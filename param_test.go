package webkit

import (
	"testing"
)

func TestParamNil(t *testing.T) {
	var ps Params
	_, err := ps.Int64("key")
	if err != ErrNilParam {
		t.Fatalf("want error %s, but got %s", ErrNilParam, err)
	}
}
