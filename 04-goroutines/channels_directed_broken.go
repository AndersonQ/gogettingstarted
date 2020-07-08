package main

// start OMIT
func main() {
	var rch <-chan int
	var sch chan<- int

	ch := make(chan int, 1)

	rch = ch
	sch = ch

	rch <- 42
	<-sch
}

// end OMIT
