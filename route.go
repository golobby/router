package router

import "strings"

// Route holds Route information.
type Route struct {
	method string
	path   string
	name   string
	stack  []Handler
}

// Method returns route method.
func (r *Route) Method() string {
	return r.method
}

// Path returns route method.
func (r *Route) Path() string {
	return r.path
}

// Name returns route method.
func (r *Route) Name() string {
	return r.name
}

// SetName sets/updates route name.
func (r *Route) SetName(name string) {
	r.name = name
}

// URL generate URL from route path with given parameters.
func (r *Route) URL(parameters map[string]string) string {
	uri := r.path
	if parameters != nil {
		for name, value := range parameters {
			uri = strings.Replace(uri, ":"+name, value, 1)
		}
	}
	return uri
}

// newRoute creates a new Route instance.
func newRoute(method, path string, stack []Handler) *Route {
	return &Route{method, path, "", stack}
}
