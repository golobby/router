package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Handler is an interface for Route handlers (controllers).
// When the router finds a Route for the incoming HTTP request, it calls the Route's handler.
type Handler func(c Context) error

func MakeHandler[IN any, OUT any](h func(IN) (OUT, error)) Handler {
	return func(c Context) error {
		bs, _ := ioutil.ReadAll(c.Request().Body)
		defer c.Request().Body.Close()

		body := new(IN)
		_ = json.Unmarshal(bs, body)

		out, err := h(*body)

		if err != nil {
			fmt.Fprintf(c.Response(), "err: %s", err.Error())
			return nil
		}

		bs, _ = json.Marshal(out)
		c.Response().Write(bs)
		return nil
	}
}
