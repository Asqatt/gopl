package main


import (
	"net"
	"log"
	"bufio"
	"fmt"
)

type client chan<-string // an outgoing message channel

var(
	entering =make (chan client, 0)
	leaving = make (chan client, 0)
	messages= make(chan string, 0)
)

func main()  {
	listener ,err := net.Listen("tcp","localhost:8000")
	if err!=nil {
		log.Fatal(err)
	}

	go broadcaster()

	for{
		conn,err:=listener.Accept()
		if err!=nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func clientWriter(conn net.Conn, msgCh <-chan string)  {
	for{
		msg:=<-msgCh
		fmt.Fprintf(conn,msg)
	}
}


func handleConn(conn net.Conn)  {
	ch := make(chan string)
	go clientWriter(conn,ch)

	who:=conn.RemoteAddr().String()

	ch<-"You are "+who

	messages<-who+" has arrived!\t\n"

	entering<-ch

	input:=bufio.NewScanner(conn)

	for input.Scan(){
		messages<-who+": "+input.Text()
	}
	leaving<-ch

	messages<-who +" has left."

	conn.Close()

}

func broadcaster()  {
	clients:=make(map[client]bool) // all connected clients

	for{
		select {
		case msg:=<-messages:
			//broadcast to all
			//clients' outgoing message channels

			for cli := range clients{
				cli<-msg
			}
		case cli:=<-entering:
			clients[cli]=true
		case cli:=<-leaving:
			delete(clients,cli)
			close(cli)
		}
	}

}