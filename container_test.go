package fn_test

import (
	"testing"

	"github.com/pyrorhythm/fn"
	"github.com/pyrorhythm/fn/opt"
	"github.com/pyrorhythm/fn/res"
)

func TestElse(t *testing.T) {
	if got := fn.Else(opt.Some(42), 0); got != 42 {
		t.Errorf("Else(Some(42), 0) = %d, want 42", got)
	}
	if got := fn.Else(opt.Nil[int](), 99); got != 99 {
		t.Errorf("Else(Nil(), 99) = %d, want 99", got)
	}
	if got := fn.Else(res.OK(42), 0); got != 42 {
		t.Errorf("Else(OK(42), 0) = %d, want 42", got)
	}
	if got := fn.Else(res.Err[int](nil), 99); got != 99 {
		t.Errorf("Else(Err(), 99) = %d, want 99", got)
	}
}

func TestMust(t *testing.T) {
	if got := fn.Must(opt.Some(42)); got != 42 {
		t.Errorf("Must(Some(42)) = %d, want 42", got)
	}
	if got := fn.Must(res.OK(42)); got != 42 {
		t.Errorf("Must(OK(42)) = %d, want 42", got)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Must(Nil()) should panic")
		}
	}()
	fn.Must(opt.Nil[int]())
}

func TestFold(t *testing.T) {
	onNil := func() string { return "empty" }
	onVal := func(n int) string { return "got value" }

	if got := fn.Match(opt.Some(42), onNil, onVal); got != "got value" {
		t.Errorf("Match(Some(42)) = %q, want \"got value\"", got)
	}
	if got := fn.Match(opt.Nil[int](), onNil, onVal); got != "empty" {
		t.Errorf("Match(Nil()) = %q, want \"empty\"", got)
	}
	if got := fn.Match(res.OK(42), onNil, onVal); got != "got value" {
		t.Errorf("Match(OK(42)) = %q, want \"got value\"", got)
	}
	if got := fn.Match(res.Err[int](nil), onNil, onVal); got != "empty" {
		t.Errorf("Match(Err()) = %q, want \"empty\"", got)
	}
}
