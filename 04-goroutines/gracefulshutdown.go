package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// start_main OMIT
func main() {
	sigChan := make(chan os.Signal, 1) // we don't want the sender to block
	exit := make(chan struct{})

	go fakeCtrlC(sigChan)

	go gracefulShutdown(sigChan, exit) // HL

	go doCleverThings()
	fmt.Println("application started to do clever things")

	// block until the shutdown is complete // HL
	<-exit // HL
	fmt.Println("bye o/")
}

// end_main OMIT

// start_gracefulShutdown OMIT
func gracefulShutdown(exitSig <-chan os.Signal, done chan struct{}) {
	// It'll block here until a exit signal is received // HL
	sig := <-exitSig // HL

	fmt.Printf("received signal:\n\t%qstarting graceful shutdown...\n",
		sig.String())

	time.Sleep((time.Duration(rand.Intn(10) * 10)) * time.Millisecond)
	fmt.Println("shutdown complete")

	// signal the shutdown is complete // HL
	close(done) // HL
}

// end_gracefulShutdown OMIT

func fakeCtrlC(sigChan chan os.Signal) {
	time.Sleep(time.Second)
	sigChan <- os.Interrupt
}

func doCleverThings() {}
