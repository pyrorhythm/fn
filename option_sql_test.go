package fn

import (
	"database/sql/driver"
	"errors"
	"math"
	"testing"
	"time"
)

func TestOption_Scan_Nil(t *testing.T) {
	var opt Option[int]
	if err := opt.Scan(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opt.Valid() {
		t.Error("scanning nil should produce invalid Option")
	}
}

func TestOption_Scan_NilString(t *testing.T) {
	var opt Option[string]
	if err := opt.Scan(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opt.Valid() {
		t.Error("scanning nil should produce invalid Option")
	}
}

func TestOption_Scan_DirectMatch_Int64(t *testing.T) {
	var opt Option[int64]
	if err := opt.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 42 {
		t.Errorf("expected valid Option with value 42, got %v", opt)
	}
}

func TestOption_Scan_DirectMatch_String(t *testing.T) {
	var opt Option[string]
	if err := opt.Scan("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != "hello" {
		t.Errorf("expected valid Option with value 'hello', got %v", opt)
	}
}

func TestOption_Scan_DirectMatch_Bool(t *testing.T) {
	var opt Option[bool]
	if err := opt.Scan(true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || !opt.Val() {
		t.Errorf("expected valid Option with value true, got %v", opt)
	}
}

func TestOption_Scan_DirectMatch_Float64(t *testing.T) {
	var opt Option[float64]
	if err := opt.Scan(3.14); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 3.14 {
		t.Errorf("expected valid Option with value 3.14, got %v", opt)
	}
}

func TestOption_Scan_DirectMatch_Time(t *testing.T) {
	now := time.Now()
	var opt Option[time.Time]
	if err := opt.Scan(now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || !opt.Val().Equal(now) {
		t.Errorf("expected valid Option with time %v, got %v", now, opt)
	}
}

func TestOption_Scan_DirectMatch_Bytes(t *testing.T) {
	data := []byte("hello")
	var opt Option[[]byte]
	if err := opt.Scan(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || string(opt.Val()) != "hello" {
		t.Errorf("expected valid Option with value 'hello', got %v", opt)
	}
}

func TestOption_Scan_Int64ToInt(t *testing.T) {
	var opt Option[int]
	if err := opt.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 42 {
		t.Errorf("expected 42, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToInt8(t *testing.T) {
	var opt Option[int8]
	if err := opt.Scan(int64(127)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 127 {
		t.Errorf("expected 127, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToInt8_Negative(t *testing.T) {
	var opt Option[int8]
	if err := opt.Scan(int64(-128)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != -128 {
		t.Errorf("expected -128, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToInt16(t *testing.T) {
	var opt Option[int16]
	if err := opt.Scan(int64(32767)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 32767 {
		t.Errorf("expected 32767, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToInt32(t *testing.T) {
	var opt Option[int32]
	if err := opt.Scan(int64(2147483647)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 2147483647 {
		t.Errorf("expected 2147483647, got %v", opt.Val())
	}
}

func TestOption_Scan_Int8_Overflow(t *testing.T) {
	var opt Option[int8]
	err := opt.Scan(int64(128))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
	if opt.Valid() {
		t.Error("Option should be invalid after overflow error")
	}
}

func TestOption_Scan_Int8_Underflow(t *testing.T) {
	var opt Option[int8]
	err := opt.Scan(int64(-129))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
}

func TestOption_Scan_Int16_Overflow(t *testing.T) {
	var opt Option[int16]
	err := opt.Scan(int64(32768))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
}

func TestOption_Scan_Int32_Overflow(t *testing.T) {
	var opt Option[int32]
	err := opt.Scan(int64(2147483648))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
}

func TestOption_Scan_Int64ToUint(t *testing.T) {
	var opt Option[uint]
	if err := opt.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 42 {
		t.Errorf("expected 42, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToUint8(t *testing.T) {
	var opt Option[uint8]
	if err := opt.Scan(int64(255)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 255 {
		t.Errorf("expected 255, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToUint16(t *testing.T) {
	var opt Option[uint16]
	if err := opt.Scan(int64(65535)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 65535 {
		t.Errorf("expected 65535, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToUint32(t *testing.T) {
	var opt Option[uint32]
	if err := opt.Scan(int64(4294967295)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 4294967295 {
		t.Errorf("expected 4294967295, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToUint64(t *testing.T) {
	var opt Option[uint64]
	if err := opt.Scan(int64(math.MaxInt64)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != math.MaxInt64 {
		t.Errorf("expected %d, got %v", int64(math.MaxInt64), opt.Val())
	}
}

func TestOption_Scan_Int64ToUintptr(t *testing.T) {
	var opt Option[uintptr]
	if err := opt.Scan(int64(12345)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 12345 {
		t.Errorf("expected 12345, got %v", opt.Val())
	}
}

func TestOption_Scan_Uint_Negative(t *testing.T) {
	var opt Option[uint]
	err := opt.Scan(int64(-1))
	if err == nil {
		t.Fatal("expected error for negative value")
	}
	if !errors.Is(err, ErrScanNegativeUnsigned) {
		t.Errorf("expected ErrScanNegativeUnsigned, got %v", err)
	}
}

func TestOption_Scan_Uint8_Negative(t *testing.T) {
	var opt Option[uint8]
	err := opt.Scan(int64(-1))
	if err == nil {
		t.Fatal("expected error for negative value")
	}
	if !errors.Is(err, ErrScanNegativeUnsigned) {
		t.Errorf("expected ErrScanNegativeUnsigned, got %v", err)
	}
}

func TestOption_Scan_Uint8_Overflow(t *testing.T) {
	var opt Option[uint8]
	err := opt.Scan(int64(256))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
}

func TestOption_Scan_Uint16_Overflow(t *testing.T) {
	var opt Option[uint16]
	err := opt.Scan(int64(65536))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
}

func TestOption_Scan_Uint32_Overflow(t *testing.T) {
	var opt Option[uint32]
	err := opt.Scan(int64(4294967296))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
}

func TestOption_Scan_Float64ToFloat32(t *testing.T) {
	var opt Option[float32]
	if err := opt.Scan(float64(3.14)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() < 3.13 || opt.Val() > 3.15 {
		t.Errorf("expected ~3.14, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToFloat32(t *testing.T) {
	var opt Option[float32]
	if err := opt.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 42.0 {
		t.Errorf("expected 42.0, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToFloat64(t *testing.T) {
	var opt Option[float64]
	if err := opt.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 42.0 {
		t.Errorf("expected 42.0, got %v", opt.Val())
	}
}

func TestOption_Scan_Float32_Overflow(t *testing.T) {
	var opt Option[float32]
	err := opt.Scan(float64(math.MaxFloat64))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
}

func TestOption_Scan_Float32_NegativeOverflow(t *testing.T) {
	var opt Option[float32]
	err := opt.Scan(float64(-math.MaxFloat64))
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrScanOverflow) {
		t.Errorf("expected ErrScanOverflow, got %v", err)
	}
}

func TestOption_Scan_BytesToString(t *testing.T) {
	var opt Option[string]
	if err := opt.Scan([]byte("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != "hello" {
		t.Errorf("expected 'hello', got %v", opt.Val())
	}
}

func TestOption_Scan_StringToBytes(t *testing.T) {
	var opt Option[[]byte]
	if err := opt.Scan("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || string(opt.Val()) != "hello" {
		t.Errorf("expected 'hello', got %v", opt.Val())
	}
}

func TestOption_Scan_StringToRunes(t *testing.T) {
	var opt Option[[]rune]
	if err := opt.Scan("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || string(opt.Val()) != "hello" {
		t.Errorf("expected 'hello', got %v", string(opt.Val()))
	}
}

func TestOption_Scan_BytesToRunes(t *testing.T) {
	var opt Option[[]rune]
	if err := opt.Scan([]byte("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || string(opt.Val()) != "hello" {
		t.Errorf("expected 'hello', got %v", string(opt.Val()))
	}
}

func TestOption_Scan_Int64ToBool_True(t *testing.T) {
	var opt Option[bool]
	if err := opt.Scan(int64(1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || !opt.Val() {
		t.Errorf("expected true, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToBool_False(t *testing.T) {
	var opt Option[bool]
	if err := opt.Scan(int64(0)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() {
		t.Errorf("expected false, got %v", opt.Val())
	}
}

func TestOption_Scan_Int64ToBool_NonZero(t *testing.T) {
	var opt Option[bool]
	if err := opt.Scan(int64(42)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || !opt.Val() {
		t.Errorf("expected true for non-zero, got %v", opt.Val())
	}
}

func TestOption_Scan_StringToTime_RFC3339(t *testing.T) {
	var opt Option[time.Time]
	if err := opt.Scan("2024-01-15T10:30:00Z"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() {
		t.Error("expected valid Option")
	}
	if opt.Val().Year() != 2024 || opt.Val().Month() != 1 || opt.Val().Day() != 15 {
		t.Errorf("unexpected time: %v", opt.Val())
	}
}

func TestOption_Scan_StringToTime_DateTime(t *testing.T) {
	var opt Option[time.Time]
	if err := opt.Scan("2024-01-15 10:30:00"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() {
		t.Error("expected valid Option")
	}
}

func TestOption_Scan_StringToTime_DateOnly(t *testing.T) {
	var opt Option[time.Time]
	if err := opt.Scan("2024-01-15"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() {
		t.Error("expected valid Option")
	}
}

func TestOption_Scan_BytesToTime(t *testing.T) {
	var opt Option[time.Time]
	if err := opt.Scan([]byte("2024-01-15T10:30:00Z")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() {
		t.Error("expected valid Option")
	}
}

func TestOption_Scan_Int64ToTime_Unix(t *testing.T) {
	var opt Option[time.Time]
	timestamp := int64(1705315800)
	if err := opt.Scan(timestamp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() {
		t.Error("expected valid Option")
	}
	if opt.Val().Unix() != timestamp {
		t.Errorf("expected unix %d, got %d", timestamp, opt.Val().Unix())
	}
}

func TestOption_Scan_StringToTime_Invalid(t *testing.T) {
	var opt Option[time.Time]
	err := opt.Scan("not a time")
	if err == nil {
		t.Fatal("expected error for invalid time string")
	}
	if !errors.Is(err, ErrScanType) {
		t.Errorf("expected ErrScanType, got %v", err)
	}
}

func TestOption_Scan_Int64ToDuration(t *testing.T) {
	var opt Option[time.Duration]
	if err := opt.Scan(int64(5000000000)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 5*time.Second {
		t.Errorf("expected 5s, got %v", opt.Val())
	}
}

func TestOption_Scan_Float64ToDuration(t *testing.T) {
	var opt Option[time.Duration]
	if err := opt.Scan(float64(2.5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 2500*time.Millisecond {
		t.Errorf("expected 2.5s, got %v", opt.Val())
	}
}

func TestOption_Scan_StringToDuration(t *testing.T) {
	var opt Option[time.Duration]
	if err := opt.Scan("1h30m"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 90*time.Minute {
		t.Errorf("expected 1h30m, got %v", opt.Val())
	}
}

func TestOption_Scan_BytesToDuration(t *testing.T) {
	var opt Option[time.Duration]
	if err := opt.Scan([]byte("500ms")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val() != 500*time.Millisecond {
		t.Errorf("expected 500ms, got %v", opt.Val())
	}
}

func TestOption_Scan_StringToDuration_Invalid(t *testing.T) {
	var opt Option[time.Duration]
	err := opt.Scan("not a duration")
	if err == nil {
		t.Fatal("expected error for invalid duration string")
	}
	if !errors.Is(err, ErrScanType) {
		t.Errorf("expected ErrScanType, got %v", err)
	}
}

func TestOption_Scan_TypeMismatch_StringToInt(t *testing.T) {
	var opt Option[int]
	err := opt.Scan("hello")
	if err == nil {
		t.Fatal("expected type error")
	}
	if !errors.Is(err, ErrScanType) {
		t.Errorf("expected ErrScanType, got %v", err)
	}
}

func TestOption_Scan_TypeMismatch_BoolToInt(t *testing.T) {
	var opt Option[int]
	err := opt.Scan(true)
	if err == nil {
		t.Fatal("expected type error")
	}
	if !errors.Is(err, ErrScanType) {
		t.Errorf("expected ErrScanType, got %v", err)
	}
}

func TestOption_Scan_TypeMismatch_IntToString(t *testing.T) {
	var opt Option[string]
	err := opt.Scan(int64(42))
	if err == nil {
		t.Fatal("expected type error")
	}
	if !errors.Is(err, ErrScanType) {
		t.Errorf("expected ErrScanType, got %v", err)
	}
}

type customScanner struct {
	value string
}

func (c *customScanner) Scan(src any) error {
	if s, ok := src.(string); ok {
		c.value = "scanned:" + s
		return nil
	}
	return errors.New("customScanner only accepts strings")
}

func TestOption_Scan_CustomScanner(t *testing.T) {
	var opt Option[customScanner]
	if err := opt.Scan("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opt.Valid() || opt.Val().value != "scanned:test" {
		t.Errorf("expected 'scanned:test', got %v", opt.Val().value)
	}
}

func TestOption_Scan_CustomScanner_Error(t *testing.T) {
	var opt Option[customScanner]
	err := opt.Scan(int64(42))
	if err == nil {
		t.Fatal("expected error from custom scanner")
	}
}

type unsupportedType struct {
	field int
}

func TestOption_Scan_UnsupportedType(t *testing.T) {
	var opt Option[unsupportedType]
	err := opt.Scan("anything")
	if err == nil {
		t.Fatal("expected unsupported type error")
	}
	if !errors.Is(err, ErrScanUnsupported) {
		t.Errorf("expected ErrScanUnsupported, got %v", err)
	}
}

func TestOption_Value_Invalid(t *testing.T) {
	opt := Nil[int]()
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != nil {
		t.Errorf("expected nil for invalid Option, got %v", val)
	}
}

func TestOption_Value_Int64(t *testing.T) {
	opt := SomeAny(int64(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestOption_Value_Float64(t *testing.T) {
	opt := SomeAny(float64(3.14))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != float64(3.14) {
		t.Errorf("expected 3.14, got %v", val)
	}
}

func TestOption_Value_Bool(t *testing.T) {
	opt := SomeAny(true)
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != true {
		t.Errorf("expected true, got %v", val)
	}
}

func TestOption_Value_String(t *testing.T) {
	opt := Some("hello")
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}
}

func TestOption_Value_Bytes(t *testing.T) {
	opt := SomeAny([]byte("hello"))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(val.([]byte)) != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}
}

func TestOption_Value_Time(t *testing.T) {
	now := time.Now()
	opt := SomeAny(now)
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !val.(time.Time).Equal(now) {
		t.Errorf("expected %v, got %v", now, val)
	}
}

func TestOption_Value_IntToInt64(t *testing.T) {
	opt := SomeAny(42)
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_Int8ToInt64(t *testing.T) {
	opt := SomeAny(int8(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_Int16ToInt64(t *testing.T) {
	opt := SomeAny(int16(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_Int32ToInt64(t *testing.T) {
	opt := SomeAny(int32(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_UintToInt64(t *testing.T) {
	opt := SomeAny(uint(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_Uint8ToInt64(t *testing.T) {
	opt := SomeAny(uint8(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_Uint16ToInt64(t *testing.T) {
	opt := SomeAny(uint16(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_Uint32ToInt64(t *testing.T) {
	opt := SomeAny(uint32(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_Uint64ToInt64(t *testing.T) {
	opt := SomeAny(uint64(42))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", val, val)
	}
}

func TestOption_Value_Uint64_Overflow(t *testing.T) {
	opt := SomeAny(uint64(math.MaxUint64))
	_, err := opt.Value()
	if err == nil {
		t.Fatal("expected overflow error for uint64 > MaxInt64")
	}
}

func TestOption_Value_Float32ToFloat64(t *testing.T) {
	opt := SomeAny(float32(3.14))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := val.(float64); !ok {
		t.Errorf("expected float64, got %T", val)
	}
}

func TestOption_Value_DurationToInt64(t *testing.T) {
	opt := SomeAny(5 * time.Second)
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != int64(5*time.Second) {
		t.Errorf("expected %d, got %v", 5*time.Second, val)
	}
}

func TestOption_Value_RunesToString(t *testing.T) {
	opt := SomeAny([]rune("hello"))
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}
}

type customValuer struct {
	value string
}

func (c customValuer) Value() (driver.Value, error) {
	return "custom:" + c.value, nil
}

func TestOption_Value_CustomValuer(t *testing.T) {
	opt := SomeAny(customValuer{value: "test"})
	val, err := opt.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "custom:test" {
		t.Errorf("expected 'custom:test', got %v", val)
	}
}

type customValuerError struct{}

func (c customValuerError) Value() (driver.Value, error) {
	return nil, errors.New("custom valuer error")
}

func TestOption_Value_CustomValuer_Error(t *testing.T) {
	opt := SomeAny(customValuerError{})
	_, err := opt.Value()
	if err == nil {
		t.Fatal("expected error from custom valuer")
	}
}

func TestOption_Value_UnsupportedType(t *testing.T) {
	opt := SomeAny(unsupportedType{field: 42})
	_, err := opt.Value()
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestScanError_Error(t *testing.T) {
	err := &ScanError{
		Src:    int64(256),
		Target: "int8",
		Cause:  ErrScanOverflow,
	}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
	// Should contain type info and cause
	if !errors.Is(err, ErrScanOverflow) {
		t.Error("ScanError should unwrap to ErrScanOverflow")
	}
}

func TestScanError_Unwrap(t *testing.T) {
	err := &ScanError{
		Src:    "test",
		Target: "int",
		Cause:  ErrScanType,
	}
	if !errors.Is(err, ErrScanType) {
		t.Error("errors.Is should work with ScanError")
	}
}
