package main

import (
	"fmt"
)

// start_main OMIT
func main() {
	buffSize := 10
	ch := make(chan int, buffSize)

	for i := 0; i < buffSize; i++ {
		ch <- 42 + i
	}
	close(ch) // the channel is closed, but there is data on its buffer to be consumed // HL

	for v := range ch {
		fmt.Println(v)
	}
}

// end_main OMIT
