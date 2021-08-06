package main

import (
	"embed"
	"net/http"

	"github.com/mivinci/webkit/v2"
)

//go:embed html
var fs embed.FS

func main() {
	h := webkit.New(
		webkit.TemplateFS(fs, "html/*.html"),
	)
	h.Get("/", func(c webkit.Ctx) error {
		data := map[string]string{
			"name": "Webkit",
		}
		return c.Template(http.StatusOK, "index.html", data)
	})
	h.Run(":8080")
}
