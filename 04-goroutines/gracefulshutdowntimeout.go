package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// start_main OMIT
func main() {
	sigChan := make(chan os.Signal, 1) // we don't want the sender to block

	// use signal.Notify(sigChan, os.Interrupt) on a real application. signal.Notify
	// does not block when relaying the signal, so the channel must have enough buffer
	exit := make(chan error)

	go fakeCtrlC(sigChan)

	go gracefulShutdown(sigChan, exit, 200*time.Millisecond)

	go doCleverThings()
	fmt.Println("application started to do clever things")

	fmt.Println("bye o/:", <-exit)
}

// end_main OMIT

// start_gracefulShutdown OMIT
func gracefulShutdown(exitSig <-chan os.Signal, done chan error, timeout time.Duration) {
	// start_shutdown1 OMIT
	sig := <-exitSig // same as before // HL
	fmt.Printf("received signal: %q, starting graceful shutdown...\n",
		sig.String())

	// start timeout countdown // HL
	ctx, cancel := context.WithTimeout(context.Background(), timeout) // HL
	defer cancel()
	// end_shutdown1 OMIT

	// start_shutdown_timeout OMIT
	go func() {
		<-ctx.Done() // wait until context timeout or it's cancelled // HL
		if err := ctx.Err(); err != nil {
			done <- fmt.Errorf("graceful shutdown failed: %w", err) // HL
		}
	}()
	// end_shutdown_timeout OMIT

	// start_shutdown_wg OMIT
	// tl;dr: a WaitGroup is a counter which can increase and decrease. // HL
	// Calling WaitGroup.Wait() will block until the counter is ZERO // HL
	wg := &sync.WaitGroup{}

	wg.Add(4) // adds 4 to the counter // HL
	go closeDependency(ctx, 0, wg)
	go closeDependency(ctx, 1, wg)
	go closeDependency(ctx, 2, wg)
	go closeDependency(ctx, 3, wg)

	wg.Wait() // Wait() will return when WaitGroup counter == 0 // HL
	// end_shutdown_wg OMIT
	fmt.Println("shutdown complete")
	done <- nil
}

// end_gracefulShutdown OMIT

// start_closeDependency OMIT
func closeDependency(ctx context.Context, n int, wg *sync.WaitGroup) {
	time.Sleep((time.Duration(rand.Intn(10) * 25)) * time.Millisecond)
	fmt.Printf("[%d] closed\n", n)
	wg.Done() // decrements the WaitGroup counter by one, same as wg.Add(-1) // HL
}

// end_closeDependency OMIT

func fakeCtrlC(sigChan chan os.Signal) {
	time.Sleep(time.Second)
	sigChan <- os.Interrupt
}

func doCleverThings() {}
