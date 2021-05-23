package webkit

import "net/http"

type Plugin func(next http.Handler) http.Handler

// chain is the magic.
func chain(endpoint http.Handler, plugins ...Plugin) http.Handler {
	for i := len(plugins) - 1; i >= 0; i-- {
		endpoint = plugins[i](endpoint)
	}
	return endpoint
}
