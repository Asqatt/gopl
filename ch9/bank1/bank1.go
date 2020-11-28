package bank1

import (
	"fmt"
	"time"
)


var deposits = make(chan int)
var balances = make(chan int)

func deposit( amount int)  {
	deposits<-amount
}

func balance()int  {
	return  <- balances
}

func teller()  {
	var balance int  //balance is confined to teller go routine ,others can't directly write it

	for {
		
		select{
			case amount:= <-deposits:
				balance+=amount

			case balances<-balance:
		}
	}
}

func init()  {
	go teller()
}

func main()  {
	done := make(chan struct{}, 0)

	//alice

	go func ()  {
		deposit(200)
		fmt.Println(balance())
		done<-struct{}{}
	}()

	go func ()  {
		time.Sleep(20*time.Millisecond)
		deposit(100)
		done<-struct{}{}
	}()

  <-done
  <-done

  fmt.Println(balance())
}