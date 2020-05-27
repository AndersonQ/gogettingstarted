# Handling HTTP requests

A `http.Handler` is the interface which will receives a request, process it and returns, better saying, _writes_ 
a response. Before diving into the `http.Handler` type, let's have a look at the http request and response, they are
`http.Request` and `http.ResponseWriter` respectively.

## `http.Request`

It is the http incoming or outcoming request. It's a struct, and some of its fields are shown below. Check the documentation
(a.k.a. the source code) for the formal definition, here I'll summarise what we need for now.

```go
type Request struct {
	// Method is the HTTP method (GET, POST, PUT, etc.).
	// There are constants in the http package for the http methods and status. Thus use them instead of typing them in.
	Method string

	// URL is, well, the request's URL. We can call String() on a *url.URL and get it as a string
	URL *url.URL

	// Header represents request headers, it's build on top of a `map[string][]string` and has helper methods such as
	// Get and Set. These methods will canonicalize the header key, so prefer to use them instead of accessing the map
	// directely.
	Header http.Header

	// Body is the request body, it's important to remember that for requests received by the server there is no need to
	// close the body after to have read it.
	Body io.ReadCloser
}

// Context returns the request's context. To change the request's context, use WithContext.
func (r *Request) Context() context.Context

// WithContext returns a copy of the request with the new context.
func (r *Request) WithContext(ctx context.Context) *http.Request
```

## `http.ResponseWriter`

It's an interface, different from the request, and again I'll summarise what we need for now, check the documentation
for the proper definitions.

```go
type ResponseWriter interface {
	// Header represents the response headers. It's important to keep in mind the header map should not be changed after
	// either Write or WriteHeader have been called.
	Header() http.Header

	// Write writes the response back, if WriteHeader has not been called it'll set the response status to http.StatusOK.
	// If the Content-Type header is not set, Go will set it, trying to guess the right one.
	// The Content-Length header is added automatically.
	Write([]byte) (int, error)

	// WriteHeader sets the response status code and it must be a valid HTTP 1xx-5xx status code. Use the constants
	// defined in the http package instead of typing them in.
	WriteHeader(statusCode int)
}
```

## Handling a request

Given we have an incoming request, and we can write a response, we can write a function which receives a request and 
writes a response. Let's define `func(w http.ResponseWriter, r *http.Request)` as the signature of our function. As you
might have noticed it does return a value. We have to call `Write` on `w` to send the response and, ~ideally~
after to have handled the error `Write` might return, finish the function.

Now let's implement a handler which logs the request http method, URL, and headers, reply with a status code `418` and a
body `{"hello":"world"}`. On [handlers.go](handlers.go) there is the skeleton of a handler and on
[handlers_test.go](handlers_test.go) there is the test for it. 
For simplicity, we'll only assert the response's status code and body.

## The `http.Handler`

We have defined our http handler as any function with the following signature `func(w http.ResponseWriter, r *http.Request)`.
Let's be honest, it's neither convenient nor quite intuitive. In Go the `http.Handler` interface is:

```go
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
```

Any type with a method `ServeHTTP(http.ResponseWriter, *http.Request)` automatically implements
`http.Handler`. Components dealing with http handlers receive or return a `http.Handler`.
However, it isn't exactly convenient to create a type just to implement a one interface.

Go have you covered and brings a helper type, the `http.HandlerFunc`, as its documentation
states, _is an adapter to allow the use of ordinary functions as HTTP handlers_.
Using `http.HandlerFunc` we can make any function `func(http.ResponseWriter, *http.Request)`
to implement `http.Handler`. Let's see some code.

```go
type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`{"hello":"handler.ServeHTTP"}`)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`{"hello":"handlerFunc"}`)
}

func main() {
	var h1, h2, h3 http.Handler
	
	h1 = handler{} // it works, the type handler implements http.Handler
	
	h2 = handlerFunc // it does not work, the function handlerFunc does not implements http.Handler
	
	h2 = http.HandlerFunc(handlerFunc) // now it works!

	h3 = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(`{"hello":"anonymous func"}`)
	}) // it also works!
}
```
You can run `go run demo/00/main.go` to see it in action.
A question arises, how does `http.HandlerFunc` does its magic? First let's see what it is:

```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
```

`http.HandlerFunc` type is a function receiving two parameters, a `http.ResponseWriter`
and `*http.Request`, and does not return any value. It is possible to cast any 
function with the same signature to `http.HandlerFunc`, as our `handlerFunc` has 
this signature, we can cast it to `http.HandlerFunc`.

Now let's build our own `http.HandlerFunc` from our `handler`, to fully 
understand the trick. We had:

```go
type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`{"hello":"handler.ServeHTTP"}`)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`{"hello":"handlerFunc"}`)
}
```

a wee tweak, and the following work just as good as our old `handler`:

```go
type handler func(http.ResponseWriter, *http.Request)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`{"hello":"handler.ServeHTTP"}`)
}
```

Now `handler` is a function, and any function with the same signature can be cast
to `handler`:

```go
    // Both work, but they do not produce the same result
	h2 = http.HandlerFunc(handlerFunc)
	h3 = handler(handlerFunc)

	// Output:
	// {"hello":"handlerFunc"}
	// {"hello":"handler.ServeHTTP"}
```

Run `go run demo/01/main.go` to see for yourself.

Why? Simple, we are ignoring our receiver! We need to invoke the function 
represented by our receiver `h`. Let's fix it.


```go
type handler func(http.ResponseWriter, *http.Request)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}
```

Finally, now our `handler` works just as `http.Handler` does!

```go
	var h2, h3 http.Handler

	h2 = http.HandlerFunc(handlerFunc)
	h3 = handler(handlerFunc)

	handlerFunc(nil, nil)
	h2.ServeHTTP(nil, nil)
	h3.ServeHTTP(nil, nil)

	// Output:
	// {"hello":"handlerFunc"}
	// {"hello":"handlerFunc"}
	// {"hello":"handlerFunc"}
```

Try running `go run demo/02/main.go` to see it in action.
