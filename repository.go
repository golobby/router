package router

import (
	"errors"
	"regexp"
	"strings"
)

// repository holds the defined routes, groups, and patterns.
type repository struct {
	routes     []*Route
	parameters map[string]string
	collector  *collector
}

// addRoute adds a new route to the repository.
func (r *repository) addRoute(method, path string, handler Handler) {
	r.routes = append(r.routes, newRoute(method, r.path(path), r.stack(handler)))
}

// path appends route path with group path and converts parameters to patterns.
func (r *repository) path(path string) string {
	return r.collector.active.prefix + r.pattern(path)
}

// stack merges handler and middleware to create a stack of callables the route needs to call.
func (r *repository) stack(handler Handler) []Handler {
	stack := make([]Handler, len(r.collector.active.middleware)+1)
	stack = append(stack, handler)

	for i := len(r.collector.active.middleware); i > 0; i-- {
		stack = append(stack, r.collector.active.middleware[i-1](stack[len(stack)-1]))
	}

	return stack
}

// pattern converts path (or prefix) parameters to regular expression patterns
func (r *repository) pattern(path string) string {
	parameters := regexp.MustCompile(`{[^}]+}`).FindAllString(path, -1)
	for _, parameter := range parameters {
		name := parameter[1 : len(parameter)-1]
		optional := false
		if name[len(name)-1:] == "?" {
			name = name[0 : len(name)-1]
			optional = true
		}

		pattern := "(?P<" + name + ">[^/]+?)"
		if definedPattern, exist := r.parameters[name]; exist {
			pattern = "(?P<" + name + ">" + definedPattern + ")"
		}

		if optional {
			pattern += "?"
		}

		path = strings.Replace(path, parameter, pattern, -1)
	}

	return path
}

// addGroup adds a new group of routes to the repository.
func (r *repository) addGroup(prefix string, middleware []Middleware, body func()) {
	r.collector.update(r.pattern(prefix), middleware)
	body()
	r.collector.rollback()
}

// updateGroup update the active group without rollback.
func (r *repository) updateGroup(prefix string, middleware []Middleware) {
	r.collector.update(r.pattern(prefix), middleware)
}

// addParameter adds a new route parameter pattern to the repository.
func (r *repository) addParameter(name, pattern string) {
	r.parameters[name] = pattern
}

// find searches for a route that matches the given HTTP method and URI.
// It returns the route and its parameters.
func (r *repository) find(method, uri string) (*Route, map[string]string, error) {
	for _, route := range r.routes {
		if method == route.Method {
			if ok, parameters := r.match(route.Path, uri); ok {
				return route, parameters, nil
			}
		}
	}

	return nil, nil, errors.New("router: cannot find a route for the request")
}

// match compares the route path with the request URI.
// It will return the boolean result and route parameters if the comparison is successful.
func (r *repository) match(path, uri string) (bool, map[string]string) {
	parameters := map[string]string{}
	pathPattern := regexp.MustCompile("^" + path + "$")
	if pathPattern.MatchString(uri) {
		names := pathPattern.SubexpNames()
		for i, value := range pathPattern.FindAllStringSubmatch(uri, -1)[0] {
			if i > 0 {
				parameters[names[i]] = value
			}
		}

		return true, parameters
	}

	return false, nil
}

// newRepository creates a new repository instance.
func newRepository() *repository {
	return &repository{
		parameters: map[string]string{},
		collector:  newCollector(),
	}
}
