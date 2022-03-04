package router

import (
	"errors"
)

// repository holds the radix tree and current stateStack.
type repository struct {
	tree  *tree
	state *stateStack
}

// addRoute adds a new Route to the repository.
func (r *repository) addRoute(method, path string, handler Handler) {
	r.tree.add(newRoute(method, r.state.prefix()+path, r.stack(handler)))
}

// stack merges handler and middlewares to create a stack of callables the Route is going to call.
func (r *repository) stack(handler Handler) []Handler {
	stack := make([]Handler, len(r.state.middlewares())+1)
	stack = append(stack, handler)

	for i := len(r.state.middlewares()); i > 0; i-- {
		stack = append(stack, r.state.middlewares()[i-1](stack[len(stack)-1]))
	}

	return stack
}

// addGroup adds a new group of routes to the repository.
func (r *repository) addGroup(prefix string, middleware []Middleware, body func()) {
	r.state.push(prefix, middleware)
	body()
	r.state.pop()
}

// updateGroup push the current group without pop.
func (r *repository) updateGroup(prefix string, middleware []Middleware) {
	r.state.push(prefix, middleware)
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
		state: newStateStack(),
		tree:  newTree(),
	}
}
