package webkit

import (
	"html/template"
	"io/fs"
)

type Options struct {
	group       string
	encoder     Encoder
	binder      Binder
	handleError ErrorHandler
	template    *template.Template
	plugins     []Plugin
}

type Option func(*Options)

func Encode(e Encoder) Option {
	return func(o *Options) {
		o.encoder = e
	}
}

func Bind(b Binder) Option {
	return func(o *Options) {
		o.binder = b
	}
}

func TemplateFS(fs fs.FS, fns template.FuncMap, patterns ...string) Option {
	return func(o *Options) {
		o.template = template.Must(
			template.ParseFS(fs, patterns...)).Funcs(fns)
	}
}

func HandleError(h ErrorHandler) Option {
	return func(o *Options) {
		o.handleError = h
	}
}

func Group(name string) Option {
	return func(o *Options) {
		o.group = name
	}
}

func Plugins(plugins ...Plugin) Option {
	return func(o *Options) {
		o.plugins = plugins
	}
}
