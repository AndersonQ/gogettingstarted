package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		fmt.Println("Hello, Gophers!")
	}() // HL

	time.Sleep(time.Millisecond)
}
