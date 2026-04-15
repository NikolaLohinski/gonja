package tokens

import (
	"fmt"
	"strings"
	"testing"
)

func generateTemplate(macroCount int) string {
	var b strings.Builder
	for i := 0; i < macroCount; i++ {
		fmt.Fprintf(&b, "{%%- macro test_macro_%d(model, column_name, value) -%%}\n", i)
		fmt.Fprintf(&b, "select * from {{ model }} where {{ column_name }} = '%d'\n", i)
		b.WriteString("{%- endmacro -%}\n\n")
	}
	b.WriteString("select id, name, created_at from {{ ref('stg_orders') }} where id > 0\n")
	return b.String()
}

func BenchmarkReadablePosition(b *testing.B) {
	for _, size := range []int{100, 500, 1000, 2000} {
		template := generateTemplate(size)
		tokenPositions := make([]int, 0, 50)
		step := len(template) / 50
		for i := step; i < len(template); i += step {
			tokenPositions = append(tokenPositions, i)
		}

		b.Run(fmt.Sprintf("macros=%d/bytes=%d", size, len(template)), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				for _, pos := range tokenPositions {
					ReadablePosition(pos, template)
				}
			}
		})
	}
}

func BenchmarkReadablePositionFromOffsets(b *testing.B) {
	for _, size := range []int{100, 500, 1000, 2000} {
		template := generateTemplate(size)
		offsets := PrecomputeLineOffsets(template)

		tokenPositions := make([]int, 0, 50)
		step := len(template) / 50
		for i := step; i < len(template); i += step {
			tokenPositions = append(tokenPositions, i)
		}

		b.Run(fmt.Sprintf("macros=%d/bytes=%d", size, len(template)), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				for _, pos := range tokenPositions {
					ReadablePositionFromOffsets(pos, offsets)
				}
			}
		})
	}
}

func BenchmarkAllocComparison(b *testing.B) {
	template := generateTemplate(1000)
	offsets := PrecomputeLineOffsets(template)
	pos := len(template) * 3 / 4

	b.Run("Current_StringsSplit", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ReadablePosition(pos, template)
		}
	})

	b.Run("Fixed_BinarySearch", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ReadablePositionFromOffsets(pos, offsets)
		}
	})

	b.Run("PrecomputeLineOffsets", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			PrecomputeLineOffsets(template)
		}
	})
}

func TestReadablePositionFromOffsetsMatchesOriginal(t *testing.T) {
	template := generateTemplate(200)
	offsets := PrecomputeLineOffsets(template)

	for pos := 0; pos < len(template); pos++ {
		origLine, origCol := ReadablePosition(pos, template)
		fixedLine, fixedCol := ReadablePositionFromOffsets(pos, offsets)
		if origLine != fixedLine || origCol != fixedCol {
			t.Fatalf("mismatch at pos %d: original=(%d,%d) fixed=(%d,%d)",
				pos, origLine, origCol, fixedLine, fixedCol)
		}
	}
}

func TestReadablePositionFromOffsetsEdgeCases(t *testing.T) {
	input := "abc\ndef\nghi"
	offsets := PrecomputeLineOffsets(input)

	cases := []struct {
		pos      int
		wantLine int
		wantCol  int
	}{
		{0, 1, 1},
		{3, 1, 4},
		{4, 2, 1},
		{7, 2, 4},
		{8, 3, 1},
		{10, 3, 3},
	}

	for _, tc := range cases {
		line, col := ReadablePositionFromOffsets(tc.pos, offsets)
		origLine, origCol := ReadablePosition(tc.pos, input)

		if line != tc.wantLine || col != tc.wantCol {
			t.Errorf("pos=%d: got (%d,%d), want (%d,%d)", tc.pos, line, col, tc.wantLine, tc.wantCol)
		}
		if line != origLine || col != origCol {
			t.Errorf("pos=%d: fixed (%d,%d) != original (%d,%d)", tc.pos, line, col, origLine, origCol)
		}
	}
}
