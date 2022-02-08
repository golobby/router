package router_test

import (
	"github.com/golobby/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// Testing HTTP response writer

type responseWriter struct {
	status  int
	body    []byte
	headers http.Header
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *responseWriter) Write(body []byte) (int, error) {
	r.body = body
	return 0, nil
}

func (r *responseWriter) Header() http.Header {
	return r.headers
}

func newResponse() *responseWriter {
	return &responseWriter{headers: http.Header{}}
}

// Testing HTTP request builder

func newRequest(method, path string) *http.Request {
	return &http.Request{
		Method:     method,
		RequestURI: path,
	}
}

// Integrated tests

func TestRouter_With_Different_HTTP_Methods(t *testing.T) {
	handler := func(c router.Context) error {
		return c.Text(200, c.Request().Method)
	}

	r := router.New()
	r.GET("/", handler)
	r.POST("/", handler)
	r.PUT("/", handler)
	r.PATCH("/", handler)
	r.DELETE("/", handler)
	r.HEAD("/", handler)
	r.OPTIONS("/", handler)
	r.Map("CUSTOM", "/", handler)

	var rw *responseWriter
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "CUSTOM"}

	for _, m := range methods {
		rw = newResponse()
		r.Director().ServeHTTP(rw, newRequest(m, "/"))
		assert.Equal(t, 200, rw.status)
		assert.Equal(t, []byte(m), rw.body)
	}
}
