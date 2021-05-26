package plugins

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/mathoj/webkit"
)

func RealIP(r *http.Request) string {
	xRealIP := r.Header.Get("X-Real-Ip")
	xFFor := r.Header.Get("X-Forwarded-For")

	if len(xRealIP) == 0 && len(xFFor) == 0 {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		return ip
	}

	for _, ip := range strings.Split(xFFor, ",") {
		return strings.TrimSpace(ip)
	}

	return xRealIP
}

type ipCtxKey struct{}

func IPFromContext(ctx context.Context) (ip string, ok bool) {
	if i := ctx.Value(ipCtxKey{}); i != nil {
		ip, ok = i.(string)
		return
	}
	return
}

func IP() webkit.Plugin {
	f := func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := RealIP(r)
			ctx := context.WithValue(r.Context(), ipCtxKey{}, ip)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}

	return func(next http.Handler) http.Handler { return f(next) }
}
