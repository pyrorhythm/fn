package fn

import (
	"errors"
	"testing"
)

func TestRunOps_AllSuccess(t *testing.T) {
	var calls []int
	op1 := ErrFunc(func() error { calls = append(calls, 1); return nil })
	op2 := ErrFunc(func() error { calls = append(calls, 2); return nil })
	op3 := ErrFunc(func() error { calls = append(calls, 3); return nil })

	err := RunOps(op1, op2, op3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if len(calls) != 3 {
		t.Errorf("expected 3 calls, got %d", len(calls))
	}
	for i, c := range calls {
		if c != i+1 {
			t.Errorf("expected call %d, got %d", i+1, c)
		}
	}
}

func TestRunOps_StopsOnError(t *testing.T) {
	var calls []int
	testErr := errors.New("op2 failed")
	op1 := ErrFunc(func() error { calls = append(calls, 1); return nil })
	op2 := ErrFunc(func() error { calls = append(calls, 2); return testErr })
	op3 := ErrFunc(func() error { calls = append(calls, 3); return nil })

	err := RunOps(op1, op2, op3)
	if err != testErr {
		t.Errorf("expected testErr, got %v", err)
	}
	if len(calls) != 2 {
		t.Errorf("expected 2 calls (stop on error), got %d", len(calls))
	}
}

func TestRunOps_Empty(t *testing.T) {
	err := RunOps()
	if err != nil {
		t.Errorf("expected nil for empty ops, got %v", err)
	}
}

type testOp struct {
	runCalled      bool
	rollbackCalled bool
	runErr         error
	rollbackErr    error
}

func (o *testOp) Run() error {
	o.runCalled = true
	return o.runErr
}

func (o *testOp) Rollback() error {
	o.rollbackCalled = true
	return o.rollbackErr
}

func TestTransactOps_AllSuccess(t *testing.T) {
	op1 := &testOp{}
	op2 := &testOp{}
	op3 := &testOp{}

	err := TransactOps(op1, op2, op3)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if !op1.runCalled || !op2.runCalled || !op3.runCalled {
		t.Error("all ops should have Run called")
	}
	if op1.rollbackCalled || op2.rollbackCalled || op3.rollbackCalled {
		t.Error("no rollbacks should be called on success")
	}
}

func TestTransactOps_RollbackOnError(t *testing.T) {
	op1 := &testOp{}
	op2 := &testOp{runErr: errors.New("op2 failed")}
	op3 := &testOp{}

	err := TransactOps(op1, op2, op3)
	if err == nil {
		t.Error("expected error")
	}
	if !op1.runCalled {
		t.Error("op1 Run should be called")
	}
	if !op2.runCalled {
		t.Error("op2 Run should be called")
	}
	if op3.runCalled {
		t.Error("op3 Run should NOT be called (stopped before)")
	}
	if !op1.rollbackCalled {
		t.Error("op1 should be rolled back")
	}
	if !op2.rollbackCalled {
		t.Error("op2 should be rolled back")
	}
	if op3.rollbackCalled {
		t.Error("op3 should NOT be rolled back (never ran)")
	}
}

func TestTransactOps_RollbackOrder(t *testing.T) {
	var rollbackOrder []int

	op1 := &orderTrackOp{id: 1, order: &rollbackOrder}
	op2 := &orderTrackOp{id: 2, order: &rollbackOrder}
	op3 := &orderTrackOp{id: 3, runErr: errors.New("fail"), order: &rollbackOrder}

	_ = TransactOps(op1, op2, op3)

	// should rollback in reverse: 3, 2, 1
	expected := []int{3, 2, 1}
	if len(rollbackOrder) != 3 {
		t.Errorf("expected 3 rollbacks, got %d", len(rollbackOrder))
	}
	for i, id := range rollbackOrder {
		if id != expected[i] {
			t.Errorf("rollback order[%d]: expected %d, got %d", i, expected[i], id)
		}
	}
}

type orderTrackOp struct {
	id          int
	runErr      error
	rollbackErr error
	order       *[]int
}

func (o *orderTrackOp) Run() error      { return o.runErr }
func (o *orderTrackOp) Rollback() error { *o.order = append(*o.order, o.id); return o.rollbackErr }

func TestTransactOps_RollbackErrorsJoined(t *testing.T) {
	runErr := errors.New("run failed")
	rb1Err := errors.New("rollback 1 failed")
	rb2Err := errors.New("rollback 2 failed")

	op1 := &testOp{rollbackErr: rb1Err}
	op2 := &testOp{runErr: runErr, rollbackErr: rb2Err}

	err := TransactOps(op1, op2)
	if err == nil {
		t.Error("expected error")
	}
	// error should contain all three errors joined
	errStr := err.Error()
	if !contains(errStr, "run failed") {
		t.Error("should contain run error")
	}
	if !contains(errStr, "rollback 1 failed") {
		t.Error("should contain rollback 1 error")
	}
	if !contains(errStr, "rollback 2 failed") {
		t.Error("should contain rollback 2 error")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestTransactOps_Empty(t *testing.T) {
	err := TransactOps()
	if err != nil {
		t.Errorf("expected nil for empty ops, got %v", err)
	}
}

func TestErrFunc_Run(t *testing.T) {
	called := false
	ef := ErrFunc(func() error { called = true; return nil })
	err := ef.Run()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if !called {
		t.Error("Run should call the function")
	}
}

func TestErrFunc_Rollback(t *testing.T) {
	called := false
	ef := ErrFunc(func() error { called = true; return nil })
	err := ef.Rollback()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if !called {
		t.Error("Rollback should call the function")
	}
}
