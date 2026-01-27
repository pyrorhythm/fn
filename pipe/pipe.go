// Package pipe implements Raw, Err and [fn.Result] pipelines
package pipe

import "github.com/pyrorhythm/fn"

func Two[A, B any](a A, fb func(A) B) B {
	return fb(a)
}

func Three[A, B, C any](a A, fb func(A) B, fc func(B) C) C {
	return fc(Two(a, fb))
}

func Four[A, B, C, D any](a A, fb func(A) B, fc func(B) C, fd func(C) D) D {
	return fd(Three(a, fb, fc))
}

func Five[A, B, C, D, E any](a A, fb func(A) B, fc func(B) C, fd func(C) D, fe func(D) E) E {
	return fe(Four(a, fb, fc, fd))
}

func TwoErr[A, B any](a fn.Result[A], fb func(A) (B, error)) (z B, _ error) {
	if a.Exc() {
		return z, a.Err()
	}
	return fb(a.Val())
}

func ThreeErr[A, B, C any](a fn.Result[A], fb func(A) (B, error), fc func(B) (C, error)) (z C, _ error) {
	b, eb := TwoErr(a, fb)
	if eb != nil {
		return z, eb
	}
	return fc(b)
}

func FourErr[A, B, C, D any](a fn.Result[A], fb func(A) (B, error), fc func(B) (C, error), fd func(C) (D, error)) (z D, _ error) {
	c, ec := ThreeErr(a, fb, fc)
	if ec != nil {
		return z, ec
	}
	return fd(c)
}

func FiveErr[A, B, C, D, E any](a fn.Result[A], fb func(A) (B, error), fc func(B) (C, error), fd func(C) (D, error), fe func(D) (E, error)) (z E, _ error) {
	d, ec := FourErr(a, fb, fc, fd)
	if ec != nil {
		return z, ec
	}
	return fe(d)
}

func TwoRes[A, B any](ra fn.Result[A], fb func(A) fn.Result[B]) fn.Result[B] {
	if ra.Exc() {
		return fn.Err[B](ra.Err())
	}
	return fb(ra.Val())
}

func ThreeRes[A, B, C any](ra fn.Result[A], fb func(A) fn.Result[B], fc func(B) fn.Result[C]) fn.Result[C] {
	rb := TwoRes(ra, fb)
	if rb.Exc() {
		return fn.Err[C](rb.Err())
	}
	return fc(rb.Val())
}

func FourRes[A, B, C, D any](ra fn.Result[A], fb func(A) fn.Result[B], fc func(B) fn.Result[C], fd func(C) fn.Result[D]) fn.Result[D] {
	rc := ThreeRes(ra, fb, fc)
	if rc.Exc() {
		return fn.Err[D](rc.Err())
	}
	return fd(rc.Val())
}

func FiveRes[A, B, C, D, E any](ra fn.Result[A], fb func(A) fn.Result[B], fc func(B) fn.Result[C], fd func(C) fn.Result[D], fe func(D) fn.Result[E]) fn.Result[E] {
	rd := FourRes(ra, fb, fc, fd)
	if rd.Exc() {
		return fn.Err[E](rd.Err())
	}
	return fe(rd.Val())
}
