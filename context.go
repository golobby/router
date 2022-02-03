package router

import (
	"encoding/json"
	"net/http"
)

// Context holds the HTTP request, the HTTP response writer, the route, and the route parameters.
type Context interface {
	// Route returns the dispatched route
	Route() *Route

	// SetRoute sets the dispatched route
	SetRoute(route *Route)

	// Request returns the HTTP request.
	Request() *http.Request

	// SetRequest sets the HTTP request.
	SetRequest(request *http.Request)

	// RW return the HTTP response writer.
	RW() http.ResponseWriter

	// SetRW sets the HTTP response writer.
	SetRW(rw http.ResponseWriter)

	// Parameters returns route parameters.
	Parameters() map[string]string

	// SetParameters sets the router parameters.
	SetParameters(parameters map[string]string)

	// Parameter returns a router parameter by name.
	Parameter(name string) string

	// HasParameter checks if router parameter exists.
	HasParameter(name string) bool

	// Empty creates and sends an HTTP empty response
	Empty(status int) error

	// Json creates and sends an HTTP JSON response
	Json(status int, body interface{}) error

	// Text creates and sends an HTTP text response
	Text(status int, body string) error
}

// DefaultContext is the default implementation of Context
type DefaultContext struct {
	route      *Route
	request    *http.Request
	rw         http.ResponseWriter
	parameters map[string]string
}

func (d *DefaultContext) Route() *Route {
	return d.route
}

func (d *DefaultContext) SetRoute(route *Route) {
	d.route = route
}

func (d *DefaultContext) Request() *http.Request {
	return d.request
}

func (d *DefaultContext) SetRequest(request *http.Request) {
	d.request = request
}

func (d *DefaultContext) RW() http.ResponseWriter {
	return d.rw
}

func (d *DefaultContext) SetRW(rw http.ResponseWriter) {
	d.rw = rw
}

func (d *DefaultContext) Parameters() map[string]string {
	return d.parameters
}

func (d *DefaultContext) SetParameters(parameters map[string]string) {
	d.parameters = parameters
}

func (d *DefaultContext) Parameter(name string) string {
	if value, exist := d.parameters[name]; exist {
		return value
	}
	return ""
}

func (d *DefaultContext) HasParameter(name string) bool {
	_, exist := d.parameters[name]
	return exist
}

func (d *DefaultContext) Empty(status int) error {
	d.rw.WriteHeader(status)
	return nil
}

func (d *DefaultContext) Json(status int, body interface{}) error {
	v, err := json.Marshal(body)
	if err != nil {
		return err
	}

	d.rw.WriteHeader(status)
	d.rw.Header().Set("Content-Type", "application/json")
	_, err = d.rw.Write(v)
	return err
}

func (d *DefaultContext) Text(status int, body string) error {
	d.rw.WriteHeader(status)
	d.rw.Header().Set("Content-Type", "text/plain")
	_, err := d.rw.Write([]byte(body))
	return err
}
