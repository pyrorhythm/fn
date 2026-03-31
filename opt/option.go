// Package opt implements the Of monad for Go.
package opt

import "github.com/pyrorhythm/fn"

// Of[T] is either a present value (Some) or absent (Nil).
type Of[T any] struct {
	t T
	v bool
}

func opt[T any](t T, b bool) Of[T] { return Of[T]{t, b} }
func some[T any](t T) Of[T]        { return opt(t, true) }

// Some returns Of[T], valid only if t is non-zero.
func Some[T comparable](t T) Of[T] { return opt(t, fn.Valid(t)) }

// SomePtr returns Of[T] from a pointer. Valid if and only if the pointer is non-nil.
func SomePtr[T any](t *T) Of[T] { return opt(fn.OrZero(t), t != nil) }

// SomeAny returns an always-valid Of[T], bypassing the zero check.
func SomeAny[T any](t T) Of[T] { return some(t) }

// SomeAnyReflect returns Of[T] using reflection or known interfaces to determine validity.
// Prefer [Some] or [SomeAny] when possible; reflection is 25–50× slower.
func SomeAnyReflect[T any](t T) Of[T] { return opt(t, fn.ValidReflect(t)) }

// Nil returns the zero Of[T], representing absence.
func Nil[T any]() (_ Of[T]) { return }

// Valid reports whether the option holds a value.
func (o Of[T]) Valid() bool { return o.v }

// Val returns the held value. Zero if invalid.
func (o Of[T]) Val() T { return o.t }

// Ptr returns a pointer to the value, or nil if invalid.
func (o Of[T]) Ptr() *T {
	if !o.v {
		return nil
	}
	return &o.t
}

// Unpack returns the value and validity flag.
func (o Of[T]) Unpack() (T, bool) { return o.t, o.v }

// Fold returns onVal(t) if valid, otherwise onNil().
func (o Of[T]) Fold(onNil func() T, onVal func(T) T) T {
	if o.v {
		return onVal(o.t)
	}
	return onNil()
}

// FlatMap applies f to the value if valid, propagating Nil otherwise.
func (o Of[T]) FlatMap(f func(T) Of[T]) Of[T] {
	if o.v {
		return f(o.t)
	}
	return Nil[T]()
}
