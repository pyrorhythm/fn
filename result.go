package fn

import (
	"errors"
	"fmt"
)

// Result [T] implements the result monad. It can be represented either by:
//
// - [OK]([Option][T]) (which still could be zero or nil, even if err == nil)
//
// or
//
// - [Err](error)
type Result[T any] struct {
	v Option[T]
	e error
}

type AnyResult = Result[any]

var (
	ErrNilInOKPtr = errors.New("nil passed to okptr")
	ErrNoTFInTo   = errors.New("no tfunc passd in to(); result isnt exc")
)

func res[T any](r Option[T], e error) Result[T] { return Result[T]{r, e} }
func ok[T any](r Option[T]) Result[T]           { return Result[T]{v: r} }
func okv[T any](r T) Result[T]                  { return ok(some(r)) }
func okp[T any](r *T) Result[T]                 { return ok(some(OrZero(r))) }
func err[T any](e error) Result[T]              { return Result[T]{e: e} }
func errn[T any](s string) Result[T]            { return err[T](errors.New(s)) }

func nilok[T any]() Result[T] { return err[T](ErrNilInOKPtr) }
func notf[T any]() Result[T]  { return err[T](ErrNoTFInTo) }

// From creates [Result][T] from standard function output,
// where function yields a value of T and T implements [comparable]
func From[T comparable](val T, err error) Result[T] { return res(Some(val), err) }

// FromPtr creates [Result][T] from standard function output, where function yields a pointer of T
func FromPtr[T any](val *T, err error) Result[T] { return res(SomePtr(val), err) }

// FromAny creates [Result][T] from standard function output,
// bypassing validity checks (option set as valid)
func FromAny[T any](val T, err error) Result[T] { return res(some(val), err) }

// FromOpt creates [Result][T] from given option.
func FromOpt[T any](o Option[T], err error) Result[T] { return res(o, err) }

// FromAnyReflect creates [Result][T] from standard function output,
// where function yields a value of T and uses reflection to validate it.
//
// Highly unrecommended for use, see [ValidReflect]
func FromAnyReflect[T any](val T, err error) Result[T] { return res(SomeAnyReflect(val), err) }

// Err creates [Result][T] from error.
//
// Implements Err(error) monad.
func Err[T any](e error) Result[T] { return err[T](e) }

// Errn creates [Result][T] from error string.
//
// Implements Err(error) monad.
func Errn[T any](errStr string) Result[T] { return err[T](errors.New(errStr)) }

// Errw creates [Result][T] from wrapped error.
//
// Implements Err(error) monad.
func Errw[T any](e error, s string) Result[T] { return err[T](fmt.Errorf("%s: %w", s, e)) }

// OK creates [Result][T] from a value of T, where T implements [comparable]
func OK[T comparable](v T) Result[T] { return ok(Some(v)) }

// OKOpt creates [Result][T] from given option.
func OKOpt[T any](o Option[T]) Result[T] { return ok(o) }

// OKPtr creates [Result][T] from T pointer.
//
// When passing pointer to this function, you need to be sure that pointer is valid.
// Else, returns wrapped [ErrNilInOKPtr] as [Result][T]
func OKPtr[T any](v *T) Result[T] { return If(v != nil, okp(v), nilok[T]()) }

// OKAny creates [Result][T] from a value of T, bypassing validity checks.
func OKAny[T any](v T) Result[T] { return okv(v) }

// OKAnyReflect creates [Result][T] from T value, using reflection.
//
// Highly unrecommended for use, see [ValidReflect]
func OKAnyReflect[T any](v T) Result[T] { return ok(SomeAnyReflect(v)) }

// Opt returns underlying [Option][T] of [Result][T]
func (r Result[T]) Opt() Option[T] { return r.v }

// Err returns underlying error of [Result][T]
func (r Result[T]) Err() error { return r.e }

// Val returns underlying T value from [Option][T]
func (r Result[T]) Val() T { return r.v.t }

// Ptr returns underlying T pointer from [Option][T]
func (r Result[T]) Ptr() *T { return &r.v.t }

// OK returns whether if [Result][T] is not containing error and underlying [Option][T] is valid.
func (r Result[T]) OK() bool { return r.e == nil && r.v.v }

// Exc is a negation of [Result.OK]
func (r Result[T]) Exc() bool { return !r.OK() }

// Into sets underlying [Option][T] value to giver pointer of T
func (r Result[T]) Into(d *T) error { r.v.into(d); return r.e }

// Unpack ...unpacks?
func (r Result[T]) Unpack() (T, error) { return r.v.t, r.e }

// Any returns type-unsafe [AnyResult]
func (r Result[T]) Any() AnyResult { return AnyResult{r.v.Any(), r.e} }

// Inspect calls provided function with result itself.
func (r Result[T]) Inspect(f func(r Result[T])) {
	f(r)
}

// MapResult maps result (LOL) to the result of  the same type by calling provided function.
func (r *Result[T]) MapResult(f func(r Result[T]) Result[T]) {
	*r = f(*r)
}

// Map maps underlying optional T to the option of the same type by calling provided function.
//
// If [Result][T] contains error, function is not being called.
func (r *Result[T]) Map(f func(T, bool) (T, bool)) {
	if r.e != nil {
		return
	}
	r.v = opt(f(r.v.Unpack()))
}

// MapTo maps underlying optional T to the [Result][T] by calling provided function.
//
// If [Result][T] contains error, function is not being called.
func (r *Result[T]) MapTo(f func(T, bool) Result[T]) {
	if r.e != nil {
		return
	}
	*r = f(r.v.Unpack())
}

// Valid returns true if Result is OK (no error and valid option).
// Implements [Container] interface.
func (r Result[T]) Valid() bool { return r.OK() }

// Fold pattern matches over Result, calling onFail if Exc, onVal if OK.
func (r Result[T]) Fold(onValue func(T) T, onFail func() T) T {
	if r.OK() {
		return onValue(r.v.t)
	}
	return onFail()
}

// FlatMap chains Result operations, propagating error if Exc.
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T] {
	if r.OK() {
		return f(r.v.t)
	}
	return err[T](r.e)
}
