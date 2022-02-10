[![GoDoc](https://godoc.org/github.com/golobby/router/?status.svg)](https://godoc.org/github.com/golobby/router)
[![CI](https://github.com/golobby/router/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/router/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/router/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/golobby/router/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/router)](https://goreportcard.com/report/github.com/golobby/router)
[![Coverage Status](https://coveralls.io/repos/github/golobby/router/badge.svg?v=dev)](https://coveralls.io/github/golobby/router?branch=master)

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

You can use the `Map()` method to declare routes. It gets HTTP methods and paths (URIs).
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

### Route Parameters

To specify route parameters, put them inside curly braces like `{id}` or `{id?}` where it's optional.
In default, the regular expression pattern for parameters is "`[^/]+`".
You can change it using the `Define()` method.
To catch and check route parameters in your handlers, you'll have the `Parameters()`, `Parameter()`, and `HasParameter()` methods.

```go
func main() {
    r := router.New()
    
    // "id" parameters must be numeric
    r.Define("id", "[0-9]+")
   
    // a required parameter
    r.Get("/posts/{id}", func(c router.Context) error {
    	return c.Text(200, c.Parameter("id"))
    })
    
    // multiple required parameters
    r.Get("/posts/{pid}/comments/{cid}", func(c router.Context) error {
    	return c.Json(200, c.Parameters())
    })
    
    // an optional parameter
    r.Get("/posts/{id?}", func(c router.Context) error {
    	if c.HasParameter("id") {
    	    return c.Text(200, c.Parameter("id"))
	} else {
	    return c.Text(200, "No Parameter")
	}
    })
    
    // an optional parameter after an optional slash!
    r.Get("/posts/?{id?}", func(c router.Context) error {
   	// It runs for "/posts/1", "/posts/", and "/posts"
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### Groups

You may put routes with similar attributes in groups.
Currently, prefix and middleware attributes are supported.

#### Group by prefix

The example below demonstrates how to group routes with the same prefix.

```go
func main() {
    r := router.New()
    
    r.WithPrefix("/blog", func() {
      r.Get("/posts", PostsIndexHandler)
      r.Get("/posts/{1}", PostsShowHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Group by middleware

The example below demonstrates how to group routes with the same middleware.

```go
func AdminMiddleware(next router.Handler) router.Handler {
    return func(c router.Context) error {
        // Check user roles...
        return next(c)
    }
}

func main() {
    r := router.New()
    
    r.WithMiddleware(AdminMiddleware, func() {
      r.Get("/admin/users", UsersHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Group by middlewares

The example below demonstrates how to group routes with the same middlewares.

```go
func main() {
    r := router.New()
    
    r.WithMiddlewares([]router.Middleware{Middleware1, Middleware2, Middleware3}, func() {
        r.Get("/posts", PostsHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Group by multiple attributes

The `group()` method helps you create a group of routes with the same prefix and middlewares.

```go
func main() {
    r := router.New()
    
    r.Group("/blog", []router.Middleware{Middleware1, Middleware2}, func() {
      r.Get("/posts", PostsHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### Basic Attributes

Your application might need a base prefix or global middlewares.
In this case, you can set up these base attributes before defining routes.

#### Base prefix

The following example shows how to set a base prefix for all routes.

```go
func main() {
    r := router.New()
    r.AddPrefix("/blog")

    r.Get("/posts", PostsHandler)
    r.Get("/posts/{id}/comments", CommentsHandler)
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Base middlewares

The following example shows how to set a base middlewares for all routes.

```go
func main() {
    r := router.New()
    
    // Add a single middleware
    r.AddMiddleware(LoggerMiddleware)
    
    // Add multiple middlewares
    r.AddMiddlewares([]router.Middleware{AuthMiddleware, ThrottleMiddleware})

    r.Get("/users", UsersHandler)
    r.Get("/users/{id}/files", FilesHandler)
    
    log.Fatalln(r.Start(":8000"))
}
```

### 404 Handler

In default, the router returns the following HTTP response when a requested URI doesn't match any route.

```http
HTTP/1.1 404 Not Found
Content-Type: application/json

{"message": "Not found."}
```

You can set your custom handler like the following example.

```go
func main() {
    r := router.New()
    
    r.SetNotFoundHandler(func(c router.Context) error {
	return c.Html(404, "<p>404 Not Found</p>")
    })

    r.GET("/", Handler)
    
    log.Fatalln(r.Start(":8000"))
}
```

### Error Handling

Your handlers might return an error while processing the HTTP request.
This error can be produced by your application logic or failure in the HTTP response.
By default, the router logs it using Golang's built-in logger into the standard output and returns the HTTP response below.

```http
HTTP/1.1 500 Internal Error
Content-Type: application/json

{"message": "Internal error."}
```

It's a good practice to add a global middleware to catch all these errors, log and handler them the way you need.
The example below demonstrates how to add middleware for handling errors to the router.

```go
func main() {
    r := router.New()
    
    // Error Handler
    r.AddMiddleware(func (next router.Handler) router.Handler {
        return func(c router.Context) error {
	    if err := next(c); err != nil {
	    	myLogger.log(err);
		retrun c.Html(500, "<p>Something went wrong</p>")
	    }
	    
	    return nil
        }
    })

    r.GET("/", Handler)
    
    log.Fatalln(r.Start(":8000"))
}
```

## License
GoLobby Router is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
