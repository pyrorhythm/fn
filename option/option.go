// Package option is an alias for [fn.Option] for the sake of lesser verbosity.
package option

import "github.com/pyrorhythm/fn"

type Of[T any] = fn.Option[T]

func Some[T comparable](t T) Of[T] { return fn.Some(t) }
func SomePtr[T any](t *T) Of[T]    { return fn.SomePtr(t) }
func SomeAny[T any](t T) Of[T]     { return fn.SomeAny(t) }
func Nil[T any]() Of[T]            { return fn.Nil[T]() }
