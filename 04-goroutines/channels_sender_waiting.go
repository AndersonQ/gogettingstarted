package main

import (
	"fmt"
	"time"
)

// start_main OMIT
func main() {
	ch := make(chan int)

	go func() {
		fmt.Println("goroutine: waiting to send...")
		ch <- 42
		fmt.Println("goroutine: data sent")
	}()

	time.Sleep(time.Millisecond) // blocking this goroutine so the other can run
	fmt.Println("done")
}

// end_main OMIT
