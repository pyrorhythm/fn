package pipe

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/pyrorhythm/fn"
)

// Pure pipes — test type transformations

func TestTwo(t *testing.T) {
	got := Two(42, strconv.Itoa)
	if got != "42" {
		t.Errorf("Two = %q, want %q", got, "42")
	}
}

func TestThree(t *testing.T) {
	got := Three(42,
		strconv.Itoa,
		func(s string) []byte { return []byte(s) },
	)
	if string(got) != "42" {
		t.Errorf("Three = %v, want []byte(\"42\")", got)
	}
}

func TestFour(t *testing.T) {
	got := Four(42,
		strconv.Itoa,
		func(s string) []byte { return []byte(s) },
		func(b []byte) int { return len(b) },
	)
	if got != 2 {
		t.Errorf("Four = %d, want 2", got)
	}
}

func TestFive(t *testing.T) {
	got := Five(42,
		strconv.Itoa,
		func(s string) []byte { return []byte(s) },
		func(b []byte) int { return len(b) },
		func(n int) bool { return n > 1 },
	)
	if got != true {
		t.Errorf("Five = %v, want true", got)
	}
}

// Error pipes — test type transformations with error handling

func TestTwoErr(t *testing.T) {
	got, err := TwoErr(42, nil, func(n int) (string, error) { return strconv.Itoa(n), nil })
	if err != nil || got != "42" {
		t.Errorf("TwoErr = %q, %v, want \"42\", nil", got, err)
	}

	_, err = TwoErr(42, errors.New("init"), func(n int) (string, error) { return strconv.Itoa(n), nil })
	if err == nil {
		t.Error("TwoErr should propagate initial error")
	}

	_, err = TwoErr(-1, nil, func(n int) (string, error) {
		if n < 0 {
			return "", errors.New("negative")
		}
		return strconv.Itoa(n), nil
	})
	if err == nil {
		t.Error("TwoErr should propagate func error")
	}
}

func TestThreeErr(t *testing.T) {
	got, err := ThreeErr("42", nil,
		strconv.Atoi,
		func(n int) (float64, error) { return float64(n) * 1.5, nil },
	)
	if err != nil || got != 63.0 {
		t.Errorf("ThreeErr = %v, %v, want 63.0, nil", got, err)
	}

	_, err = ThreeErr("bad", nil,
		strconv.Atoi,
		func(n int) (float64, error) { return float64(n), nil },
	)
	if err == nil {
		t.Error("ThreeErr should short-circuit on first error")
	}
}

func TestFourErr(t *testing.T) {
	got, err := FourErr("42", nil,
		strconv.Atoi,
		func(n int) (float64, error) { return float64(n), nil },
		func(f float64) (string, error) { return strconv.FormatFloat(f, 'f', 1, 64), nil },
	)
	if err != nil || got != "42.0" {
		t.Errorf("FourErr = %q, %v, want \"42.0\", nil", got, err)
	}

	_, err = FourErr("42", nil,
		strconv.Atoi,
		func(n int) (float64, error) { return 0, errors.New("fail") },
		func(f float64) (string, error) { return "", nil },
	)
	if err == nil {
		t.Error("FourErr should short-circuit on middle error")
	}
}

func TestFiveErr(t *testing.T) {
	got, err := FiveErr("  42  ", nil,
		func(s string) (string, error) { return strings.TrimSpace(s), nil },
		strconv.Atoi,
		func(n int) (float64, error) { return float64(n) * 2.5, nil },
		func(f float64) (int, error) { return int(f), nil },
	)
	if err != nil || got != 105 {
		t.Errorf("FiveErr = %d, %v, want 105, nil", got, err)
	}
}

// Result pipes — test type transformations with Result monad

func TestTwoRes(t *testing.T) {
	got := TwoRes(
		fn.OKAny(42),
		func(n int) fn.Result[string] { return fn.OKAny(strconv.Itoa(n)) },
	)
	if !got.OK() || got.Val() != "42" {
		t.Errorf("TwoRes = %v, want OK(\"42\")", got)
	}

	got = TwoRes(
		fn.Err[int](errors.New("init")),
		func(n int) fn.Result[string] { return fn.OKAny(strconv.Itoa(n)) },
	)
	if got.OK() {
		t.Error("TwoRes should propagate initial error")
	}
}

func TestThreeRes(t *testing.T) {
	got := ThreeRes(
		fn.OKAny("42"),
		func(s string) fn.Result[int] { return fn.FromAny(strconv.Atoi(s)) },
		func(n int) fn.Result[float64] { return fn.OKAny(float64(n) * 1.5) },
	)
	if !got.OK() || got.Val() != 63.0 {
		t.Errorf("ThreeRes = %v, want OK(63.0)", got)
	}

	got = ThreeRes(
		fn.OKAny("bad"),
		func(s string) fn.Result[int] { return fn.FromAny(strconv.Atoi(s)) },
		func(n int) fn.Result[float64] { return fn.OKAny(float64(n)) },
	)
	if got.OK() {
		t.Error("ThreeRes should short-circuit on first error")
	}
}

func TestFourRes(t *testing.T) {
	got := FourRes(
		fn.OKAny("42"),
		func(s string) fn.Result[int] { return fn.FromAny(strconv.Atoi(s)) },
		func(n int) fn.Result[float64] { return fn.OKAny(float64(n)) },
		func(f float64) fn.Result[string] { return fn.OKAny(strconv.FormatFloat(f, 'f', 1, 64)) },
	)
	if !got.OK() || got.Val() != "42.0" {
		t.Errorf("FourRes = %v, want OK(\"42.0\")", got)
	}
}

func TestFiveRes(t *testing.T) {
	got := FiveRes(
		fn.OKAny("  42  "),
		func(s string) fn.Result[string] { return fn.OKAny(strings.TrimSpace(s)) },
		func(s string) fn.Result[int] { return fn.FromAny(strconv.Atoi(s)) },
		func(n int) fn.Result[float64] { return fn.OKAny(float64(n) * 2.5) },
		func(f float64) fn.Result[int] { return fn.OKAny(int(f)) },
	)
	if !got.OK() || got.Val() != 105 {
		t.Errorf("FiveRes = %v, want OK(105)", got)
	}
}
