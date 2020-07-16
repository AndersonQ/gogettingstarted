package main

import (
	"fmt"
	"time"
)

// start_p OMIT
func produce(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
		time.Sleep(420 * time.Millisecond)
	}
	close(ch) // So the consumer know nothing else will be produced // HL
}

// end_p OMIT

// start_c OMIT
func consumer(id int, ch <-chan int) {
	// only finishes the loop when the channel is drained and closed // HL
	for v := range ch {
		fmt.Printf("[%d] received: %d\n", id, v)
	}
}

// end_c OMIT

// start_main OMIT
func main() {
	ch := make(chan int)

	for i := 0; i < 5; i++ {
		go consumer(i, ch)
	}

	produce(ch)

	fmt.Println("done :)")
}

// end_main OMIT
