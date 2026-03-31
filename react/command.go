package react

import (
	"sync"
)

type Command struct {
	mu          sync.RWMutex
	execute     func()
	canExec     func() bool
	nextID      uint64
	listeners   map[uint64]func(bool)
	lastCanExec *bool
}

func NewCommand(execute func(), canExecute func() bool) *Command {
	return &Command{
		execute:   execute,
		canExec:   canExecute,
		listeners: make(map[uint64]func(bool)),
	}
}

func (c *Command) Execute() {
	if !c.CanExecute() {
		return
	}
	c.execute()
}

func (c *Command) CanExecute() bool {
	if c.canExec == nil {
		return true
	}
	cur := c.canExec()
	if c.lastCanExec != nil && *c.lastCanExec == cur {
		return cur
	}
	*c.lastCanExec = cur
	defer c.notify(cur)
	return cur
}

func (c *Command) Refresh() {
	c.CanExecute()
}

func (c *Command) OnCanExecuteChanged(fn func(canExec bool)) func() {
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

func (c *Command) notify(cur bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, fn := range c.listeners {
		fn(cur)
	}
}
