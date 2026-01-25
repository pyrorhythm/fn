package fn

// To is similar to [Morph], but specializes on error propagation.
// If [Result.Exc] == true, propagates error without calling (and needing) the function.
// Else, returns either a tf-call [Result], or a [ErrNoTFInTo] error.
func To[T, U any](r Result[T], tf ...func(T) U) Result[U] {
	if r.Exc() {
		return err[U](r.e)
	}

	if len(tf) == 0 {
		return notf[U]()
	}

	return okv(tf[0](r.Val()))
}

// Morph casts [Result][T] to [Result][U] using func(T) [Result][U]. If [Result.Exc] == true, propagates error without calling the function.
func Morph[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.OK() {
		return f(r.Val())
	}

	return err[U](r.Err())
}
