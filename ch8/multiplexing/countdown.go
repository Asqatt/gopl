package main

import (
	"fmt"
	"time"
	"os"
)

func main()  {
	fmt.Println("Commencing countdown,press return key to abort.")

	// tick := time.Tick(1*time.Second)
	ticker := time.NewTicker(1*time.Second)

	abort := make(chan struct{})

	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort<-struct{}{}
	}()

	for countdown := 10; countdown >0; countdown-- {
		select {
		case <-ticker.C:
			fmt.Println(countdown)
		case <-abort:
			fmt.Println("Launch aborted! ")
			ticker.Stop() //avoid go routine leak
			return
		}
	}

	launch()
}

func launch()  {
	fmt.Println("Lauched!")
}