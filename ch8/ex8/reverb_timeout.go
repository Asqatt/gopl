package main

import(
	"bufio"
	"log"
	"net"
	"time"
	"fmt"
	"strings"
	"sync"
)

func main(){
	listener,err := net.Listen("tcp","localhost:8000")

	if err!=nil {
		log.Fatal(err)
	}
	for{
		conn,err :=listener.Accept()
		if err!=nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func echo(c net.Conn, shout string, delay time.Duration,wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Fprintln(c,"\t",strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c,"\t",shout)
	time.Sleep(delay)
	fmt.Fprintln(c,"\t",strings.ToLower(shout))
}

func handleConn(c net.Conn)  {

	ticker := time.NewTicker(1*time.Second)

	shouts:= make(chan string, 0)
	input:=bufio.NewScanner(c)
	go func() {
		for input.Scan(){
			shouts<-input.Text()
		}
	}()

	var countdown int


	var wg  sync.WaitGroup

	for countdown=10;countdown>0;countdown--{

		select {
		case txt:= <-shouts:
			countdown=10
			wg.Add(1)
			go echo(c,txt,1*time.Second,&wg)
		case <-ticker.C:
		}
	}

	fmt.Fprintln(c,"Connection closed: idle too long!")
	
	ticker.Stop()
	wg.Wait()
	if tcpConn,ok:=c.(*net.TCPConn); ok{
		tcpConn.CloseWrite()
	}
}