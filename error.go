package webkit

import (
	"fmt"
	"net/http"
	"strings"
)

type ErrorHandler func(Ctx, error)

type HTTPError struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Values  interface{} `json:"values"`
}

func Error(code int, messages ...string) *HTTPError {
	var e HTTPError
	if messages != nil {
		e.Message = strings.Join(messages, "\n")
	}
	e.Code = code
	return &e
}

func (e *HTTPError) Wrap(err error) *HTTPError {
	if err != nil {
		e.Message = err.Error()
	}
	return e
}

func (e *HTTPError) WithValues(o interface{}) *HTTPError {
	e.Values = o
	return e
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

func (e *HTTPError) Reserved() bool {
	return e.Code >= http.StatusBadRequest && e.Code <= http.StatusNetworkAuthenticationRequired
}
