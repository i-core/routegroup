# routegroup [![GoDoc][doc-img]][doc] [![Build Status][build-img]][build]

routegroup provides a modular approach to organizing HTTP handlers and encapsulates an HTTP router's implementation.

## Installation

```bash
go get github.com/i-core/routegroup
```

## Usage

### Grouping of HTTP handlers

Create a Go module that contains a resuable group of HTTP handlers. For example,
a group of HTTP handlers that provides health checking and the application's version over HTTP:

```go
package httpstat // import "example.org/httpstat""

import (
    "encoding/json"
    "net/http"
)

type Handler struct {
    version string
}

func NewHandler(version string) *Handler {
    return &Handler{version: version}
}

func (h *Handler) AddRoutes(apply func(m, p string, h http.Handler, mws ...func(http.Handler) http.Handler)) {
    apply(http.MethodGet, "/health/alive", newHealthHandler())
    apply(http.MethodGet, "/health/ready", newHealthHandler())
    apply(http.MethodGet, "/version", newVersionHandler(h.version))
}

func newHealthHandler() http.Handler {
    return http.HandlerFunc(w http.ResponseWriter, r *http.Request) {
        r.WriteHeader(http.StatusOK)
    }
}

func newVersionHandler(version string) http.Handler {
    return http.HandlerFunc(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, version)
    }
}
```

Now we can use this module in our applications using the routegroup router:

```go
package main

import (
    "fmt"
    "net/http"

    "example.org/httpstat"
    "github.com/i-core/routegroup"
)

var version = "<APP_VERSION>"

func main() {
    router := routegroup.NewRouter()
    router.AddRoutes(httpstat.NewHandler(version), "/stat")

    // Start a web server.
    fmt.Println("Server started")
    log.Printf("Server finished; error: %s\n", http.ListenAndServe(":8080", router))
}
```

### Common middlewares

In some cases, you may need to use some middlewares for all routes in an application.
To do it you should specify the middlewares as arguments in the router's constructor.

For example, we add a logging middleware to the router in the example above:

```go
package main

import (
    "fmt"
    "net/http"

    "example.org/httpstat"
    "github.com/i-core/routegroup"
)

var version = "<APP_VERSION>"

func main() {
    router := routegroup.NewRouter(logw)
    router.AddRoutes(httpstat.NewHandler(version), "/stat")

    // Start a web server.
    fmt.Println("Server started")
    log.Printf("Server finished; error: %s\n", http.ListenAndServe(":8080", router))
}

func logw(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Incomming request %s %s\n", r.Method, r. URL)
        h.ServeHTTP(w, r)
    })
}
```

### Path parameters

The routegroup router support path parameters.
Path parameters are defined in the format `:param_name`, for example, `/employee/:id`.
routegroup provides the function `routegroup.PathParam` to get a value of a path parameter by a name:

```go
package api

import (
    "encoding/json"
    "net/http"

    "github.com/i-core/routegroup"
)

type Handler struct {}

func NewHandler() *Handler {
    return &Handler{}
}

func (h *Handler) AddRoutes(apply func(m, p string, h http.Handler, mws ...func(http.Handler) http.Handler)) {
    // Register a route with a path parameter.
    apply(http.MethodGet, "/employee/:id", newEmployeeHandler())
}

func newEmployeeHandler() http.Handler {
    return http.HandlerFunc(w http.ResponseWriter, r *http.Request) {
        // Get a value of the path parameter registered above.
        empId := routegroup.PathParam(r.Context(), "id")
        var emp map[string]interface{}

        // Load employee's data from the DB...

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(emp); err != nil {
            panic(err)
        }
    }
}
```

## Implementation

routegroup uses [github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) as an HTTP router.

## Contributing

Thanks for your interest in contributing to this project.
Get started with our [Contributing Guide][contrib].

## License

The code in this project is licensed under [MIT license][license].

[doc-img]: https://godoc.org/github.com/i-core/routegroup?status.svg
[doc]: https://godoc.org/github.com/i-core/routegroup

[build-img]: https://travis-ci.com/i-core/routegroup.svg?branch=master
[build]: https://travis-ci.com/i-core/routegroup

[contrib]: https://github.com/i-core/.github/blob/master/CONTRIBUTING.md
[license]: LICENSE