package main

import (
	"github.com/golobby/router"
	"log"
	"net/http"
)

func MyHandler(c router.Context) error {
	return c.Text(http.StatusOK, "It is a "+c.Request().Method+" Request")
}

func main() {
	r := router.New()

	r.Get("/", MyHandler)
	r.Post("/", MyHandler)
	r.Put("/", MyHandler)
	r.Patch("/", MyHandler)
	r.Delete("/", MyHandler)
	r.Map("CUSTOM", "/", MyHandler)

	log.Fatalln(r.Start(":8000"))
}
