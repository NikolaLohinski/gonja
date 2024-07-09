package pystring

import (
	"testing"
)

func BenchmarkFomatSingleReplacement(b *testing.B) {
	b.StopTimer()
	rawString := "{0}"
	vargs := []any{"foo"}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PyString(rawString).Format(vargs, nil)
	}
}

func BenchmarkFomatSingleAutoReplacement(b *testing.B) {
	b.StopTimer()
	rawString := "{}"
	vargs := []any{"foo"}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PyString(rawString).Format(vargs, nil)
	}
}

func BenchmarkFomatDoubleReplacement(b *testing.B) {
	b.StopTimer()
	rawString := "{0} {1}"
	vargs := []any{"foo", "bar"}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PyString(rawString).Format(vargs, nil)
	}
}

func BenchmarkFomatDoubleAutoReplacement(b *testing.B) {
	b.StopTimer()
	rawString := "{} {}"
	vargs := []any{"foo", "bar"}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PyString(rawString).Format(vargs, nil)
	}
}

func BenchmarkFomatSingleReplacementWithPadding(b *testing.B) {
	b.StopTimer()
	rawString := "{:<10}"
	vargs := []any{"foo"}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PyString(rawString).Format(vargs, nil)
	}
}

func BenchmarkFomatDoubleReplacementWithPadding(b *testing.B) {
	b.StopTimer()
	rawString := "{:<10} {:<10}"
	vargs := []any{"foo", "bar"}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PyString(rawString).Format(vargs, nil)
	}
}

func BenchmarkComplexSpec(b *testing.B) {
	b.StopTimer()
	rawString := "{:<} hello {:{}.4%} world {} {m.sub} {m[sub] } {:.2f} "
	vargs := []any{
		"foo", 123, "{<60", "buz", 123.321321312,
	}
	kwargs := map[string]any{
		"m": map[string]any{
			"sub": "subvalue",
		},
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		PyString(rawString).Format(vargs, kwargs)
	}
}
