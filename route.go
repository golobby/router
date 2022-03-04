package router

import "strings"

// Route holds Route information.
type Route struct {
	Method string
	Path   string
	Name   string
	stack  []Handler
}

func (r *Route) ToURL(parameters map[string]string) string {
	uri := r.Path
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
