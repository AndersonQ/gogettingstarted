package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// start_declarations OMIT
	c1 := make(chan string)
	c2 := make(chan []float64)
	c42 := make(chan int)

	go func() {
		s := rand.Intn(500)
		time.Sleep(time.Duration(s) * time.Millisecond)
		c1 <- "sent to c1"
	}()

	go func() {
		s := rand.Intn(5)
		time.Sleep(time.Duration(s) * time.Millisecond)
		c2 <- []float64{3.14, 1.618}
	}()

	go func() {
		s := rand.Intn(5)
		time.Sleep(time.Duration(s) * time.Millisecond)
		fmt.Println("received from c42:", <-c42)
	}()
	// end_declarations OMIT

	// start_select OMIT
	select {
	case v1 := <-c1:
		fmt.Println("received from c1:", v1)
	case v2 := <-c2:
		fmt.Println("received from c2:", v2)
	case c42 <- 42:
		fmt.Println("sent 42 to c42")
	default: // HL
		fmt.Println("no one is ready") // HL
	}
	// end_select OMIT
}
