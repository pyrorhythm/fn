package option

import "testing"

func TestSome(t *testing.T) {
	o := Some(42)
	if !o.Valid() {
		t.Error("Some(42) should be valid")
	}
	if o.Val() != 42 {
		t.Errorf("Val() = %d, want 42", o.Val())
	}
}

func TestSomeZero(t *testing.T) {
	o := Some(0)
	if o.Valid() {
		t.Error("Some(0) should be invalid (zero value)")
	}
}

func TestSomePtr(t *testing.T) {
	v := 42
	o := SomePtr(&v)
	if !o.Valid() {
		t.Error("SomePtr(&v) should be valid")
	}
	if o.Val() != 42 {
		t.Errorf("Val() = %d, want 42", o.Val())
	}
}

func TestSomePtrNil(t *testing.T) {
	o := SomePtr[int](nil)
	if o.Valid() {
		t.Error("SomePtr(nil) should be invalid")
	}
}

func TestSomeAny(t *testing.T) {
	o := SomeAny(0)
	if !o.Valid() {
		t.Error("SomeAny(0) should be valid (bypasses zero check)")
	}
}

func TestNil(t *testing.T) {
	o := Nil[int]()
	if o.Valid() {
		t.Error("Nil() should be invalid")
	}
}

func TestFold(t *testing.T) {
	some := Some(10)
	got := some.Fold(func() int { return -1 }, func(v int) int { return v * 2 })
	if got != 20 {
		t.Errorf("Fold on Some = %d, want 20", got)
	}

	none := Nil[int]()
	got = none.Fold(func() int { return -1 }, func(v int) int { return v * 2 })
	if got != -1 {
		t.Errorf("Fold on Nil = %d, want -1", got)
	}
}

func TestFlatMap(t *testing.T) {
	double := func(v int) Of[int] { return SomeAny(v * 2) }

	some := Some(10)
	got := some.FlatMap(double)
	if !got.Valid() || got.Val() != 20 {
		t.Errorf("FlatMap on Some = %v, want Some(20)", got)
	}

	none := Nil[int]()
	got = none.FlatMap(double)
	if got.Valid() {
		t.Error("FlatMap on Nil should return Nil")
	}
}
