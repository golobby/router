package router

import (
	"log"
	"net/http"
	"net/url"
)

// director is the base HTTP handler.
// It receives the request, and the response writer objects then pass them to the route through the middleware.
type director struct {
	repository *repository
	notFound   Handler
}

// ServeHTTP serves HTTP requests and uses other modules to handle them.
func (d *director) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	c := &DefaultContext{}
	c.SetRequest(request)
	c.SetRW(rw)

	uri, err := url.ParseRequestURI(request.RequestURI)
	if err != nil {
		log.Println("router: fail to parse request uri. err=" + err.Error())
		_ = c.Empty(http.StatusNotFound)
		return
	}

	route, parameters, err := d.repository.find(request.Method, uri.Path)
	if err != nil {
		if err = d.notFound(c); err != nil {
			log.Println("router: 404 handler error=" + err.Error())
			_ = c.Empty(http.StatusNotFound)
		}
		return
	}

	c.SetRoute(route)
	c.SetParameters(parameters)

	if err = d.Run(route, c); err != nil {
		log.Println("router: handler or middleware error=" + err.Error())
		_ = c.Empty(http.StatusInternalServerError)
	}
}

// Run calls the middleware and handler stack.
func (d *director) Run(route *Route, c Context) error {
	stack := make([]Handler, len(route.Middleware)+1)
	stack = append(stack, route.Handler)

	for i := len(route.Middleware); i > 0; i-- {
		stack = append(stack, route.Middleware[i-1](stack[len(stack)-1]))
	}

	return stack[len(stack)-1](c)
}

// newDirector creates a new instance of director.
func newDirector(repository *repository) *director {
	return &director{
		repository: repository,
		notFound: func(c Context) error {
			return c.Text(http.StatusNotFound, "404 Not Found.")
		},
	}
}
