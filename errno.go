package webkit

import (
	"fmt"
	"net/http"
	"strings"
)

type ErrorHandler func(Ctx, error)

type Error struct {
	Code    int
	Message string
	Values  interface{}
}

func HTTPError(code int, messages ...string) *Error {
	var e Error
	if messages != nil {
		e.Message = strings.Join(messages, "\n")
	}
	e.Code = code
	return &e
}

func (e *Error) Wrap(err error) *Error {
	if err != nil {
		e.Message = err.Error()
	}
	return e
}

func (e *Error) WithValues(o interface{}) *Error {
	e.Values = o
	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

func (e *Error) Reserved() bool {
	return e.Code >= http.StatusBadRequest && e.Code <= http.StatusNetworkAuthenticationRequired
}
