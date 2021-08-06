package main

import (
	"net/http"
	"strconv"

	"github.com/mivinci/webkit/v2"
)

type arg struct {
	ID int `json:"id"`
}

func (a *arg) Bind(c webkit.Ctx) (err error) {
	qs := c.Query()
	a.ID, err = strconv.Atoi(qs.Get("id"))
	return
}

func main() {
	h := webkit.New()
	h.Get("/", func(c webkit.Ctx) error {
		var a arg
		if err := c.BindCtx(&a); err != nil {
			return webkit.HTTPError(http.StatusBadRequest)
		}
		return c.JSON(http.StatusOK, a)
	})
	http.ListenAndServe(":8080", h)
}
