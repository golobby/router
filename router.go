package router

import (
	"net/http"
)

type Router struct {
	repository *Repository
	director   *Director
}

func (r Router) Define(method, path string, handler Handler) {
	r.repository.Save(method, path, handler)
}

func (r Router) Start(address string) error {
	return http.ListenAndServe(address, r.director)
}

func New() *Router {
	repository := &Repository{}

	return &Router{
		repository: repository,
		director:   &Director{repository},
	}
}
