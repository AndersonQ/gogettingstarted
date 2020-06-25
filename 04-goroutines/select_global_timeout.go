package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	// start_declarations OMIT
	c1 := make(chan string)
	c2 := make(chan []float64)
	c42 := make(chan int)

	go func() {
		for {
			time.Sleep(randDurationMs(700, 1000))
			c1 <- "sent to c1"
		}
	}()

	go func() {
		for {
			time.Sleep(randDurationMs(700, 1000))
			c2 <- []float64{3.14, 1.618}
		}
	}()

	go func() {
		for {
			time.Sleep(randDurationMs(700, 1000))
			fmt.Println("received from c42:", <-c42)
		}
	}()
	// end_declarations OMIT

	// start_select OMIT
	timeout := time.After(3 * time.Second) // HL
	for i := 0; ; i++ {
		fmt.Printf("[%d]: ", i)
		select {
		case v1 := <-c1:
			fmt.Println("received from c1:", v1)
		case v2 := <-c2:
			fmt.Println("received from c2:", v2)
		case c42 <- 42:
			fmt.Println("sent 42 to c42")
		case <-timeout: // HL
			fmt.Println("too slow, quiting")
			return
		}
	}
	// end_select OMIT
}

func randDurationMs(min, max int) time.Duration {
	return time.Duration(rand.Intn(max-min+1)+min) * time.Millisecond
}
