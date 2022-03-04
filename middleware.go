package router

// Middleware is an interface for Route middlewares.
// It returns a Handler that receives HTTP Context to watch or manipulate and calls the next middlewares/handler.
type Middleware func(next Handler) Handler
