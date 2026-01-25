package fn

import (
	"fmt"

	"github.com/bytedance/sonic"
)

type resultJSON[T any] struct {
	Value      *T     `json:"result_value,omitempty"`
	ValueValid bool   `json:"result_value_valid"`
	Error      string `json:"result_error,omitempty"`
}

// UnmarshalJSON implements [json.Unmarshaler].
func (r *Result[T]) UnmarshalJSON(ba []byte) error {
	var rj resultJSON[T]

	if err := sonic.Unmarshal(ba, &rj); err != nil {
		*r = Err[T](fmt.Errorf("failed to unmarshal Result[%T]: %w", *new(T), err))

		return err
	}

	// restore error if present
	if rj.Error != "" {
		*r = Err[T](fmt.Errorf("%s", rj.Error))

		return nil
	}

	// restore option based on validity and value
	if rj.Value == nil {
		*r = ok(Nil[T]())
	} else if rj.ValueValid {
		*r = ok(some(*rj.Value))
	} else {
		*r = ok(Option[T]{t: *rj.Value, v: false})
	}

	return nil
}

// MarshalJSON implements [json.Marshaler].
func (r Result[T]) MarshalJSON() ([]byte, error) {
	rj := resultJSON[T]{
		ValueValid: r.v.Valid(),
	}

	if r.e != nil {
		rj.Error = r.e.Error()
	}

	if r.v.Valid() {
		v := r.v.Value()
		rj.Value = &v
	}

	return sonic.Marshal(rj)
}
