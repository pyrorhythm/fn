package chain

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/pyrorhythm/fn/res"
)

func TestWrap_OK(t *testing.T) {
	got := Wrap(res.OK(42), func(n int) (string, error) { return strconv.Itoa(n), nil })
	if !got.OK() || got.Val() != "42" {
		t.Errorf("Wrap = %v, want OK(\"42\")", got)
	}
}

func TestWrap_InitialError(t *testing.T) {
	got := Wrap(
		res.Err[int](errors.New("init")),
		func(n int) (string, error) { return strconv.Itoa(n), nil },
	)
	if got.OK() {
		t.Error("Wrap should propagate initial error")
	}
}

func TestWrap_FuncError(t *testing.T) {
	got := Wrap(res.OK(-1), func(n int) (string, error) {
		if n < 0 {
			return "", errors.New("negative")
		}
		return strconv.Itoa(n), nil
	})
	if got.OK() {
		t.Error("Wrap should wrap function error into Result")
	}
}

func TestWrap_Chain(t *testing.T) {
	rInt := Wrap(res.OK("42"), strconv.Atoi)
	rFloat := Wrap(rInt, func(n int) (float64, error) { return float64(n) * 1.5, nil })
	if !rFloat.OK() || rFloat.Val() != 63.0 {
		t.Errorf("Wrap chain = %v, want OK(63.0)", rFloat)
	}
}

func TestWrap_ChainShortCircuit(t *testing.T) {
	rInt := Wrap(res.OK("bad"), strconv.Atoi)
	rFloat := Wrap(rInt, func(n int) (float64, error) { return float64(n), nil })
	if rFloat.OK() {
		t.Error("Wrap chain should short-circuit on first error")
	}
}

func TestWrap_LongChain(t *testing.T) {
	rStr := Wrap(res.OK("  42  "), func(s string) (string, error) { return strings.TrimSpace(s), nil })
	rInt := Wrap(rStr, strconv.Atoi)
	rFloat := Wrap(rInt, func(n int) (float64, error) { return float64(n) * 2.5, nil })
	rFinal := Wrap(rFloat, func(f float64) (int, error) { return int(f), nil })
	if !rFinal.OK() || rFinal.Val() != 105 {
		t.Errorf("Wrap long chain = %v, want OK(105)", rFinal)
	}
}

func TestWrap3(t *testing.T) {
	got := Wrap3(res.OK("42"), strconv.Atoi, func(n int) (float64, error) { return float64(n) * 1.5, nil })
	if !got.OK() || got.Val() != 63.0 {
		t.Errorf("Wrap3 = %v, want OK(63.0)", got)
	}

	got2 := Wrap3(res.OK("bad"), strconv.Atoi, func(n int) (float64, error) { return float64(n), nil })
	if got2.OK() {
		t.Error("Wrap3 should short-circuit on first error")
	}
}

func TestWrap4(t *testing.T) {
	got := Wrap4(res.OK("42"),
		strconv.Atoi,
		func(n int) (float64, error) { return float64(n), nil },
		func(f float64) (string, error) { return strconv.FormatFloat(f, 'f', 1, 64), nil },
	)
	if !got.OK() || got.Val() != "42.0" {
		t.Errorf("Wrap4 = %v, want OK(\"42.0\")", got)
	}
}

func TestWrap5(t *testing.T) {
	got := Wrap5(res.OK("  42  "),
		func(s string) (string, error) { return strings.TrimSpace(s), nil },
		strconv.Atoi,
		func(n int) (float64, error) { return float64(n) * 2.5, nil },
		func(f float64) (int, error) { return int(f), nil },
	)
	if !got.OK() || got.Val() != 105 {
		t.Errorf("Wrap5 = %v, want OK(105)", got)
	}
}

func TestWrapPtr_OK(t *testing.T) {
	v := 42
	got := WrapPtr(res.OK("42"), func(s string) (*int, error) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		v = n
		return &v, nil
	})
	if !got.OK() || got.Val() != 42 {
		t.Errorf("WrapPtr = %v, want OK(42)", got)
	}
}

func TestWrapPtr_NilReturn(t *testing.T) {
	got := WrapPtr(res.OK("42"), func(s string) (*int, error) {
		return nil, nil
	})
	if got.OK() {
		t.Error("WrapPtr with nil pointer result should not be OK")
	}
}

func TestWrapPtr_FuncError(t *testing.T) {
	got := WrapPtr(res.OK("bad"), func(s string) (*int, error) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return &n, nil
	})
	if got.OK() {
		t.Error("WrapPtr should propagate function error")
	}
}
