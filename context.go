package webkit

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

const (
	headerContentType = "Content-Type"
	charsetUTF8       = "charset=utf-8"
)

type Ctx interface {
	Request() *http.Request
	Writer() http.ResponseWriter
	Reset(http.ResponseWriter, *http.Request)
	Params() Params
	Query() url.Values
	Form() url.Values
	Bind(interface{}) error
	BindCtx(BinderCtx) error
	RealIP() string
	Text(int, string) error
	Blob(int, string, []byte) error
	JSON(int, interface{}) error
	JSONBlob(int, []byte) error
	XML(int, interface{}) error
	XMLBlob(int, []byte) error
	NoContent(int) error
}

var _ Ctx = (*ctx)(nil)

type ctx struct {
	r *http.Request
	w http.ResponseWriter
	p *Router
}

func (c *ctx) Request() *http.Request {
	return c.r
}

func (c *ctx) Writer() http.ResponseWriter {
	return c.w
}

func (c *ctx) Reset(w http.ResponseWriter, r *http.Request) {
	c.r = r
	c.w = w
}

// Params represents parameters in the request path
func (c *ctx) Params() Params {
	return c.r.Context().Value(paramCtxKey{}).(Params)
}

// Query is an alias of http.Request.URL.Query
func (c *ctx) Query() url.Values {
	return c.r.URL.Query()
}

// Form is an alias of http.Request.PostForm
func (c *ctx) Form() url.Values {
	if c.r.PostForm == nil {
		c.r.ParseForm()
	}
	return c.r.PostForm
}

func (c *ctx) Bind(o interface{}) error {
	return c.p.opts.binder.Bind(o)
}

func (c *ctx) BindCtx(b BinderCtx) error {
	return b.Bind(c)
}

func (c *ctx) RealIP() string {
	r := c.Request()
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

func (c *ctx) Blob(code int, contentType string, b []byte) error {
	c.contentType(contentType)
	c.status(code)
	_, err := c.w.Write(b)
	return err
}

func (c *ctx) Text(code int, text string) error {
	return c.Blob(code, "text/plain", []byte(text))
}

func (c *ctx) HTML(code int, s string) error {
	return c.Blob(code, "text/html", []byte(s))
}

func (c *ctx) JSON(code int, o interface{}) error {
	c.contentType("application/json")
	c.status(code)
	return c.p.opts.encoder.JSON(c.w, o)
}

func (c *ctx) JSONBlob(code int, b []byte) error {
	return c.Blob(code, "application/json", b)
}

func (c *ctx) XML(code int, o interface{}) error {
	c.contentType("application/xml")
	c.status(code)
	return c.p.opts.encoder.XML(c.w, o)
}

func (c *ctx) XMLBlob(code int, b []byte) error {
	return c.Blob(code, "application/xml", b)
}

func (c *ctx) contentType(value string) {
	h := c.w.Header()
	if h.Get(headerContentType) == "" {
		h.Set(headerContentType, value)
	}
}

func (c *ctx) status(code int) {
	c.w.WriteHeader(code)
}

func (c *ctx) NoContent(code int) error {
	c.w.WriteHeader(code)
	return nil
}
