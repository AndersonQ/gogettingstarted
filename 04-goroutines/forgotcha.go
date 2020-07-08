package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		i := i // HL
		go func() { fmt.Println(i) }()
	}

	time.Sleep(time.Millisecond)
}
