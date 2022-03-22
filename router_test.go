package router_test

import (
	"errors"
	"github.com/golobby/router"
	"github.com/golobby/router/pkg/response"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
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
		URL:        &url.URL{Path: path},
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

	r.GET("/products/:id", func(c router.Context) error {
		return c.Text(200, c.Parameter("id"))
	})
	r.GET("/poly/:id", func(c router.Context) error {
		if c.HasParameter("id") {
			return c.Text(200, "has")
		} else {
			return c.Text(200, "has not")
		}
	})
	r.GET("/poly/:word", func(c router.Context) error {
		return c.Text(200, c.Parameter("word"))
	})
	r.GET("/word/:word/before", func(c router.Context) error {
		return c.Text(200, c.Parameter("word")+" before")
	})
	r.GET("/word/:word", func(c router.Context) error {
		return c.Text(200, c.Parameter("word"))
	})
	r.GET("/word/:word/after", func(c router.Context) error {
		return c.Text(200, c.Parameter("word")+" after")
	})
	r.GET("/multiple/:a/:b", func(c router.Context) error {
		return c.Text(200, c.Parameter("a")+c.Parameter("b"))
	})
	r.GET("/multiple/:a/:b/:c", func(c router.Context) error {
		return c.Text(200, strconv.Itoa(len(c.Parameters())))
	})
	r.GET("/no-parameter", func(c router.Context) error {
		return c.Text(200, c.Parameter("id"))
	})

	var rw *responseWriter

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/products/13"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "13", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/products/test"))
	assert.Equal(t, 404, rw.status)

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/poly/13"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "has", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/poly/test"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "test", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/word/test/before"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "test before", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/word/test"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "test", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/word/test/after"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "test after", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/multiple/1/2"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "12", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/multiple/1/2/3"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "3", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/no-parameter"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "", rw.stringBody())
}

func TestRouter_With_Both_Static_Part_And_Parameter(t *testing.T) {
	r := router.New()

	r.GET("/", func(c router.Context) error {
		return c.Text(200, "home")
	})
	r.GET("/page", func(c router.Context) error {
		return c.Text(200, "page")
	})
	r.GET("/:id", func(c router.Context) error {
		return c.Text(200, "id")
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "home", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/page"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "page", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/13"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "id", rw.stringBody())
}

func TestRouter_With_Wildcard(t *testing.T) {
	r := router.New()

	r.GET("/", func(c router.Context) error {
		return c.Text(200, "home")
	})
	r.GET("/page", func(c router.Context) error {
		return c.Text(200, "page")
	})
	r.GET("/files/*", func(c router.Context) error {
		return c.Text(200, "wildcard-files")
	})
	r.GET("/*", func(c router.Context) error {
		return c.Text(200, "wildcard-all")
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "home", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/page"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "page", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/files/abc"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "wildcard-files", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/files/abc/123"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "wildcard-files", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/abc"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "wildcard-all", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/abc/123"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "wildcard-all", rw.stringBody())
}

func TestRouter_With_Context_Parameters(t *testing.T) {
	r := router.New()
	r.GET("/", func(c router.Context) error {
		return c.Text(200, c.Route().Method()+" "+c.Route().Path())
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "GET /", rw.stringBody())
}

func TestRouter_With_File_Response(t *testing.T) {
	r := router.New()
	r.GET("/", func(c router.Context) error {
		return c.File(200, "text/plain", "assets/text.txt")
	})
	r.GET("/404", func(c router.Context) error {
		return c.File(200, "text/plain", "assets/no-file")
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "This is a text file.", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/404"))
	assert.Equal(t, 500, rw.status)
}

func TestRouter_With_Serving_Static_Files(t *testing.T) {
	r := router.New()
	r.Files("/files/notes/*", "assets/notes")
	r.Files("/*", "assets")

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/files/notes/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "<p>This is notes index.</p>", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/files/notes/note1.txt"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "This is note 1.", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "<p>This is root index.</p>", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/text.txt"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "This is a text file.", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/notes/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "<p>This is notes index.</p>", rw.stringBody())
}

func TestRouter_With_Route_Names(t *testing.T) {
	r := router.New()
	r.GET("/", func(c router.Context) error {
		return c.Text(200, c.URL("home", nil))
	}).SetName("home")
	r.GET("/single/:id", func(c router.Context) error {
		return c.Text(200, c.URL("single", map[string]string{"id": "13"}))
	}).SetName("single")
	r.GET("/multi/:one/:two", func(c router.Context) error {
		return c.Text(200, c.URL("multi", map[string]string{"one": "13", "two": "33"}))
	}).SetName("multi")
	r.GET("/else/:id", func(c router.Context) error {
		return c.Text(200, c.URL("other", map[string]string{"id": "13"}))
	}).SetName("else")
	r.GET("/name", func(c router.Context) error {
		return c.Text(200, c.Route().Name())
	}).SetName("name")

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "/", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/single/1"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "/single/13", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/multi/1/2"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "/multi/13/33", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/else/1"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/name"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "name", rw.stringBody())
}

func TestRouter_WithPrefix(t *testing.T) {
	r := router.New()
	r.WithPrefix("/path", func() {
		r.WithPrefix("/to", func() {
			r.GET("/page", func(c router.Context) error {
				return c.Text(200, "Page1")
			})
		})
		r.GET("/page", func(c router.Context) error {
			return c.Text(200, "Page2")
		})
	})
	r.GET("/page", func(c router.Context) error {
		return c.Text(200, "Page3")
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/path/to/page"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "Page1", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/path/page"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "Page2", rw.stringBody())

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/page"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "Page3", rw.stringBody())
}

func TestRouter_AddPrefix(t *testing.T) {
	r := router.New()
	r.AddPrefix("/content")
	r.GET("/page", func(c router.Context) error {
		return c.Text(200, "OK")
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/content/page"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "OK", rw.stringBody())
}

func TestRouter_WithMiddleware(t *testing.T) {
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

func TestRouter_AddMiddleware(t *testing.T) {
	r := router.New()
	r.AddMiddleware(Middleware1)
	r.GET("/", func(c router.Context) error {
		return c.Text(200, c.Response().Header().Get("Middleware1"))
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "Middleware1", rw.stringBody())
}

func TestRouter_WithMiddlewares(t *testing.T) {
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

func TestRouter_AddMiddlewares(t *testing.T) {
	r := router.New()
	r.AddMiddlewares([]router.Middleware{Middleware1, Middleware2})
	r.GET("/", func(c router.Context) error {
		b := c.Response().Header().Get("Middleware1") + " " + c.Response().Header().Get("Middleware2")
		return c.Text(200, b)
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "Middleware1 Middleware2", rw.stringBody())
}

func TestRouter_SetNotFoundHandler(t *testing.T) {
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

func TestRouter_Internal_Error(t *testing.T) {
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

func TestRouter_With_Different_Responses(t *testing.T) {
	r := router.New()
	r.GET("/empty", func(c router.Context) error {
		return c.Empty(200)
	})
	r.GET("/redirect", func(c router.Context) error {
		return c.Redirect(301, "https://miladrahimi.com")
	})
	r.GET("/text", func(c router.Context) error {
		return c.Text(200, "Text")
	})
	r.GET("/html", func(c router.Context) error {
		return c.HTML(200, "<p>HTML</p>")
	})
	r.GET("/json", func(c router.Context) error {
		return c.JSON(200, response.M{"message": "JSON"})
	})
	r.GET("/json-fail", func(c router.Context) error {
		return c.JSON(200, func() {})
	})
	r.GET("/json-pretty", func(c router.Context) error {
		return c.PrettyJSON(200, response.M{"message": "JSON"})
	})
	r.GET("/json-pretty-fail", func(c router.Context) error {
		return c.PrettyJSON(200, func() {})
	})
	r.GET("/xml", func(c router.Context) error {
		return c.XML(200, struct {
			XMLName struct{} `xml:"User"`
		}{})
	})
	r.GET("/xml-fail", func(c router.Context) error {
		return c.XML(200, func() {})
	})
	r.GET("/xml-pretty", func(c router.Context) error {
		return c.PrettyXML(200, struct {
			XMLName struct{} `xml:"User"`
		}{})
	})
	r.GET("/xml-pretty-fail", func(c router.Context) error {
		return c.PrettyXML(200, func() {})
	})

	rw := newResponse()
	r.Serve(rw, newRequest("GET", "/empty"))
	assert.Equal(t, 200, rw.status)
	assert.Equal(t, "", rw.stringBody())
	assert.Equal(t, "", rw.Header().Get("Content-Type"))

	rw = newResponse()
	r.Serve(rw, newRequest("GET", "/redirect"))
	assert.Equal(t, 301, rw.status)
	assert.Equal(t, "<a href=\"https://miladrahimi.com\">Moved Permanently</a>.\n\n", rw.stringBody())
	assert.Equal(t, "https://miladrahimi.com", rw.Header().Get("Location"))

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
