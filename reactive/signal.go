package reactive

import (
	"sync"
	"sync/atomic"
)

var idCounter uint64

func nextID() uint64 {
	return atomic.AddUint64(&idCounter, 1)
}

type Signal[T any] struct {
	mu        sync.RWMutex
	listeners map[uint64]func(T)
}

func NewSignal[T any]() *Signal[T] {
	return &Signal[T]{listeners: make(map[uint64]func(T))}
}

func (s *Signal[T]) Emit(v T) {
	s.mu.RLock()
	fns := make([]func(T), 0, len(s.listeners))
	for _, fn := range s.listeners {
		fns = append(fns, fn)
	}
	s.mu.RUnlock()
	for _, fn := range fns {
		fn(v)
	}
}

func (s *Signal[T]) Subscribe(fn func(T)) func() {
	id := nextID()
	s.mu.Lock()
	s.listeners[id] = fn
	s.mu.Unlock()
	return func() {
		s.mu.Lock()
		delete(s.listeners, id)
		s.mu.Unlock()
	}
}

// VoidSignal is a Signal that carries no value.
type VoidSignal struct {
	inner Signal[struct{}]
}

func NewVoidSignal() *VoidSignal {
	return &VoidSignal{inner: Signal[struct{}]{listeners: make(map[uint64]func(struct{}))}}
}

func (s *VoidSignal) Emit() {
	s.inner.Emit(struct{}{})
}

func (s *VoidSignal) Subscribe(fn func()) func() {
	return s.inner.Subscribe(func(struct{}) { fn() })
}
