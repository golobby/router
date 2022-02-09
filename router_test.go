package router_test

import (
	"errors"
	"github.com/golobby/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// Testing HTTP response writer

type responseWriter struct {
	status int
	body   []byte
	header http.Header
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *responseWriter) Write(body []byte) (int, error) {
	r.body = body
	return 0, nil
}

func (r *responseWriter) Header() http.Header {
	return r.header
}

func (r *responseWriter) stringBody() string {
	return string(r.body)
}

func newResponse() *responseWriter {
	return &responseWriter{
		status: 0,
		header: http.Header{},
		body:   []byte(""),
	}
}

// Testing HTTP request builder

func newRequest(method, path string) *http.Request {
	return &http.Request{
		Method:     method,
		RequestURI: path,
	}
}

// Common values

const InternalErrorJson = "{\"message\":\"Internal error.\"}"

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
		return c.Text(200, strconv.Itoa(len(c.Parameters())))
	})
	r.GET("/optional/page/{id?}", func(c router.Context) error {
		return c.Text(200, c.Parameter("id"))
	})
	r.GET("/optional/page2/?{id?}", func(c router.Context) error {
		return c.Text(200, c.Parameter("id"))
	})
	r.GET("/else/no-parameter", func(c router.Context) error {
		if c.HasParameter("id") {
			return c.Text(200, "Yes and "+c.Parameter("id"))
		} else {
			return c.Text(200, "No but "+c.Parameter("id"))
		}
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
	assert.Equal(t, "3", rw.stringBody())

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

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/else/no-parameter"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "No but ", rw.stringBody())
}

func TestRouter_With_Context_Parameters(t *testing.T) {
	r := router.New()
	r.GET("/", func(c router.Context) error {
		return c.Text(200, c.Route().Method+" "+c.Route().Path)
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "GET /", rw.stringBody())
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

func TestRouter_With_404_Error(t *testing.T) {
	r := router.New()

	r.SetNotFoundHandler(func(c router.Context) error {
		return c.Text(504, "New Not Found")
	})

	r.GET("/%", func(c router.Context) error {
		return c.Text(200, "OK")
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 504, rw.status)
	assert.Equal(t, "New Not Found", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/%"))
	assert.Equal(t, 504, rw.status)
	assert.Equal(t, "New Not Found", rw.stringBody())
}

func TestRouter_With_Internal_Error(t *testing.T) {
	r := router.New()

	r.SetNotFoundHandler(func(c router.Context) error {
		return errors.New("error inside error handler")
	})

	r.GET("/error", func(c router.Context) error {
		return errors.New("error")
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/error"))
	assert.Equal(t, 500, rw.status)
	assert.Equal(t, InternalErrorJson, rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 500, rw.status)
	assert.Equal(t, InternalErrorJson, rw.stringBody())

}

func TestRouter_Start(t *testing.T) {
	r := router.New()
	r.GET("/", func(c router.Context) error {
		return c.Text(200, "OK")
	})

	ec := make(chan error)
	go func() {
		ec <- r.Start(":8585")
	}()

	select {
	case err := <-ec:
		assert.Fail(t, err.Error())
	case <-time.After(3 * time.Second):
		assert.True(t, true)
	}
}

func TestRouter_With_Different_Responses(t *testing.T) {
	r := router.New()
	r.GET("/empty", func(c router.Context) error {
		return c.Empty(200)
	})
	r.GET("/text", func(c router.Context) error {
		return c.Text(200, "Text")
	})
	r.GET("/html", func(c router.Context) error {
		return c.Html(200, "<p>HTML</p>")
	})
	r.GET("/json", func(c router.Context) error {
		return c.Json(200, router.S{"message": "JSON"})
	})
	r.GET("/json-fail", func(c router.Context) error {
		return c.Json(200, func() {})
	})
	r.GET("/json-pretty", func(c router.Context) error {
		return c.JsonPretty(200, router.S{"message": "JSON"})
	})
	r.GET("/json-pretty-fail", func(c router.Context) error {
		return c.JsonPretty(200, func() {})
	})
	r.GET("/xml", func(c router.Context) error {
		return c.Xml(200, struct {
			XMLName struct{} `xml:"User"`
		}{})
	})
	r.GET("/xml-fail", func(c router.Context) error {
		return c.Xml(200, func() {})
	})
	r.GET("/xml-pretty", func(c router.Context) error {
		return c.XmlPretty(200, struct {
			XMLName struct{} `xml:"User"`
		}{})
	})
	r.GET("/xml-pretty-fail", func(c router.Context) error {
		return c.XmlPretty(200, func() {})
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/empty"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "", rw.stringBody())
	assert.Equal(t, "", rw.Header().Get("Content-Type"))

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/text"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "Text", rw.stringBody())
	assert.Equal(t, "text/plain", rw.Header().Get("Content-Type"))

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/html"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "<p>HTML</p>", rw.stringBody())
	assert.Equal(t, "text/html", rw.Header().Get("Content-Type"))

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/json"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "{\"message\":\"JSON\"}", rw.stringBody())
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/json-fail"))
	assert.Equal(t, 500, rw.status)
	assert.Equal(t, InternalErrorJson, rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/json-pretty"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "{\n  \"message\": \"JSON\"\n}", rw.stringBody())
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/json-pretty-fail"))
	assert.Equal(t, 500, rw.status)
	assert.Equal(t, InternalErrorJson, rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/xml"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "<User></User>", rw.stringBody())
	assert.Equal(t, "application/xml", rw.Header().Get("Content-Type"))

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/xml-fail"))
	assert.Equal(t, 500, rw.status)
	assert.Equal(t, InternalErrorJson, rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/xml-pretty"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "<User></User>", rw.stringBody())
	assert.Equal(t, "application/xml", rw.Header().Get("Content-Type"))

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/xml-pretty-fail"))
	assert.Equal(t, 500, rw.status)
	assert.Equal(t, InternalErrorJson, rw.stringBody())
}
