package plugin

import (
	"net/http"

	"github.com/mivinci/webkit/v2"
	"github.com/mivinci/webkit/v2/store"
)

const (
	CSRFTokenKey       = "authenticity_token"
	CSRFTokenKeyHeader = "X-CSRF-Token"

	CSRFTokenHeader CSRFTokenType = iota
	CSRFTokenQuery
	CSRFTokenForm
)

type CSRFTokenType int8

func (t CSRFTokenType) Token(c webkit.Ctx) (token string) {
	switch t {
	case CSRFTokenForm:
		token = c.Form().Get(CSRFTokenKey)
	case CSRFTokenQuery:
		token = c.Query().Get(CSRFTokenKey)
	default:
		token = c.Request().Header.Get(CSRFTokenKeyHeader)
	}
	return
}

func (t CSRFTokenType) Key() string {
	if t == CSRFTokenHeader {
		return CSRFTokenKeyHeader
	}
	return CSRFTokenKey
}

func CSRF(s store.Store, t CSRFTokenType) webkit.Plugin {
	return func(next webkit.Handler) webkit.Handler {
		return func(c webkit.Ctx) error {
			var token string
			key := t.Key()
			// compare token with that in cookie or session.
			if v, err := s.Get(key); err == nil {
				token = v.(string)
			} else {
				token = randStr(32)
			}

			// TODO: non-zero token checking

			// validate csrf token only in unsafe requests (RFC7231).
			switch c.Request().Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			default:
				incoming := t.Token(c)
				if incoming == "" {
					return webkit.Error(http.StatusBadRequest)
				}
				if incoming != token {
					return webkit.Error(http.StatusForbidden)
				}
			}

			if err := s.Set(key, token, store.Forever); err != nil {
				return webkit.Error(http.StatusInternalServerError).Wrap(err)
			}

			c.SetParams(key, token)
			return next(c)
		}
	}
}
