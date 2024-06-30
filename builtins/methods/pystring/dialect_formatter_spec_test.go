package pystring

import (
	"fmt"
	"math"
	"testing"
	"unicode"
)

func TestNewFormatterSpecFromStr(t *testing.T) {
	tests := []struct {
		format   string
		expected FormatSpec
		errMsg   string
	}{
		{
			format: "<10.2f",
			expected: FormatSpec{
				Fill:        0,
				Align:       '<',
				Sign:        0,
				Alternate:   false,
				ZeroPadding: false,
				MinWidth:    10,
				Precision:   2,
				Type:        'f',
			},
			errMsg: "",
		},
		{
			format: "+10.2f",
			expected: FormatSpec{
				Fill:        0,
				Align:       0,
				Sign:        '+',
				Alternate:   false,
				ZeroPadding: false,
				MinWidth:    10,
				Precision:   2,
				Type:        'f',
			},
			errMsg: "",
		},
		{
			format: "+010.2f",
			expected: FormatSpec{
				Fill:        0,
				Align:       0,
				Sign:        '+',
				Alternate:   false,
				ZeroPadding: true,
				MinWidth:    10,
				Precision:   2,
				Type:        'f',
			},
			errMsg: "",
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Format: %s", test.format), func(t *testing.T) {
			spec, err := NewFormatterSpecFromStr(test.format)

			if test.errMsg != "" {
				if err == nil {
					t.Errorf("Expected an error but got nil")
				} else if err.Error() != test.errMsg {
					t.Errorf("Expected error message '%s' but got '%s'", test.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if !isEqual(spec, test.expected) {
					t.Errorf("Expected '%+v' but got '%+v' ", test.expected, spec)
				}
			}
		})
	}
}

func TestValidExpressionsCanFormatValues(t *testing.T) {
	expressions := getValidExpressions(DefaultDialect)
	t.Logf("found %d valid expressions", len(expressions))
	for _, expr := range expressions {
		if expr.ExpectFloatType() {
			if _, err := expr.Format(1.123456789); err != nil {
				t.Errorf("Expected no error but got: %v for expr: %s", err, expr)
				return
			}
		} else if expr.ExpectIntType() {
			if _, err := expr.Format(16789); err != nil {
				t.Errorf("Expected no error but got: %v for expr: %s", err, expr)
				return
			}
		} else if expr.ExpectNumericType() {
			if _, err := expr.Format(167800119); err != nil {
				t.Errorf("Expected no error but got: %v for expr: %s", err, expr)
				return
			}
		} else if expr.Sign == 0 && expr.Fill != ' ' {
			if _, err := expr.Format("foobar"); err != nil {
				t.Errorf("Expected no error but got: %v for expr: %s", err, expr)
				return
			}
		}
	}
}

func TestValidExpressionsCanBeSerializedAndDeserialized(t *testing.T) {
	expressions := getValidExpressions(DefaultDialect)
	for _, expr := range expressions {
		exprStr := expr.String()
		expr2, err := NewFormatterSpecFromStr(exprStr)
		if err != nil {
			t.Fatalf("Expected to be able to reconstruct valid expression but got err: %v from spec %#v => %s", err, expr, exprStr)
		}
		if !isEqual(expr, expr2) {
			t.Fatalf("Expected to be able to reconstruct valid expression but spec %#v != %#v", expr2, expr)
		}
	}
}

func TestCoarceNegativeZeroWithZ(t *testing.T) {
	// Works after 3.11
	mustEvalToMatch(t, DialectPython3_11, "z.2f", math.Copysign(0.0, -1), "0.00")

	// not supported before 3.11
	mustEvalToError(t, DialectPython3_10, "z.2f", 0) // z not supported
	mustEvalToMatch(t, DialectPython3_10, ".2f", math.Copysign(0.0, -1), "-0.00")
}

func TestPastFailures(t *testing.T) {
	// Special cases we don't want to forget about.
	mustEvalToMatch(t, DefaultDialect, " #1b", 16789, " 0b100000110010101")
	mustEvalToMatch(t, DefaultDialect, "#1b", 16789,  "0b100000110010101")

	mustEvalToMatch(t, DefaultDialect, "2.1", "foobar", "f ")

	mustEvalToMatch(t, DefaultDialect, "02.1", "foobar", "f0")

	mustEvalToMatch(t, DefaultDialect, "#10o", 16789, "   0o40625")

	mustEvalToMatch(t, DefaultDialect, "#030_x", 100000000, "0x000_0000_0000_0000_05f5_e100")

	mustEvalToMatch(t, DefaultDialect, "#_b", 16789, "0b100_0001_1001_0101")

	mustEvalToMatch(t, DefaultDialect, "#10,d", 16789, "    16,789")
	mustEvalToMatch(t, DefaultDialect, "#010_o", 0o0040625, "0o004_0625")

	mustEvalToMatch(t, DefaultDialect, "#0_b", 167, "0b1010_0111")
	mustEvalToMatch(t, DefaultDialect, "#0,d", 16789, "16,789")

	mustEvalToError(t, DefaultDialect, " 0", "foobar")
	mustEvalToMatch(t, DefaultDialect, "010", "foobar", "foobar0000")

	mustEvalToMatch(t, DefaultDialect, "#010d", 16789, "0000016789")
	mustEvalToMatch(t, DefaultDialect, ",=-10.5G", 77.11121111111112, ",,,,77.111")
	mustEvalToMatch(t, DefaultDialect, "A=#x", -53, "-0x35")
	mustEvalToMatch(t, DefaultDialect, "`=87", 10, "`````````````````````````````````````````````````````````````````````````````````````10")
	mustEvalToError(t, DefaultDialect, "}<1", "no_used")
}

func FuzzFormatSpec(t *testing.F) {
	t.Add('<', rune(0), rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add('>', rune(0), rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add('^', rune(0), rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add('=', rune(0), rune(0), false, false, uint(0), uint(0), rune(0))

	t.Add(rune(0), ' ', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '>', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '.', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), 'g', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '0', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), 'O', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '#', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '<', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '^', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '=', rune(0), false, false, uint(0), uint(0), rune(0))

	t.Add(rune(0), rune(0), '+', false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), '-', false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), ' ', false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), true, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, true, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(1), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(2), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(3), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(4), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(5), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(6), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(7), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(8), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(9), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(10), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(1), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(2), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(3), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(4), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(5), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(6), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(7), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(8), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(9), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(10), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'b')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'c')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'd')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'o')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'x')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'X')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'e')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'E')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'f')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'F')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'g')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'G')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), '%')

	t.Fuzz(func(
		t *testing.T,
		Fill rune,
		Align rune,
		Sign rune,
		Alternate bool,
		ZeroPadding bool,
		MinWidth uint,
		Precision uint,
		Type rune,
	) {

		// Some fuzz responses aren't really valid test cases.
		if !unicode.IsPrint(Fill) {
			return
		}
		if Fill != 0 && Align != '<' && Align != '>' && Align != '^' && Align != '=' && Align != 0 {
			return
		}
		if Sign != 0 && Sign != '+' && Sign != '-' && Sign != ' ' {
			return
		}
		if Type != 0 && Type != 'b' && Type != 'c' && Type != 'd' && Type != 'o' && Type != 'x' && Type != 'X' && Type != 'e' && Type != 'E' && Type != 'f' && Type != 'F' && Type != 'g' && Type != 'G' && Type != '%' {
			return
		}
		if Fill != 0 && (Fill != '<' && Fill != '>' && Fill != '^' && Fill != '=') {
			return
		}
		if Align == 0 && (Fill == '<' || Fill == '>' || Fill == '^' || Fill == '=') {
			return
		}

		spec := FormatSpec{
			Fill:        Fill,
			Align:       Align,
			Sign:        Sign,
			Alternate:   Alternate,
			ZeroPadding: ZeroPadding,
			MinWidth:    MinWidth,
			Precision:   Precision,
			Type:        Type,
		}
		orig := spec.String()

		specRecovered, err := NewFormatterSpecFromStr(orig)
		if err != nil {
			t.Errorf("Expected to parse '%s' but got error: %v on data %#v", orig, err, spec)
		}
		if orig != specRecovered.String() {
			t.Errorf("Expected '%s' = '%s' on data %#v", specRecovered.String(), orig, spec)
		}
	})
}

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

/*
	Test utility functions
*/

func isEqual(spec1, spec2 FormatSpec) bool {
	return spec1.Fill == spec2.Fill &&
		spec1.Align == spec2.Align &&
		spec1.Sign == spec2.Sign &&
		spec1.Alternate == spec2.Alternate &&
		spec1.ZeroPadding == spec2.ZeroPadding &&
		spec1.MinWidth == spec2.MinWidth &&
		spec1.Precision == spec2.Precision &&
		spec1.Type == spec2.Type
}

func getValidExpressions(d Dialect) []FormatSpec {
	// printable runes are valid but let's ensure that our special characters aren't messing anything up.
	validFills := []rune{0, '<', '>', '^', '=', '+', '-', ' ', '0', 'O', '#', '%', 'b', 'c', 'd', 'e', 'E', 'f', 'F', 'g', 'G', 'n', 'o', 's', 'x', 'X'}
	validAligns := []rune{0, '<', '>', '^', '='}
	validSigns := []rune{0, ' ', '-', '+'}
	ValidTypes := []rune{0, '%', 'b', 'c', 'd', 'e', 'E', 'f', 'F', 'g', 'G', 'n', 'o', 's', 'x', 'X'}
	validAlternate := []bool{true, false}
	validZeroPadding := []bool{true, false}
	validWidths := []uint{0, 1, 2, 3, 10, 20, 30, 100}
	validGroupings := []rune{0, ',', '_'}
	validPrecisions := []uint{0, 1, 2, 5, 10, 15, 20, 30}

	res := make([]FormatSpec, 0, len(validFills)*len(validAligns)*len(validSigns)*len(validAlternate)*len(validZeroPadding)*len(validWidths)*len(validPrecisions)*len(ValidTypes))
	for _, fill := range validFills {
		for _, align := range validAligns {
			for _, sign := range validSigns {
				for _, alternate := range validAlternate {
					for _, zeroPadding := range validZeroPadding {
						for _, width := range validWidths {
							for _, _ = range validGroupings {
								for _, precision := range validPrecisions {
									for _, t := range ValidTypes {
										spec := FormatSpec{
											dialect:        d,
											Fill:           fill,
											Align:          align,
											Sign:           sign,
											Alternate:      alternate,
											ZeroPadding:    zeroPadding,
											MinWidth:       width,
											GroupingOption: 0,
											Precision:      precision,
											Type:           t,
										}
										if err := spec.Validate(); err == nil {
											res = append(res, spec)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return res
}

func mustFormatSpec(t *testing.T, d Dialect, template string) FormatSpec {
	spec, err := d.NewFormatterSpecFromStr(template)
	if err != nil {
		t.Fatalf("Failed to parse format spec: %s => %v", template, err)
	}
	return spec
}

func mustEvalToMatch(t *testing.T, d Dialect, template string, val any, expected string) {
	spec := mustFormatSpec(t, d, template)
	s, err := spec.Format(val)
	if err != nil {
		t.Fatalf("Failed to format spec: %s", err)
	}
	if s != expected {
		t.Fatalf("Expected: '%s', got: '%s'", expected, s)
	}
}

func mustEvalToError(t *testing.T, d Dialect, template string, val any) {
	spec, err := d.NewFormatterSpecFromStr(template)
	if err == nil {
		_, err2 := spec.Format(val)
		if err2 == nil {
			t.Fatalf("should fail to parse but got %#v: %s", template, spec)
		}
	}
}
