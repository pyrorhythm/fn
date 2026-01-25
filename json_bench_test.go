package fn

import (
	"encoding/json"
	"testing"
)

type benchPerson struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	intData    = []byte(`42`)
	nullData   = []byte(`null`)
	structData = []byte(`{"name":"Alice","age":30}`)
	errData    = []byte(`{"_ERROR":"something went wrong"}`)
)

// Option benchmarks

func BenchmarkOption_Marshal_Int(b *testing.B) {
	opt := Some(42)
	b.SetBytes(int64(len(intData)))
	b.ResetTimer()
	for b.Loop() {
		_, _ = json.Marshal(opt)
	}
}

func BenchmarkOption_Marshal_Struct(b *testing.B) {
	opt := SomeAny(benchPerson{Name: "Alice", Age: 30})
	b.SetBytes(int64(len(structData)))
	b.ResetTimer()
	for b.Loop() {
		_, _ = json.Marshal(opt)
	}
}

func BenchmarkOption_Unmarshal_Int(b *testing.B) {
	b.SetBytes(int64(len(intData)))
	b.ResetTimer()
	for b.Loop() {
		var opt Option[int]
		_ = json.Unmarshal(intData, &opt)
	}
}

func BenchmarkOption_Unmarshal_Null(b *testing.B) {
	b.SetBytes(int64(len(nullData)))
	b.ResetTimer()
	for b.Loop() {
		var opt Option[int]
		_ = json.Unmarshal(nullData, &opt)
	}
}

func BenchmarkOption_Unmarshal_Struct(b *testing.B) {
	b.SetBytes(int64(len(structData)))
	b.ResetTimer()
	for b.Loop() {
		var opt Option[benchPerson]
		_ = json.Unmarshal(structData, &opt)
	}
}

// Result benchmarks

func BenchmarkResult_Marshal_OK(b *testing.B) {
	r := OK(42)
	b.SetBytes(int64(len(intData)))
	b.ResetTimer()
	for b.Loop() {
		_, _ = json.Marshal(r)
	}
}

func BenchmarkResult_Marshal_Err(b *testing.B) {
	r := Errn[int]("something went wrong")
	b.SetBytes(int64(len(errData)))
	b.ResetTimer()
	for b.Loop() {
		_, _ = json.Marshal(r)
	}
}

func BenchmarkResult_Marshal_Struct(b *testing.B) {
	r := OKAny(benchPerson{Name: "Alice", Age: 30})
	b.SetBytes(int64(len(structData)))
	b.ResetTimer()
	for b.Loop() {
		_, _ = json.Marshal(r)
	}
}

func BenchmarkResult_Unmarshal_OK(b *testing.B) {
	b.SetBytes(int64(len(intData)))
	b.ResetTimer()
	for b.Loop() {
		var r Result[int]
		_ = json.Unmarshal(intData, &r)
	}
}

func BenchmarkResult_Unmarshal_Err(b *testing.B) {
	b.SetBytes(int64(len(errData)))
	b.ResetTimer()
	for b.Loop() {
		var r Result[int]
		_ = json.Unmarshal(errData, &r)
	}
}

func BenchmarkResult_Unmarshal_Struct(b *testing.B) {
	b.SetBytes(int64(len(structData)))
	b.ResetTimer()
	for b.Loop() {
		var r Result[benchPerson]
		_ = json.Unmarshal(structData, &r)
	}
}

// Parallel benchmarks

func BenchmarkOption_Parallel_Unmarshal_Struct(b *testing.B) {
	b.SetBytes(int64(len(structData)))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var opt Option[benchPerson]
			_ = json.Unmarshal(structData, &opt)
		}
	})
}

func BenchmarkResult_Parallel_Unmarshal_OK(b *testing.B) {
	b.SetBytes(int64(len(intData)))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var r Result[int]
			_ = json.Unmarshal(intData, &r)
		}
	})
}

func BenchmarkResult_Parallel_Unmarshal_Struct(b *testing.B) {
	b.SetBytes(int64(len(structData)))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var r Result[benchPerson]
			_ = json.Unmarshal(structData, &r)
		}
	})
}
