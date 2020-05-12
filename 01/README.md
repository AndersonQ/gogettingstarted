# HTTP handler and middlewares

Go provides out of the box a production ready and easy to use http server, 
the [`http.Server`](https://golang.org/pkg/net/http/#Server). Tl;dr it's a struct, we'll
focus on how to add a handler for requests, a few of its fields to configure the server and on 
`func (*Server) ListenAndServe` which starts the http server.

### `http.Handler`

`http.Handler` is the entity which will receives a request, process it and returns, better saying, _writes_ 
a response. Before diving into the `http.Handler` type, let's have a look at the request and response, they are
`http.Request` and `http.ResponseWriter` respectively.

#### `http.Request`

It is the incoming or outcoming request. It's a struct, and some of its fields are shown below. Check the documentation
(a.k.a. the source code) for the formal definition, here I'll summarise what we need for now.

```go
type Request struct {
	// Method is the HTTP method (GET, POST, PUT, etc.).
	// there are constants in the http package for the http methods and status. Thus use them instead of type then in.
	Method string
	
    // URL is, well, the request's URL. We can call String() on a *url.URL and get it as, well, string
	URL *url.URL
	
    // Header represents request headers, it's build on top of a `map[string][]string` and has helper methods such as Get and Set.
	Header http.Header
	
    // Body is the request body, it's important to remember that for requests received by the server there is no need to close
	// the body after to have read it.
	Body io.ReadCloser
}

// Context returns the request's context. To change the request's context, use WithContext.
func (r *Request) Context() context.Context

// WithContext returns a copy of the request with the new context.
func (r *Request) WithContext(ctx context.Context) *http.Request
```

#### `http.ResponseWriter`

It's an interface, different from the request, and again I'll summarise what we need for now, check the documentation
for the proper definition.

```go
type ResponseWriter interface {
    // Header represents the response headers. It's important to keep in mind the header map should not be changed after
    // either Write or WriteHeader have been called.
	Header() http.Header
    
    // Write writes the response back, if WriteHeader has not been called it'll set the response status to http.StatusOK.
    // If the Content-Type header is not set, Go will set it, trying to guess the right one.
    // The Content-Length header is added automatically.
	Write([]byte) (int, error)

    // WriteHeader sets the response status code and it must be a valid HTTP 1xx-5xx status code.
	WriteHeader(statusCode int)
}
```

#### Handling a request

Given we have an incoming request, and we can write a response, we can write a function which receives a request and 
writes a response. Let's define `func(w http.ResponseWriter, r *http.Request)` as the signature of our function. As you
might have noticed it does not have a return value. We have to call `Write` on `w` to send the response and ~ideally~
after to have handled the error `Write` might return, finish the function.

Now let's implement a handler which logs the request http method, URL, and headers, reply with a status code `418` and a
body `{"hello":"world"}`. On [handlers_test.go](handlers_test.go) there is the skeleton of a test for our handler. 
For simplicity, we'll only assert the response's status code and body.

### (WIP) http.Server
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