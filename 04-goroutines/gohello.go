package main

import (
	"fmt"
)

func hello() {
	fmt.Println("Hello, Gophers!")
}

func main() {
	go hello()
}
