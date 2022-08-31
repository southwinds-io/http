# HTTP Server Library

The HTTP server library helps to quickly create a consistent HTTP server
with some useful middlewares, OpenAPI and standard access controls.

The server is focused on serving REST style APIs and is based on the [Gorilla Mux](https://github.com/gorilla/mux) request 
router and dispatcher.

### Getting started

Add the module to your project:

```bash
$  go get southwinds.dev/http
```

Set up the server:

```go
package main

import (
    "github.com/gorilla/mux"
    "net/http"
    h "southwinds.dev/http"
)

func main() {
    // create a new server
    server := h.New("test")
    // configure http handlers
    server.Http = func(router *mux.Router) {
        // add basic authentication
        router.Use(server.AuthenticationMiddleware)
        // add handler
        router.HandleFunc("/", doSomething).Methods("GET")
    }
    server.Serve()
}

func doSomething(writer http.ResponseWriter, request *http.Request) {
    writer.WriteHeader(http.StatusOK)
}
```

For a more complete example [see here](server_test.go)
