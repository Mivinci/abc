package webkit

import "context"

type Params map[string]string

type ParamsCtxKey struct{}

func ParamsFromContext(ctx context.Context) Params {
	ps, _ := ctx.Value(ParamsCtxKey{}).(Params)
	return ps
}
