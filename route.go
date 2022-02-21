package router

// Route holds Route information.
type Route struct {
	Method string
	Path   string
	stack  []Handler
}

// newRoute creates a new Route instance.
func newRoute(method, path string, stack []Handler) *Route {
	return &Route{method, path, stack}
}
