package fn

import (
	"errors"
	"testing"
)

func TestSome_ValidValue(t *testing.T) {
	opt := Some(42)
	if !opt.Valid() {
		t.Error("Some(42) should be valid")
	}
	if opt.Val() != 42 {
		t.Errorf("expected 42, got %d", opt.Val())
	}
}

func TestSome_ZeroValue(t *testing.T) {
	opt := Some(0)
	if opt.Valid() {
		t.Error("Some(0) should not be valid (zero value)")
	}
}

func TestSome_EmptyString(t *testing.T) {
	opt := Some("")
	if opt.Valid() {
		t.Error("Some(\"\") should not be valid (zero value)")
	}
}

func TestSome_NonEmptyString(t *testing.T) {
	opt := Some("hello")
	if !opt.Valid() {
		t.Error("Some(\"hello\") should be valid")
	}
	if opt.Val() != "hello" {
		t.Errorf("expected \"hello\", got %q", opt.Val())
	}
}

func TestSomePtr_ValidPointer(t *testing.T) {
	v := 42
	opt := SomePtr(&v)
	if !opt.Valid() {
		t.Error("SomePtr(&v) should be valid")
	}
	if opt.Val() != 42 {
		t.Errorf("expected 42, got %d", opt.Val())
	}
}

func TestSomePtr_NilPointer(t *testing.T) {
	var p *int
	opt := SomePtr(p)
	if opt.Valid() {
		t.Error("SomePtr(nil) should not be valid")
	}
}

func TestSomePtr_PointerToZero(t *testing.T) {
	v := 0
	opt := SomePtr(&v)
	if !opt.Valid() {
		t.Error("SomePtr(&0) should be valid (pointer is non-nil)")
	}
	if opt.Val() != 0 {
		t.Errorf("expected 0, got %d", opt.Val())
	}
}

func TestNil(t *testing.T) {
	opt := Nil[int]()
	if opt.Valid() {
		t.Error("Nil[int]() should not be valid")
	}
}

func TestOption_Ptr(t *testing.T) {
	opt := Some(42)
	p := opt.Ptr()
	if p == nil {
		t.Error("Ptr() should not return nil")
	}
	if p != nil && *p != 42 {
		t.Errorf("expected *p == 42, got %d", *p)
	}
}

func TestOption_PtrInvalid(t *testing.T) {
	opt := Nil[int]()
	p := opt.Ptr()
	if p != nil {
		t.Error("Ptr() should return nil if Option is invalid")
	}
}

func TestOption_Any(t *testing.T) {
	opt := Some(42)
	anyOpt := opt.Any()
	if !anyOpt.Valid() {
		t.Error("Any() should preserve validity")
	}
	if anyOpt.Val() != 42 {
		t.Errorf("expected 42, got %v", anyOpt.Val())
	}
}

func TestSomeAnyReflect_Slice(t *testing.T) {
	opt := SomeAnyReflect([]int{1, 2, 3})
	if !opt.Valid() {
		t.Error("SomeAnyReflect(non-empty slice) should be valid")
	}
}

func TestSomeAnyReflect_EmptySlice(t *testing.T) {
	opt := SomeAnyReflect([]int{})
	if opt.Valid() {
		t.Error("SomeAnyReflect(empty slice) should not be valid")
	}
}

func TestSomeAnyReflect_Map(t *testing.T) {
	opt := SomeAnyReflect(map[string]int{"a": 1})
	if !opt.Valid() {
		t.Error("SomeAnyReflect(non-empty map) should be valid")
	}
}

func TestSomeAnyReflect_EmptyMap(t *testing.T) {
	opt := SomeAnyReflect(map[string]int{})
	if opt.Valid() {
		t.Error("SomeAnyReflect(empty map) should not be valid")
	}
}

type boolValidator struct{ valid bool }

func (b boolValidator) Bool() bool { return b.valid }

func TestSomeAnyReflect_BoolInterface(t *testing.T) {
	opt := SomeAnyReflect(boolImpl{true})
	if !opt.Valid() {
		t.Error("SomeAnyReflect with Bool() true should be valid")
	}

	opt2 := SomeAnyReflect(boolImpl{false})
	if opt2.Valid() {
		t.Error("SomeAnyReflect with Bool() false should be invalid")
	}
}

func TestSomeAnyReflect_OKInterface(t *testing.T) {
	opt := SomeAnyReflect(okImpl{true})
	if !opt.Valid() {
		t.Error("SomeAnyReflect with Ok() true should be valid")
	}

	opt2 := SomeAnyReflect(okImpl{false})
	if opt2.Valid() {
		t.Error("SomeAnyReflect with Ok() false should be invalid")
	}
}

func TestSomeAnyReflect_ValidateErrInterface(t *testing.T) {
	opt := SomeAnyReflect(validateErrImpl{nil})
	if !opt.Valid() {
		t.Error("SomeAnyReflect with Validate() error == nil should be valid")
	}

	opt2 := SomeAnyReflect(validateErrImpl{errors.New("asdasd")})
	if opt2.Valid() {
		t.Error("SomeAnyReflect with Validate() error != nil should be invalid")
	}
}
