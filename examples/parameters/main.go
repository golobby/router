package main

import (
	"github.com/golobby/router"
	"log"
	"net/http"
)

func MyHandler(c router.Context) error {
	return c.Json(http.StatusOK, c.Parameters())
}

func main() {
	r := router.New()

	r.Define("id", "[0-9]+")

	r.Get("/{id}", MyHandler)
	r.Get("/{name}", MyHandler)
	r.Get("/{id}/more/{p2}/{p3}", MyHandler)

	log.Fatalln(r.Start(":8000"))
}
