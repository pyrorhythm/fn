package fn

import "errors"

type (
	// rollbackable can undo its own work
	rollbackable interface {
		Rollback() error
	}

	// runnable can do some work
	runnable interface {
		Run() error
	}

	// operation can do work and undo it
	operation interface {
		rollbackable
		runnable
	}

	// ErrFunc is a simple adapter for func() error to [runnable] and [rollbackable]
	ErrFunc func() error
)

// Run implements [runnable]
func (ef ErrFunc) Run() error {
	return ef()
}

// Rollback implements [rollbackable]
func (ef ErrFunc) Rollback() error {
	return ef()
}

// RunOps runs [runnable]s sequentially, stops on first error
func RunOps(rbls ...runnable) (e error) {
	for _, r := range rbls {
		e = r.Run()
		if e != nil {
			return
		}
	}
	return
}

type stack []rollbackable

func (q *stack) push(n rollbackable) {
	*q = append(*q, n)
}

func (q *stack) pop() (n rollbackable) {
	x := q.c() - 1
	n = (*q)[x]
	*q = (*q)[:x]
	return
}

func (q *stack) c() int {
	return len(*q)
}

// TransactOps runs [operation]s sequentially.
// On error, rollbacks all previously executed ops in reverse order.
// Returns joined error if rollbacks also fail.
func TransactOps(ops ...operation) (e error) {
	st := new(stack)

	for _, op := range ops {
		st.push(op)

		e = op.Run()
		if e != nil {
			goto rb
		}
	}

	return

rb:
	errs := make([]error, 0)
	errs = append(errs, e)

	for st.c() != 0 {
		rbe := st.pop().Rollback()
		if rbe != nil {
			errs = append(errs, rbe)
		}
	}

	return errors.Join(errs...)
}
