package main

import (
	"fmt"
	"net/http"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`{"hello":"handler.ServeHTTP"}`)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print(`{"hello":"handlerFunc"}`, "\n")
}

func main() {
	var h1, h2, h3 http.Handler

	h1 = handler{} // it works, the type handler implements http.Handler

	// h2 = handlerFunc // it does not work, the function handlerFunc does not implements http.Handler

	h2 = http.HandlerFunc(handlerFunc) // now it works!

	h3 = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(`{"hello":"anonymous func"}`, "\n")
	}) // it also works!

	h1.ServeHTTP(nil, nil)
	h2.ServeHTTP(nil, nil)
	h3.ServeHTTP(nil, nil)

	// Output:
	// {"hello":"handler.ServeHTTP"}
	// {"hello":"handlerFunc"}
	// {"hello":"anonymous func"}
}
