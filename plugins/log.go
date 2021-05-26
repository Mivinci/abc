package plugins

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mathoj/webkit"
)

func Log(w io.Writer) webkit.Plugin {
	log.SetOutput(w)

	f := func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			next.ServeHTTP(w, r)
			ip, _ := IPFromContext(r.Context())
			log.Printf("%s %s %s %v", ip, r.Method, r.URL.Path, time.Since(t))
		}
	}

	return func(next http.Handler) http.Handler { return f(next) }
}
