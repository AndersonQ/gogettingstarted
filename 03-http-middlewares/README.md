# HTTP Middleware

A loose definition for a http middleware is some function which sits before or after your handler being invoked, which 
is able to access the request and the response. Bear in mind only one call to the `ResponseWriter.Write` is allowed
(there are exceptions, but let's keep it simple), therefore a middleware should not call `Write`. There are ways
to bypass it, but again, let's keep it simple. A visual representation of a middleware is:

```text
incoming request => middleware => your handler => middleware
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
func JSONResponse(w http.ResponseWriter, r *http.Request, h http.Hnadler) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	h.ServeHTTP(w, r)
}
```

Nice, it works with any handler! However now it cannot replace a http handler as the signature is different. The main
requirement is we need to end up with a `http.Handler` in our hands. Thus, we need a
function which receives a http handler and returns another http handler, what lead us to the following signature:
```go
func( http.Hnadler ) http.Hnadler
```

Now head to [middlewares.go](middlewares.go)
and implement `TrackingID`, a middleware which reads the http header `X-TrackingId` and adds the tracking id to the context. For now,
you don't need to worry about working with contexts, on [context.go](context.go) you'll find 
`func ContextWithID(id string) context.Context` which will take care of dealing with the context for you. On 
[middlewares_test.go](middlewares_test.go) there is a test for the middleware you'll build.
