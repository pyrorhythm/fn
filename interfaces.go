package fn

// Container provides value extraction from monadic types.
type Container[T any] interface {
	Val() T
	Valid() bool
}

// Monad represents a chainable container.
type Monad[T any, Self any] interface {
	Container[T]
	FlatMap(func(T) Self) Self
}
