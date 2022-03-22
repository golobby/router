package router

import (
	"net/http"
	"strings"
)

// Handler is an interface for Route handlers (controllers).
// When the router finds a Route for the incoming HTTP request, it calls the Route's handler.
type Handler func(c Context) error

// FilesHandler is a special handler for serving static files.
// It returns files stored in the given root directory that matches the stripped request URI (path).
func FilesHandler(path, directory string) Handler {
	return func(c Context) error {
		fs := http.FileServer(http.Dir(directory))
		h := http.StripPrefix(strings.TrimRight(path, "*"), fs)
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
