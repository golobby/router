package router

import "net/http"

// Context holds the HTTP request, the HTTP response writer, the route, and the route parameters.
type Context interface {
	// Request returns the HTTP request.
	Request() *http.Request

	// SetRequest sets the HTTP request.
	SetRequest(*http.Request)

	// ResponseWriter return the HTTP response writer.
	ResponseWriter() http.ResponseWriter

	// SetResponseWriter sets the HTTP response writer.
	SetResponseWriter(http.ResponseWriter)

	// Parameters returns route parameters.
	Parameters() map[string]string

	// SetParameters sets the router parameters.
	SetParameters(map[string]string)

	// Parameter returns a router parameter by name.
	Parameter(name string) string

	// HasParameter checks if router parameter exists.
	HasParameter(name string) bool
}

// DefaultContext is the default implementation of Context
type DefaultContext struct {
	request    *http.Request
	rw         http.ResponseWriter
	parameters map[string]string
}

func (d *DefaultContext) Request() *http.Request {
	return d.request
}

func (d *DefaultContext) SetRequest(request *http.Request) {
	d.request = request
}

func (d *DefaultContext) ResponseWriter() http.ResponseWriter {
	return d.rw
}

func (d *DefaultContext) SetResponseWriter(rw http.ResponseWriter) {
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
