package router

import (
	"fmt"
	"net/http"
	"net/url"
)

type Director struct {
	repository *Repository
}

func (d *Director) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	uri, err := url.ParseRequestURI(request.RequestURI)
	if err != nil {
		fmt.Println(err)
		return
	}

	if route, err := d.repository.Find(request.Method, uri.Path); err == nil {
		c := &context{}
		c.SetRequest(request)
		c.SetResponseWriter(rw)

		err = route.Handler(c)
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	fmt.Println("404", uri)
}
