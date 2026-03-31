package fn

// Container provides value extraction from monadic types.
type Container[T any] interface {
	Val() T
	Valid() bool
}
