package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mivinci/webkit/v2"
	"github.com/mivinci/webkit/v2/plugin"
)

func basicAuthValidator(user, passwd string) error {
	fmt.Printf("user=%s passwd=%s", user, passwd)
	if user == "a" && passwd == "123" {
		return nil
	}
	return errors.New("invalid credentials")
}

func main() {
	h := webkit.New(
		webkit.Encode(webkit.StdEncoder()),
		webkit.Plugins(
			plugin.BasicAuth("example", basicAuthValidator),
		),
	)

	h.Get("/", func(c webkit.Ctx) error {
		return nil
	})

	http.ListenAndServe(":8080", h)
}
