package main

import (
	"fmt"
	"time"
)

// start_observer OMIT

type Observer interface {
	Observe(int)
}

type ObserverFunc func(int)

func (f ObserverFunc) Observe(s int) {
	f(s)
}

// end_observer OMIT

// start_observable OMIT

type observable struct {
	observers []Observer
}

func (o *observable) Subscribe(observer Observer) {
	o.observers = append(o.observers, observer)
}

func (o *observable) notify(msg int) {
	for _, obs := range o.observers {
		go obs.Observe(msg) // calling the observers on a different goroutine // HL
	}
}

// end_observable OMIT

// start OMIT
func main() {
	ch := make(chan int)
	go produce(ch)

	some := &observable{}

	for i := 0; i < 3; i++ {
		i := i
		some.Subscribe(ObserverFunc(func(n int) {
			fmt.Printf("[%d]: %d\n", i, n)
		}))
	}

	for msg := range ch {
		some.notify(msg)
	}
}

// end OMIT

func produce(ch chan int) {
	var i int
	for {
		fmt.Println("producing: ", i)
		ch <- i
		i++
		time.Sleep(300 * time.Millisecond)
	}
}
