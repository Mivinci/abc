package main

import (
	"net/http"
	"strconv"

	"github.com/mivinci/webkit/v2"
)

func main() {
	h := webkit.New(
		webkit.Encode(webkit.StdEncoder()),
	)

	h.Get("/:id", func(c webkit.Ctx) error {
		id, err := c.Params().Int("id")
		if err != nil {
			return webkit.HTTPError(http.StatusBadRequest)
		}
		return c.Text(http.StatusOK, strconv.Itoa(id))
	})

	http.ListenAndServe(":8080", h)
}
