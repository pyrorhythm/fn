// Package chain provides helpers for chaining fallible functions in a [res.Of] pipeline.
package chain

import "pyrorhythm.dev/fn/res"

// Wrap chains a (T, error) function in a Result pipeline.
// If r has an error, f is not called and the error propagates.
func Wrap[T, U any](r res.Of[T], f func(T) (U, error)) res.Of[U] {
	return res.Morph(r, func(v T) res.Of[U] {
		return res.FromAny(f(v))
	})
}

// WrapPtr chains a (*T, error) function in a Result pipeline.
// If r has an error, f is not called and the error propagates.
func WrapPtr[T, U any](r res.Of[T], f func(T) (*U, error)) res.Of[U] {
	return res.Morph(r, func(v T) res.Of[U] {
		return res.FromPtr(f(v))
	})
}

// Wrap3 composes two (T, error) steps in one call.
func Wrap3[In, M, Out any](r res.Of[In], f1 func(In) (M, error), f2 func(M) (Out, error)) res.Of[Out] {
	return Wrap(Wrap(r, f1), f2)
}

// Wrap4 composes three (T, error) steps in one call.
func Wrap4[In, M1, M2, Out any](r res.Of[In], f1 func(In) (M1, error), f2 func(M1) (M2, error), f3 func(M2) (Out, error)) res.Of[Out] {
	return Wrap(Wrap3(r, f1, f2), f3)
}

// Wrap5 composes four (T, error) steps in one call.
func Wrap5[In, M1, M2, M3, Out any](r res.Of[In], f1 func(In) (M1, error), f2 func(M1) (M2, error), f3 func(M2) (M3, error), f4 func(M3) (Out, error)) res.Of[Out] {
	return Wrap(Wrap4(r, f1, f2, f3), f4)
}
