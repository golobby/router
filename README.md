[![Go Reference](https://pkg.go.dev/badge/github.com/golobby/router.svg)](https://pkg.go.dev/github.com/golobby/router)
[![CI](https://github.com/golobby/router/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/router/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/router/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/golobby/router/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/router)](https://goreportcard.com/report/github.com/golobby/router)
[![Coverage Status](https://coveralls.io/repos/github/golobby/router/badge.svg?r=2)](https://coveralls.io/github/golobby/router?branch=master)

# GoLobby Router
GoLobby Router is a lightweight yet powerful HTTP router for Go projects.
It's built on top of the Go HTTP package and uses radix tree to provide the following features:
* Routing based on HTTP method and URI
* Route parameters and parameter patterns
* Route wildcards
* Middleware
* HTTP Responses (such as JSON, XML, Text, Empty, File, and Redirect)
* Static file serving
* No footprint!
* Zero-dependency!

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
        // c.Request() is original http.Request
        // c.Response() is original http.ResponseWriter
        return c.Text(http.StatusOK, "Hello from GoLobby Router!")
    })
    
    r.PUT("/products/:id", func(c router.Context) error {
        return c.Text(http.StatusOK, "Update product with ID: "+c.Parameter("id"))
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### HTTP Methods

You can use the `Map()` method to declare routes. It gets HTTP methods and paths (URIs).
There are also some methods available for the most used HTTP methods.
These methods are `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, and `OPTIONS`.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

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

To specify route parameters, prepend a colon like `:id`.
In default, parameters could be anything but you can determine a regex pattern using the `Define()` method. Of course, regex patterns slow down your application, and it is recommended not to use them if possible.
To catch and check route parameters in your handlers, you'll have the `Parameters()`, `Parameter()`, and `HasParameter()` methods.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    // "id" parameters must be numeric
    r.Define("id", "[0-9]+")
   
    // a route with one parameter
    r.GET("/posts/:id", func(c router.Context) error {
        return c.Text(http.StatusOK, c.Parameter("id"))
    })
    
    // a route with multiple parameters
    r.GET("/posts/:id/comments/:cid", func(c router.Context) error {
        return c.JSON(http.StatusOK, c.Parameters())
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### Wildcard Routes

Wildcard routes match any URI with the specified prefix.
The following example shows how it works.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    // Exception routes must come first.
    r.GET("/pages/contact", ContactHandler)
    
    r.GET("/pages/*", PagesHandler)
    // It matches:
    // - /pages/
    // - /pages/about
    // - /pages/about/us
    // - /pages/help
    
    log.Fatalln(r.Start(":8000"))
}
```

### Serving Static Files

The `FileHandler` and `FileHandlerWithStripper` handlers are provided to serve static files directly.
The examples below demonstrate how to use them.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    r.GET("/api", YourApiHandler)
    
    r.GET("/*", router.FileHandler("./files"))
    // example.com/            ==> ./files/index.html
    // example.com/photo.jpg   ==> ./files/photo.jpg
    // example.com/notes/1.txt ==> ./files/notes/1.txt
    
    log.Fatalln(r.Start(":8000"))
}
```

You might serve static files with different URI.
In this case, you must strip the extra URI prefix like this example.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    r.GET("/api", YourApiHandler)
    
    r.GET("/files/*", router.FileHandlerWithStripper("./files", "/files/"))
    // example.com/files/            ==> ./files/index.html
    // example.com/files/photo.jpg   ==> ./files/photo.jpg
    // example.com/files/notes/1.txt ==> ./files/notes/1.txt
    
    log.Fatalln(r.Start(":8000"))
}
```

### Named Routes

Named routes allow the convenient generation of URLs or redirects for specific routes.
You may specify a name for a route by chaining the `SetName()` method onto the route definition:

```go
package main

import (
    "github.com/golobby/router"
    "github.com/golobby/router/pkg/response"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    r.GET("/", func(c router.Context) error {
        return c.Text(http.StatusOK, "I am the home!")
    }).SetName("home")
    
    r.GET("/posts/:id", func(c router.Context) error {
        return c.Text(http.StatusOK, "I am a post!")
    }).SetName("post")
    
    r.GET("/links", func(c router.Context) error {
        return c.JSON(http.StatusOK, response.M{
            "home": c.URL("home", nil), // "/"
            "post-1": c.URL("post", map[string]string{"id": "1"}), // "/posts/1"
            "post-2": c.URL("post", map[string]string{"id": "2"}), // "/posts/2"
        })
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### Responses

The router comes with `Empty`, `Redirect`, `Text`, `HTML`, `JSON`, `PrettyJSON`, `XML`, `PrettyXML`, and `Bytes` responses out of the box.
The examples below demonstrate how to use built-in and custom responses.

```go
package main

import (
    "github.com/golobby/router"
    "github.com/golobby/router/pkg/response"
    "log"
    "net/http"
)

func main() {
    r := router.New()

    r.GET("/empty", func(c router.Context) error {
        return c.Empty(204)
    })

    r.GET("/redirect", func(c router.Context) error {
        return c.Redirect(301, "https://github.com/golobby/router")
    })

    r.GET("/text", func(c router.Context) error {
        return c.Text(200, "A text response")
    })

    r.GET("/html", func(c router.Context) error {
        return c.HTML(200, "<p>A HTML response</p>")
    })

    r.GET("/json", func(c router.Context) error {
        return c.JSON(200, User{"id": 13})
    })

    r.GET("/json", func(c router.Context) error {
        return c.JSON(200, response.M{"message": "Using response.M helper"})
    })

    r.GET("/json-pretty", func(c router.Context) error {
        return c.PrettyJSON(200, response.M{"message": "A pretty JSON response!"})
    })

    r.GET("/xml", func(c router.Context) error {
        return c.XML(200, User{"id": 13})
    })

    r.GET("/xml-pretty", func(c router.Context) error {
        return c.PrettyXML(200, User{"id": 13})
    })

    r.GET("/bytes", func(c router.Context) error {
        return c.Bytes(200, []bytes("Some bytes!"))
    })

    r.GET("/file", func(c router.Context) error {
	return c.File(200, "text/plain", "text.txt")
    })

    r.GET("/custom", func(c router.Context) error {
        c.Response().Header().Set("Content-Type", "text/csv")
        return c.Bytes(200, []bytes("Column 1, Column 2, Column 3"))
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
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    r.WithPrefix("/blog", func() {
        r.GET("/posts", PostsHandler)    // "/blog/posts"
        r.GET("/posts/:id", PostHandler) // "/blog/posts/:id"
        r.WithPrefix("/pages", func() {
            r.GET("/about", AboutHandler)     // "/blog/pages/about"
            r.GET("/contact", ContactHandler) // "/blog/pages/contact"
        })
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Group by middleware

The example below demonstrates how to group routes with the same middleware.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func AdminMiddleware(next router.Handler) router.Handler {
    return func(c router.Context) error {
        // Check user roles...
        return next(c)
    }
}

func main() {
    r := router.New()
    
    r.WithMiddleware(AdminMiddleware, func() {
        r.GET("/admin/users", UsersHandler)
        r.GET("/admin/products", ProductsHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Group by middlewares

The example below demonstrates how to group routes with the same middlewares.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    middlewares := []router.Middleware{Middleware1, Middleware2, Middleware3}
    r.WithMiddlewares(middlewares, func() {
        r.GET("/posts", PostsIndexHandler)
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Group by multiple attributes

The `group()` method helps you create a group of routes with the same prefix and middlewares.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    r.Group("/blog", []router.Middleware{Middleware1, Middleware2}, func() {
        r.GET("/posts", PostsHandler)
        r.GET("/posts/:id/comments", CommentsHandler)
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
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    // Add a prefix to all routes
    r.AddPrefix("/blog")

    r.GET("/posts", PostsHandler)
    r.GET("/posts/:id/comments", CommentsHandler)
    
    log.Fatalln(r.Start(":8000"))
}
```

#### Base middlewares

The following example shows how to set a base middlewares for all routes.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    // Add a single middleware
    r.AddMiddleware(LoggerMiddleware)
    
    // Add multiple middlewares at once
    r.AddMiddlewares([]router.Middleware{AuthMiddleware, ThrottleMiddleware})

    r.GET("/users", UsersHandler)
    r.GET("/users/:id/files", FilesHandler)
    
    log.Fatalln(r.Start(":8000"))
}
```

### 404 Handler

In default, the router returns the following HTTP 404 response when a requested URI doesn't match any route.

```json
{"message": "Not found."}
```

You can set your custom handler like the following example.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    // Custom (HTML) Not Found Handler
    r.SetNotFoundHandler(func(c router.Context) error {
        return c.HTML(404, "<p>404 Not Found</p>")
    })

    r.GET("/", Handler)
    
    log.Fatalln(r.Start(":8000"))
}
```

### Error Handling

Your handlers might return an error while processing the HTTP request.
This error can be produced by your application logic or failure in the HTTP response.
By default, the router logs it using Golang's built-in logger into the standard output and returns the HTTP 500 response below.

```json
{"message": "Internal error."}
```

It's a good practice to add a global middleware to catch all these errors, log and handle them the way you need.
The example below demonstrates how to add middleware for handling errors.

```go
package main

import (
    "github.com/golobby/router"
    "log"
    "net/http"
)

func main() {
    r := router.New()
    
    // Error Handler
    r.AddMiddleware(func (next router.Handler) router.Handler {
        return func(c router.Context) error {
            if err := next(c); err != nil {
                myLogger.log(err)
                return c.HTML(500, "<p>Something went wrong</p>")
            }
            
            // No error will raise to the router base handler
            return nil
        }
    })

    r.GET("/", Handler)
    
    log.Fatalln(r.Start(":8000"))
}
```

## License
GoLobby Router is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
