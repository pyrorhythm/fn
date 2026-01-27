package result

import (
	"errors"
	"testing"

	"github.com/pyrorhythm/fn"
)

func TestOK(t *testing.T) {
	r := OK(42)
	if !r.OK() {
		t.Error("OK(42) should be OK")
	}
	if r.Val() != 42 {
		t.Errorf("Val() = %d, want 42", r.Val())
	}
}

func TestOKZero(t *testing.T) {
	r := OK(0)
	if r.OK() {
		t.Error("OK(0) should not be OK (zero value)")
	}
}

func TestOKAny(t *testing.T) {
	r := OKAny(0)
	if !r.OK() {
		t.Error("OKAny(0) should be OK (bypasses zero check)")
	}
}

func TestOKPtr(t *testing.T) {
	v := 42
	r := OKPtr(&v)
	if !r.OK() {
		t.Error("OKPtr(&v) should be OK")
	}
}

func TestOKPtrNil(t *testing.T) {
	r := OKPtr[int](nil)
	if r.OK() {
		t.Error("OKPtr(nil) should not be OK")
	}
}

func TestOKOpt(t *testing.T) {
	o := fn.Some(42)
	r := OKOpt(o)
	if !r.OK() {
		t.Error("OKOpt(Some(42)) should be OK")
	}
}

func TestErr(t *testing.T) {
	e := errors.New("fail")
	r := Err[int](e)
	if r.OK() {
		t.Error("Err should not be OK")
	}
	if r.Err() != e {
		t.Error("Err() should return original error")
	}
}

func TestErrn(t *testing.T) {
	r := Errn[int]("fail")
	if r.OK() {
		t.Error("Errn should not be OK")
	}
	if r.Err().Error() != "fail" {
		t.Errorf("Err().Error() = %q, want %q", r.Err().Error(), "fail")
	}
}

func TestFrom(t *testing.T) {
	r := From(42, nil)
	if !r.OK() {
		t.Error("From(42, nil) should be OK")
	}

	r = From(42, errors.New("fail"))
	if r.OK() {
		t.Error("From(42, err) should not be OK")
	}
}

func TestFromPtr(t *testing.T) {
	v := 42
	r := FromPtr(&v, nil)
	if !r.OK() {
		t.Error("FromPtr(&v, nil) should be OK")
	}
}

func TestFromAny(t *testing.T) {
	r := FromAny(0, nil)
	if !r.OK() {
		t.Error("FromAny(0, nil) should be OK")
	}
}

func TestFromOpt(t *testing.T) {
	o := fn.Some(42)
	r := FromOpt(o, nil)
	if !r.OK() {
		t.Error("FromOpt(Some(42), nil) should be OK")
	}
}

func TestValid(t *testing.T) {
	r := OK(42)
	if !r.Valid() {
		t.Error("Valid() should match OK()")
	}
}

func TestFold(t *testing.T) {
	ok := OK(10)
	got := ok.Fold(func() int { return -1 }, func(v int) int { return v * 2 })
	if got != 20 {
		t.Errorf("Fold on OK = %d, want 20", got)
	}

	fail := Err[int](errors.New("fail"))
	got = fail.Fold(func() int { return -1 }, func(v int) int { return v * 2 })
	if got != -1 {
		t.Errorf("Fold on Err = %d, want -1", got)
	}
}

func TestFlatMap(t *testing.T) {
	double := func(v int) Of[int] { return OKAny(v * 2) }

	ok := OK(10)
	got := ok.FlatMap(double)
	if !got.OK() || got.Val() != 20 {
		t.Errorf("FlatMap on OK = %v, want OK(20)", got)
	}

	fail := Err[int](errors.New("fail"))
	got = fail.FlatMap(double)
	if got.OK() {
		t.Error("FlatMap on Err should return Err")
	}
}
