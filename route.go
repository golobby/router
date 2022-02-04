package router

// Route holds route information.
type Route struct {
	Method string
	Path   string
	Stack  []Handler
}

// newRoute creates a new route
func newRoute(method, path string, stack []Handler) *Route {
	return &Route{method, path, stack}
}
