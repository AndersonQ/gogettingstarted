package main

import (
	"fmt"
	"math/rand"
	"time"
)

// start OMIT
func main() {
	done := make(chan struct{})

	go run(done)

	<-done

	fmt.Println("we're done :)")
}

func run(done chan struct{}) {
	for i := 0; i < 5; i++ {
		time.Sleep(randDurationMs(300, 800))
		fmt.Println(i)
	}

	close(done)
	return
}

// end OMIT

func randDurationMs(min, max int) time.Duration {
	return time.Duration(rand.Intn(max-min+1)+min) * time.Millisecond
}
