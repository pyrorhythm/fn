package fn

// If is a generic ternary operator.
func If[T any](cond bool, then T, or T) T {
	if cond {
		return then
	}
	return or
}
