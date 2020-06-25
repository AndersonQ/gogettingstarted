package main

import (
	"fmt"
)

// start_main OMIT
func main() {
	ch := make(chan int)

	fmt.Println(<-ch)
	fmt.Println("done")
}

// end_main OMIT
