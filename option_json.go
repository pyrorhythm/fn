package fn

import (
	"fmt"

	"github.com/bytedance/sonic"
)

var jsonNull = []byte("null")

// UnmarshalJSON implements [json.Unmarshaler].
func (o *Option[T]) UnmarshalJSON(ba []byte) error {
	var dest *T

	// error while unmarshalling, set zero value
	if err := sonic.Unmarshal(ba, &dest); err != nil {
		*o = Nil[T]()

		return fmt.Errorf("failed to unmarshal Option[%T]: %w", o.t, err)
	}

	// option is invalid, set zero value
	if dest == nil {
		*o = Nil[T]()

		return nil
	}

	// pointer non-nil, can say that option is valid
	*o = some(*dest)

	return nil
}

// MarshalJSON implements [json.Marshaler].
//
// Marshal value if [Option.Valid]. Else - nil.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.Valid() {
		return sonic.Marshal(o.t)
	}

	return jsonNull, nil
}
