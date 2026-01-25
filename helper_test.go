package fn

import "testing"

func TestValid_Int(t *testing.T) {
	if Valid(0) {
		t.Error("0 should not be valid")
	}
	if !Valid(42) {
		t.Error("42 should be valid")
	}
	if !Valid(-1) {
		t.Error("-1 should be valid")
	}
}

func TestValid_String(t *testing.T) {
	if Valid("") {
		t.Error("empty string should not be valid")
	}
	if !Valid("hello") {
		t.Error("non-empty string should be valid")
	}
}

func TestValid_Pointer(t *testing.T) {
	var p *int
	if Valid(p) {
		t.Error("nil pointer should not be valid")
	}
	v := 42
	if !Valid(&v) {
		t.Error("non-nil pointer should be valid")
	}
}

func TestValidReflect_Slice(t *testing.T) {
	if ValidReflect([]int{}) {
		t.Error("empty slice should not be valid")
	}
	if !ValidReflect([]int{1, 2, 3}) {
		t.Error("non-empty slice should be valid")
	}
}

func TestValidReflect_Map(t *testing.T) {
	if ValidReflect(map[string]int{}) {
		t.Error("empty map should not be valid")
	}
	if !ValidReflect(map[string]int{"a": 1}) {
		t.Error("non-empty map should be valid")
	}
}

func TestValidReflect_Int(t *testing.T) {
	if ValidReflect(0) {
		t.Error("0 should not be valid")
	}
	if !ValidReflect(42) {
		t.Error("42 should be valid")
	}
}

type boolImpl struct{ v bool }

func (b boolImpl) Bool() bool { return b.v }

func TestValidReflect_BoolInterface(t *testing.T) {
	if ValidReflect(boolImpl{false}) {
		t.Error("Bool() false should not be valid")
	}
	if !ValidReflect(boolImpl{true}) {
		t.Error("Bool() true should be valid")
	}
}

type okImpl struct{ v bool }

func (o okImpl) Ok() bool { return o.v }

func TestValidReflect_OkInterface(t *testing.T) {
	if ValidReflect(okImpl{false}) {
		t.Error("Ok() false should not be valid")
	}
	if !ValidReflect(okImpl{true}) {
		t.Error("Ok() true should be valid")
	}
}

type validateErrImpl struct{ err error }

func (v validateErrImpl) Validate() error { return v.err }

func TestValidReflect_ValidateErrorInterface(t *testing.T) {
	if ValidReflect(validateErrImpl{err: nil}) != true {
		t.Error("Validate() nil should be valid")
	}
	if ValidReflect(validateErrImpl{err: errTest}) {
		t.Error("Validate() error should not be valid")
	}
}

var errTest = &testError{}

type testError struct{}

func (e *testError) Error() string { return "test" }

type validateBoolImpl struct{ v bool }

func (v validateBoolImpl) Validate() bool { return v.v }

func TestValidReflect_ValidateBoolInterface(t *testing.T) {
	if ValidReflect(validateBoolImpl{false}) {
		t.Error("Validate() false should not be valid")
	}
	if !ValidReflect(validateBoolImpl{true}) {
		t.Error("Validate() true should be valid")
	}
}

type isZeroImpl struct{ zero bool }

func (i isZeroImpl) IsZero() bool { return i.zero }

func TestValidReflect_IsZeroInterface(t *testing.T) {
	if ValidReflect(isZeroImpl{true}) {
		t.Error("IsZero() true should not be valid")
	}
	if !ValidReflect(isZeroImpl{false}) {
		t.Error("IsZero() false should be valid")
	}
}

func TestOr_NonNil(t *testing.T) {
	v := 42
	result := Or(&v, 99)
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestOr_Nil(t *testing.T) {
	var p *int
	result := Or(p, 99)
	if result != 99 {
		t.Errorf("expected 99, got %d", result)
	}
}

func TestOrZero_NonNil(t *testing.T) {
	v := 42
	result := OrZero(&v)
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestOrZero_Nil(t *testing.T) {
	var p *int
	result := OrZero(p)
	if result != 0 {
		t.Errorf("expected 0, got %d", result)
	}
}

func TestOrZero_String(t *testing.T) {
	var p *string
	result := OrZero(p)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
