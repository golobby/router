package router

import (
	"errors"
)

type Repository struct {
	routes []*Route
}

func (r *Repository) Save(method, path string, handler Handler) {
	r.routes = append(r.routes, &Route{path, method, handler})
}

func (r *Repository) Find(method, path string) (*Route, error) {
	for _, v := range r.routes {
		if method == v.Method && Match(v.Path, path) {
			return v, nil
		}
	}

	return nil, errors.New("router: cannot find a route for the request")
}

func (r *Repository) All() []*Route {
	return r.routes
}
