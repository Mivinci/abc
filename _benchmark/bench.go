package main

import (
	"fmt"
	"net/http"

	"github.com/mathoj/webkit"
)

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	ps := webkit.ParamsFromContext(r.Context())
	fmt.Fprintf(w, "Hello %v\n", ps["name"])
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func main() {
	h := webkit.New()
	h.Get("/some/page/:name", dynamicHandler)
	h.Get("/other/page/path", staticHandler)
	http.ListenAndServe(":8080", h)
}
