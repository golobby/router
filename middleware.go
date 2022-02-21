package router

// Middleware is an interface for Route middleware.
// It returns a Handler that receives HTTP Context to watch or manipulate and calls the next middleware/handler.
type Middleware func(next Handler) Handler
