package main

import (
	"fmt"
)

// start OMIT
func main() {
	var rch <-chan int
	var sch chan<- int

	ch := make(chan int, 1)

	rch = ch
	sch = ch

	sch <- 42
	fmt.Println(<-rch)

	ch <- 43
	fmt.Println(<-rch)
}

// end OMIT
