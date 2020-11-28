package main

import (
	"fmt"
	"time"
)


var (
	sema =make(chan struct{}, 1)
	balance int
)

func Deposit( amount int)  {
	sema <- struct{}{} //acquire token
	balance=balance+amount
	<-sema //realese token
}

func Balance()int  {
	
	sema <- struct{}{} //acquire token
	defer func ()  {
		<-sema //realese token
	}()
	return balance
}




func main()  {
	done := make(chan struct{}, 0)

	//alice

	go func ()  {
		Deposit(200)
		fmt.Println(Balance())
		done<-struct{}{}
	}()

	go func ()  {
		time.Sleep(1*time.Millisecond)
		Deposit(100)
		done<-struct{}{}
	}()

  <-done
  <-done

  fmt.Println(Balance())
}