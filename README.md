# Router
A lightweight yet powerful HTTP router for Go projects.
It's built on top of the built-in Golang HTTP library and added real-world requirements to it.

## Documentation
### Required Go Version
It requires Go `v1.11` or newer versions.

### Installation
To install this package run the following command in the root of your project.

```bash
go get github.com/golobby/router
```

### Quick Start

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.Get("/", func(c router.Context) error {
      _, err := c.ResponseWriter().Write([]byte("Hello World!"))
      return err
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### HTTP Methods

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.Get("/", Handler)
    r.Post("/", Handler)
    r.Put("/", Handler)
    r.Patch("/", Handler)
    r.Delete("/", Handler)
    r.Head("/", Handler)
    r.Options("/", Handler)
    
    r.Map("GET", "/", Handler)
    r.Map("CUSTOM", "/", Handler)
    
    log.Fatalln(r.Start(":8000"))
}
```

### Route Paramters

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.Get("/posts/{pid}/comments/{cid}", func(c router.Context) error {
      postId := c.Parameter("pid")
      commentId := c.Parameter("cid")
      
      _, err := c.ResponseWriter().Write([]byte("Hello Comment!"))
      return err
    })
    
    log.Fatalln(r.Start(":8000"))
}
```

### Groups

#### WithPrefix

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.WithPrefix("/blog", func() {
      r.Get("/post", PostHandler)
	  })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### WithMiddleware

```go
import 	"github.com/golobby/router"

func main() {
    r := router.New()
    
    r.WithMiddleware(myMiddleware, func() {
      r.Get("/post", PostHandler)
	  })
    
    log.Fatalln(r.Start(":8000"))
}
```

#### WithMiddlewareList

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

## License
GoLobby Router is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
