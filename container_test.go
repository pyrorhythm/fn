package fn

import "testing"

func TestOrElse(t *testing.T) {
	// Option
	got := Else(Some(42), 0)
	if got != 42 {
		t.Errorf("OrElse(Some(42), 0) = %d, want 42", got)
	}

	got = Else(Nil[int](), 99)
	if got != 99 {
		t.Errorf("OrElse(Nil(), 99) = %d, want 99", got)
	}

	// Result
	got = Else(OK(42), 0)
	if got != 42 {
		t.Errorf("OrElse(OK(42), 0) = %d, want 42", got)
	}

	got = Else(Err[int](nil), 99)
	if got != 99 {
		t.Errorf("OrElse(Err(), 99) = %d, want 99", got)
	}
}

func TestMust(t *testing.T) {
	// Option
	got := Must(Some(42))
	if got != 42 {
		t.Errorf("Must(Some(42)) = %d, want 42", got)
	}

	// Result
	got = Must(OK(42))
	if got != 42 {
		t.Errorf("Must(OK(42)) = %d, want 42", got)
	}

	// Panic on invalid
	defer func() {
		if r := recover(); r == nil {
			t.Error("Must(Nil()) should panic")
		}
	}()
	Must(Nil[int]())
}

func TestFold(t *testing.T) {
	onNil := func() string { return "empty" }
	onVal := func(n int) string { return "got value" }

	// Option
	got := Fold(Some(42), onNil, onVal)
	if got != "got value" {
		t.Errorf("Fold(Some(42)) = %q, want %q", got, "got value")
	}

	got = Fold(Nil[int](), onNil, onVal)
	if got != "empty" {
		t.Errorf("Fold(Nil()) = %q, want %q", got, "empty")
	}

	// Result
	got = Fold(OK(42), onNil, onVal)
	if got != "got value" {
		t.Errorf("Fold(OK(42)) = %q, want %q", got, "got value")
	}

	got = Fold(Err[int](nil), onNil, onVal)
	if got != "empty" {
		t.Errorf("Fold(Err()) = %q, want %q", got, "empty")
	}
}

func TestChain(t *testing.T) {
	double := func(n int) Option[int] { return SomeAny(n * 2) }

	// Option
	got := Chain(SomeAny(21), double)
	if !got.Valid() || got.Val() != 42 {
		t.Errorf("Chain(Some(21), double) = %v, want Some(42)", got)
	}

	got = Chain(Nil[int](), double)
	if got.Valid() {
		t.Error("Chain(Nil, double) should be Nil")
	}

	// Result
	doubleRes := func(n int) Result[int] { return OKAny(n * 2) }
	gotRes := Chain(OKAny(21), doubleRes)
	if !gotRes.OK() || gotRes.Val() != 42 {
		t.Errorf("Chain(OK(21), double) = %v, want OK(42)", gotRes)
	}
}

func TestAndThen(t *testing.T) {
	// Option
	got := AndThen(SomeAny(1), SomeAny(42))
	if !got.Valid() || got.Val() != 42 {
		t.Errorf("AndThen(Some(1), Some(42)) = %v, want Some(42)", got)
	}

	got = AndThen(Nil[int](), SomeAny(42))
	if got.Valid() {
		t.Error("AndThen(Nil, Some(42)) should be Nil")
	}

	// Result
	gotRes := AndThen(OKAny(1), OKAny(42))
	if !gotRes.OK() || gotRes.Val() != 42 {
		t.Errorf("AndThen(OK(1), OK(42)) = %v, want OK(42)", gotRes)
	}
}
