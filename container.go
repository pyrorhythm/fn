package fn

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

// Match pattern matches over a container, calling onNil if invalid, onVal if valid.
func Match[T, U any, M Container[T]](m M, onNil func() U, onVal func(T) U) U {
	if m.Valid() {
		return onVal(m.Val())
	}
	return onNil()
}
