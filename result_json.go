package fn

import (
	"bytes"
	"fmt"

	"github.com/bytedance/sonic"
)

type jerr struct {
	Error string `json:"_ERROR"`
}

var errPrefix = []byte(`{"_ERROR"`)

// UnmarshalJSON implements [json.Unmarshaler].
func (r *Result[T]) UnmarshalJSON(ba []byte) error {
	// Fast path: check prefix to detect error objects
	if bytes.HasPrefix(ba, errPrefix) {
		var je jerr
		if err := sonic.Unmarshal(ba, &je); err == nil && je.Error != "" {
			*r = errn[T](je.Error)
			return nil
		}
	}

	// Otherwise unmarshal as Option[T]
	var opt Option[T]
	if err := opt.UnmarshalJSON(ba); err != nil {
		return fmt.Errorf("failed to unmarshal Result[%T]: %w", opt.t, err)
	}

	*r = ok(opt)
	return nil
}

// MarshalJSON implements [json.Marshaler].
func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.e != nil {
		return sonic.Marshal(jerr{r.e.Error()})
	}

	return r.v.MarshalJSON()
}
