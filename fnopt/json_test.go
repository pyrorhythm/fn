package fnopt

import "testing"

func TestOf_MarshalJSON_Valid(t *testing.T) {
	o := Some(42)
	b, err := o.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	if string(b) != "42" {
		t.Errorf("MarshalJSON = %s, want 42", b)
	}
}

func TestOf_MarshalJSON_Nil(t *testing.T) {
	o := Nil[int]()
	b, err := o.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	if string(b) != "null" {
		t.Errorf("MarshalJSON(Nil) = %s, want null", b)
	}
}

func TestOf_MarshalJSON_String(t *testing.T) {
	o := Some("hello")
	b, err := o.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	if string(b) != `"hello"` {
		t.Errorf("MarshalJSON = %s, want \"hello\"", b)
	}
}

func TestOf_UnmarshalJSON_Value(t *testing.T) {
	var o Of[int]
	if err := o.UnmarshalJSON([]byte("42")); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}
	if !o.Valid() || o.Val() != 42 {
		t.Errorf("UnmarshalJSON = %v, want Some(42)", o)
	}
}

func TestOf_UnmarshalJSON_Null(t *testing.T) {
	var o Of[int]
	if err := o.UnmarshalJSON([]byte("null")); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}
	if o.Valid() {
		t.Error("UnmarshalJSON(null) should be Nil")
	}
}

func TestOf_UnmarshalJSON_Invalid(t *testing.T) {
	var o Of[int]
	if err := o.UnmarshalJSON([]byte("notjson{")); err == nil {
		t.Error("UnmarshalJSON should return error on invalid JSON")
	}
	if o.Valid() {
		t.Error("UnmarshalJSON on error should leave Of Nil")
	}
}

func TestOf_RoundTrip(t *testing.T) {
	original := Some("roundtrip")
	b, err := original.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	var restored Of[string]
	if err := restored.UnmarshalJSON(b); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	if !restored.Valid() || restored.Val() != "roundtrip" {
		t.Errorf("RoundTrip = %v, want Some(\"roundtrip\")", restored)
	}
}
