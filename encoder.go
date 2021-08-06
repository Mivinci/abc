package webkit

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type Encoder interface {
	JSON(io.Writer, interface{}) error
	XML(io.Writer, interface{}) error
}

type stdEncoder struct{}

// TODO: use github.com/goccy/go-json instead
func (stdEncoder) JSON(w io.Writer, o interface{}) error {
	return json.NewEncoder(w).Encode(o)
}

func (stdEncoder) XML(w io.Writer, o interface{}) error {
	return xml.NewEncoder(w).Encode(o)
}

func StdEncoder() Encoder {
	return &stdEncoder{}
}
