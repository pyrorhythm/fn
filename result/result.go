// Package result is an alias for [fn.Result] for the sake of lesser verbosity.
package result

import "github.com/pyrorhythm/fn"

type Of[T any] = fn.Result[T]

func OK[T comparable](v T) Of[T]                   { return fn.OK(v) }
func OKPtr[T any](v *T) Of[T]                      { return fn.OKPtr(v) }
func OKAny[T any](v T) Of[T]                       { return fn.OKAny(v) }
func OKOpt[T any](o fn.Option[T]) Of[T]            { return fn.OKOpt(o) }
func Err[T any](e error) Of[T]                     { return fn.Err[T](e) }
func Errn[T any](s string) Of[T]                   { return fn.Errn[T](s) }
func Errw[T any](e error, s string) Of[T]          { return fn.Errw[T](e, s) }
func From[T comparable](v T, e error) Of[T]        { return fn.From(v, e) }
func FromPtr[T any](v *T, e error) Of[T]           { return fn.FromPtr(v, e) }
func FromAny[T any](v T, e error) Of[T]            { return fn.FromAny(v, e) }
func FromOpt[T any](o fn.Option[T], e error) Of[T] { return fn.FromOpt(o, e) }
