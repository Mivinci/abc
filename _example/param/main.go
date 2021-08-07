package main

import (
	"net/http"

	"github.com/mivinci/webkit/v2"
)

func f(next webkit.Handler) webkit.Handler {
	return func(c webkit.Ctx) error {
		ps := c.Params()
		ps["name"] = "webkit"
		return next(c)
	}
}

func main() {
	h := webkit.New()
	h.Use(f)
	h.Get("/params/:id", func(c webkit.Ctx) error {
		ps := c.Params()
		return c.JSON(http.StatusOK, ps)
		// {
		//   "id": 123,
		//   "name": "webkit"
		// }
	})
	h.Run(":8080")
}
