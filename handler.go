package router

import (
	"fmt"
	"net/http"
)

// Handler is an interface for Route handlers (controllers).
// When the router finds a Route for the incoming HTTP request, it calls the Route's handler.
type Handler func(c Context) error

// FileHandler is a special handler for serving static files.
// It returns files stored in the given root path that matches request URI.
func FileHandler(path string) Handler {
	return func(c Context) error {
		fs := http.FileServer(http.Dir(path))
		fs.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// FileHandlerWithStripper is a special handler for serving static files.
// It returns files stored in the given root path that matches the stripped request URI.
func FileHandlerWithStripper(path, strip string) Handler {
	return func(c Context) error {
		fs := http.FileServer(http.Dir(path))
		fmt.Println("P", strip, path)
		h := http.StripPrefix(strip, fs)
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
