package null

import (
	"encoding/json/v2"
	"testing"
)

func BenchmarkIntUnmarshalJSON(b *testing.B) {
	input := []byte("123456")
	var nullable Int
	b.ReportAllocs()
	for b.Loop() {
		if err := json.Unmarshal(input, &nullable); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIntStringUnmarshalJSON(b *testing.B) {
	input := []byte(`"123456"`)
	var nullable String
	b.ReportAllocs()
	for b.Loop() {
		if err := json.Unmarshal(input, &nullable); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNullIntUnmarshalJSON(b *testing.B) {
	input := nullLiteral
	var nullable Int
	b.ReportAllocs()
	for b.Loop() {
		if err := json.Unmarshal(input, &nullable); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStringUnmarshalJSON(b *testing.B) {
	input := []byte(`"hello"`)
	var nullable String
	b.ReportAllocs()
	for b.Loop() {
		if err := json.Unmarshal(input, &nullable); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNullStringUnmarshalJSON(b *testing.B) {
	input := nullLiteral
	var nullable String
	b.ReportAllocs()
	for b.Loop() {
		if err := json.Unmarshal(input, &nullable); err != nil {
			b.Fatal(err)
		}
	}
}
