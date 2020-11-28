package du1

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"time"
	"path/filepath"
)


func walkDir(dir string, fileSizes chan<- int64)   {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir,entry.Name())
			walkDir(subdir,fileSizes)
		}else{
			fileSizes<-entry.Size()
		}
	}
}


// dirents returns the entries of the directory
func dirents(dir string) []os.FileInfo {
	entries,err:= ioutil.ReadDir(dir)
	if err!=nil {
		fmt.Fprintf(os.Stderr,"du1:%v\n",err)
		return nil
	}
	return entries
}


var verbose = flag.Bool("v",false,"show verbose progress message ")

func main()  {
	
	//Determine initial dirs

	flag.Parse()

	

	roots:=flag.Args()

	

	if len(roots)==0 {
		roots=[]string{"."}
	}

	fileSizes := make(chan int64, 0)

	go func() {
		for _, root := range roots {
			walkDir(root,fileSizes)
		}
		close(fileSizes)
	}()
	var nfiles , nbytes int64

	var tick <-chan time.Time

	if *verbose {
		tick =time.Tick(500*time.Millisecond)
	}

	// for size:=range fileSizes{
	// 	nfiles++
	// 	nbytes+=size
	// }
loop:
	for{
		select {
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