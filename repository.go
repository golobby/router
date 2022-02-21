package router

import (
	"errors"
)

// repository holds the defined routes, groups, and patterns.
type repository struct {
	tree       *tree
	routes     []*Route
	parameters map[string]string
	collector  *collector
}

// addRoute adds a new Route to the repository.
func (r *repository) addRoute(method, path string, handler Handler) {
	r.tree.add(newRoute(method, r.collector.active.prefix+path, r.stack(handler)))
}

// stack merges handler and middleware to create a stack of callables the Route needs to call.
func (r *repository) stack(handler Handler) []Handler {
	stack := make([]Handler, len(r.collector.active.middleware)+1)
	stack = append(stack, handler)

	for i := len(r.collector.active.middleware); i > 0; i-- {
		stack = append(stack, r.collector.active.middleware[i-1](stack[len(stack)-1]))
	}

	return stack
}

// addGroup adds a new group of routes to the repository.
func (r *repository) addGroup(prefix string, middleware []Middleware, body func()) {
	r.collector.update(prefix, middleware)
	body()
	r.collector.rollback()
}

// updateGroup update the active group without rollback.
func (r *repository) updateGroup(prefix string, middleware []Middleware) {
	r.collector.update(prefix, middleware)
}

// addParameter adds a new Route parameter pattern to the repository.
func (r *repository) addParameter(name, pattern string) {
	r.tree.patterns[name] = pattern
}

// find searches for a Route that matches the given HTTP method and URI.
// It returns the Route and its parameters.
func (r *repository) find(method, uri string) (*Route, map[string]string, error) {
	if route, parameters := r.tree.find(method, uri); route != nil {
		return route, parameters, nil
	}

	return nil, nil, errors.New("router: cannot find a Route for the request")
}

// newRepository creates a new repository instance.
func newRepository() *repository {
	return &repository{
		parameters: map[string]string{},
		collector:  newCollector(),
		tree:       newTree(),
	}
}
