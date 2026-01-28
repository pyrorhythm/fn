// Package pipe implements Raw, Err and [fn.Result] pipelines
package pipe

import "github.com/pyrorhythm/fn"

func Two[In, Out any](
	in In,
	fin fn.Transformer[In, Out],
) Out {
	return fin(in)
}

func Three[In, inter, Out any](
	in In,
	fin fn.Transformer[In, inter],
	finter fn.Transformer[inter, Out],
) Out {
	return finter(Two(in, fin))
}

func Four[In, inter1, inter2, Out any](
	in In,
	fin fn.Transformer[In, inter1],
	finter1 fn.Transformer[inter1, inter2],
	finter2 fn.Transformer[inter2, Out],
) Out {
	return finter2(Three(in, fin, finter1))
}

func Five[In, inter1, inter2, inter3, Out any](
	in In,
	fin fn.Transformer[In, inter1],
	finter1 fn.Transformer[inter1, inter2],
	finter2 fn.Transformer[inter2, inter3],
	finter3 fn.Transformer[inter3, Out],
) Out {
	return finter3(Four(in, fin, finter1, finter2))
}

func TwoErr[In, Out any](a fn.Result[In], fin func(In) (Out, error)) (z Out, _ error) {
	if a.Exc() {
		return z, a.Err()
	}
	return fin(a.Val())
}

func ThreeErr[In, inter, Out any](
	a fn.Result[In],
	fin func(In) (inter, error),
	finter1 func(inter) (Out, error),
) (z Out, _ error) {
	b, eb := TwoErr(a, fin)
	if eb != nil {
		return z, eb
	}
	return finter1(b)
}

func FourErr[In, inter1, inter2, Out any](
	a fn.Result[In],
	fin func(In) (inter1, error),
	finter1 func(inter1) (inter2, error),
	finter2 func(inter2) (Out, error),
) (z Out, _ error) {
	c, ec := ThreeErr(a, fin, finter1)
	if ec != nil {
		return z, ec
	}
	return finter2(c)
}

func FiveErr[In, inter1, inter2, inter3, Out any](
	a fn.Result[In],
	fin func(In) (inter1, error),
	finter1 func(inter1) (inter2, error),
	finter2 func(inter2) (inter3, error),
	finter3 func(inter3) (Out, error),
) (z Out, _ error) {
	d, ec := FourErr(a, fin, finter1, finter2)
	if ec != nil {
		return z, ec
	}
	return finter3(d)
}

func TwoRes[In, Out any](
	rin fn.Result[In],
	fin fn.Transformer[In, fn.Result[Out]],
) fn.Result[Out] {
	return fn.Morph(rin, func(i In) fn.Result[Out] {
		return fin(i)
	})
}

func ThreeRes[In, inter, Out any](
	rin fn.Result[In],
	fin fn.Transformer[In, fn.Result[inter]],
	finter1 fn.Transformer[inter, fn.Result[Out]],
) fn.Result[Out] {
	return fn.Morph(TwoRes(rin, fin), func(i inter) fn.Result[Out] {
		return finter1(i)
	})
}

func FourRes[In, inter1, inter2, Out any](
	rin fn.Result[In],
	fin fn.Transformer[In, fn.Result[inter1]],
	finter1 fn.Transformer[inter1, fn.Result[inter2]],
	finter2 fn.Transformer[inter2, fn.Result[Out]],
) fn.Result[Out] {
	return fn.Morph(ThreeRes(rin, fin, finter1), func(i2 inter2) fn.Result[Out] {
		return finter2(i2)
	})
}

func FiveRes[In, inter1, inter2, inter3, Out any](
	ra fn.Result[In],
	fin fn.Transformer[In, fn.Result[inter1]],
	finter1 fn.Transformer[inter1, fn.Result[inter2]],
	finter2 fn.Transformer[inter2, fn.Result[inter3]],
	finter3 fn.Transformer[inter3, fn.Result[Out]],
) fn.Result[Out] {
	return fn.Morph(FourRes(ra, fin, finter1, finter2), func(i3 inter3) fn.Result[Out] {
		return finter3(i3)
	})
}

func TwoWrap[In, Out any](
	in fn.Result[In],
	fin func(In) (Out, error),
) fn.Result[Out] {
	return fn.Morph(in, func(i In) fn.Result[Out] {
		return fn.FromAny(fin(i))
	})
}

func ThreeWrap[In, inter, Out any](
	in fn.Result[In],
	fin func(In) (inter, error),
	finter1 func(inter) (Out, error),
) fn.Result[Out] {
	return fn.Morph(TwoWrap(in, fin), func(i inter) fn.Result[Out] {
		return fn.FromAny(finter1(i))
	})
}

func FourWrap[In, inter1, inter2, Out any](
	in fn.Result[In],
	fin func(In) (inter1, error),
	finter1 func(inter1) (inter2, error),
	finter2 func(inter2) (Out, error),
) fn.Result[Out] {
	return fn.Morph(ThreeWrap(in, fin, finter1), func(i2 inter2) fn.Result[Out] {
		return fn.FromAny(finter2(i2))
	})
}

func FiveWrap[In, inter1, inter2, inter3, Out any](
	in fn.Result[In],
	fin func(In) (inter1, error),
	finter1 func(inter1) (inter2, error),
	finter2 func(inter2) (inter3, error),
	finter3 func(inter3) (Out, error),
) fn.Result[Out] {
	return fn.Morph(FourWrap(in, fin, finter1, finter2),
		func(i3 inter3) fn.Result[Out] {
			return fn.FromAny(finter3(i3))
		},
	)
}

func TwoWrapPtr[In, Out any](
	in fn.Result[In],
	fin func(In) (*Out, error),
) fn.Result[Out] {
	return fn.Morph(in, func(i In) fn.Result[Out] {
		return fn.FromPtr(fin(i))
	})
}

func ThreeWrapPtr[In, inter, Out any](
	in fn.Result[In],
	fin func(In) (inter, error),
	finter1 func(inter) (*Out, error),
) fn.Result[Out] {
	return fn.Morph(TwoWrap(in, fin), func(i inter) fn.Result[Out] {
		return fn.FromPtr(finter1(i))
	})
}

func FourWrapPtr[In, inter1, inter2, Out any](
	in fn.Result[In],
	fin func(In) (inter1, error),
	finter1 func(inter1) (inter2, error),
	finter2 func(inter2) (*Out, error),
) fn.Result[Out] {
	return fn.Morph(ThreeWrap(in, fin, finter1), func(i2 inter2) fn.Result[Out] {
		return fn.FromPtr(finter2(i2))
	})
}

func FiveWrapPtr[In, inter1, inter2, inter3, Out any](
	in fn.Result[In],
	fin func(In) (inter1, error),
	finter1 func(inter1) (inter2, error),
	finter2 func(inter2) (inter3, error),
	finter3 func(inter3) (*Out, error),
) fn.Result[Out] {
	return fn.Morph(FourWrap(in, fin, finter1, finter2),
		func(i3 inter3) fn.Result[Out] {
			return fn.FromPtr(finter3(i3))
		},
	)
}
