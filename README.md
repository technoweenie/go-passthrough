# Passthrough

Passes a white list of headers and a response body from a client HTTP request
to a ResponseWriter.  This can be useful for little proxy services around
internal APIs.

```go
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
  res, _ := http.Get("internal/service")
  
  pass := passthrough.New([]string{"Content-Type", "ETag", "Last-Modified"})
  pass.Pass(res, w, 200)
  
  // maybe you just want to pass the headers, but not the body
  pass.PassHeaders(res.Header, w)
}
```

_Barely enough code to extract._
