package router

import (
	"encoding/json"
	"net/http"
)

// Handler is a type/interface for route handlers (controllers).
// Whenever the router finds a route for the incoming HTTP request, it calls the route's handler.
type Handler func(c Context) error
type Response func(w http.ResponseWriter) error
type Handler2 func(r *http.Request) Response

func Json(status int, obj interface{}) func(writer http.ResponseWriter) error {
	return func(writer http.ResponseWriter) error {
		bs, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		writer.Write(bs)
		writer.WriteHeader(status)
		return nil
	}
}
func wrap(h Handler2) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		res := h(request)
		res(writer)
	}
}
func myHandler(r *http.Request) Response {
	return Json(200, map[string]string{"message": "hello world"})
}

func main() {

}
