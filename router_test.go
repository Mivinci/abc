package webkit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

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

var testEncoder = stdEncoder{}

func TestGroupRouter(t *testing.T) {
	s := New(
		Group("a"),
		Encode(testEncoder))

	w := httptest.NewRecorder()
	h := func(Ctx) error { return nil }

	s.Get("/b", h)
	r, _ := http.NewRequest(http.MethodGet, "/b", nil)
	s.ServeHTTP(w, r)

	s.Group("c", func(r *Router) {
		r.Get("/d", h)
	})
	r, _ = http.NewRequest(http.MethodGet, "/c/d", nil)
	s.ServeHTTP(w, r)
}

func ExamplePlugin() {
	f := make([]Plugin, 4)
	for i := 0; i < 4; i++ {
		f[i] = func(i int) Plugin {
			return func(next Handler) Handler {
				return func(c Ctx) (err error) {
					fmt.Println(i)
					if err = next(c); err != nil {
						return
					}
					fmt.Println(i + 10)
					return
				}
			}
		}(i)
	}

	h := New(Plugins(f...), Encode(testEncoder))
	h.Get("/", func(c Ctx) (err error) {
		fmt.Println("endpoint")
		return
	})

	w := httptest.NewRecorder()
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

func ExampleRouter_Use() {
	h := New(Plugins(func(next Handler) Handler {
		return func(c Ctx) error {
			fmt.Println(1)
			return next(c)
		}
	}))

	p := func(next Handler) Handler {
		return func(c Ctx) error {
			fmt.Println(2)
			return next(c)
		}
	}

	f := func(c Ctx) error {
		fmt.Println(3)
		return nil
	}

	h.Group("/group", func(r *Router) {
		r.Use(p)
		r.Get("/", f)
	})

	h.Mount("/mount", func(r *Router) {
		r.Use(p)
		r.Get("/", f)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/group", nil)
	h.ServeHTTP(w, r)

	r, _ = http.NewRequest(http.MethodGet, "/mount", nil)
	h.ServeHTTP(w, r)

	// Output:
	// 1
	// 2
	// 3
	// 2
	// 3
}
