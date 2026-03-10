package main

func ellipsis() {}

var ellipsisFunc = func() {}
var variadicFunc = func(int, ...any) {}

func B[T comparable](a, b T) bool {
	return a == b
}

type b[T comparable] struct {
	x T
}

func main() {
	x := 1
	var y = 4
	z := 2.3 * float64(y)

	zp := &z
	zpp := &zp

	v := x

	s := b[int]{v}

	c := make(chan int)
	rc := func(c chan int) <-chan int { return c }(c)
	bf := B
	bfInt := B[int]

	select {
	case val := <-rc:
		eq := bfInt(x, y)
		print("val: ", val, " eq: ", eq)
	default:
	}

	_ = s
	_ = zpp
}
