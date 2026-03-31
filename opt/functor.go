package opt

// To maps Of[T] to Of[U] by applying tf to the value.
// Returns Nil[U] if the option is invalid.
func To[T, U any](o Of[T], tf func(T) U) Of[U] {
	if !o.v {
		return Nil[U]()
	}
	return some(tf(o.t))
}

// Morph maps Of[T] to Of[U] using a function that returns an option.
// Returns Nil[U] if the option is invalid or if f returns Nil.
func Morph[T, U any](o Of[T], f func(T) Of[U]) Of[U] {
	if o.v {
		return f(o.t)
	}
	return Nil[U]()
}
