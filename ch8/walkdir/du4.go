package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"time"
	"path/filepath"
	"sync"
)


var verbose = flag.Bool("v",false,"show verbose progress message ")

var signal = make(chan struct{})


//sema is a counting semaphore to limit concurrency dirents


func isCanceld() bool {
	select {
	case <-signal:
		return true
	default:
		return false
	}
}



func main()  {
	
	//Determine initial dirs

	flag.Parse()
	roots:=flag.Args()
	if len(roots)==0 {
		roots=[]string{"."}
	}

	

	go func ()  {
		os.Stdin.Read(make([]byte, 1))
		close(signal) //broadcasting the event: cancel
	}()

	fileSizes := make(chan int64, 0)
	var wg  sync.WaitGroup
	for _, root := range roots {
		wg.Add(1)
		go walkDir(root,fileSizes,&wg)  //pass address : wait group is  noCopy 
	}

	/*
	since fileSizes is a unbufferd channel , the second goroutine sending to this chan
	 will block and waing for the receiver to take out the previous sent value to the chan, 
	 meanwhile receiver goroutine( main) waiting for the others to finish and close chan:
	wg.Wait()
	close(fileSizes)

	*/

	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	var nfiles , nbytes int64

	var tick <-chan time.Time

	if *verbose {
		tick =time.Tick(500*time.Millisecond)
	}

loop:
	for{
		select {
		case <-signal:
			// canceld , then drain the channel fileSizes to allow existing /blocked goroutines to finish
			for range fileSizes {
				//do nothing
			}
		case size,ok:=<-fileSizes:
			if !ok {
				break	loop //fileSizes was closed.
			}
			nfiles++
			nbytes+=size

		case <-tick:
			printDiskUsage(nfiles,nbytes)
		}
		
	}

	printDiskUsage(nfiles,nbytes)
}

func printDiskUsage(nfiles,nbytes int64)  {
	fmt.Printf("%d files  %8.1f GB\n",nfiles,float64(nbytes)/1e9)
}




func walkDir(dir string, fileSizes chan<- int64,wg *sync.WaitGroup)   {
	defer wg.Done()
	if isCanceld(){
		return
	}
	for _, entry := range dirents(dir) {
		if isCanceld(){
			break
		}
		if entry.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dir,entry.Name())
			go walkDir(subdir,fileSizes,wg)
		}else{
			fileSizes<-entry.Size()
		}
	}
}

var sema =make(chan struct{}, 20)

// dirents returns the entries of the directory , and limits the opeing files at once to 20
func dirents(dir string) []os.FileInfo {
	sema <-struct{}{}  // acquire token 
	defer func (){   // release token
		<-sema
	}() 
	if isCanceld() {
		return nil
	}
	entries,err:= ioutil.ReadDir(dir)
	if err!=nil {
		fmt.Fprintf(os.Stderr,"du1:%v\n",err)
		return nil
	}
	return entries
}