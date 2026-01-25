package fn

import (
	"errors"
	"testing"

	"github.com/bytedance/sonic"
)

func TestResult_MarshalJSON_OK(t *testing.T) {
	r := OK(42)
	ba, err := sonic.Marshal(r)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"result_value":42,"result_value_valid":true}`
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestResult_MarshalJSON_Err(t *testing.T) {
	r := Err[int](errors.New("something went wrong"))
	ba, err := sonic.Marshal(r)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"result_value_valid":false,"result_error":"something went wrong"}`
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestResult_MarshalJSON_InvalidOption(t *testing.T) {
	r := OK(0) // zero value, option invalid
	ba, err := sonic.Marshal(r)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"result_value_valid":false}`
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestResult_MarshalJSON_String(t *testing.T) {
	r := OK("hello")
	ba, err := sonic.Marshal(r)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"result_value":"hello","result_value_valid":true}`
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestResult_MarshalJSON_Struct(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	r := OKAny(person{Name: "Alice", Age: 30})
	ba, err := sonic.Marshal(r)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"result_value":{"name":"Alice","age":30},"result_value_valid":true}`
	if string(ba) != expected {
		t.Errorf("expected %s, got %s", expected, string(ba))
	}
}

func TestResult_UnmarshalJSON_OK(t *testing.T) {
	data := `{"result_value":42,"result_value_valid":true}`
	var r Result[int]
	err := sonic.Unmarshal([]byte(data), &r)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !r.OK() {
		t.Error("expected OK result")
	}
	if r.Val() != 42 {
		t.Errorf("expected 42, got %d", r.Val())
	}
}

func TestResult_UnmarshalJSON_Err(t *testing.T) {
	data := `{"result_value_valid":false,"result_error":"something went wrong"}`
	var r Result[int]
	err := sonic.Unmarshal([]byte(data), &r)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if r.OK() {
		t.Error("expected error result")
	}
	if r.Err() == nil {
		t.Error("expected error to be set")
	}
	if r.Err().Error() != "something went wrong" {
		t.Errorf("expected 'something went wrong', got %q", r.Err().Error())
	}
}

func TestResult_UnmarshalJSON_InvalidOption(t *testing.T) {
	data := `{"result_value_valid":false}`
	var r Result[int]
	err := sonic.Unmarshal([]byte(data), &r)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if r.OK() {
		t.Error("expected invalid option")
	}
	if r.Err() != nil {
		t.Error("should not have error, just invalid option")
	}
}

func TestResult_UnmarshalJSON_String(t *testing.T) {
	data := `{"result_value":"hello","result_value_valid":true}`
	var r Result[string]
	err := sonic.Unmarshal([]byte(data), &r)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !r.OK() {
		t.Error("expected OK result")
	}
	if r.Val() != "hello" {
		t.Errorf("expected 'hello', got %q", r.Val())
	}
}

func TestResult_UnmarshalJSON_Struct(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	data := `{"result_value":{"name":"Alice","age":30},"result_value_valid":true}`
	var r Result[person]
	err := sonic.Unmarshal([]byte(data), &r)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if !r.OK() {
		t.Error("expected OK result")
	}
	if r.Val().Name != "Alice" || r.Val().Age != 30 {
		t.Errorf("expected {Alice, 30}, got %+v", r.Val())
	}
}

func TestResult_UnmarshalJSON_InvalidJSON(t *testing.T) {
	data := `{invalid json}`
	var r Result[int]
	err := sonic.Unmarshal([]byte(data), &r)
	if err == nil {
		t.Error("expected unmarshal error")
	}
}

func TestResult_JSON_Roundtrip_OK(t *testing.T) {
	original := OK(42)
	ba, err := sonic.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var restored Result[int]
	err = sonic.Unmarshal(ba, &restored)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if original.OK() != restored.OK() {
		t.Errorf("OK mismatch: %v vs %v", original.OK(), restored.OK())
	}
	if original.Val() != restored.Val() {
		t.Errorf("Val mismatch: %v vs %v", original.Val(), restored.Val())
	}
}

func TestResult_JSON_Roundtrip_Err(t *testing.T) {
	original := Err[int](errors.New("test error"))
	ba, err := sonic.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var restored Result[int]
	err = sonic.Unmarshal(ba, &restored)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if original.OK() != restored.OK() {
		t.Errorf("OK mismatch: %v vs %v", original.OK(), restored.OK())
	}
	if restored.Err() == nil {
		t.Error("expected error to be restored")
	}
	if restored.Err().Error() != "test error" {
		t.Errorf("error message mismatch: expected 'test error', got %q", restored.Err().Error())
	}
}

func TestResult_JSON_Roundtrip_InvalidOption(t *testing.T) {
	original := OK(0) // invalid option
	ba, err := sonic.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var restored Result[int]
	err = sonic.Unmarshal(ba, &restored)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if original.OK() != restored.OK() {
		t.Errorf("OK mismatch: %v vs %v", original.OK(), restored.OK())
	}
	if restored.Opt().Valid() {
		t.Error("option should be invalid")
	}
}

func TestResult_UnmarshalJSON_ValuePresentButInvalid(t *testing.T) {
	// Case: value is present but marked as invalid
	data := `{"result_value":42,"result_value_valid":false}`
	var r Result[int]
	err := sonic.Unmarshal([]byte(data), &r)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if r.OK() {
		t.Error("expected not OK (option invalid)")
	}
	if r.Err() != nil {
		t.Error("should not have error, just invalid option")
	}
	if r.Val() != 42 {
		t.Errorf("value should still be 42, got %d", r.Val())
	}
	if r.Opt().Valid() {
		t.Error("option should be invalid")
	}
}
