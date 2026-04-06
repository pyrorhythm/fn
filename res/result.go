// Package res implements the Of monad for Go.
package res

import (
	"errors"
	"fmt"

	"github.com/pyrorhythm/fn"
	"github.com/pyrorhythm/fn/opt"
)

// Of[T] is either an OK value or an error.
type Of[T any] struct {
	v opt.Of[T]
	e error
}

func ok[T any](r opt.Of[T]) Of[T] { return Of[T]{v: r} }
func okv[T any](r T) Of[T]        { return ok(opt.SomeAny(r)) }
func okp[T any](r *T) Of[T]       { return ok(opt.SomeAny(fn.Deref(r))) }
func rErr[T any](e error) Of[T]   { return Of[T]{e: e} }

var errNilPtr = errors.New("nil pointer")

// From creates Of[T] from a standard (T, error) return where T is comparable.
func From[T comparable](val T, err error) Of[T] { return Of[T]{v: opt.Some(val), e: err} }

// FromPtr creates Of[T] from a standard (*T, error) return.
func FromPtr[T any](val *T, err error) Of[T] { return Of[T]{v: opt.SomePtr(val), e: err} }

// FromAny creates Of[T] from a standard (T, error) return, bypassing zero checks.
func FromAny[T any](val T, err error) Of[T] { return Of[T]{v: opt.SomeAny(val), e: err} }

// FromOpt creates Of[T] from an existing option and error.
func FromOpt[T any](o opt.Of[T], err error) Of[T] { return Of[T]{v: o, e: err} }

// FromAnyReflect creates Of[T] using reflection to determine validity.
// Prefer [From] or [FromAny] when possible; reflection is 25–50× slower.
func FromAnyReflect[T any](val T, err error) Of[T] {
	return Of[T]{v: opt.SomeAnyReflect(val), e: err}
}

// Err creates Of[T] from an error.
func Err[T any](e error) Of[T] { return rErr[T](e) }

// Errn creates Of[T] from an error string.
func Errn[T any](s string) Of[T] { return rErr[T](errors.New(s)) }

// Errw creates Of[T] from a wrapped error.
func Errw[T any](e error, s string) Of[T] { return rErr[T](fmt.Errorf("%s: %w", s, e)) }

// OK creates Of[T] from a comparable value. Invalid if value is zero.
func OK[T comparable](v T) Of[T] { return ok(opt.Some(v)) }

// OKOpt creates Of[T] from an existing option.
func OKOpt[T any](o opt.Of[T]) Of[T] { return ok(o) }

// OKPtr creates Of[T] from a pointer. Returns an error if the pointer is nil.
func OKPtr[T any](v *T) Of[T] {
	if v == nil {
		return rErr[T](errNilPtr)
	}
	return okp(v)
}

// OKAny creates Of[T] from any value, bypassing zero checks.
func OKAny[T any](v T) Of[T] { return okv(v) }

// OKAnyReflect creates Of[T] using reflection to determine validity.
// Prefer [OK] or [OKAny] when possible; reflection is 25–50× slower.
func OKAnyReflect[T any](v T) Of[T] { return ok(opt.SomeAnyReflect(v)) }

// Opt returns the underlying Option[T].
func (r Of[T]) Opt() opt.Of[T] { return r.v }

// Err returns the underlying error, or nil if OK.
func (r Of[T]) Err() error { return r.e }

// Val returns the underlying value. Zero if not OK.
func (r Of[T]) Val() T { return r.v.Val() }

// Ptr returns a pointer to the value, or nil if not OK.
func (r Of[T]) Ptr() *T {
	if !r.OK() {
		return nil
	}
	return r.v.Ptr()
}

// OK reports whether the result has no error and a valid option.
func (r Of[T]) OK() bool { return r.e == nil && r.v.Valid() }

// Into writes the value to *d and returns the error.
func (r Of[T]) Into(d *T) error { *d = r.v.Val(); return r.e }

// Unpack returns the value and error.
func (r Of[T]) Unpack() (T, error) { return r.v.Val(), r.e }

// Valid reports whether the result is OK. Implements [fn.Container].
func (r Of[T]) Valid() bool { return r.OK() }

// Fold returns onValue(t) if OK, otherwise onFail().
func (r Of[T]) Fold(onValue func(T) T, onFail func() T) T {
	if r.OK() {
		return onValue(r.v.Val())
	}
	return onFail()
}

// FlatMap applies f to the value if OK, propagating the error otherwise.
func (r Of[T]) FlatMap(f func(T) Of[T]) Of[T] {
	if r.OK() {
		return f(r.v.Val())
	}
	return rErr[T](r.e)
}
