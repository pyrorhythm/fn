package fn

import (
	"errors"
	"testing"
)

func TestOK(t *testing.T) {
	r := OK(42)
	if !r.OK() {
		t.Error("OK(42) should be OK")
	}
	if r.Exc() {
		t.Error("OK(42) should not be Exc")
	}
	if r.Val() != 42 {
		t.Errorf("expected 42, got %d", r.Val())
	}
}

func TestOK_ZeroValue(t *testing.T) {
	r := OK(0)
	if r.OK() {
		t.Error("OK(0) should not be OK (zero value)")
	}
}

func TestErr(t *testing.T) {
	e := errors.New("test error")
	r := Err[int](e)
	if r.OK() {
		t.Error("Err should not be OK")
	}
	if !r.Exc() {
		t.Error("Err should be Exc")
	}
	if r.Err() != e {
		t.Error("Err() should return the error")
	}
}

func TestErrn(t *testing.T) {
	r := Errn[int]("test error")
	if r.OK() {
		t.Error("Errn should not be OK")
	}
	if r.Err() == nil {
		t.Error("Errn should have an error")
	}
	if r.Err().Error() != "test error" {
		t.Errorf("expected 'test error', got %q", r.Err().Error())
	}
}

func TestFrom_WithValue(t *testing.T) {
	r := From(42, nil)
	if !r.OK() {
		t.Error("From(42, nil) should be OK")
	}
	if r.Val() != 42 {
		t.Errorf("expected 42, got %d", r.Val())
	}
}

func TestFrom_WithError(t *testing.T) {
	e := errors.New("test error")
	r := From(42, e)
	if r.OK() {
		t.Error("From(42, err) should not be OK")
	}
	if r.Err() != e {
		t.Error("should have the error")
	}
}

func TestFromA(t *testing.T) {
	r := FromA([]int{1, 2, 3}, nil)
	if !r.OK() {
		t.Error("FromA with non-nil slice should be OK")
	}
}

func TestFromP_ValidPointer(t *testing.T) {
	v := 42
	r := FromP(&v, nil)
	if !r.OK() {
		t.Error("FromP(&v, nil) should be OK")
	}
	if r.Val() != 42 {
		t.Errorf("expected 42, got %d", r.Val())
	}
}

func TestFromP_NilPointer(t *testing.T) {
	var p *int
	r := FromP(p, nil)
	if r.OK() {
		t.Error("FromP(nil, nil) should not be OK")
	}
}

func TestFromAReflect_Slice(t *testing.T) {
	r := FromAReflect([]int{1, 2, 3}, nil)
	if !r.OK() {
		t.Error("FromAReflect with non-empty slice should be OK")
	}
}

func TestFromAReflect_EmptySlice(t *testing.T) {
	r := FromAReflect([]int{}, nil)
	if r.OK() {
		t.Error("FromAReflect with empty slice should not be OK")
	}
}

func TestFromAReflect_Map(t *testing.T) {
	r := FromAReflect(map[string]int{"a": 1}, nil)
	if !r.OK() {
		t.Error("FromAReflect with non-empty map should be OK")
	}
}

func TestFromAReflect_EmptyMap(t *testing.T) {
	r := FromAReflect(map[string]int{}, nil)
	if r.OK() {
		t.Error("FromAReflect with empty map should not be OK")
	}
}

func TestFromAReflect_WithError(t *testing.T) {
	e := errors.New("test error")
	r := FromAReflect([]int{1, 2, 3}, e)
	if r.OK() {
		t.Error("FromAReflect with error should not be OK")
	}
	if r.Err() != e {
		t.Error("should preserve the error")
	}
}

func TestFromAReflect_ZeroInt(t *testing.T) {
	r := FromAReflect(0, nil)
	if r.OK() {
		t.Error("FromAReflect(0) should not be OK")
	}
}

func TestFromAReflect_NonZeroInt(t *testing.T) {
	r := FromAReflect(42, nil)
	if !r.OK() {
		t.Error("FromAReflect(42) should be OK")
	}
	if r.Val() != 42 {
		t.Errorf("expected 42, got %d", r.Val())
	}
}

func TestOKPtr_ValidPointer(t *testing.T) {
	v := 42
	r := OKPtr(&v)
	if !r.OK() {
		t.Error("OKPtr(&v) should be OK")
	}
	if r.Val() != 42 {
		t.Errorf("expected 42, got %d", r.Val())
	}
}

func TestOKPtr_NilPointer(t *testing.T) {
	var p *int
	r := OKPtr(p)
	if r.OK() {
		t.Error("OKPtr(nil) should not be OK")
	}
	if r.Err() == nil {
		t.Error("OKPtr(nil) should have error")
	}
}

func TestResult_Into(t *testing.T) {
	r := OK(42)
	var v int
	e := r.Into(&v)
	if e != nil {
		t.Errorf("Into should return nil error, got %v", e)
	}
	if v != 42 {
		t.Errorf("expected v == 42, got %d", v)
	}
}

func TestResult_Into_WithError(t *testing.T) {
	err := errors.New("test error")
	r := Err[int](err)
	var v int
	e := r.Into(&v)
	if e != err {
		t.Error("Into should return the error")
	}
}

func TestResult_Unpack(t *testing.T) {
	r := OK(42)
	v, e := r.Unpack()
	if e != nil {
		t.Errorf("Unpack should return nil error, got %v", e)
	}
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestResult_Unpack_WithError(t *testing.T) {
	err := errors.New("test error")
	r := Err[int](err)
	_, e := r.Unpack()
	if e != err {
		t.Error("Unpack should return the error")
	}
}

func TestResult_Opt(t *testing.T) {
	r := OK(42)
	opt := r.Opt()
	if !opt.Valid() {
		t.Error("Opt() should be valid for OK result")
	}
	if opt.Value() != 42 {
		t.Errorf("expected 42, got %d", opt.Value())
	}
}

func TestResult_Ptr(t *testing.T) {
	r := OK(42)
	p := r.Ptr()
	if p == nil {
		t.Error("Ptr() should not be nil")
	}
	if *p != 42 {
		t.Errorf("expected *p == 42, got %d", *p)
	}
}

func TestResult_Any(t *testing.T) {
	r := OK(42)
	anyR := r.Any()
	if !anyR.OK() {
		t.Error("Any() should preserve OK status")
	}
	if anyR.Val() != 42 {
		t.Errorf("expected 42, got %v", anyR.Val())
	}
}
