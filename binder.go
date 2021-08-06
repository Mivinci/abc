package webkit

type BinderCtx interface {
	Bind(Ctx) error
}

type Binder interface {
	Bind(interface{}) error
}

func DefaultBinder() Binder {
	return &binder{}
}

type binder struct{}

func (b binder) Bind(o interface{}) error {
	return nil
}
