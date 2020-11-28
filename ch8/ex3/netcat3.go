package main

import (
	"os"
	"net"
	"log"
	"io"
)

func main()  {

	addr, err := net.ResolveTCPAddr("tcp", "localhost:8000")
	log.Println(addr)
	if err != nil {
		log.Fatal(err)
	}

	conn,err := net.DialTCP("tcp",nil,addr)
	if err!=nil {
		log.Fatal(err)
	}

	done:= make(chan struct{})
	go func() {
		io.Copy(os.Stdout,conn)
		log.Println("done")
		done<-struct{}{}
		}()
	mustCopy(conn,os.Stdin)
	conn.CloseWrite()
	// if tcpConn,ok:=conn.(*net.TCPConn);ok {
	// 	tcpConn.CloseWrite()
	// }
	<-done
}

func mustCopy(dst io.Writer,src io.Reader)  {
	if _,err:= io.Copy(dst,src);err!=nil {
		log.Fatal(err)
	}
}