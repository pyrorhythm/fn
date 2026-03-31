package fnopt

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"time"

	"golang.org/x/exp/constraints"
)

var (
	_ sql.Scanner   = (*Of[int])(nil)
	_ driver.Valuer = (*Of[int])(nil)
	_ driver.Valuer = Of[int]{}
)

var (
	ErrScanOverflow         = errors.New("value overflows target type")
	ErrScanNegativeUnsigned = errors.New("cannot scan negative value into unsigned type")
	ErrScanType             = errors.New("incompatible types")
	ErrScanUnsupported      = errors.New("unsupported target type")
)

// ScanError describes an error that occurred while scanning a SQL value.
type ScanError struct {
	Src    any
	Target string
	Cause  error
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("cannot scan %T into %s: %v", e.Src, e.Target, e.Cause)
}

func (e *ScanError) Unwrap() error { return e.Cause }

func scanErr[T any](src any, cause error) error {
	var zero T
	return &ScanError{Src: src, Target: fmt.Sprintf("%T", zero), Cause: cause}
}

// Scan implements [sql.Scanner] for Of[T].
func (o *Of[T]) Scan(src any) error {
	var (
		val T
		ok  bool
	)

	if src == nil {
		*o = Nil[T]()
		return nil
	}

	if scanner, ok := any(&val).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}
		*o = some(val)
		return nil
	}

	if val, ok = src.(T); ok {
		*o = some(val)
		return nil
	}

	result, err := convertTo[T](src)
	if err != nil {
		*o = Nil[T]()
		return err
	}
	*o = result
	return nil
}

// Value implements [driver.Valuer] for Of[T].
func (o Of[T]) Value() (driver.Value, error) {
	if !o.v {
		return nil, nil
	}

	if valuer, ok := any(o.t).(driver.Valuer); ok {
		return valuer.Value()
	}

	switch val := any(o.t).(type) {
	case int64, float64, bool, []byte, string, time.Time:
		return val, nil
	default:
		if converted, err := valueOf(val); err == nil {
			return converted, nil
		}
	}

	return nil, fmt.Errorf("Option.Value: unsupported type %T", o.t)
}

func convertTo[T any](src any) (Of[T], error) {
	var result T
	var err error

	switch tz := any(&result).(type) {
	case *int:
		err = likeSigned(src, tz)
	case *int8:
		err = likeSigned(src, tz)
	case *int16:
		err = likeSigned(src, tz)
	case *int32:
		err = likeSigned(src, tz)
	case *int64:
		err = likeSigned(src, tz)
	case *uint:
		err = likeUnsigned(src, tz)
	case *uint8:
		err = likeUnsigned(src, tz)
	case *uint16:
		err = likeUnsigned(src, tz)
	case *uint32:
		err = likeUnsigned(src, tz)
	case *uint64:
		err = likeUnsigned(src, tz)
	case *uintptr:
		err = likeUnsigned(src, tz)
	case *float32:
		err = likeFloat(src, tz)
	case *float64:
		err = likeFloat(src, tz)
	case *string:
		err = likeString(src, tz)
	case *[]byte:
		err = likeByteSlice(src, tz)
	case *[]rune:
		err = likeRuneSlice(src, tz)
	case *bool:
		err = likeBool(src, tz)
	case *time.Time:
		err = likeTime(src, tz)
	case *time.Duration:
		err = likeDuration(src, tz)
	default:
		return Nil[T](), scanErr[T](src, ErrScanUnsupported)
	}

	if err != nil {
		return Nil[T](), err
	}
	return some(result), nil
}

func likeSigned[S constraints.Signed](src any, dst *S) error {
	serr := scanErr[S]
	ovfErr := serr(src, ErrScanOverflow)
	i, ok := src.(int64)
	if !ok {
		return serr(src, ErrScanType)
	}
	switch p := any(dst).(type) {
	case *int:
		if i < math.MinInt || i > math.MaxInt {
			return ovfErr
		}
		*p = int(i)
	case *int8:
		if i < math.MinInt8 || i > math.MaxInt8 {
			return ovfErr
		}
		*p = int8(i)
	case *int16:
		if i < math.MinInt16 || i > math.MaxInt16 {
			return ovfErr
		}
		*p = int16(i)
	case *int32:
		if i < math.MinInt32 || i > math.MaxInt32 {
			return ovfErr
		}
		*p = int32(i)
	case *int64:
		*p = i
	}
	return nil
}

func likeUnsigned[U constraints.Unsigned](src any, dst *U) error {
	serr := scanErr[U]
	ovfErr := serr(src, ErrScanOverflow)
	i, ok := src.(int64)
	if !ok {
		return serr(src, ErrScanType)
	} else if i < 0 {
		return serr(src, ErrScanNegativeUnsigned)
	}
	switch p := any(dst).(type) {
	case *uint:
		if uint64(i) > math.MaxUint {
			return ovfErr
		}
		*p = uint(i)
	case *uint8:
		if i > math.MaxUint8 {
			return ovfErr
		}
		*p = uint8(i)
	case *uint16:
		if i > math.MaxUint16 {
			return ovfErr
		}
		*p = uint16(i)
	case *uint32:
		if i > math.MaxUint32 {
			return ovfErr
		}
		*p = uint32(i)
	case *uint64:
		*p = uint64(i)
	case *uintptr:
		if uint64(i) > math.MaxUint {
			return ovfErr
		}
		*p = uintptr(i)
	}
	return nil
}

func likeFloat[F constraints.Float](src any, dst *F) error {
	switch p := any(dst).(type) {
	case *float32:
		serr := scanErr[float32]
		switch v := src.(type) {
		case float64:
			if v < -math.MaxFloat32 || v > math.MaxFloat32 {
				return serr(src, ErrScanOverflow)
			}
			*p = float32(v)
		case int64:
			*p = float32(v)
		default:
			return serr(src, ErrScanType)
		}
	case *float64:
		switch v := src.(type) {
		case float64:
			*p = v
		case int64:
			*p = float64(v)
		default:
			return scanErr[float64](src, ErrScanType)
		}
	}
	return nil
}

func likeString(src any, dst *string) error {
	switch v := src.(type) {
	case string:
		*dst = v
	case []byte:
		*dst = string(v)
	default:
		return scanErr[string](src, ErrScanType)
	}
	return nil
}

func likeByteSlice(src any, dst *[]byte) error {
	switch v := src.(type) {
	case []byte:
		*dst = v
	case string:
		*dst = []byte(v)
	default:
		return scanErr[[]byte](src, ErrScanType)
	}
	return nil
}

func likeRuneSlice(src any, dst *[]rune) error {
	switch v := src.(type) {
	case string:
		*dst = []rune(v)
	case []byte:
		*dst = []rune(string(v))
	default:
		return scanErr[[]rune](src, ErrScanType)
	}
	return nil
}

func likeBool(src any, dst *bool) error {
	switch v := src.(type) {
	case bool:
		*dst = v
	case int64:
		*dst = v != 0
	default:
		return scanErr[bool](src, ErrScanType)
	}
	return nil
}

func likeDuration(src any, dst *time.Duration) error {
	serr := scanErr[time.Duration]
	switch v := src.(type) {
	case int64:
		*dst = time.Duration(v)
	case float64:
		*dst = time.Duration(v * float64(time.Second))
	case string:
		d, err := time.ParseDuration(v)
		if err != nil {
			return serr(src, ErrScanType)
		}
		*dst = d
	case []byte:
		d, err := time.ParseDuration(string(v))
		if err != nil {
			return serr(src, ErrScanType)
		}
		*dst = d
	default:
		return serr(src, ErrScanType)
	}
	return nil
}

func likeTime(src any, dst *time.Time) error {
	switch v := src.(type) {
	case time.Time:
		*dst = v
	case string:
		return parseTimeString(v, src, dst)
	case []byte:
		return parseTimeString(string(v), src, dst)
	case int64:
		*dst = time.Unix(v, 0)
	default:
		return scanErr[time.Time](src, ErrScanType)
	}
	return nil
}

func parseTimeString(s string, src any, dst *time.Time) error {
	for _, layout := range []string{
		time.RFC3339,
		time.DateTime,
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999999",
		"2006-01-02",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			*dst = t
			return nil
		}
	}
	return scanErr[time.Time](src, ErrScanType)
}

func valueOf(val any) (driver.Value, error) {
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return nil, fmt.Errorf("uint64 %d overflows int64", v)
		}
		return int64(v), nil
	case uintptr:
		if uint64(v) > math.MaxInt64 {
			return nil, fmt.Errorf("uintptr %d overflows int64", v)
		}
		return int64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	case time.Duration:
		return int64(v), nil
	case []rune:
		return string(v), nil
	}
	return nil, fmt.Errorf("unsupported type %T", val)
}
