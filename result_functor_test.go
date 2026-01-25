package fn

import (
	"errors"
	"strconv"
	"testing"
)

func TestTo_OK(t *testing.T) {
	r := OK(42)
	r2 := To(r, func(i int) string { return strconv.Itoa(i) })
	if !r2.OK() {
		t.Error("To on OK result should be OK")
	}
	if r2.Val() != "42" {
		t.Errorf("expected \"42\", got %q", r2.Val())
	}
}

func TestTo_Err(t *testing.T) {
	e := errors.New("test error")
	r := Err[int](e)
	r2 := To(r, func(i int) string { return strconv.Itoa(i) })
	if r2.OK() {
		t.Error("To on Err should propagate error")
	}
	if r2.Err() != e {
		t.Error("To should preserve original error")
	}
}

func TestTo_NoFunc(t *testing.T) {
	r := OK(42)
	r2 := To[int, string](r)
	if r2.OK() {
		t.Error("To without func should be Exc")
	}
	if r2.Err() == nil {
		t.Error("To without func should have error")
	}
}

func TestTo_Chain(t *testing.T) {
	r := OK(10)
	r2 := To(r, func(i int) int { return i * 2 })
	r3 := To(r2, func(i int) int { return i + 5 })
	if !r3.OK() {
		t.Error("chained To should be OK")
	}
	if r3.Val() != 25 {
		t.Errorf("expected 25, got %d", r3.Val())
	}
}

func TestTo_ChainWithError(t *testing.T) {
	e := errors.New("first error")
	r := Err[int](e)
	r2 := To(r, func(i int) int { return i * 2 })
	r3 := To(r2, func(i int) int { return i + 5 })
	if r3.OK() {
		t.Error("chained To with error should propagate")
	}
	if r3.Err() != e {
		t.Error("should preserve first error")
	}
}

func TestMorph_OK(t *testing.T) {
	r := OK(42)
	r2 := Morph(r, func(i int) Result[string] {
		return OK(strconv.Itoa(i))
	})
	if !r2.OK() {
		t.Error("Morph on OK should be OK")
	}
	if r2.Val() != "42" {
		t.Errorf("expected \"42\", got %q", r2.Val())
	}
}

func TestMorph_Err(t *testing.T) {
	e := errors.New("test error")
	r := Err[int](e)
	r2 := Morph(r, func(i int) Result[string] {
		return OK(strconv.Itoa(i))
	})
	if r2.OK() {
		t.Error("Morph on Err should propagate error")
	}
	if r2.Err() != e {
		t.Error("Morph should preserve original error")
	}
}

func TestMorph_FuncReturnsErr(t *testing.T) {
	r := OK(42)
	e := errors.New("inner error")
	r2 := Morph(r, func(i int) Result[string] {
		return Err[string](e)
	})
	if r2.OK() {
		t.Error("Morph should propagate inner error")
	}
	if r2.Err() != e {
		t.Error("Morph should have inner error")
	}
}

func TestMorph_Chain(t *testing.T) {
	r := OK(10)
	r2 := Morph(r, func(i int) Result[int] {
		if i > 5 {
			return OK(i * 2)
		}
		return Errn[int]("too small")
	})
	r3 := Morph(r2, func(i int) Result[string] {
		return OK(strconv.Itoa(i))
	})
	if !r3.OK() {
		t.Error("chained Morph should be OK")
	}
	if r3.Val() != "20" {
		t.Errorf("expected \"20\", got %q", r3.Val())
	}
}

func TestMorph_ChainShortCircuit(t *testing.T) {
	r := OK(3)
	r2 := Morph(r, func(i int) Result[int] {
		if i > 5 {
			return OK(i * 2)
		}
		return Errn[int]("too small")
	})
	r3 := Morph(r2, func(i int) Result[string] {
		return OK(strconv.Itoa(i))
	})
	if r3.OK() {
		t.Error("should short-circuit on first error")
	}
	if r3.Err().Error() != "too small" {
		t.Errorf("expected 'too small', got %q", r3.Err().Error())
	}
}
