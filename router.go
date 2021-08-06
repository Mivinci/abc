package webkit

import (
	"context"
	"net/http"
	"path"
	"sync"
)

type Handler func(Ctx) error

type Plugin func(Handler) Handler

type Router struct {
	pool  sync.Pool
	opts  Options
	trees map[string]*node
}

func New(opts ...Option) *Router {
	var opt Options
	opt.encoder = StdEncoder()
	opt.handleError = defaultErrorHandler
	for _, o := range opts {
		o(&opt)
	}

	var r Router
	r.opts = opt
	r.trees = make(map[string]*node)
	r.pool.New = func() interface{} { return &ctx{p: &r} }
	return &r
}

func parse(s string) []string {
	i, j, m, k, n := 1, 0, 1, 4, len(s)

	if n < 2 {
		return nil
	}

	if n > k {
		k = n >> 1
	}

	vs := make([]string, k)

	for i < n {
		if s[i] == '/' {
			vs[j] = s[m:i]
			m = i + 1
			j++
		} else if s[i] == '*' {
			n = i + 1
			break
		}
		i++
	}

	vs[j] = s[m:n]
	return vs[:j+1]
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.RequestURI == "*" {
		if req.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	root, ok := r.trees[req.Method]

	if !ok {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var params Params
	var target node

	root.lookup(parse(req.URL.Path), 0, &params, &target)

	h := target.handler
	if h == nil {
		http.NotFound(w, req)
		return
	}

	if params != nil {
		ctx := context.WithValue(req.Context(), paramCtxKey{}, params)
		req = req.WithContext(ctx)
	}

	c := r.pool.Get().(Ctx)
	c.Reset(w, req)

	if err := h(c); err != nil {
		r.opts.handleError(c, err)
	}

	r.pool.Put(c)
}

func (r *Router) handle(method, pattern string, handler Handler) {
	if r.trees[method] == nil {
		r.trees[method] = &node{}
	}
	r.trees[method].insert(parse(pattern), handler, 0)
}

func (r *Router) Handle(method, pattern string, handler Handler, plugins ...Plugin) {
	r.handle(method, r.join(pattern), chain(handler, r.append(plugins)))
}

func (r *Router) Get(pattern string, handler Handler, plugins ...Plugin) {
	r.Handle(http.MethodGet, pattern, handler, plugins...)
}

func (r *Router) Post(pattern string, handler Handler, plugins ...Plugin) {
	r.Handle(http.MethodPost, pattern, handler, plugins...)
}

func (r *Router) Put(pattern string, handler Handler, plugins ...Plugin) {
	r.Handle(http.MethodPut, pattern, handler, plugins...)
}

func (r *Router) Delete(pattern string, handler Handler, plugins ...Plugin) {
	r.Handle(http.MethodDelete, pattern, handler, plugins...)
}

func (r *Router) Group(name string, f func(*Router)) {
	r.mount(name, r.opts.plugins, f)
}

func (r *Router) Mount(name string, f func(*Router)) {
	r.mount(name, nil, f)
}

func (r *Router) mount(group string, plugins []Plugin, f func(*Router)) {
	dummy := Router{
		trees: r.trees,
		opts: Options{
			group:   r.join(group),
			plugins: plugins,
		},
	}
	f(&dummy)
}

func (r *Router) append(plugins []Plugin) []Plugin {
	i, j := len(r.opts.plugins), len(plugins)
	ps := make([]Plugin, i+j)
	copy(ps, r.opts.plugins)
	copy(ps[i:], plugins)
	return ps
}

func (r *Router) join(sub string) string {
	s := path.Join(r.opts.group, sub)
	if s[0] != '/' {
		s = "/" + s
	}
	return path.Clean(s)
}

func chain(endpoint Handler, plugins []Plugin) Handler {
	for i := len(plugins) - 1; i >= 0; i-- {
		endpoint = plugins[i](endpoint)
	}
	return endpoint
}

// RFC7807
func defaultErrorHandler(c Ctx, err error) {
	code := http.StatusOK
	e, ok := err.(*Error)

	if !ok {
		e = &Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if e.Reserved() {
		code = e.Code
	}

	if c.Request().Method == http.MethodHead {
		c.NoContent(code)
	}

	if err = c.JSON(code, e); err != nil {
		http.Error(c.Writer(), err.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}

func (r *Router) RunTLS(addr, cert, key string) error {
	return http.ListenAndServeTLS(addr, cert, key, r)
}
