package router

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// Context holds the HTTP request, the HTTP responseWriter, the Route, and the Route parameters.
type Context interface {
	// Route returns the dispatched Route
	Route() *Route

	// SetRoute sets the dispatched Route
	SetRoute(route *Route)

	// Request returns the HTTP request.
	Request() *http.Request

	// SetRequest sets the HTTP request.
	SetRequest(request *http.Request)

	// Response return the HTTP responseWriter.
	Response() http.ResponseWriter

	// SetResponse sets the HTTP responseWriter.
	SetResponse(rw http.ResponseWriter)

	// Parameters returns Route parameters.
	Parameters() map[string]string

	// SetParameters sets the router parameters.
	SetParameters(parameters map[string]string)

	// Parameter returns a router parameter by name.
	Parameter(name string) string

	// HasParameter checks if router parameter exists.
	HasParameter(name string) bool

	// Status sets the HTTP responseWriter status code.
	Status(status int)

	// Bytes creates and sends a custom HTTP response.
	Bytes(status int, body []byte) error

	// Empty creates and sends an HTTP empty response.
	Empty(status int) error

	// Redirect creates and sends an HTTP redirection response.
	Redirect(status int, url string) error

	// Text creates and sends an HTTP text response.
	Text(status int, body string) error

	// HTML creates and sends an HTTP HTML response.
	HTML(status int, body string) error

	// JSON creates and sends an HTTP JSON response.
	JSON(status int, body interface{}) error

	// PrettyJSON creates and sends an HTTP JSON (with indents) response.
	PrettyJSON(status int, body interface{}) error

	// XML creates and sends an HTTP XML response.
	XML(status int, body interface{}) error

	// PrettyXML creates and sends an HTTP XML (with indents) response.
	PrettyXML(status int, body interface{}) error
}

// DefaultContext is the default implementation of Context interface.
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

func (d *DefaultContext) Response() http.ResponseWriter {
	return d.rw
}

func (d *DefaultContext) SetResponse(rw http.ResponseWriter) {
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

func (d *DefaultContext) Status(status int) {
	d.rw.WriteHeader(status)
}

func (d *DefaultContext) Bytes(status int, body []byte) error {
	d.Status(status)
	_, err := d.rw.Write(body)
	return err
}

func (d *DefaultContext) Empty(status int) error {
	d.Status(status)
	return nil
}

func (d *DefaultContext) Redirect(status int, url string) error {
	http.Redirect(d.Response(), d.Request(), url, status)
	return nil
}

func (d *DefaultContext) Text(status int, body string) error {
	d.Response().Header().Set("Content-Type", "text/plain")
	return d.Bytes(status, []byte(body))
}

func (d *DefaultContext) HTML(status int, body string) error {
	d.Response().Header().Set("Content-Type", "text/html")
	return d.Bytes(status, []byte(body))
}

func (d *DefaultContext) JSON(status int, body interface{}) error {
	d.Response().Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return d.Bytes(status, bytes)
}

func (d *DefaultContext) PrettyJSON(status int, body interface{}) error {
	d.Response().Header().Set("Content-Type", "application/json")
	bytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}
	return d.Bytes(status, bytes)
}

func (d *DefaultContext) XML(status int, body interface{}) error {
	d.Response().Header().Set("Content-Type", "application/xml")
	bytes, err := xml.MarshalIndent(body, "", "")
	if err != nil {
		return err
	}
	return d.Bytes(status, bytes)
}

func (d *DefaultContext) PrettyXML(status int, body interface{}) error {
	d.Response().Header().Set("Content-Type", "application/xml")
	bytes, err := xml.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}
	return d.Bytes(status, bytes)
}
