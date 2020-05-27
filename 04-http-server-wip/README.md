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
