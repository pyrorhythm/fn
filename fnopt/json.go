package fnopt

import (
	"fmt"

	"github.com/bytedance/sonic"
)

var jsonNull = []byte("null")

// UnmarshalJSON implements [json.Unmarshaler].
func (o *Of[T]) UnmarshalJSON(ba []byte) error {
	var dest *T
	if err := sonic.Unmarshal(ba, &dest); err != nil {
		*o = Nil[T]()
		return fmt.Errorf("failed to unmarshal Option[%T]: %w", o.t, err)
	}
	if dest == nil {
		*o = Nil[T]()
		return nil
	}
	*o = some(*dest)
	return nil
}

// MarshalJSON implements [json.Marshaler].
// Marshals the value if valid, otherwise marshals null.
func (o Of[T]) MarshalJSON() ([]byte, error) {
	if o.v {
		return sonic.Marshal(o.t)
	}
	return jsonNull, nil
}
