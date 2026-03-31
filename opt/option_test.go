package opt

import (
	"errors"
	"testing"
)

func TestSome_ValidValue(t *testing.T) {
	o := Some(42)
	if !o.Valid() || o.Val() != 42 {
		t.Errorf("Some(42): valid=%v val=%d", o.Valid(), o.Val())
	}
}

func TestSome_ZeroValue(t *testing.T) {
	if Some(0).Valid() {
		t.Error("Some(0) should not be valid")
	}
}

func TestSome_EmptyString(t *testing.T) {
	if Some("").Valid() {
		t.Error("Some(\"\") should not be valid")
	}
}

func TestSomePtr_NonNil(t *testing.T) {
	o := SomePtr(new(42))
	if !o.Valid() || o.Val() != 42 {
		t.Errorf("SomePtr(&v): valid=%v val=%d", o.Valid(), o.Val())
	}
}

func TestSomePtr_Nil(t *testing.T) {
	if SomePtr[int](nil).Valid() {
		t.Error("SomePtr(nil) should not be valid")
	}
}

func TestSomePtr_PointerToZero(t *testing.T) {
	o := SomePtr(new(0))
	if !o.Valid() || o.Val() != 0 {
		t.Error("SomePtr(&0) should be valid")
	}
}

func TestNil(t *testing.T) {
	if Nil[int]().Valid() {
		t.Error("Nil[int]() should not be valid")
	}
}

func TestOption_Ptr(t *testing.T) {
	p := Some(42).Ptr()
	if p == nil || *p != 42 {
		t.Errorf("Some(42).Ptr() = %v", p)
	}
	if Nil[int]().Ptr() != nil {
		t.Error("Nil.Ptr() should be nil")
	}
}

type boolImpl struct{ v bool }

func (b boolImpl) Bool() bool { return b.v }

type okImpl struct{ v bool }

func (o okImpl) Ok() bool { return o.v }

type validateErrImpl struct{ err error }

func (v validateErrImpl) Validate() error { return v.err }

func TestSomeAnyReflect_Slice(t *testing.T) {
	if !SomeAnyReflect([]int{1, 2, 3}).Valid() {
		t.Error("non-empty slice should be valid")
	}
	if SomeAnyReflect([]int{}).Valid() {
		t.Error("empty slice should not be valid")
	}
}

func TestSomeAnyReflect_Map(t *testing.T) {
	if !SomeAnyReflect(map[string]int{"a": 1}).Valid() {
		t.Error("non-empty map should be valid")
	}
	if SomeAnyReflect(map[string]int{}).Valid() {
		t.Error("empty map should not be valid")
	}
}

func TestSomeAnyReflect_BoolInterface(t *testing.T) {
	if !SomeAnyReflect(boolImpl{true}).Valid() {
		t.Error("Bool()=true should be valid")
	}
	if SomeAnyReflect(boolImpl{false}).Valid() {
		t.Error("Bool()=false should not be valid")
	}
}

func TestSomeAnyReflect_OKInterface(t *testing.T) {
	if !SomeAnyReflect(okImpl{true}).Valid() {
		t.Error("Ok()=true should be valid")
	}
	if SomeAnyReflect(okImpl{false}).Valid() {
		t.Error("Ok()=false should not be valid")
	}
}

func TestSomeAnyReflect_ValidateErrInterface(t *testing.T) {
	if !SomeAnyReflect(validateErrImpl{nil}).Valid() {
		t.Error("Validate()=nil should be valid")
	}
	if SomeAnyReflect(validateErrImpl{errors.New("e")}).Valid() {
		t.Error("Validate()=err should not be valid")
	}
}
