package fn

import "errors"

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
	ErrNilPointerInOKPtr = errors.New("nilptr passed to okptr")
	ErrNoTFInTo          = errors.New("no tfunc passd in to(); result isnt exc")
)

func res[T any](r Option[T], e error) Result[T] { return Result[T]{r, e} }
func err[T any](e error) Result[T]              { return Result[T]{e: e} }
func ok[T any](r Option[T]) Result[T]           { return Result[T]{v: r} }
func okv[T any](r T) Result[T]                  { return Result[T]{v: some(r)} }
func okp[T any](r *T) Result[T]                 { return Result[T]{v: some(OrZero(r))} }

func nilptr[T any]() Result[T] { return err[T](ErrNilPointerInOKPtr) }
func notf[T any]() Result[T]   { return err[T](ErrNoTFInTo) }

// From creates [Result][T] from standard function output, where function yields a value of T and T implements [comparable]
func From[T comparable](val T, err error) Result[T] { return res(Some(val), err) }

// FromA creates [Result][T] from standard function output, bypassing value check (option is set as valid)
func FromA[T any](val T, err error) Result[T] { return res(some(val), err) }

// FromAReflect creates [Result][T] from standard function output, where function yields a value of T and uses reflection to validate it.
//
// Highly unrecommended for use, see [ValidReflect]
func FromAReflect[T any](val T, err error) Result[T] { return res(SomeAnyReflect(val), err) }

// FromP creates [Result][T] from standard function output, where function yields a pointer of T
func FromP[T any](val *T, err error) Result[T] { return res(SomeP(val), err) }

// Err creates [Result][T] from error.
//
// Implements Err(error) monad.
func Err[T any](e error) Result[T] { return err[T](e) }

// Errn creates [Result][T] from error string.
//
// Implements Err(error) monad.
func Errn[T any](errStr string) Result[T] { return err[T](errors.New(errStr)) }

// OK creates [Result][T] from a value of T, where T implements [comparable]
func OK[T comparable](v T) Result[T] { return ok(Some(v)) }

// OKPtr creates [Result][T] from T pointer.
//
// If pointer is nil, returns wrapped [ErrNilPointerInOKPtr] as [Result][T]
func OKPtr[T any](v *T) Result[T] { return If(v != nil, okp(v), nilptr[T]()) }

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

// OK returns whether if [Result][T] is not containing error nor underlying [Option][T] is valid.
func (r Result[T]) OK() bool { return r.e == nil && r.v.v }

// Exc is effectively a negation of [Result.OK]
func (r Result[T]) Exc() bool { return r.e != nil || !r.v.v }

// Into sets underlying [Option][T] value to giver pointer of T
func (r Result[T]) Into(d *T) error { r.v.into(d); return r.e }

// Unpack ...unpacks?
func (r Result[T]) Unpack() (T, error) { return r.v.t, r.e }

// Any returns type-unsafe [AnyResult]
func (r Result[T]) Any() AnyResult { return AnyResult{r.v.Any(), r.e} }
