package router

// Handler is an interface for route handlers (controllers).
// When the router finds a route for the incoming HTTP request, it calls the route's handler.
type Handler func(c Context) error
