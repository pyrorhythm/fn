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

func Map[T any, U any](s []T, f func(T) U) []U {
	res := make([]U, len(s))
	for i, item := range s {
		res[i] = f(item)
	}
	return res
}

func Reduce[T any, U any](s []T, f func(U, T) U, init ...U) U {
	var res U
	if len(init) > 0 {
		res = init[0]
	}
	for _, item := range s {
		res = f(res, item)
	}
	return res
}

func Filter[T any](s []T, f func(T) bool) []T {
	res := make([]T, 0, len(s))
	for _, item := range s {
		if f(item) {
			res = append(res, item)
		}
	}
	return res
}

func Find[T any](s []T, f func(T) bool) (T, bool) {
	for _, item := range s {
		if f(item) {
			return item, true
		}
	}
	var z T
	return z, false
}

func MapErr[T any, U any](s []T, f func(T) (U, error)) ([]U, error) {
	res := make([]U, len(s))
	for i, item := range s {
		var err error
		res[i], err = f(item)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func MapOk[T any, U any](s []T, f func(T) (U, bool)) []U {
	res := make([]U, 0, len(s))
	for _, item := range s {
		u, ok := f(item)
		if ok {
			res = append(res, u)
		}
	}
	return res
}
