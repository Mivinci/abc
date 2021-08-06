package webkit

type Options struct {
	group       string
	encoder     Encoder
	binder      Binder
	handleError ErrorHandler
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
