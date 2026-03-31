package react

import (
	"context"
	"sync"
)

type ContextCommand struct {
	execute     func(ctx context.Context)
	canExec     func(ctx context.Context) bool
	nextID      uint64
	listeners   map[uint64]func(bool)
	lastCanExec *bool
	mu          sync.RWMutex
}

func NewContextCommand(execute func(ctx context.Context), canExecute func(ctx context.Context) bool) *ContextCommand {
	return &ContextCommand{
		execute:   execute,
		canExec:   canExecute,
		listeners: make(map[uint64]func(bool)),
	}

}

func (c *ContextCommand) Execute(ctx context.Context) {
	if !c.CanExecute(ctx) {
		return
	}
	c.execute(ctx)
}

func (c *ContextCommand) CanExecute(ctx context.Context) bool {
	if c.canExec == nil {
		return true
	}
	cur := c.canExec(ctx)
	if c.lastCanExec != nil && *c.lastCanExec == cur {
		return cur
	}
	*c.lastCanExec = cur
	defer c.notify(cur)
	return cur
}

func (c *ContextCommand) Refresh(ctx context.Context) {
	c.CanExecute(ctx)
}

func (c *ContextCommand) OnCanExecuteChanged(fn func(canExec bool)) func() {
	id := nextID()

	c.mu.Lock()
	defer c.mu.Unlock()
	c.listeners[id] = fn

	return func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.listeners, id)
	}
}

func (c *ContextCommand) notify(cur bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, fn := range c.listeners {
		fn(cur)
	}
}
