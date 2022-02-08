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

func (r *responseWriter) stringBody() string {
	return string(r.body)
}

func newResponse() *responseWriter {
	return &responseWriter{
		status:  0,
		headers: http.Header{},
		body:    []byte(""),
	}
}

// Testing HTTP request builder

func newRequest(method, path string) *http.Request {
	return &http.Request{
		Method:     method,
		RequestURI: path,
	}
}

// Testing middlewares

func Middleware1(next router.Handler) router.Handler {
	return func(c router.Context) error {
		c.Response().Header().Set("Middleware1", "Middleware1")
		return next(c)
	}
}

func Middleware2(next router.Handler) router.Handler {
	return func(c router.Context) error {
		c.Response().Header().Set("Middleware2", "Middleware2")
		return next(c)
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

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "CUSTOM"}
	for _, m := range methods {
		rw := newResponse()
		r.Serve(rw, newRequest(m, "/"))
		assert.Equal(t, 200, rw.status)
		assert.Equal(t, m, rw.stringBody())
	}

	rw := newResponse()
	r.Serve(rw, newRequest("FAIL", "/"))
	assert.Equal(t, 404, rw.status)
}

func TestRouter_With_Route_Parameters(t *testing.T) {
	r := router.New()
	r.Define("id", "[0-9]+")

	r.GET("/{id}", func(c router.Context) error {
		return c.Text(200, c.Parameter("id"))
	})
	r.GET("/{word}/before", func(c router.Context) error {
		return c.Text(200, c.Parameter("word")+" before")
	})
	r.GET("/{word}", func(c router.Context) error {
		return c.Text(200, c.Parameter("word"))
	})
	r.GET("/{word}/after", func(c router.Context) error {
		return c.Text(200, c.Parameter("word")+" after")
	})
	r.GET("/fail/{id}", func(c router.Context) error {
		return c.Text(200, c.Parameter("id"))
	})
	r.GET("/multi/{a}/{b}/{c}", func(c router.Context) error {
		return c.Text(200, c.Parameter("a")+c.Parameter("b")+c.Parameter("c"))
	})
	r.GET("/optional/page/{id?}", func(c router.Context) error {
		return c.Text(200, c.Parameter("id"))
	})
	r.GET("/optional/page2/?{id?}", func(c router.Context) error {
		return c.Text(200, c.Parameter("id"))
	})

	var rw *responseWriter

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/13"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "13", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/test/before"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "test before", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/test"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "test", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/test/after"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "test after", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/fail/test"))
	assert.Equal(t, 404, rw.status)

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/multi/1/2/3"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "123", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/optional/page/13"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "13", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/optional/page/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/optional/page"))
	assert.Equal(t, 404, rw.status)

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/optional/page2"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "", rw.stringBody())
}

func TestRouter_With_Prefix(t *testing.T) {
	r := router.New()
	r.WithPrefix("/prefix", func() {
		r.GET("/page", func(c router.Context) error {
			return c.Text(200, "OK")
		})
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/prefix/page"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "OK", rw.stringBody())
}

func TestRouter_With_Middleware(t *testing.T) {
	r := router.New()
	r.WithMiddleware(Middleware1, func() {
		r.GET("/", func(c router.Context) error {
			return c.Text(200, c.Response().Header().Get("Middleware1"))
		})
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "Middleware1", rw.stringBody())
}

func TestRouter_With_Middlewares(t *testing.T) {
	r := router.New()
	r.WithMiddlewares([]router.Middleware{Middleware1, Middleware2}, func() {
		r.GET("/", func(c router.Context) error {
			b := c.Response().Header().Get("Middleware1") + " " + c.Response().Header().Get("Middleware2")
			return c.Text(200, b)
		})
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "Middleware1 Middleware2", rw.stringBody())
}
