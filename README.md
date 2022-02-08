[![GoDoc](https://godoc.org/github.com/golobby/router/?status.svg)](https://godoc.org/github.com/golobby/router/v3)
[![CI](https://github.com/golobby/router/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/router/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/router/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/golobby/router/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/router)](https://goreportcard.com/report/github.com/golobby/router)
[![Coverage Status](https://coveralls.io/repos/github/golobby/router/badge.svg?v=0)](https://coveralls.io/github/golobby/router?branch=master)

# Router
A lightweight yet powerful HTTP router for Go projects.
It's built on top of the built-in Golang HTTP package and provides the following features:
* HTTP routing based on HTTP method and URI (path)
* Route parameters (and parameter regular expression patterns)
* Middleware
* HTTP Responses like JSON, XML, Text, Empty and so on

## Documentation
### Required Go Version
It requires Go `v1.11` or newer versions.

### Installation
To install this package run the following command in the root of your project.

```bash
go get github.com/golobby/router
```

### Quick Start

The following example demonstrates a simple example of using the router.

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.Get("/", func(c router.Context) error {
        return c.Text(http.StatusOK, "Hello from GoLobby Router!")
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### HTTP Methods

You may use the "Map()" to declare routes for the given HTTP method.
There are also some methods available for most used HTTP methods ("GET", "POST", "PUT", "PATCH", and "DELETE").

```go
import 	"github.com/golobby/router"

func Handler(c router.Context) error {
    return c.Text(http.StatusOK, "Hello from GoLobby Router!")
}

func main() {
    r := router.New()
    
    r.Get("/", Handler)
    r.Post("/", Handler)
    r.Put("/", Handler)
    r.Patch("/", Handler)
    r.Delete("/", Handler)
    
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
