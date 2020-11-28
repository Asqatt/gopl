package main

import (
	"fmt"
	"time"
	"sync"
)


var (
	mutex sync.Mutex
	balance int
)

func Deposit( amount int)  {
	mutex.Lock()
	balance=balance+amount
	mutex.Unlock()
}

func Balance()int  {
	
	mutex.Lock()
	defer mutex.Unlock()
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