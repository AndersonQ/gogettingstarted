# HTTP handler and middlewares

## Handling http requests

A `http.Handler` is the interface which will receives a request, process it and returns, better saying, _writes_ 
a response. Before diving into the `http.Handler` type, let's have a look at the http request and response, they are
`http.Request` and `http.ResponseWriter` respectively.

### `http.Request`

It is the incoming or outcoming request. It's a struct, and some of its fields are shown below. Check the documentation
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

### `http.ResponseWriter`

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

### Handling a request

Given we have an incoming request, and we can write a response, we can write a function which receives a request and 
writes a response. Let's define `func(w http.ResponseWriter, r *http.Request)` as the signature of our function. As you
might have noticed it does not have a return value. We have to call `Write` on `w` to send the response and, ~ideally~
after to have handled the error `Write` might return, finish the function.

Now let's implement a handler which logs the request http method, URL, and headers, reply with a status code `418` and a
body `{"hello":"world"}`. On [handlers_test.go](handlers_test.go) there is the skeleton of a test for our handler. 
For simplicity, we'll only assert the response's status code and body.


## Middlewares

A loose definition for a http middleware is some function which sits before or after your handler being invoked, which 
is able to access the request and the response. Bear in mind only one call to the `ResponseWriter.Write` is allowed
(there are exceptions, but let's keep it simple), therefore a middleware should not call `Write`. There are ways
to bypass it, but again, let's keep it simple. A visual representation of a middleware is:

```text
incoming request => middleware => your handler => middleware.
```

Let's build a middleware step by step. Our middleware will set the content type of the response to `application/json`,
and we'll call it `JSONResponse`. In order to receive a request our middleware must look like a handler, so it needs to
have the same signature as a http handler.

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    // do stuff...
}

func JSONResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	Handler(w, r)
}
```

It solves our problem, however it isn't quite flexible. Ideally our middleware is able to work with any handler, so it
needs to receive the handler as well. If we receive a handler as a parameter we'd have:
```go
func JSONResponse(w http.ResponseWriter, r *http.Request, h func(w http.ResponseWriter, r *http.Request)) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	h(w, r)
}
```

Nice, it works with any handler! However now it cannot replace a http handler as the signature is different. The main
requirement is we need to end up with a `func(w http.ResponseWriter, r *http.Request)` in our hands. Thus, we need a
function which receives a http handler and returns another http handler, what lead us to the following signature:
```go
func( func(w http.ResponseWriter, r *http.Request) ) func(w http.ResponseWriter, r *http.Request)
```

It might look a bit scary, or just a way to many `func` in the same line, but it's simple. We have a function which
receives a function as its only parameter and returns another function. Now head to [middlewares.go](middlewares.go)
and implement `TrackingID`, a middleware which reads the http header `X-TrackingId` and adds the tracking id to the context. For now,
you don't need to worry about working with contexts, on [context.go](context.go) you'll find 
`ContextWithID(id string) context.Context` which will take care of dealing with the context for you. On 
[middlewares_test.go](middlewares_test.go) there is a test for the middleware you'll build.


## The `http.Handler`

We have defined our http handler as any function with the following signature `func(w http.ResponseWriter, r *http.Request)`.
Let's be honest, it's neither convenient nor quite intuitive. In Go a http handler is:

```go
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
```

anything with a method `ServeHTTP(http.ResponseWriter, *http.Request)` automatically implements
`http.Handler`. Whatever id dealing with http handlers will receive or return `http.Handler`.
However, it isn't the most convenient to create a type just to implement one method.
Go have you covered and brings a helper type, the `http.HandlerFunc`, as its documentation
states, _is an adapter to allow the use of ordinary functions as HTTP handlers_.
Using `http.HandlerFunc` we can make any function `func(http.ResponseWriter, *http.Request)`
to implement `http.Handler`. Let's see some code.

```go
type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"hello":"world"}`))
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"hello":"world"}`))
}

func main() {
	var h1, h2, h3 http.Handler
	
	h1 = handler{} // it works, the type handler implements http.Handler
	
	h2 = handlerFunc // it does not work, the function handlerFunc does not implements http.Handler
	
	h2 = http.HandlerFunc(handlerFunc) // now it works!

	h3 = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"hello":"world"}`))
	}) // it also works!
}
```

How does `http.HandlerFunc` does its magic? First let's see what it is:

```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
```

`http.HandlerFunc` is a function which receives a `http.ResponseWriter` and `*http.Request`,
and does not return any value. It is possible to cast any function with the same
signature to `http.HandlerFunc`, as our `handlerFunc` has this signature, we can
cast it to `http.HandlerFunc` as well.

Now let's build our own `http.HandlerFunc` from our `handler`, we had:

```go
type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"hello":"world"}`))
}
```

a wee tweak, and the following work just as good as our old `handler`:

```go
type handler func(http.ResponseWriter, *http.Request)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"hello":"world"}`))
}
```

We are ignoring our receiver in the snippet above, as `h handler` is a
`func(http.ResponseWriter, *http.Request)` we can do

```go
type handler func(http.ResponseWriter, *http.Request)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}
```

Finally, now any instance of `handler` is a `http.Handler`

```go

```

Also, a http middleware signature's, using it, will be `func (http.Handler) http.Handler`.
Which is definitely a lot better to read and reason about.




TODO: middleware like func (...) func(http.Handler) http.Handler
Another common need when defining a middleware is to pass some parameters to the middleware,



### (WIP) http.Server

Go provides out of the box a production ready and easy to use http server, 
the [`http.Server`](https://golang.org/pkg/net/http/#Server). Tl;dr it's a struct, we'll
focus on how to add a handler for requests, a few of its fields to configure the server and on 
`func (*Server) ListenAndServe` which starts the http server.

Start the http server is as easy as:

```go
server := http.Server{}
server.ListenAndServe()
```

It doesn't do much right now, obviously we are missing a few things, define some routes and their handlers, 
set a timeout and handle the error returned by `ListenAndServe`.

First the error, let's keep it simple

```go 
server := http.Server{}
if err := server.ListenAndServe(); err!= nil {
    panic("ListenAndServe returne")
}
```

We'll set the `ReadTimeout` and `WriteTimeout`. The first is the maximum time for the server to read the request,
including the body. The latter is the maximum time starting after the request headers have been read up to the end
of the write of the response. Cloudflare's blog has got a 
[nice article](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/) about the Go's http timeouts.

```go
server := http.Server{
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 5 * time.Second,
}
if err := server.ListenAndServe(); err != nil {
    panic("ListenAndServe returned: %v", err)
}
```
