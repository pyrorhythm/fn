package fn

import (
	"errors"
	"testing"
)

func TestIf_True(t *testing.T) {
	result := If(true, "yes", "no")
	if result != "yes" {
		t.Errorf("expected 'yes', got %q", result)
	}
}

func TestIf_False(t *testing.T) {
	result := If(false, "yes", "no")
	if result != "no" {
		t.Errorf("expected 'no', got %q", result)
	}
}

func TestIf_Int(t *testing.T) {
	result := If(10 > 5, 100, 0)
	if result != 100 {
		t.Errorf("expected 100, got %d", result)
	}
}

func TestFlatIf_True(t *testing.T) {
	opt := FlatIf(true, 42, 0)
	if !opt.Valid() {
		t.Error("FlatIf(true, ...) should be valid")
	}
	if opt.Value() != 42 {
		t.Errorf("expected 42, got %d", opt.Value())
	}
}

func TestFlatIf_False(t *testing.T) {
	opt := FlatIf(false, 42, 99)
	if !opt.Valid() {
		t.Error("FlatIf(false, ...) should be valid")
	}
	if opt.Value() != 99 {
		t.Errorf("expected 99, got %d", opt.Value())
	}
}

func TestErrIf_True(t *testing.T) {
	r := ErrIf(true, 42, errors.New("error"))
	if !r.OK() {
		t.Error("ErrIf(true, ...) should be OK")
	}
	if r.Val() != 42 {
		t.Errorf("expected 42, got %d", r.Val())
	}
}

func TestErrIf_False(t *testing.T) {
	testErr := errors.New("test error")
	r := ErrIf(false, 42, testErr)
	if r.OK() {
		t.Error("ErrIf(false, ...) should not be OK")
	}
	if r.Err() != testErr {
		t.Error("should have the error")
	}
}

func TestIfPtr_TrueUsesValue(t *testing.T) {
	fallback := 99
	opt := IfPtr(true, 42, &fallback)
	if !opt.Valid() {
		t.Error("should be valid")
	}
	if opt.Value() != 42 {
		t.Errorf("expected 42, got %d", opt.Value())
	}
}

func TestIfPtr_FalseUsesPointer(t *testing.T) {
	fallback := 99
	opt := IfPtr(false, 42, &fallback)
	if !opt.Valid() {
		t.Error("should be valid (pointer not nil)")
	}
	if opt.Value() != 99 {
		t.Errorf("expected 99, got %d", opt.Value())
	}
}

func TestIfPtr_FalseNilPointer(t *testing.T) {
	var p *int
	opt := IfPtr(false, 42, p)
	if opt.Valid() {
		t.Error("should not be valid (nil pointer)")
	}
}

func TestPtrIf_TrueUsesPointer(t *testing.T) {
	val := 42
	opt := PtrIf(true, &val, 99)
	if !opt.Valid() {
		t.Error("should be valid")
	}
	if opt.Value() != 42 {
		t.Errorf("expected 42, got %d", opt.Value())
	}
}

func TestPtrIf_TrueNilPointer(t *testing.T) {
	var p *int
	opt := PtrIf(true, p, 99)
	if opt.Valid() {
		t.Error("should not be valid (nil pointer)")
	}
}

func TestPtrIf_FalseUsesValue(t *testing.T) {
	val := 42
	opt := PtrIf(false, &val, 99)
	if !opt.Valid() {
		t.Error("should be valid")
	}
	if opt.Value() != 99 {
		t.Errorf("expected 99, got %d", opt.Value())
	}
}
