package main

import (
	"fmt"
)

type Person struct{ Name string }

func (p Person) Hello() string { return fmt.Sprint("Hello, ", p.Name) }

type Name string

func (n Name) Hello() string { return fmt.Sprint("Hello, ", n) }

type Surname func() string

func (s Surname) Hello() string { return fmt.Sprint("Hello, ", s()) }

type Helloer interface{ Hello() string }

func SayHello(h Helloer) { fmt.Println(h.Hello()) }

func main() {
	p := Person{Name: "John Doe"}
	n := Name("John")
	s := Surname(func() string { return "Doe" })
	SayHello(p)
	SayHello(n)
	SayHello(s)
	// Output:
	// Hello, John Doe
	// Hello, John
	// Hello, Doe
}
