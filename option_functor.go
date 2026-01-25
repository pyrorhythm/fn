package fn

// OptTo is similar to [OptMorph], but specializes on [Nil] propagation.
// If ![Result.Valid], propagates [Nil] without calling (and needing) the function.
// Else, returns either a [Option][U] from call of TF, or a [Nil][U] if no TF is passed as an argument.
func OptTo[T, U any](o Option[T], tf ...func(T) U) Option[U] {
	if !o.Valid() {
		return Nil[U]()
	}
	if len(tf) == 0 {
		return Nil[U]()
	}
	return some(tf[0](o.Value()))
}

// OptMorph casts [Option] of type T to type U using func(T) [Option][U].
// If [Option][T] is invalid, func(T) [Option][U] would not be called.
func OptMorph[T, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.Valid() {
		return f(o.Value())
	}
	return Nil[U]()
}
