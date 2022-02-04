// Package router is a lightweight yet powerful HTTP router.
// It's built on top of the built-in Golang HTTP library and adds real-world requirements to it.
package router

import (
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

// Get maps a GET route.
func (r Router) Get(path string, handler Handler) {
	r.Map("GET", path, handler)
}

// Post maps a POST route.
func (r Router) Post(path string, handler Handler) {
	r.Map("POST", path, handler)
}

// Put maps a PUT route.
func (r Router) Put(path string, handler Handler) {
	r.Map("PUT", path, handler)
}

// Patch maps a PATCH route.
func (r Router) Patch(path string, handler Handler) {
	r.Map("PATCH", path, handler)
}

// Delete maps a DELETE route.
func (r Router) Delete(path string, handler Handler) {
	r.Map("DELETE", path, handler)
}

// Group creates a group of routes with common attributes.
// Currently, prefix and middleware attributes are supported.
func (r Router) Group(prefix string, middleware []Middleware, body func()) {
	r.repository.addGroup(prefix, middleware, body)
}

// WithPrefix creates a group of routes with common prefix.
func (r Router) WithPrefix(prefix string, body func()) {
	r.repository.addGroup(prefix, []Middleware{}, body)
}

// WithMiddleware creates a group of routes with common middleware.
func (r Router) WithMiddleware(middleware Middleware, body func()) {
	r.repository.addGroup("", []Middleware{middleware}, body)
}

// WithMiddlewareList creates a group of routes with common set of middleware.
func (r Router) WithMiddlewareList(middleware []Middleware, body func()) {
	r.repository.addGroup("", middleware, body)
}

func (r Router) SetNotFoundHandler(handler Handler) {
	r.director.notFoundHandler = handler
}

// Start runs the HTTP listener and waits for HTTP requests.
// It should be called after definitions of routes.
func (r Router) Start(address string) error {
	return http.ListenAndServe(address, r.director)
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
