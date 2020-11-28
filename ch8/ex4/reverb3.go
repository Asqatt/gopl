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

	var wg  sync.WaitGroup

	input:=bufio.NewScanner(c)
	for input.Scan(){
		wg.Add(1)
		go echo(c,input.Text(),1*time.Second,&wg)
	}
	wg.Wait()
	if tcpConn,ok:=c.(*net.TCPConn); ok{
		tcpConn.CloseWrite()
	}
}