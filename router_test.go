package webkit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockResponseWriter struct{}

func (mockResponseWriter) Header() (h http.Header) { return }

func (mockResponseWriter) Write([]byte) (n int, err error) { return }

func (mockResponseWriter) WriteHeader(int) {}

func BenchmarkParse(b *testing.B) {
	path := "/some/page/path"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parse(path)
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		path string
		want []string
	}{
		{"/a", []string{"a"}},
		{"/a/b/c", []string{"a", "b", "c"}},
		{"/a/:b/:c", []string{"a", ":b", ":c"}},
		{"/a/b/*", []string{"a", "b", "*"}},
		{"/a/*/b/c", []string{"a", "*"}},
		{"/*", []string{"*"}},
		{"/", nil},
		{"/a/b/c/d/e/f", []string{"a", "b", "c", "d", "e", "f"}},
		{"/apple/banana/orange", []string{"apple", "banana", "orange"}},
	}

	for _, test := range testCases {
		actual := parse(test.path)
		if !reflect.DeepEqual(actual, test.want) {
			t.Fatalf("path %s want %v, but got %v", test.path, test.want, actual)
		}
	}
}

func TestRootHandler(t *testing.T) {
	h := New()
	h.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	w := new(mockResponseWriter)
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
}

func TestParamsFromContext(t *testing.T) {
	want := Params{"name": "webkit"}
	h := New()
	h.Get("/user/:name", func(_ http.ResponseWriter, r *http.Request) {
		ps := ParamsFromContext(r.Context())
		if !reflect.DeepEqual(want, ps) {
			t.Fatalf("want: %v, got: %v", want, ps)
		}
	})
	w := new(mockResponseWriter)
	r, _ := http.NewRequest(http.MethodGet, "/user/webkit", nil)
	h.ServeHTTP(w, r)
}

func TestRouter404(t *testing.T) {
	h := New()
	h.Get("/a", nil)
	r, _ := http.NewRequest(http.MethodGet, "/b", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("404 expected")
	}
}

func TestRouterGroup(t *testing.T) {
	h := New(Group("a"))
	h.Get("/b", func(_ http.ResponseWriter, _ *http.Request) {})
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/a/b", nil)
	h.ServeHTTP(w, r)

	h.Group("c", func(r *Router) {
		r.Get("/d", func(_ http.ResponseWriter, _ *http.Request) {})
	})
	r, _ = http.NewRequest(http.MethodGet, "/a/c/d", nil)
	h.ServeHTTP(w, r)

	h.Group("d", func(r *Router) {
		r.Get("/", func(_ http.ResponseWriter, _ *http.Request) {})
	})

	r, _ = http.NewRequest(http.MethodGet, "/a/d", nil)
	h.ServeHTTP(w, r)
}

func ExamplePlugin() {
	h := New()
	f := make([]Plugin, 4)
	for i := 0; i < 4; i++ {
		f[i] = func(i int) Plugin {
			return func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Println(i)
					next.ServeHTTP(w, r)
					fmt.Println(i + 10)
				})
			}
		}(i)
	}
	for i := 0; i < 4; i++ {
		h.Use(f[i])
	}
	h.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("endpoint")
	})
	w := new(mockResponseWriter)
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	// Output:
	// 0
	// 1
	// 2
	// 3
	// endpoint
	// 13
	// 12
	// 11
	// 10
}
