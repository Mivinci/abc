package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mivinci/webkit"
	"github.com/mivinci/webkit/plugins"
)

func main() {
	h := webkit.New()
	h.Use(plugins.Log(os.Stdout))
	h.Use(plugins.IP())
	h.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", r.URL.Path)
	})
	http.ListenAndServe(":8080", h)
}
