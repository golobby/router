package router

import (
	"github.com/golobby/router/pkg/response"
	"log"
	"net/http"
	"net/url"
)

// director is the base HTTP handler.
// It receives the request, and the responseWriter objects then pass them to the Route through the middlewares.
type director struct {
	repository      *repository
	notFoundHandler Handler
}

// ServeHTTP serves HTTP requests and uses other modules to handle them.
func (d *director) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	c := &DefaultContext{}
	c.SetRequest(request)
	c.SetResponse(rw)

	uri, err := url.ParseRequestURI(request.RequestURI)
	if err != nil {
		d.serveNotFoundError(c)
		return
	}

	route, parameters, err := d.repository.find(request.Method, uri.Path)
	if err != nil {
		d.serveNotFoundError(c)
		return
	}

	c.SetRoute(route)
	c.SetParameters(parameters)

	if err = route.stack[len(route.stack)-1](c); err != nil {
		d.serveInternalError(c, err)
	}
}

// serveInternalError handles internal errors.
func (d *director) serveInternalError(c Context, err error) {
	log.Println("router: uncaught error=" + err.Error())
	_ = c.Json(http.StatusInternalServerError, response.M{"message": "Internal error."})
}

// serveNotFoundError handles 404 errors.
func (d *director) serveNotFoundError(c Context) {
	err := d.notFoundHandler(c)
	if err != nil {
		d.serveInternalError(c, err)
	}
}

// newDirector creates a new director instance.
func newDirector(repository *repository) *director {
	return &director{
		repository: repository,
		notFoundHandler: func(c Context) error {
			return c.Json(http.StatusNotFound, response.M{"message": "Not found."})
		},
	}
}
