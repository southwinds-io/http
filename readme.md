# HTTP Server Library

The HTTP server library helps to quickly create a consistent HTTP server
with some useful middlewares, OpenAPI, Prometheus metrics and standard access controls.

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

For a more complete example [see here](server_test.go).

### Configuratiion

The server can be configured via the following environment variables:

| name | description | default |
|---|---|---|
| `HTTP_METRICS_ENABLED` | enables Prometheus endpoint at /metrics | true |
|`HTTP_OPEN_API_ENABLED` | enables Swagger enpoint at /api/ | true |
|`HTTP_PORT` | port on which the server listen | 8080 |
|`HTTP_USER` | basic authentication user name  | admin|
|`HTTP_PASSWORD` | basic authentication password  | auto-generated |
|`HTTP_REALM` | the realm name for browser WWW-Authenticate prompt | none |

