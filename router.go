// Package router is a lightweight yet powerful HTTP router.
// It's built on top of the built-in Golang HTTP library and adds real-world requirements to it.
package router

import (
	"log"
	"net/http"
)

// Router is the entry point of the package.
// It gets routes, route parameter patterns, and middleware then dispatch them.
// It receives every HTTP request, finds the related route then runs the route handler through its middleware.
type Router struct {
	repository *repository
	director   *director
}

// Define assigns a regular expression pattern to a route parameter.
// After the definition, the router only chooses the related route if the request URI matches the pattern.
func (r Router) Define(parameter, pattern string) {
	r.repository.addParameter(parameter, pattern)
}

// Map defines a new route (using HTTP method and path/URI) and assigns a handler.
// The path may contain route parameters.
func (r Router) Map(method, path string, handler Handler) {
	r.repository.addRoute(method, path, handler)
}

// Group creates a group of routes with common attributes.
// Currently, prefix and middleware attributes are supported.
func (r Router) Group(prefix string, middleware []Middleware, body func()) {
	r.repository.addGroup(prefix, middleware, body)
}

// WithPrefix creates a group of routes with common prefix.
func (r Router) WithPrefix(prefix string, body func()) {
	r.Group(prefix, []Middleware{}, body)
}

// WithMiddleware creates a group of routes with common middleware.
func (r Router) WithMiddleware(middleware Middleware, body func()) {
	r.Group("", []Middleware{middleware}, body)
}

// WithMiddlewares creates a group of routes with common set of middleware.
func (r Router) WithMiddlewares(middleware []Middleware, body func()) {
	r.Group("", middleware, body)
}

// AddPrefix adds a global prefix for next routes.
func (r Router) AddPrefix(prefix string) {
	r.repository.updateGroup(prefix, []Middleware{})
}

// AddMiddleware adds a global middleware for next routes.
func (r Router) AddMiddleware(middleware Middleware) {
	r.repository.updateGroup("", []Middleware{middleware})
}

// AddMiddlewares adds global middlewares for next routes.
func (r Router) AddMiddlewares(middlewares []Middleware) {
	r.repository.updateGroup("", middlewares)
}

func (r Router) SetNotFoundHandler(handler Handler) {
	r.director.notFoundHandler = handler
}

// Start runs the HTTP listener and waits for HTTP requests.
// It should be called after definitions of routes.
func (r Router) Start(address string) error {
	log.Println("http router listening to " + address)
	return http.ListenAndServe(address, r.director)
}

// Serve handles the request manually with a given request and a response writer.
func (r Router) Serve(rw http.ResponseWriter, request *http.Request) {
	r.director.ServeHTTP(rw, request)
}

// GET maps a GET route.
func (r Router) GET(path string, handler Handler) {
	r.Map("GET", path, handler)
}

// POST maps a POST route.
func (r Router) POST(path string, handler Handler) {
	r.Map("POST", path, handler)
}

// PUT maps a PUT route.
func (r Router) PUT(path string, handler Handler) {
	r.Map("PUT", path, handler)
}

// PATCH maps a PATCH route.
func (r Router) PATCH(path string, handler Handler) {
	r.Map("PATCH", path, handler)
}

// DELETE maps a DELETE route.
func (r Router) DELETE(path string, handler Handler) {
	r.Map("DELETE", path, handler)
}

// HEAD maps a HEAD route.
func (r Router) HEAD(path string, handler Handler) {
	r.Map("HEAD", path, handler)
}

// OPTIONS maps a OPTIONS route.
func (r Router) OPTIONS(path string, handler Handler) {
	r.Map("OPTIONS", path, handler)
}

// New creates a new instance of the HTTP router.
func New() *Router {
	repository := newRepository()
	director := newDirector(repository)

	return &Router{
		repository: repository,
		director:   director,
	}
}
