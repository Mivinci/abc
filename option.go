package webkit

type Options struct {
	group string
}

type Option func(*Options)

func Group(group string) Option {
	return func(o *Options) {
		o.group = group
	}
}
