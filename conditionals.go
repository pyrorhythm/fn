package fn

// If is a generic ternary operator
func If[T any](cond bool, then T, or T) T {
	if cond {
		return then
	}
	return or
}

// FlatIf is [If] but returns [Option][T], always valid
func FlatIf[T any](cond bool, then T, or T) Option[T] {
	if cond {
		return some(then)
	}
	return some(or)
}

// ErrIf returns [Result][T] with value if cond is true, else wraps error
func ErrIf[T any](cond bool, then T, or error) Result[T] {
	if cond {
		return okv(then)
	}
	return err[T](or)
}

// IfPtr returns [Option] with value if cond is true, else dereferences pointer
func IfPtr[T any](cond bool, then T, or *T) Option[T] {
	if cond {
		return some(then)
	}
	return SomePtr(or)
}

// PtrIf returns [Option] by dereferencing pointer if cond is true, else uses value
func PtrIf[T any](cond bool, then *T, or T) Option[T] {
	if cond {
		return SomePtr(then)
	}
	return some(or)
}
