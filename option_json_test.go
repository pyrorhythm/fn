package fn

import (
	"testing"

	"github.com/bytedance/sonic"
)

func TestOption_MarshalJSON_Valid(t *testing.T) {
	opt := Some(42)
	ba, err := sonic.Marshal(opt)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := "42"
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestOption_MarshalJSON_Invalid(t *testing.T) {
	opt := Nil[int]()
	ba, err := sonic.Marshal(opt)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := "null"
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestOption_MarshalJSON_String(t *testing.T) {
	opt := Some("hello")
	ba, err := sonic.Marshal(opt)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `"hello"`
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestOption_MarshalJSON_Struct(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	opt := SomeAny(person{Name: "Alice", Age: 30})
	ba, err := sonic.Marshal(opt)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"name":"Alice","age":30}`
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestOption_UnmarshalJSON_Valid(t *testing.T) {
	data := "42"
	var opt Option[int]
	err := sonic.Unmarshal([]byte(data), &opt)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !opt.Valid() {
		t.Error("expected valid option")
	}
	if opt.Value() != 42 {
		t.Errorf("expected 42, got %d", opt.Value())
	}
}

func TestOption_UnmarshalJSON_Null(t *testing.T) {
	data := "null"
	var opt Option[int]
	err := sonic.Unmarshal([]byte(data), &opt)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if opt.Valid() {
		t.Error("expected invalid option for null")
	}
}

func TestOption_UnmarshalJSON_String(t *testing.T) {
	data := `"hello"`
	var opt Option[string]
	err := sonic.Unmarshal([]byte(data), &opt)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !opt.Valid() {
		t.Error("expected valid option")
	}
	if opt.Value() != "hello" {
		t.Errorf("expected 'hello', got %q", opt.Value())
	}
}

func TestOption_UnmarshalJSON_Struct(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	data := `{"name":"Alice","age":30}`
	var opt Option[person]
	err := sonic.Unmarshal([]byte(data), &opt)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !opt.Valid() {
		t.Error("expected valid option")
	}
	if opt.Value().Name != "Alice" || opt.Value().Age != 30 {
		t.Errorf("expected {Alice, 30}, got %+v", opt.Value())
	}
}

func TestOption_UnmarshalJSON_InvalidJSON(t *testing.T) {
	data := `{invalid json}`
	var opt Option[int]
	err := sonic.Unmarshal([]byte(data), &opt)
	if err == nil {
		t.Error("expected unmarshal error")
	}
	if opt.Valid() {
		t.Error("option should be invalid after unmarshal error")
	}
}

func TestOption_JSON_Roundtrip_Valid(t *testing.T) {
	original := Some(42)
	ba, err := sonic.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var restored Option[int]
	err = sonic.Unmarshal(ba, &restored)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if original.Valid() != restored.Valid() {
		t.Errorf("Valid mismatch: %v vs %v", original.Valid(), restored.Valid())
	}
	if original.Value() != restored.Value() {
		t.Errorf("Value mismatch: %v vs %v", original.Value(), restored.Value())
	}
}

func TestOption_JSON_Roundtrip_Invalid(t *testing.T) {
	original := Nil[int]()
	ba, err := sonic.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var restored Option[int]
	err = sonic.Unmarshal(ba, &restored)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if original.Valid() != restored.Valid() {
		t.Errorf("Valid mismatch: %v vs %v", original.Valid(), restored.Valid())
	}
}
