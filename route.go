package router

// Route holds route information.
type Route struct {
	Path       string
	Method     string
	Middleware []Middleware
	Handler    Handler
}
