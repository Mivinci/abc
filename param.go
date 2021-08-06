package webkit

import (
	"errors"
	"strconv"
)

var (
	ErrNilParams = errors.New("nil params")
	ErrNoKey     = errors.New("no such key in params")
)

type paramCtxKey struct{}

type Params map[string]string

func (ps Params) Int(key string) (int, error) {
	if ps == nil {
		return 0, ErrNilParams
	}
	i, ok := ps[key]
	if !ok {
		return 0, ErrNoKey
	}
	return strconv.Atoi(i)
}

func (ps Params) Int64(key string) (int64, error) {
	if ps == nil {
		return 0, ErrNilParams
	}
	i, ok := ps[key]
	if !ok {
		return 0, ErrNoKey
	}
	return strconv.ParseInt(i, 10, 64)
}

func (ps Params) Float64(key string) (float64, error) {
	return ps.float(key, 64)
}

func (ps Params) Float32(key string) (float32, error) {
	f, err := ps.float(key, 32)
	return float32(f), err
}

func (ps Params) float(key string, bs int) (float64, error) {
	if ps == nil {
		return 0, ErrNilParams
	}
	f, ok := ps[key]
	if !ok {
		return 0, ErrNoKey
	}
	return strconv.ParseFloat(f, bs)
}

func (ps Params) String(key string) (string, error) {
	if ps == nil {
		return "", ErrNilParams
	}
	s, ok := ps[key]
	if !ok {
		return "", ErrNoKey
	}
	return s, nil
}
