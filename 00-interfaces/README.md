# Methods and Interfaces

## Methods

Lacking formality, a method is a function attached to a type and the method has access to the type instance
on which it was called. It's possible to add methods to almost any type in Go. The type which the method is attached to
is called receiver, and the receiver is, by convention, named after the first letter of the type. Let's see them to make
it clear.

Given the following function
```go
func Hello(name string) string {
    return fmt.Sprintf("Hello,", name)
}
```

We can have it on a type holding the name and just call `Hello()`

```go
type Person struct {
    Name string
}

func (p Person) Hello() string {
    return fmt.Sprintf("Hello,", p.Name)
}
```

We can think `(p Person)` as one parameter passed to the method `Hello()`, the `Person` instance in which the method
`Hello()` is invoked. In Python, we see it quite clear, a class method has got `self` as its first parameter, which 
is the object on which the method is being called. In Java there is the `this` which is not explicitly declared, but is
also the obeject in which the method is invoked. Go, I'd say, sits in between them, there is an explicit declaration of
the instance in which the method is called, but it doesn't appear in the method's parameter list.

Below both implementations of `Hello` need an instance of `Person` and yield the same result.

```go

func (p Person) Hello() string {
    return fmt.Sprintf("Hello,", p.Name)
}

func Hello(p Person) string {
    return fmt.Sprintf("Hello,", p.Name)
}
```

Looking at these two implementations of `Hello` helps to understand the reason any change a method does on its receiver
only affects the instance in which the method is invoked if, and only if, the receiver is a pointer. Bellow we see it
in practice.

```go
func (p Person) Hello() {
	p.Name = "Hello, " + p.Name
}

func Hello(p Person) {
	p.Name = "Hello, " + p.Name
}

func (p *Person) HelloP() {
	p.Name = "Hello, " + p.Name
}

func HelloP(p *Person) {
	p.Name = "Hello, " + p.Name
}

func main() {
	p1 := Person{Name: "John p1"}
	p1.Hello()
	fmt.Println(p1.Name)
	// Output: John p1

	p2 := Person{Name: "John p2"}
	Hello(p2)
	fmt.Println(p2.Name)
	// Output: John p2

	p3 := &Person{Name: "John p3"}
	p3.HelloP()
	fmt.Println(p3.Name)
	// Output: Hello, John p3

	p4 := &Person{Name: "John p4"}
	HelloP(p4)
	fmt.Println(p4.Name)
	// Output: Hello, John p4
}
```

Go is nice and let us add functions to almost any type, so we can do:

```go
type Name string

func (n Name) Hello() string {
    return fmt.Sprintf("Hello,", n)
}

func main() {
    n := Name("John") // here we're casting the string "John" to the type "Name"
    fmt.Println(n)
    // Output: John
}
```

For the `Person` struct the receiver `p` is the struct itself, the same for `Name`, the receiver `n` is the instance of
`Name` in which `Hello()` is invoked.

It might feel odd, but even functions can have methods:

```go
// Surname is a function which receives zero parameters and returns a string
type Surname func() string

// To create a function of type Surname we can:
func Doe() string {
    return "Doe"
}

func main() {
	// In the same way we cast the string "John" to Name, here we cast the `func() string` Doe to Surname
	a := Surname(Doe)
	fmt.Print(a())
	//Output: Doe

	// The same is valid for an implicit function
	d := Surname(func() string { return "Doe" })
	fmt.Print(d())
	//Output: Doe
}
```

With this all and interfaces, which is explained next, we'll be able to do some quite nice stuff.
 
## Interfaces

Having defined a method we can talk about interfaces, quoting Effective Go:
> Interfaces in Go provide a way to specify the behavior of an object: 
if something can do _this_, then it can be used _here_.

If you understand a method is a behaviour, the quote above makes a lot of sense. An Interface defines some methods and
whichever type implements these methods, implements the interface. 

Our `Hello() string` can be an interface:

```go
type Helloer interface {
    Hello() string
}
```

A side note, one method interfaces in Go are usually named after its method with a `er` suffix, such as `io.Writer`,
`io.Reader`, `driver.Execer`. Not always quite right english, but it's done anyway.

Now using the `Helloer` interface, we can say hello to `Person`, `Name` and `Surname`:

```go
func SayHello(h Helloer) {
    fmt.Println(h.Hello())
}
```

The fully functional code would be:

```go
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
```