package router

// Middleware is an interface for route middleware.
// It returns a Handler that receives HTTP DefaultContext to watch or manipulate and calls the next middleware/handler.
type Middleware func(next Handler) Handler
