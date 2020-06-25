package main

import (
	"fmt"
	"time"
)

func hello() {
	fmt.Println("Hello, Gophers!")
}

func main() {
	go hello()
	time.Sleep(time.Millisecond)
}
