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

	r.GET("/", MyHandler)
	r.POST("/", MyHandler)
	r.PUT("/", MyHandler)
	r.PATCH("/", MyHandler)
	r.DELETE("/", MyHandler)
	r.Map("CUSTOM", "/", MyHandler)

	log.Fatalln(r.Start(":8000"))
}
