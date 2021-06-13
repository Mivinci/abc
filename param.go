package webkit

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

var (
	ErrNilParam = errors.New("nil params")
	ErrNoKey    = errors.New("no key")
)

type Params map[string]string

type ParamsCtxKey struct{}

func ParamsFromContext(ctx context.Context) Params {
	ps, _ := ctx.Value(ParamsCtxKey{}).(Params)
	return ps
}

func (ps Params) Int64(key string) (int64, error) {
	if ps == nil {
		return 0, ErrNilParam
	}
	i64, ok := ps[key]
	if !ok {
		return 0, ErrNoKey
	}
	return strconv.ParseInt(i64, 10, 64)
}

func (ps Params) Uint64(key string) (uint64, error) {
	if ps == nil {
		return 0, ErrNilParam
	}
	u64, ok := ps[key]
	if !ok {
		return 0, ErrNoKey
	}
	return strconv.ParseUint(u64, 10, 64)
}

func (ps Params) String(key string) (string, error) {
	if ps == nil {
		return "", ErrNilParam
	}
	s, ok := ps[key]
	if !ok {
		return "", ErrNoKey
	}
	return s, nil
}

func (ps Params) StringSlice(key string) ([]string, error) {
	s, err := ps.String(key)
	if err != nil {
		return nil, err
	}
	return strings.Split(s, ","), nil
}
