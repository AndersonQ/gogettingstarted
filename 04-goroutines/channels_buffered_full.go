package main

import "fmt"

func main() {
	var ch chan int
	ch = make(chan int, 1)

	ch <- 42
	ch <- 43 // the channels is full, now it blocks // HL

	fmt.Println(<-ch)
}
