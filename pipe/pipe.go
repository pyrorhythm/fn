// Package pipe implements Raw, Err and [fn.Result] pipelines
package pipe

import "github.com/pyrorhythm/fn"

func Two[In, Out any](
	in In,
	fin func(In) Out,
) Out {
	return fin(in)
}

func Three[In, inter, Out any](
	in In,
	fin func(In) inter,
	finter func(inter) Out,
) Out {
	return finter(Two(in, fin))
}

func Four[In, inter1, inter2, Out any](
	in In,
	fin func(In) inter1,
	finter1 func(inter1) inter2,
	finter2 func(inter2) Out,
) Out {
	return finter2(Three(in, fin, finter1))
}

func Five[In, inter1, inter2, inter3, Out any](
	in In,
	fin func(In) inter1,
	finter1 func(inter1) inter2,
	finter2 func(inter2) inter3,
	finter3 func(inter3) Out,
) Out {
	return finter3(Four(in, fin, finter1, finter2))
}

func TwoErr[In, Out any](a fn.Result[In], fb func(In) (Out, error)) (z Out, _ error) {
	if a.Exc() {
		return z, a.Err()
	}
	return fb(a.Val())
}

func ThreeErr[In, inter, Out any](
	a fn.Result[In],
	fb func(In) (inter, error),
	fc func(inter) (Out, error),
) (z Out, _ error) {
	b, eb := TwoErr(a, fb)
	if eb != nil {
		return z, eb
	}
	return fc(b)
}

func FourErr[In, inter1, inter2, Out any](
	a fn.Result[In],
	fb func(In) (inter1, error),
	fc func(inter1) (inter2, error),
	fd func(inter2) (Out, error),
) (z Out, _ error) {
	c, ec := ThreeErr(a, fb, fc)
	if ec != nil {
		return z, ec
	}
	return fd(c)
}

func FiveErr[In, inter1, inter2, inter3, Out any](
	a fn.Result[In],
	fb func(In) (inter1, error),
	fc func(inter1) (inter2, error),
	fd func(inter2) (inter3, error),
	fe func(inter3) (Out, error),
) (z Out, _ error) {
	d, ec := FourErr(a, fb, fc, fd)
	if ec != nil {
		return z, ec
	}
	return fe(d)
}

func TwoRes[In, Out any](
	ra fn.Result[In],
	fb func(In) fn.Result[Out],
) fn.Result[Out] {
	return fn.Morph(ra, func(i In) fn.Result[Out] {
		return fb(i)
	})
}

func ThreeRes[In, inter, Out any](
	ra fn.Result[In],
	fb func(In) fn.Result[inter],
	fc func(inter) fn.Result[Out],
) fn.Result[Out] {
	return fn.Morph(TwoRes(ra, fb), func(i inter) fn.Result[Out] {
		return fc(i)
	})
}

func FourRes[In, inter1, inter2, Out any](
	ra fn.Result[In],
	fb func(In) fn.Result[inter1],
	fc func(inter1) fn.Result[inter2],
	fd func(inter2) fn.Result[Out],
) fn.Result[Out] {
	return fn.Morph(ThreeRes(ra, fb, fc), func(i2 inter2) fn.Result[Out] {
		return fd(i2)
	})
}

func FiveRes[In, inter1, inter2, inter3, Out any](
	ra fn.Result[In],
	fb func(In) fn.Result[inter1],
	fc func(inter1) fn.Result[inter2],
	fd func(inter2) fn.Result[inter3],
	fe func(inter3) fn.Result[Out],
) fn.Result[Out] {
	return fn.Morph(FourRes(ra, fb, fc, fd), func(i3 inter3) fn.Result[Out] {
		return fe(i3)
	})
}

func TwoWrap[In, Out any](
	a fn.Result[In],
	fb func(In) (Out, error),
) fn.Result[Out] {
	return fn.Morph(a, func(i In) fn.Result[Out] {
		return fn.FromAny(fb(i))
	})
}

func ThreeWrap[In, inter, Out any](
	a fn.Result[In],
	fb func(In) (inter, error),
	fc func(inter) (Out, error),
) fn.Result[Out] {
	return fn.Morph(TwoWrap(a, fb), func(i inter) fn.Result[Out] {
		return fn.FromAny(fc(i))
	})
}

func FourWrap[In, inter1, inter2, Out any](
	a fn.Result[In],
	fb func(In) (inter1, error),
	fc func(inter1) (inter2, error),
	fd func(inter2) (Out, error),
) fn.Result[Out] {
	return fn.Morph(ThreeWrap(a, fb, fc), func(i2 inter2) fn.Result[Out] {
		return fn.FromAny(fd(i2))
	})
}

func FiveWrap[In, inter1, inter2, inter3, Out any](
	a fn.Result[In],
	fb func(In) (inter1, error),
	fc func(inter1) (inter2, error),
	fd func(inter2) (inter3, error),
	fe func(inter3) (Out, error),
) fn.Result[Out] {
	return fn.Morph(FourWrap(a, fb, fc, fd),
		func(i3 inter3) fn.Result[Out] {
			return fn.FromAny(fe(i3))
		},
	)
}
