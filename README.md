# fn 

[![Go Reference](https://pkg.go.dev/badge/github.com/pyrorhythm/fn.svg)](https://pkg.go.dev/github.com/pyrorhythm/fn)
[![Go Report Card](https://goreportcard.com/badge/github.com/pyrorhythm/fn)](https://goreportcard.com/report/github.com/pyrorhythm/fn)
[![Coverage Status](https://coveralls.io/repos/github/pyrorhythm/fn/badge.svg?branch=main)](https://coveralls.io/github/pyrorhythm/fn?branch=main)

Functional programming primitives for Go. Option, Result, composition, transactional operations.

## Install

```bash
go get github.com/pyrorhythm/fn
```

## Option

```go
opt := fn.Some(42)        // Option[int], valid
opt := fn.Some(0)         // not valid (zero value)
opt := fn.SomeP(&value)   // from pointer, nil-safe
opt := fn.Nil[int]()      // explicit nil

if opt.Valid() {
    fmt.Println(opt.Value())
}
```

## Result

```go
r := fn.OK(42)
r := fn.Err[int](errors.New("failed"))
r := fn.From(someFunc())  // wraps (T, error)

if r.OK() {
    fmt.Println(r.Val())
} else {
    fmt.Println(r.Err())
}

// unpack back to (T, error)
val, err := r.Unpack()
```

## Map / FlatMap

```go
r := fn.OK(10)

// Map: transform value
r2 := fn.To(r, func(i int) string {
    return strconv.Itoa(i)
})

// FlatMap: chain operations that return Result
r3 := fn.Morph(r, func(i int) fn.Result[int] {
    if i > 5 {
        return fn.OK(i * 2)
    }
    return fn.Errn[int]("too small")
})
```

## Transactional Operations

```go
// RunOps: run sequentially, stop on error
err := fn.RunOps(op1, op2, op3)

// TransactOps: run with automatic rollback on failure
err := fn.TransactOps(op1, op2, op3)
// if op2 fails: rollback op2, then op1
```

Operations implement `Run() error` and `Rollback() error`:

```go
type myOp struct{}

func (o *myOp) Run() error      { /* do work */ }
func (o *myOp) Rollback() error { /* undo work */ }
```

Or use `ErrFunc` for simple cases:

```go
fn.RunOps(
    fn.FuncErr(func() error { return step1() }),
    fn.FuncErr(func() error { return step2() }),
)
```

## Helpers

```go
// ternary
val := fn.If(cond, "yes", "no")

// nil-safe pointer unwrap
val := fn.Or(ptr, defaultValue)
val := fn.OrZero(ptr)

// zero value check
if fn.Valid(someValue) { ... }
```

## License

See [LICENSE](./LICENSE)
