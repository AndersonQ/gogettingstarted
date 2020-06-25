package main

import (
	"fmt"
	"time"
)

// start_main OMIT
func main() {
	var ch chan int
	ch = make(chan int)

	go func() {
		fmt.Println("sending 42")
		ch <- 42
		fmt.Println("sent 42")
	}()

	go func() {
		fmt.Println("sending 43")
		ch <- 43
		fmt.Println("sent 43")
	}()

	time.Sleep(time.Second)
	fmt.Println("received:", <-ch)
}

// end_main OMIT
