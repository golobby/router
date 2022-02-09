[![GoDoc](https://godoc.org/github.com/golobby/router/?status.svg)](https://godoc.org/github.com/golobby/router)
[![CI](https://github.com/golobby/router/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/router/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/router/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/golobby/router/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/router)](https://goreportcard.com/report/github.com/golobby/router)
[![Coverage Status](https://coveralls.io/repos/github/golobby/router/badge.svg)](https://coveralls.io/github/golobby/router?branch=master)

# GoLobby Router
GoLobby Router is a lightweight yet powerful HTTP router for Go projects.
It's built on top of the Golang HTTP package and adds  the following features to it:
* Routing based on HTTP method and URI
* Route parameters and parameter patterns
* Middleware
* HTTP Responses (such as JSON, XML, Text, Empty, and Redirect)

## Documentation
### Required Go Version
It requires Go `v1.11` or newer versions.

### Installation
To install this package, run the following command in your project directory.

```bash
go get github.com/golobby/router
```

### Quick Start

The following example demonstrates a simple example of using the router package.

```go
package main

import (
	"github.com/golobby/router"
	"log"
	"net/http"
)

func main() {
    r := router.New()
    
    r.GET("/", func(c router.Context) error {
        return c.Text(http.StatusOK, "Hello from GoLobby Router!")
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### HTTP Methods

You can use the "Map()" method to declare routes. It gets HTTP methods and paths (URIs).
There are also some methods available for the most used HTTP methods.
These methods are `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, and `OPTIONS`.

```go
func Handler(c router.Context) error {
    return c.Text(http.StatusOK, "Hello from GoLobby Router!")
}

func main() {
    r := router.New()
    
    r.GET("/", Handler)
    r.POST("/", Handler)
    r.PUT("/", Handler)
    r.PATCH("/", Handler)
    r.DELETE("/", Handler)
    r.HEAD("/", Handler)
    r.OPTIONS("/", Handler)
    
    r.Map("GET", "/", Handler)
    r.Map("CUSTOM", "/", Handler)
    
    log.Fatalln(r.Start(":8000"))
}
```

### Route Paramters

You can put route parameters inside curly braces like `{id}`.
To fetch them inside your handler, call the `Parameter()` method of the context.
You can also get all parameters at once, using the `Parameters()` method.

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.Get("/posts/{pid}/comments/{cid}", func(c router.Context) error {
      postId := c.Parameter("pid")
      commentId := c.Parameter("cid")
      // To get all parameters: c.Parameters()
      return c.Write([]byte("Hello Comment!"))
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### Groups

You may put routes with similar attributes into groups.
Currently, groups support prefix and middleware lists.

#### WithPrefix

You can group routes with common prefixes like the example below.

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.WithPrefix("/blog", func() {
      r.Get("/posts", PostsIndexHandler)
      r.Get("/posts/{1}", PostsShowHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### WithMiddleware

You can group routes with common middleware like the example below.

```go
import 	"github.com/golobby/router"

func AdminMiddleware(next router.Handler) router.Handler {
	return func(c router.Context) error {
		// Check user roles...
		return next(c)
	}
}

func main() {
    r := router.New()
    
    r.WithMiddleware(AdminMiddleware, func() {
      r.Get("/admin", PostHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### WithMiddlewareList

You can also assign multiple middlewares like the following example.

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.WithMiddlewareList([]router.Middleware{M1, M2}, func() {
      r.Get("/post", PostHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Alltogether

You can also create a group with both prefix and middlewares like the following sample.

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.Group("/blog", []router.Middleware{M1, M2}, func() {
      r.Get("/post", PostHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

## License
GoLobby Router is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
