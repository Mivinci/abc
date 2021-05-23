package webkit

import (
	"context"
	"net/http"
)

type Router struct {
	trees map[string]*node
}

func New() *Router {
	return &Router{trees: make(map[string]*node)}
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

// /a/b/c   -> [a,b,c]
// /a/:b/:c -> [a,:b,:c]
// /a/b/*   -> [a,b,*]
// /a/*/b/c -> [a,*]
// /*       -> [*]
func parse(s string) []string {
	i, j, m := 1, 0, 1
	n := len(s)
	vs := make([]string, 4) // 4 pre-allocs

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

func (r *Router) Handle(method, pattern string, endpoint http.Handler, plugins ...Plugin) {
	r.handle(method, pattern, chain(endpoint, plugins...))
}

func (r *Router) HandleFunc(method, pattern string, endpoint http.HandlerFunc, plugins ...Plugin) {
	r.Handle(method, pattern, chain(endpoint, plugins...))
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
