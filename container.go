package fn

type (
	Zeroer[T any]         func() T
	Transformer[T, U any] func(T) U
)

// Else returns the value if valid, otherwise returns fallback.
func Else[T any, M Container[T]](m M, fallback T) T {
	if m.Valid() {
		return m.Val()
	}
	return fallback
}

// Must returns the value if valid, otherwise panics.
func Must[T any, M Container[T]](m M) T {
	if !m.Valid() {
		panic("fn.Must: invalid container")
	}
	return m.Val()
}

// Fold pattern matches over a container, calling onNil if invalid, onVal if valid.
func Fold[T, U any, M Container[T]](m M, onNil Zeroer[U], onVal Transformer[T, U]) U {
	if m.Valid() {
		return onVal(m.Val())
	}
	return onNil()
}

// Chain applies f to the value if valid, otherwise returns the invalid container.
// This is FlatMap/Bind as a free function.
func Chain[T any, Self Monad[T, Self]](m Self, f Transformer[T, Self]) Self {
	return m.FlatMap(f)
}

// AndThen applies f if valid, ignoring the current value.
func AndThen[T any, M Monad[T, M]](m M, next M) M {
	return m.FlatMap(func(T) M { return next })
}
