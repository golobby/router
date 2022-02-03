package main

import (
	"github.com/golobby/router"
	"log"
)

func main() {
	r := router.New()

	r.Get("/", func(c router.Context) error {
		return c.Empty(204)
	})

	r.Get("/text", func(c router.Context) error {
		return c.Text(200, "It's a text response")
	})

	r.Get("/json", func(c router.Context) error {
		s := struct {
			Message string `json:"message"`
		}{"It's a JSON response!"}

		return c.Json(200, s)
	})

	r.Get("/json-pretty", func(c router.Context) error {
		s := struct {
			Message string `json:"message"`
		}{"It's a JSON response!"}

		return c.JsonPretty(200, s)
	})

	r.Get("/xml", func(c router.Context) error {
		return c.Xml(200, []int{1, 2, 3})
	})

	r.Get("/xml", func(c router.Context) error {
		return c.Xml(200, []int{1, 2, 3})
	})

	log.Fatalln(r.Start(":8000"))
}
