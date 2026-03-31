package react

import (
	"context"
	"sync"
	"time"
)

type valueListener[T any] struct {
	id     uint64
	wanted T
	fn     func()
}

type Prop[T any] struct {
	mu             sync.RWMutex
	value          T
	eq             func(T, T) bool
	listeners      map[uint64]func(T)
	valueListeners []valueListener[T]
}

// NewProp creates a Prop for comparable types, using == for equality.
func NewProp[T comparable](initial T) *Prop[T] {
	return &Prop[T]{
		value:     initial,
		eq:        func(a, b T) bool { return a == b },
		listeners: make(map[uint64]func(T)),
	}
}

// NewPropEq creates a Prop for any type with a custom equality function.
func NewPropEq[T any](initial T, eq func(T, T) bool) *Prop[T] {
	return &Prop[T]{
		value:     initial,
		eq:        eq,
		listeners: make(map[uint64]func(T)),
	}
}

func (p *Prop[T]) Get() T {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.value
}

func (p *Prop[T]) Set(v T) {
	p.mu.Lock()
	if p.eq(p.value, v) {
		p.mu.Unlock()
		return
	}
	p.value = v

	var valueFns []func()
	for _, vl := range p.valueListeners {
		if p.eq(vl.wanted, v) {
			valueFns = append(valueFns, vl.fn)
		}
	}
	allFns := make([]func(T), 0, len(p.listeners))
	for _, fn := range p.listeners {
		allFns = append(allFns, fn)
	}
	p.mu.Unlock()

	for _, fn := range valueFns {
		fn()
	}
	for _, fn := range allFns {
		fn(v)
	}
}

func (p *Prop[T]) OnChange(fn func(T)) func() {
	id := nextID()
	p.mu.Lock()
	p.listeners[id] = fn
	p.mu.Unlock()
	return func() {
		p.mu.Lock()
		delete(p.listeners, id)
		p.mu.Unlock()
	}
}

// OnValue registers fn to be called whenever the value becomes equal to wanted.
func (p *Prop[T]) OnValue(wanted T, fn func()) func() {
	id := nextID()
	p.mu.Lock()
	p.valueListeners = append(p.valueListeners, valueListener[T]{id: id, wanted: wanted, fn: fn})
	p.mu.Unlock()
	return func() {
		p.mu.Lock()
		for i, vl := range p.valueListeners {
			if vl.id == id {
				p.valueListeners = append(p.valueListeners[:i], p.valueListeners[i+1:]...)
				break
			}
		}
		p.mu.Unlock()
	}
}

// Poll starts a goroutine that calls pollFn on the given interval and sets the result.
// Stops when ctx is cancelled.
func (p *Prop[T]) Poll(ctx context.Context, every time.Duration, pollFn func() T) {
	go func() {
		t := time.NewTicker(every)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				p.Set(pollFn())
			}
		}
	}()
}
