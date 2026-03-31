package react

import "sync"

type EnumCommand[T comparable] struct {
	mu          sync.RWMutex
	executables map[T]*Command
}

func NewEnumCommand[T comparable]() *EnumCommand[T] {
	return &EnumCommand[T]{executables: make(map[T]*Command)}
}

func (c *EnumCommand[T]) Register(on T, cmd *Command) func() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executables[on] = cmd

	return func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.executables, on)
	}
}

func (c *EnumCommand[T]) Execute(val T) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if cmd, ok := c.executables[val]; ok {
		cmd.Execute()
	}
}

func (c *EnumCommand[T]) On(val T) *Command {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.executables[val]
}
