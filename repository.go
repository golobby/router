package router

import (
	"errors"
)

// repository holds the defined routes, groups, and patterns.
type repository struct {
	routes    []*Route
	matcher   *matcher
	collector *collector
}

// addRoute adds a new route to the repository.
func (r *repository) addRoute(method, path string, handler Handler) {
	r.routes = append(r.routes, &Route{
		r.collector.active.prefix + path,
		method,
		r.collector.active.middleware,
		handler,
	})
}

// addGroup adds a new collector of routes to the repository.
func (r *repository) addGroup(prefix string, middleware []Middleware, body func()) {
	r.collector.update(prefix, middleware)
	body()
	r.collector.rollback()
}

// addParameter adds a new route parameter pattern to the repository.
func (r *repository) addParameter(name, pattern string) {
	r.matcher.addParameter(name, pattern)
}

// find searches for a route that matches the given HTTP method and URI.
// It returns the route and its parameters.
func (r *repository) find(method, uri string) (*Route, map[string]string, error) {
	for _, route := range r.routes {
		if method == route.Method {
			if ok, parameters := r.matcher.match(route.Path, uri); ok {
				return route, parameters, nil
			}
		}
	}

	return nil, nil, errors.New("router: cannot find a route for the request")
}

// newRepository creates a new repository for the new router.
func newRepository() *repository {
	return &repository{
		matcher:   newMatcher(),
		collector: newCollector(),
	}
}
