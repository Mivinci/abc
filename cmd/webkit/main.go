package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"text/template"
)

const Repo = `package repo

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "{{ .Dsn }}"

var DB *sql.DB

func Init() {
	repo, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	DB = repo
}

func Close() error {
	return DB.Close()
}
`

const Main = `package main

import (
	"net/http"
	"os"

	"{{ .Repo }}/handle"
	"{{ .Repo }}/repo"
	"github.com/mivinci/webkit"
	"github.com/mivinci/webkit/plugins"
)

func common(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	repo.Init()

	h := webkit.New()
	h.Use(plugins.Log(os.Stdout))
	h.Use(common)

	h.Get("/", handle.Hello)

	http.ListenAndServe(":8080", h)
}
`

const Makefile = `TARGET={{ .Name }}
IMAGE=mivinci/$(TARGET)
BIN=bin/$(TARGET)

.PHONY: build
build:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o $(BIN) main.go 


.PHONY: image
image:
	docker image build -t $(IMAGE) .


.PHONY: run
run:
	go run main.go


.PHONY: run-image
run-image:
	docker run -d --rm -p 8080:8080 $(IMAGE)


.PHONY: publish
publish:
	docker push $(IMAGE)


.PHONY: tidy
tidy:
	go mod tidy


.PHONY: clean
clean:
	rm $(BIN)
`

const Handle = `package handle

import "net/http"

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Webkit!"))
}
`

const Dockerfile = `FROM alpine:3

COPY bin/{{ .Name }} {{ .Name }}

ENTRYPOINT [ "./{{ .Name }}" ]
`

const GoMod = `module {{ .Repo }}

go 1.16

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/mivinci/webkit v0.0.3
)
`

type Config struct {
	Name   string
	Author string
	Repo   string
	Dsn    string
	Port   int
}

type Source struct {
	filename string
	template string
}

var (
	root      string
	repoDir   string
	handleDir string
)

const mode = 0755

func main() {

	if len(os.Args) < 2 {
		dieError(fmt.Errorf("you need to provide a directory name as follows:\n\n    %s hello", os.Args[0]))
	}

	root = os.Args[1]

	c := &Config{
		Name:   "hello",
		Author: "mivinci",
		Port:   8080,
	}

	fmt.Println("Initializing Webkit project")
	fmt.Printf("name (%s): ", c.Name)
	fmt.Scanln(&c.Name)
	fmt.Printf("author (%s): ", c.Author)
	fmt.Scanln(&c.Author)

	for c.Repo == "" {
		fmt.Print("repo: ")
		fmt.Scanln(&c.Repo)
	}

	fmt.Print("dsn (empty string): ")
	fmt.Scanln(&c.Dsn)
	fmt.Print("port (8080): ")
	fmt.Scanln(&c.Port)

	confirm := 'Y'
	b, _ := json.MarshalIndent(c, "", "    ")
	fmt.Printf("Your config is:\n%s\n", string(b))
	fmt.Print("Is this OK? (Y/n): ")
	fmt.Scanln(&confirm)
	if confirm != 'Y' {
		return
	}

	repoDir = path.Join(root, "repo")
	handleDir = path.Join(root, "handle")

	dieError(os.MkdirAll(repoDir, mode))
	dieError(os.MkdirAll(handleDir, mode))

	t := template.New("")

	sources := [...]Source{
		{path.Join(repoDir, "repo.go"), Repo},
		{path.Join(handleDir, "handle.go"), Handle},
		{path.Join(root, "Makefile"), Makefile},
		{path.Join(root, "Dockerfile"), Dockerfile},
		{path.Join(root, "go.mod"), GoMod},
		{path.Join(root, "main.go"), Main},
	}

	for _, s := range sources {
		f, err := os.OpenFile(s.filename, os.O_CREATE|os.O_WRONLY, mode)
		dieError(err)

		t, _ := t.Parse(s.template)
		dieError(t.Execute(f, c))
	}

	fmt.Print("Completed. You can start by running:\n\n")
	if root != "" {
		fmt.Printf("    cd %s\n", root)
	}

	fmt.Print("    make tidy\n")
	fmt.Print("    make run\n\n")
	fmt.Println("Make sure you have `make` installed.")

}

func dieError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.RemoveAll(root)
		os.Exit(1)
	}
}
