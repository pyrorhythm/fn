package fnres

// To maps Of[T] to Of[U] by applying tf to the value.
// If the result has an error, tf is not called and the error propagates.
func To[T, U any](r Of[T], tf func(T) U) Of[U] {
	if !r.OK() {
		return rErr[U](r.e)
	}
	return okv(tf(r.Val()))
}

// ToPtr maps Of[T] to Of[U] by applying tf and unwrapping the pointer.
// If the result has an error, tf is not called and the error propagates.
func ToPtr[T, U any](r Of[T], tf func(T) *U) Of[U] {
	if !r.OK() {
		return rErr[U](r.e)
	}
	return okp(tf(r.Val()))
}

// Morph maps Of[T] to Of[U] using f, which returns a Of[U].
// If the result has an error, f is not called and the error propagates.
func Morph[T, U any](r Of[T], f func(T) Of[U]) Of[U] {
	if r.OK() {
		return f(r.Val())
	}
	return rErr[U](r.e)
}
