package main

import (
	"github.com/golobby/router"
	"log"
	"time"
)

func M1(next router.Handler) router.Handler {
	return func(c router.Context) error {
		c.Response().Header().Add("M1", time.Now().String())
		return next(c)
	}
}

func M2(next router.Handler) router.Handler {
	return func(c router.Context) error {
		c.Response().Header().Add("M2", time.Now().String())
		return next(c)
	}
}

func M3(next router.Handler) router.Handler {
	return func(c router.Context) error {
		c.Response().Header().Add("M3", time.Now().String())
		return next(c)
	}
}

func main() {
	r := router.New()

	r.WithMiddleware(M1, func() {
		r.GET("/single", func(c router.Context) error {
			return c.Text(200, "OK")
		})
	})
	r.WithMiddlewareList([]router.Middleware{M1, M2, M3}, func() {
		r.GET("/multiple", func(c router.Context) error {
			return c.Text(200, "OK")
		})
	})

	log.Fatalln(r.Start(":8000"))
}
