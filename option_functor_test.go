package fn

import (
	"strconv"
	"testing"
)

func TestOptTo_Valid(t *testing.T) {
	opt := Some(42)
	opt2 := OptTo(opt, func(i int) string { return strconv.Itoa(i) })
	if !opt2.Valid() {
		t.Error("OptTo on valid option should be valid")
	}
	if opt2.Value() != "42" {
		t.Errorf("expected \"42\", got %q", opt2.Value())
	}
}

func TestOptTo_Invalid(t *testing.T) {
	opt := Nil[int]()
	opt2 := OptTo(opt, func(i int) string { return strconv.Itoa(i) })
	if opt2.Valid() {
		t.Error("OptTo on invalid option should be invalid")
	}
}

func TestOptTo_NoFunc(t *testing.T) {
	opt := Some(42)
	opt2 := OptTo[int, string](opt)
	if opt2.Valid() {
		t.Error("OptTo without func should be invalid")
	}
}

func TestOptTo_NoFunc_InvalidOption(t *testing.T) {
	opt := Nil[int]()
	opt2 := OptTo[int, string](opt)
	if opt2.Valid() {
		t.Error("OptTo without func on invalid option should be invalid")
	}
}

func TestOptTo_Chain(t *testing.T) {
	opt := Some(10)
	opt2 := OptTo(opt, func(i int) int { return i * 2 })
	opt3 := OptTo(opt2, func(i int) int { return i + 5 })
	if !opt3.Valid() {
		t.Error("chained OptTo should be valid")
	}
	if opt3.Value() != 25 {
		t.Errorf("expected 25, got %d", opt3.Value())
	}
}

func TestOptTo_ChainPropagatesInvalid(t *testing.T) {
	opt := Nil[int]()
	opt2 := OptTo(opt, func(i int) int { return i * 2 })
	opt3 := OptTo(opt2, func(i int) int { return i + 5 })
	if opt3.Valid() {
		t.Error("chained OptTo with invalid should propagate invalid")
	}
}

func TestOptMorph_Valid(t *testing.T) {
	opt := Some(42)
	opt2 := OptMorph(opt, func(i int) Option[string] {
		return Some(strconv.Itoa(i))
	})
	if !opt2.Valid() {
		t.Error("OptMorph on valid option should be valid")
	}
	if opt2.Value() != "42" {
		t.Errorf("expected \"42\", got %q", opt2.Value())
	}
}

func TestOptMorph_Invalid(t *testing.T) {
	opt := Nil[int]()
	opt2 := OptMorph(opt, func(i int) Option[string] {
		return Some(strconv.Itoa(i))
	})
	if opt2.Valid() {
		t.Error("OptMorph on invalid option should be invalid")
	}
}

func TestOptMorph_FuncReturnsNil(t *testing.T) {
	opt := Some(42)
	opt2 := OptMorph(opt, func(i int) Option[string] {
		return Nil[string]()
	})
	if opt2.Valid() {
		t.Error("OptMorph should propagate inner Nil")
	}
}

func TestOptMorph_Chain(t *testing.T) {
	opt := Some(10)
	opt2 := OptMorph(opt, func(i int) Option[int] {
		if i > 5 {
			return Some(i * 2)
		}
		return Nil[int]()
	})
	opt3 := OptMorph(opt2, func(i int) Option[string] {
		return Some(strconv.Itoa(i))
	})
	if !opt3.Valid() {
		t.Error("chained OptMorph should be valid")
	}
	if opt3.Value() != "20" {
		t.Errorf("expected \"20\", got %q", opt3.Value())
	}
}

func TestOptMorph_ChainShortCircuit(t *testing.T) {
	opt := Some(3)
	opt2 := OptMorph(opt, func(i int) Option[int] {
		if i > 5 {
			return Some(i * 2)
		}
		return Nil[int]()
	})
	opt3 := OptMorph(opt2, func(i int) Option[string] {
		return Some(strconv.Itoa(i))
	})
	if opt3.Valid() {
		t.Error("should short-circuit on first Nil")
	}
}
