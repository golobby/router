package router

type Route struct {
	Path    string
	Method  string
	Handler Handler
}
