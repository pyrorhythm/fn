package fn

// Option implements option monad for type T: either [Nil]() or [Some]([SomePtr])(T)
type Option[T any] struct {
	t T
	v bool
}

type AnyOption = Option[any]

func opt[T any](t T, b bool) Option[T] { return Option[T]{t, b} }
func some[T any](t T) Option[T]        { return opt(t, true) }

// Some returns [Option][T] for a value of T, where T implements [comparable]
func Some[T comparable](t T) Option[T] { return opt(t, Valid(t)) }

// SomePtr returns [Option][T] for a pointer of T
func SomePtr[T any](t *T) Option[T] { return opt(OrZero(t), t != nil) }

// SomeAny returns always valid [Option][T] for a value of T, bypassing validity checks.
func SomeAny[T any](t T) Option[T] { return some(t) }

// SomeAnyReflect Sreturns [Option][T] for a value of T.
//
// Highly unrecommended for use, see [ValidReflect].
func SomeAnyReflect[T any](t T) Option[T] { return opt(t, ValidReflect(t)) }

// Nil is an alias for zero [Option][T]
func Nil[T any]() (_ Option[T]) { return }

func (o Option[T]) into(d *T) { *d = o.t }

// Valid tells if [Option] is valid
//
// P.S. (it may be zero inside but it was nil on creation)
func (o Option[T]) Valid() bool { return o.v }

// Value returns held value in [Option]
func (o Option[T]) Value() T { return o.t }

// Ptr returns pointer of held value in [Option].
//
// Always (i think (and pray)) non-nil.
func (o Option[T]) Ptr() *T { return &o.t }

// Any returns type-unsafe [Option]
func (o Option[T]) Any() AnyOption { return AnyOption{o.t, o.v} }
