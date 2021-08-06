package main

import (
	"net/http"

	"github.com/mivinci/webkit/v2"
)

func main() {
	h := webkit.New()

	h.Get("/ip", func(c webkit.Ctx) error {
		return c.Text(http.StatusOK, c.RealIP())
	})

	h.Run(":8080")
}
