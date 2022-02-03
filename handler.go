package router

// Handler is a type/interface for route handlers (controllers).
// Whenever the router finds a route for the incoming HTTP request, it calls the route's handler.
type Handler func(c Context) error
