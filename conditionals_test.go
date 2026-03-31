package fn

import "testing"

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
