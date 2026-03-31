package opt

import "testing"

func TestTo_Valid(t *testing.T) {
	got := To(Some(42), func(n int) string { return "value" })
	if !got.Valid() || got.Val() != "value" {
		t.Errorf("To(Some(42)) = %v, want Some(\"value\")", got)
	}
}

func TestTo_Nil(t *testing.T) {
	got := To(Nil[int](), func(n int) string { return "value" })
	if got.Valid() {
		t.Error("To(Nil) should return Nil")
	}
}

func TestMorph_Valid(t *testing.T) {
	got := Morph(Some(42), func(n int) Of[string] { return Some("ok") })
	if !got.Valid() || got.Val() != "ok" {
		t.Errorf("Morph(Some(42)) = %v, want Some(\"ok\")", got)
	}
}

func TestMorph_Nil(t *testing.T) {
	got := Morph(Nil[int](), func(n int) Of[string] { return Some("ok") })
	if got.Valid() {
		t.Error("Morph(Nil) should return Nil without calling f")
	}
}

func TestMorph_FReturnsNil(t *testing.T) {
	got := Morph(Some(42), func(n int) Of[string] { return Nil[string]() })
	if got.Valid() {
		t.Error("Morph should propagate Nil returned by f")
	}
}

func TestFlatMap_Valid(t *testing.T) {
	got := Some(10).FlatMap(func(n int) Of[int] { return Some(n * 2) })
	if !got.Valid() || got.Val() != 20 {
		t.Errorf("FlatMap = %v, want Some(20)", got)
	}
}

func TestFlatMap_Nil(t *testing.T) {
	got := Nil[int]().FlatMap(func(n int) Of[int] { return Some(n * 2) })
	if got.Valid() {
		t.Error("FlatMap(Nil) should return Nil")
	}
}

func TestFold_Valid(t *testing.T) {
	got := Some(5).Fold(func() int { return 0 }, func(n int) int { return n * 2 })
	if got != 10 {
		t.Errorf("Fold(Some(5)) = %d, want 10", got)
	}
}

func TestFold_Nil(t *testing.T) {
	got := Nil[int]().Fold(func() int { return -1 }, func(n int) int { return n })
	if got != -1 {
		t.Errorf("Fold(Nil) = %d, want -1", got)
	}
}
