package router

import "net/http"

type Context interface {
	Request() *http.Request
	SetRequest(*http.Request)
	ResponseWriter() http.ResponseWriter
	SetResponseWriter(http.ResponseWriter)
}

type context struct {
	request *http.Request
	rw      http.ResponseWriter
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) SetRequest(request *http.Request) {
	c.request = request
}

func (c *context) ResponseWriter() http.ResponseWriter {
	return c.rw
}

func (c *context) SetResponseWriter(rw http.ResponseWriter) {
	c.rw = rw
}
