package router

// repository holds the radix tree and current stateStack.
type repository struct {
	tree  *tree
	state *stateStack
}

// addRoute adds a new Route to the repository.
func (r *repository) addRoute(method, path string, handler Handler) *Route {
	route := newRoute(method, r.state.prefix()+path, r.stack(handler))
	r.tree.add(route)
	return route
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

// addParameterPattern adds a new Route parameter pattern to the radix tree.
func (r *repository) addParameterPattern(name, pattern string) {
	r.tree.patterns[name] = pattern
}

// findByRequest searches for a Route that matches the given HTTP method and URI.
// It returns the Route and its parameters.
func (r *repository) findByRequest(method, uri string) (*Route, map[string]string) {
	return r.tree.findByRequest(method, uri)
}

// findByName searches for a Route with the give name.
func (r *repository) findByName(name string) *Route {
	return r.tree.findByName(name)
}

// newRepository creates a new repository instance.
func newRepository() *repository {
	return &repository{newTree(), newStateStack()}
}
