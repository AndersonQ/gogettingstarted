package main

import "fmt"

// start_main OMIT
func main() {
	buffSize := 10
	ch := make(chan int, buffSize)

	for i := 0; i < buffSize; i++ {
		ch <- 42 + i
	}

	for v := range ch {
		fmt.Println(v)
	}
}

// end_main OMIT
