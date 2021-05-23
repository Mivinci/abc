package webkit

import (
	"net/http"
	"reflect"
	"testing"
)

type mockResponseWriter struct{}

func (mockResponseWriter) Header() (h http.Header) { return }

func (mockResponseWriter) Write([]byte) (n int, err error) { return }

func (mockResponseWriter) WriteHeader(int) {}

func BenchmarkRouterLookup(b *testing.B) {
	var ps Params
	h := New()
	h.Get("/some/page/path", nil)
	h.Get("/some/page/:id", nil)
	root := h.trees[http.MethodGet]
	if root == nil {
		b.Fatal("failed to find root")
	}

	b.ReportAllocs()
	b.Run("static", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.lookup("/some/page/path", root, &ps)
		}
	})

	b.Run("dynamic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.lookup("/some/page/123", root, &ps)
		}
	})
}

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
	}

	for _, test := range testCases {
		actual := parse(test.path)
		if !reflect.DeepEqual(actual, test.want) {
			t.Fatalf("want %v, but got %v", test.want, actual)
		}
	}
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
