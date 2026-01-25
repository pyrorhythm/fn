# fn

> [![Go Reference](https://pkg.go.dev/badge/github.com/pyrorhythm/fn.svg)](https://pkg.go.dev/github.com/pyrorhythm/fn)
> [![Go Report Card](https://goreportcard.com/badge/github.com/pyrorhythm/fn)](https://goreportcard.com/report/github.com/pyrorhythm/fn)
> [![Coverage Status](https://coveralls.io/repos/github/pyrorhythm/fn/badge.svg?branch=main)](https://coveralls.io/github/pyrorhythm/fn?branch=main)
> ##### (yet another) functional primitives for Go -> Option with sql.Scanner / driver.Valuer functionality, JSON-serializable Result, composition, transactions primitive and moooore...

## Install

```bash
go get github.com/pyrorhythm/fn
```

## Option

```go
opt := fn.Some(42)           // valid if non-zero
opt := fn.Some(0)            // invalid (zero value)
opt := fn.SomeAny(0)         // valid (bypasses zero check)
opt := fn.SomePtr(&value)    // from pointer, nil-safe
opt := fn.Nil[int]()         // explicit nil

if opt.Valid() {
    v := opt.Val()           // get value
    p := opt.Ptr()           // get pointer (nil if invalid)
}
```

## Result

```go
r := fn.OK(42)                        // from value
r := fn.OKPtr(&value)                 // from pointer
r := fn.Err[int](err)                 // from error
r := fn.Errn[int]("failed")           // from string
r := fn.From(os.Open("file"))         // wrap (T, error)

if r.OK() {
    v := r.Val()
} else {
    e := r.Err()
}

val, err := r.Unpack()                // back to (T, error)
err = r.Into(&dest)                   // assign + return error
```

## Map / FlatMap

```go
// Option
opt2 := fn.OptTo(opt, strconv.Itoa)              // map
opt3 := fn.OptMorph(opt, func(i int) Option[string] { ... })  // flatmap

// Result
r2 := fn.To(r, strconv.Itoa)                     // map (propagates error)
r3 := fn.Morph(r, func(i int) Result[string] { ... })  // flatmap
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
