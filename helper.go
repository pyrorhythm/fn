package fn

import "reflect"

// Valid returns whether T is non-zero, where T implements [comparable]
func Valid[T comparable](v T) bool {
	return v != *new(T)
}

// ValidReflect - valid through type-casting and [reflectValue]
//
// NOT recommended if not implementing interfaces below, as reflection in Go
// may be from 25 to 50 times slower than stdlib.
//
// If you really want to, you can check casted interfaces below,
// and implement one of them to avoid using reflection.
func ValidReflect[T any](v T) bool {
	switch m := any(v).(type) {
	case interface{ Bool() bool }:
		return m.Bool()
	case interface{ Ok() bool }:
		return m.Ok()
	case interface{ Validate() error }:
		return m.Validate() == nil
	case interface{ Validate() bool }:
		return m.Validate()
	case interface{ Valid() bool }:
		return m.Valid()
	case interface{ IsZero() bool }:
		return !m.IsZero()
	}

	return reflectValue(&v)
}

func reflectValue(vp any) bool {
	switch rv := reflect.ValueOf(vp).Elem(); rv.Kind() {
	case reflect.Map, reflect.Slice:
		return rv.Len() != 0
	default:
		return !rv.IsZero()
	}
}

// Or returns first pointer as a value if pointer is not nil.
//
// Else returns fallback T value, which was passed as a second argument.
func Or[T any](p *T, v T) T {
	if p != nil {
		return *p
	}

	return v
}

// OrZero validates pointer and returns pointer as a value if pointer is not nil.
//
// Else returns zero T value.
func OrZero[T any](p *T) T {
	return Or(p, *new(T))
}

func Is[T any](v any) bool {
	_, ok := v.(T)
	return ok
}

func Cast[T any](v any) T {
	t, ok := v.(T)
	if !ok {
		return Z[T]()
	}
	return t
}

// Z ...wait... OH! It's Zero Value!
func Z[T any]() T { return *new(T) }
