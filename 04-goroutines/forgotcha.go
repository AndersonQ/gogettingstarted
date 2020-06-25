package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 1; i++ {
		go func() { fmt.Println(i) }()
	}

	time.Sleep(time.Millisecond)
}
