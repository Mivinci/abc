package plugins

import (
	"net/http"

	"github.com/mivinci/webkit"
)

func RateLimit() webkit.Plugin {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
}
