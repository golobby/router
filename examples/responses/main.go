package main

import (
	"github.com/golobby/router"
	"log"
)

func main() {
	r := router.New()

	r.GET("/", func(c router.Context) error {
		return c.Empty(204)
	})

	r.GET("/text", func(c router.Context) error {
		return c.Text(200, "It's a text response")
	})

	r.GET("/html", func(c router.Context) error {
		return c.Html(200, "<h1>HTML</h1><p>This is paragraph</p>")
	})

	r.GET("/json", func(c router.Context) error {
		s := struct {
			Message string `json:"message"`
		}{"It's a JSON response!"}

		return c.Json(200, s)
	})

	r.GET("/json-pretty", func(c router.Context) error {
		s := struct {
			Message string `json:"message"`
		}{"It's a JSON response!"}

		return c.JsonPretty(200, s)
	})

	r.GET("/xml", func(c router.Context) error {
		return c.Xml(200, []int{1, 2, 3})
	})

	r.GET("/xml-pretty", func(c router.Context) error {
		return c.XmlPretty(200, []int{1, 2, 3})
	})

	log.Fatalln(r.Start(":8000"))
}
