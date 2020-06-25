package main

import "fmt"

func main() {
	var ch chan int
	ch = make(chan int, 1)

	ch <- 42 // sending does not block // HL

	fmt.Println(<-ch)
}
