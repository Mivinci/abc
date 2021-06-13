package webkit

import (
	"context"
	"net/http"
	"path"
)

type Router struct {
	options Options
	plugins []Plugin
	trees   map[string]*node
}

func New(opts ...Option) *Router {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return &Router{
		options: opt,
		plugins: make([]Plugin, 0),
		trees:   make(map[string]*node),
	}
}

func (r *Router) handle(method, pattern string, handler http.Handler) {
	if r.trees[method] == nil {
		r.trees[method] = &node{}
	}
	segments := parse(pattern)
	r.trees[method].insert(segments, handler, 0)
}

func (r *Router) lookup(path string, root *node, ps *Params) http.Handler {
	segments := parse(path)
	node := root.lookup(segments, 0, ps)
	if node != nil && node.handler != nil {
		return node.handler
	}
	return nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.RequestURI == "*" {
		if req.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ps Params
	path := req.URL.Path

	if root, ok := r.trees[req.Method]; ok {
		if h := r.lookup(path, root, &ps); h != nil {
			if ps != nil {
				ctx := context.WithValue(req.Context(), ParamsCtxKey{}, ps)
				req = req.WithContext(ctx)
			}
			h.ServeHTTP(w, req)
			return
		}
	}

	http.NotFound(w, req)
}

func (r *Router) compose(subs []Plugin) []Plugin {
	gLen, pLen := len(r.plugins), len(subs)
	plugins := make([]Plugin, gLen+pLen)
	copy(plugins, r.plugins)
	copy(plugins[gLen:], subs)
	return plugins
}

func (r *Router) concat(sub string) string {
	s := path.Join(r.options.group, sub)
	if s[0] != '/' {
		s = "/" + s
	}
	return s
}

// Group is a hackable method of Router that creates a new router that inherits all the plugins
// and handlers from itself. The new router will be released after calling this method and leaving
// newly registered group of handlers and plugins on the root router.
func (r *Router) Group(name string, f func(*Router)) {
	r.mount(name, r.plugins, f)
}

// Mount unlike Group, it only inherits handlers from the root router and has its own stack of plugins.
func (r *Router) Mount(name string, f func(*Router)) {
	r.mount(name, nil, f)
}

func (r *Router) mount(name string, plugins []Plugin, f func(*Router)) {
	dummy := Router{
		trees:   r.trees,
		plugins: plugins,
		options: Options{
			group: r.concat(name),
		},
	}
	f(&dummy)
}

// /a/b/c   -> [a,b,c]
// /a/:b/:c -> [a,:b,:c]
// /a/b/*   -> [a,b,*]
// /a/*/b/c -> [a,*]
// /*       -> [*]
// /        -> []
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

func (r *Router) Use(plugin Plugin) {
	r.plugins = append(r.plugins, plugin)
}

func (r *Router) Handle(method, pattern string, endpoint http.Handler, plugins ...Plugin) {
	r.handle(method, r.concat(pattern), chain(endpoint, r.compose(plugins)))
}

func (r *Router) HandleFunc(method, pattern string, endpoint http.HandlerFunc, plugins ...Plugin) {
	r.Handle(method, pattern, endpoint, plugins...)
}

func (r *Router) Get(pattern string, endpoint http.HandlerFunc, plugins ...Plugin) {
	r.Handle(http.MethodGet, pattern, endpoint, plugins...)
}

func (r *Router) Post(pattern string, endpoint http.HandlerFunc, plugins ...Plugin) {
	r.Handle(http.MethodPost, pattern, endpoint, plugins...)
}

func (r *Router) Put(pattern string, endpoint http.HandlerFunc, plugins ...Plugin) {
	r.Handle(http.MethodPut, pattern, endpoint, plugins...)
}

func (r *Router) Delete(pattern string, endpoint http.HandlerFunc, plugins ...Plugin) {
	r.Handle(http.MethodDelete, pattern, endpoint, plugins...)
}

type Plugin func(next http.Handler) http.Handler

// chain is the magic.
func chain(endpoint http.Handler, plugins []Plugin) http.Handler {
	for i := len(plugins) - 1; i >= 0; i-- {
		endpoint = plugins[i](endpoint)
	}
	return endpoint
}
