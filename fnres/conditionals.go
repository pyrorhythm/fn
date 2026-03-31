package fnres

// ErrIf returns Of[T] with the value if cond is true, otherwise wraps the error.
func ErrIf[T any](cond bool, then T, or error) Of[T] {
	if cond {
		return okv(then)
	}
	return rErr[T](or)
}

// IfErr dispatches between two handlers based on whether err is nil.
// Calls f_ok(v) if err == nil, f_err(err) otherwise.
func IfErr[T any](v T, err error) func(func(T), func(error)) {
	return func(f_ok func(T), f_err func(error)) {
		if err != nil {
			f_err(err)
		} else {
			f_ok(v)
		}
	}
}
