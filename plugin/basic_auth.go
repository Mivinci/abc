package plugin

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/mivinci/webkit/v2"
)

const (
	basic   = "Basic"
	keyUser = "user_id"
)

type BasicAuthValidator func(user, passwd string) error

func BasicAuth(realm string, va BasicAuthValidator) webkit.Plugin {
	return func(next webkit.Handler) webkit.Handler {
		return func(c webkit.Ctx) error {
			auth := c.Request().Header.Get("Authorization")
			n := len(basic)
			if len(auth) > n+1 && auth[:n] == basic {
				b, err := base64.StdEncoding.DecodeString(auth[n+1:])
				if err != nil {
					return err
				}

				cred := string(b)
				i := strings.IndexByte(cred, ':')
				if i != -1 {
					if err := va(cred[:i], cred[i+1:]); err == nil {
						c.SetParams(keyUser, cred[:i])
						return next(c)
					}
				}
			}

			c.Writer().Header().Set("WWW-Authenticate", basic+" realm="+realm)
			return webkit.HTTPError(http.StatusUnauthorized)
		}
	}
}
