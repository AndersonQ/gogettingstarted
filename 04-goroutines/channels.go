package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	// start_Sleep_OMIT
	go func() {
		time.Sleep(time.Second) // HLsleep
		ch <- 42                // HL
	}()
	// end_Sleep_OMIT

	fmt.Printf(
		"The answer to the Ultimate Question of Life, \n" +
			"the Universe, and Everything is ")
	fmt.Print(<-ch) // HL
}
