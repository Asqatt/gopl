package main

import (
	"fmt"
	"time"
	"sync"
)


var (

	//RWMutex is multiple readers , single writer lock
	mutex sync.RWMutex
	balance int
)


//basic un-exported deposit function
func deposit(amount int)  {
	balance=balance+amount
}




func Withdrawal(amount int) bool {
	mutex.Lock()
	defer mutex.Unlock()
	deposit(-amount)
	if balance<0 {
		deposit(amount)
		return false
	}
	return true
}


func Deposit( amount int)  {

	//acquire exclusive/writer lock
	mutex.Lock()
	deposit(amount)
	mutex.Unlock()
}

func Balance()int  {
	//acquire reader lock ,
	mutex.RLock()
	defer mutex.RUnlock()
	return balance
}




func main()  {
	done := make(chan struct{}, 0) //signal main go routine to finish

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

	go func ()  {
		b:=Withdrawal(500)
		fmt.Println(b)
		done<-struct{}{}
	}()
	go func ()  {
		time.Sleep(4*time.Millisecond)
		b:=Withdrawal(50)
		fmt.Println(b)
		done<-struct{}{}
	}()


	//waiting for worker go routines
  <-done
  <-done
  <-done
  <-done


  fmt.Println(Balance())
}