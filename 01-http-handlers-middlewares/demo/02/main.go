package main

import (
	"fmt"
	"net/http"
)

type handler func(http.ResponseWriter, *http.Request)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`{"hello":"handlerFunc"}`)
}

func main() {
	var h2, h3 http.Handler

	h2 = http.HandlerFunc(handlerFunc)
	h3 = handler(handlerFunc)

	handlerFunc(nil, nil)
	h2.ServeHTTP(nil, nil)
	h3.ServeHTTP(nil, nil)

	// Output:
	// {"hello":"handlerFunc"}
	// {"hello":"handlerFunc"}
	// {"hello":""handlerFunc"}
}
