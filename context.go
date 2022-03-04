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

	// Request returns the HTTP request.
	Request() *http.Request

	// Response return the HTTP responseWriter.
	Response() http.ResponseWriter

	// Parameters returns Route parameters.
	Parameters() map[string]string

	// Parameter returns a router parameter by name.
	Parameter(name string) string

	// HasParameter checks if router parameter exists.
	HasParameter(name string) bool

	// URL generates a URL for given route name and actual parameters.
	// It returns an empty string if it cannot find any route.
	URL(route string, parameters map[string]string) string

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
	repository *repository
	request    *http.Request
	rw         http.ResponseWriter
	parameters map[string]string
}

// Route returns the dispatched Route
func (d *DefaultContext) Route() *Route {
	return d.route
}

// Request returns the HTTP request.
func (d *DefaultContext) Request() *http.Request {
	return d.request
}

// Response return the HTTP responseWriter.
func (d *DefaultContext) Response() http.ResponseWriter {
	return d.rw
}

// Parameters returns Route parameters.
func (d *DefaultContext) Parameters() map[string]string {
	return d.parameters
}

// Parameter returns a router parameter by name.
func (d *DefaultContext) Parameter(name string) string {
	if value, exist := d.parameters[name]; exist {
		return value
	}
	return ""
}

// HasParameter checks if router parameter exists.
func (d *DefaultContext) HasParameter(name string) bool {
	_, exist := d.parameters[name]
	return exist
}

// URL generates a URL for given route name and actual parameters.
// It returns an empty string if it cannot find any route.
func (d *DefaultContext) URL(route string, parameters map[string]string) string {
	if route := d.repository.findByName(route); route != nil {
		return route.URL(parameters)
	}
	return ""
}

// Bytes creates and sends a custom HTTP response.
func (d *DefaultContext) Bytes(status int, body []byte) error {
	d.rw.WriteHeader(status)
	_, err := d.rw.Write(body)
	return err
}

// Empty creates and sends an HTTP empty response.
func (d *DefaultContext) Empty(status int) error {
	d.rw.WriteHeader(status)
	return nil
}

// Redirect creates and sends an HTTP redirection response.
func (d *DefaultContext) Redirect(status int, url string) error {
	http.Redirect(d.Response(), d.Request(), url, status)
	return nil
}

// Text creates and sends an HTTP text response.
func (d *DefaultContext) Text(status int, body string) error {
	d.Response().Header().Set("Content-Type", "text/plain")
	return d.Bytes(status, []byte(body))
}

// HTML creates and sends an HTTP HTML response.
func (d *DefaultContext) HTML(status int, body string) error {
	d.Response().Header().Set("Content-Type", "text/html")
	return d.Bytes(status, []byte(body))
}

// JSON creates and sends an HTTP JSON response.
func (d *DefaultContext) JSON(status int, body interface{}) error {
	d.Response().Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return d.Bytes(status, bytes)
}

// PrettyJSON creates and sends an HTTP JSON (with indents) response.
func (d *DefaultContext) PrettyJSON(status int, body interface{}) error {
	d.Response().Header().Set("Content-Type", "application/json")
	bytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}
	return d.Bytes(status, bytes)
}

// XML creates and sends an HTTP XML response.
func (d *DefaultContext) XML(status int, body interface{}) error {
	d.Response().Header().Set("Content-Type", "application/xml")
	bytes, err := xml.MarshalIndent(body, "", "")
	if err != nil {
		return err
	}
	return d.Bytes(status, bytes)
}

// PrettyXML creates and sends an HTTP XML (with indents) response.
func (d *DefaultContext) PrettyXML(status int, body interface{}) error {
	d.Response().Header().Set("Content-Type", "application/xml")
	bytes, err := xml.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}
	return d.Bytes(status, bytes)
}
