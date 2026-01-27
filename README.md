# fn

> [![Go Reference](https://pkg.go.dev/badge/github.com/pyrorhythm/fn.svg)](https://pkg.go.dev/github.com/pyrorhythm/fn)
> [![Go Report Card](https://goreportcard.com/badge/github.com/pyrorhythm/fn)](https://goreportcard.com/report/github.com/pyrorhythm/fn)
> [![Coverage Status](https://coveralls.io/repos/github/pyrorhythm/fn/badge.svg?branch=main)](https://coveralls.io/github/pyrorhythm/fn?branch=main)
> ##### (yet another) functional primitives for Go -> Option with sql.Scanner / driver.Valuer functionality, JSON-serializable Result, composition, transactions primitive and moooore...

## Install

```bash
go get github.com/pyrorhythm/fn
```

## Packages

```go
import "github.com/pyrorhythm/fn/option"  // option.Of[T]
import "github.com/pyrorhythm/fn/result"  // result.Of[T]
import "github.com/pyrorhythm/fn/pipe"    // pipe.Two, pipe.Three...
```

## Option

```go
import "github.com/pyrorhythm/fn/option"

var o option.Of[int]

o = option.Some(42)           // valid if non-zero
o = option.Some(0)            // invalid (zero value)
o = option.SomeAny(0)         // valid (bypasses zero check)
o = option.SomePtr(&value)    // from pointer, nil-safe
o = option.Nil[int]()         // explicit nil

if o.Valid() {
    v := o.Val()              // get value
    p := o.Ptr()              // get pointer (nil if invalid)
}

// FlatMap
o2 := o.FlatMap(func(n int) option.Of[int] { return option.SomeAny(n * 2) })

// Fold
v := o.Fold(func() int { return 0 }, func(n int) int { return n * 2 })
```

## Result

```go
import "github.com/pyrorhythm/fn/result"

var r result.Of[int]

r = result.OK(42)                        // from value
r = result.OKPtr(&value)                 // from pointer
r = result.OKAny(0)                      // bypasses zero check
r = result.Err[int](err)                 // from error
r = result.Errn[int]("failed")           // from string
r = result.From(os.Open("file"))         // wrap (T, error)

if r.OK() {
    v := r.Val()
} else {
    e := r.Err()
}

val, err := r.Unpack()                   // back to (T, error)
err = r.Into(&dest)                      // assign + return error

// FlatMap
r2 := r.FlatMap(func(n int) result.Of[int] { return result.OKAny(n * 2) })

// Fold
v := r.Fold(func() int { return -1 }, func(n int) int { return n })
```

## Pipe

```go
import "github.com/pyrorhythm/fn/pipe"

// Pure composition
out := pipe.Three(input,
    parseJSON,
    validate,
)

// Error tuple composition (short-circuits on error)
out, err := pipe.ThreeErr(input, nil,
    func(s string) (int, error) { return strconv.Atoi(s) },
    func(n int) (float64, error) { return float64(n) * 1.5, nil },
)

// Result monad composition (short-circuits on Exc)
out := pipe.ThreeRes(
    result.From(fetchData()),
    transform,
    save,
)
```

## Generic Functions

Works with both `Option` and `Result`:

```go
import "github.com/pyrorhythm/fn"

// Extract value or fallback
v := fn.OrElse(o, 0)
v := fn.OrElse(r, "default")

// Extract or panic
v := fn.Must(o)
v := fn.Must(r)

// Pattern match with type change
s := fn.Fold(o, func() string { return "empty" }, func(n int) string { return strconv.Itoa(n) })

// Check validity
ok := fn.IsValid(o)
ok := fn.IsValid(r)

// Monadic chain (FlatMap as free function)
o2 := fn.Chain(o, func(n int) option.Of[int] { return option.SomeAny(n * 2) })
r2 := fn.Chain(r, func(n int) result.Of[int] { return result.OKAny(n * 2) })

// Sequence (ignore current value)
o2 := fn.AndThen(o, option.Some(42))
r2 := fn.AndThen(r, result.OK("next"))
```

## Interfaces

```go
// Container - value extraction
type Container[T any] interface {
    Val() T
    Valid() bool
}

// Monad - chainable container
type Monad[T, Self any] interface {
    Container[T]
    FlatMap(func(T) Self) Self
}
```

Both `Option` and `Result` implement these interfaces.

## Map / Morph (type-changing)

```go
// Option
opt2 := fn.OptTo(opt, strconv.Itoa)              // map T -> U
opt3 := fn.OptMorph(opt, func(i int) fn.Option[string] { ... })  // flatmap

// Result
r2 := fn.To(r, strconv.Itoa)                     // map (propagates error)
r3 := fn.Morph(r, func(i int) fn.Result[string] { ... })  // flatmap
```

## Conditionals

```go
v := fn.If(cond, "yes", "no")                    // ternary
o := fn.FlatIf(cond, valA, valB)                 // returns Option
r := fn.ErrIf(cond, value, err)                  // returns Result
```

## Helpers

```go
v := fn.Or(ptr, defaultVal)      // nil-safe pointer unwrap
v := fn.OrZero(ptr)              // unwrap or zero value
ok := fn.Valid(value)            // non-zero check
ok := fn.Is[T](v)                // type assertion check
v := fn.Cast[T](v)               // safe cast (zero on fail)
z := fn.Z[T]()                   // zero value
```

## Transactions

```go
// Sequential execution, stops on first error
err := fn.RunOps(op1, op2, op3)

// With automatic rollback on failure
err := fn.TransactOps(op1, op2, op3)

// Simple func adapter
fn.RunOps(
    fn.FuncErr(step1),
    fn.FuncErr(step2),
)
```

Operations implement `Run() error` and `Rollback() error`.

## JSON

Option and Result implement `json.Marshaler` / `json.Unmarshaler` (via sonic):

```go
type User struct {
    Name  string           `json:"name"`
    Email fn.Option[string] `json:"email"`  // null if invalid
}

// Result marshals error as {"_ERROR": "message"}
```

## SQL

Option implements `sql.Scanner` and `driver.Valuer`:

```go
var email fn.Option[string]
row.Scan(&email)  // NULL -> Nil[string]()

db.Exec("INSERT ...", fn.Some("test@example.com"))  // Some -> value, Nil -> NULL
```

Supported conversions:
- `int64` → `int`, `int8`, `int16`, `int32`, `int64` (with overflow check)
- `int64` → `uint`, `uint8`, `uint16`, `uint32`, `uint64` (with negative/overflow check)
- `float64` → `float32`, `float64` (with overflow check)
- `string` ↔ `[]byte` ↔ `[]rune`
- `int64` → `bool` (0 = false)
- `string`/`[]byte`/`int64` → `time.Time` (RFC3339, DateTime, unix)
- `int64`/`float64`/`string` → `time.Duration`

## License

MIT
