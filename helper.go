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

// Deref safely dereferences p and returns the value.
// If p == nil, returns **first** passed value in variadic chain.
// If no additional values passed, returns zero value.
// Else returns fallback T value, which was passed as a second argument.
func Deref[T any](p *T, v ...T) (z T) {
	if p != nil {
		return *p
	}
	if len(v) > 0 {
		return v[0]
	}
	return z
}

// Cast tries to cast v to type T, if failes - returns zero value.
func Cast[T any](v any) (z T) {
	t, ok := v.(T)
	if !ok {
		return z
	}
	return t
}

// Or returns first non-zero value, or zero if all passed values are zero.
func Or[T comparable](values ...T) (z T) {
	for _, v := range values {
		if Valid(v) {
			return v
		}
	}

	return z
}

// OrDef returns first non-zero value, or default if all passed values are zero.
func OrDef[T comparable](def T, values ...T) T {
	for _, v := range values {
		if Valid(v) {
			return v
		}
	}

	return def
}
