package router

// Handler is an interface for Route handlers (controllers).
// When the router finds a Route for the incoming HTTP request, it calls the Route's handler.
type Handler func(c Context) error
