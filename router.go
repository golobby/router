// Package router is a lightweight yet powerful HTTP router.
// It's built on top of the Golang HTTP package and uses the radix tree to provide routing requirements for modern
// applications.
package router

import (
	"net/http"
)

// Router is the entry point of the package.
// It gets routes, Route parameter patterns, and middlewares then dispatches them.
// It receives HTTP requests, finds the related Route then runs the Route handler through its middlewares.
type Router struct {
	repository *repository
	director   *director
}

// Define assigns a regular expression pattern to a Route parameter.
// After the definition, the router only dispatches the related Route if the request URI matches the pattern.
func (r Router) Define(parameter, pattern string) {
	r.repository.addParameterPattern(parameter, pattern)
}

// Files defines a new Route by HTTP method and path and assigns a handler.
// The path (URI) may contain Route parameters.
func (r Router) Files(path, directory string) *Route {
	return r.GET(path, FilesHandler(path, directory))
}

// Map defines a new Route by HTTP method and path and assigns a handler.
// The path (URI) may contain Route parameters.
func (r Router) Map(method, path string, handler Handler) *Route {
	return r.repository.addRoute(method, path, handler)
}

// Group creates a group of routes with common attributes.
// Currently, content and middlewares attributes are supported.
func (r Router) Group(prefix string, middleware []Middleware, body func()) {
	r.repository.addGroup(prefix, middleware, body)
}

// WithPrefix creates a group of routes with common content.
func (r Router) WithPrefix(prefix string, body func()) {
	r.Group(prefix, []Middleware{}, body)
}

// WithMiddleware creates a group of routes with common middlewares.
func (r Router) WithMiddleware(middleware Middleware, body func()) {
	r.Group("", []Middleware{middleware}, body)
}

// WithMiddlewares creates a group of routes with common set of middlewares.
func (r Router) WithMiddlewares(middleware []Middleware, body func()) {
	r.Group("", middleware, body)
}

// AddPrefix adds a global content for next or all routes.
func (r Router) AddPrefix(prefix string) {
	r.repository.updateGroup(prefix, []Middleware{})
}

// AddMiddleware adds a global middlewares for next or all routes.
func (r Router) AddMiddleware(middleware Middleware) {
	r.repository.updateGroup("", []Middleware{middleware})
}

// AddMiddlewares adds set of global middlewares for next or all routes.
func (r Router) AddMiddlewares(middlewares []Middleware) {
	r.repository.updateGroup("", middlewares)
}

// SetNotFoundHandler receives a handler and runs it when user request won't lead to any declared Route.
// It is the application 404 error handler, indeed.
func (r Router) SetNotFoundHandler(handler Handler) {
	r.director.notFoundHandler = handler
}

// Start runs the HTTP listener and waits for HTTP requests.
// It should be called after definitions of routes.
func (r Router) Start(address string) error {
	return http.ListenAndServe(address, r.director)
}

// Serve handles the request manually with a given request and a response writer.
func (r Router) Serve(rw http.ResponseWriter, request *http.Request) {
	r.director.ServeHTTP(rw, request)
}

// GET maps a GET Route.
func (r Router) GET(path string, handler Handler) *Route {
	return r.Map("GET", path, handler)
}

// POST maps a POST Route.
func (r Router) POST(path string, handler Handler) *Route {
	return r.Map("POST", path, handler)
}

// PUT maps a PUT Route.
func (r Router) PUT(path string, handler Handler) *Route {
	return r.Map("PUT", path, handler)
}

// PATCH maps a PATCH Route.
func (r Router) PATCH(path string, handler Handler) *Route {
	return r.Map("PATCH", path, handler)
}

// DELETE maps a DELETE Route.
func (r Router) DELETE(path string, handler Handler) *Route {
	return r.Map("DELETE", path, handler)
}

// HEAD maps a HEAD Route.
func (r Router) HEAD(path string, handler Handler) *Route {
	return r.Map("HEAD", path, handler)
}

// OPTIONS maps a OPTIONS Route.
func (r Router) OPTIONS(path string, handler Handler) *Route {
	return r.Map("OPTIONS", path, handler)
}

// New creates a new Router instance.
func New() *Router {
	repository := newRepository()
	director := newDirector(repository)
	return &Router{repository, director}
}
