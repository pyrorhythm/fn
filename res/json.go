package res

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"
	"pyrorhythm.dev/fn/opt"
)

type jerr struct {
	Error string `json:"_ERROR"`
}

var errPrefix = []byte(`{"_ERROR"`)

// UnmarshalJSON implements [json.Unmarshaler].
func (r *Of[T]) UnmarshalJSON(ba []byte) error {
	if bytes.HasPrefix(ba, errPrefix) {
		var je jerr
		if err := json.Unmarshal(ba, &je); err == nil && je.Error != "" {
			*r = Errn[T](je.Error)
			return nil
		}
	}

	var o opt.Of[T]
	if err := o.UnmarshalJSON(ba); err != nil {
		return fmt.Errorf("failed to unmarshal Result[%T]: %w", r.Val(), err)
	}
	*r = ok(o)
	return nil
}

// MarshalJSON implements [json.Marshaler].
func (r Of[T]) MarshalJSON() ([]byte, error) {
	if r.e != nil {
		return json.Marshal(jerr{r.e.Error()})
	}
	return r.v.MarshalJSON()
}
