package pystring

import (
	"fmt"
	"math"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewFormatterSpecFromStr", func() {
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
		format := test.format
		expected := test.expected
		errMsg := test.errMsg

		It(fmt.Sprintf("should handle format: %s", format), func() {
			spec, err := NewFormatterSpecFromStr(format)

			if errMsg != "" {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(errMsg))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(isEqual(spec, expected)).To(Equal(true))
			}
		})
	}
})

/*
	var _ = Describe("ValidExpressions", func() {
		expressions := getValidExpressions(DefaultDialect)
		for _, expr := range expressions {
			expr := expr

			It(fmt.Sprintf("should format values correctly for expr: %s", expr), func() {
				if expr.ExpectFloatType() {
					_, err := expr.Format(1.123456789)
					Expect(err).NotTo(HaveOccurred(), "for expr: %s", expr)
				} else if expr.ExpectIntType() {
					_, err := expr.Format(16789)
					Expect(err).NotTo(HaveOccurred(), "for expr: %s", expr)
				} else if expr.ExpectNumericType() {
					_, err := expr.Format(167800119)
					Expect(err).NotTo(HaveOccurred(), "for expr: %s", expr)
				} else if expr.Sign == 0 && expr.Fill != ' ' {
					_, err := expr.Format("foobar")
					Expect(err).NotTo(HaveOccurred(), "for expr: %s", expr)
				}
			})
		}
	})
*/

// GinKo doesn't seem to support large table tests. It is several of
// order of magnitude slower for these tests iterating over all valid
// expressions. So, we will just use a regular test for this.
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

/*
var _ = Describe("ValidExpressions serialization and deserialization", func() {
	return
	expressions := getValidExpressions(DefaultDialect)

	for _, expr := range expressions {
		expr := expr

		It(fmt.Sprintf("should serialize and deserialize correctly for expr: %s", expr), func() {
			exprStr := expr.String()
			expr2, err := NewFormatterSpecFromStr(exprStr)
			Expect(err).NotTo(HaveOccurred(), "for expr: %s", exprStr)
			Expect(isEqual(expr, expr2)).To(BeTrue(), "Expected %#v to equal %#v", expr2, expr)
		})
	}
})
*/

// GinKo doesn't seem to support large table tests. It is several of
// order of magnitude slower for these tests iterating over all valid
// expressions. So, we will just use a regular test for this.
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

var _ = Describe("Coerce Negative Zero With Z", func() {
	It("should work correctly after Python 3.11", func() {
		mustEvalToMatch(DialectPython3_11, "z.2f", math.Copysign(0.0, -1), "0.00")
	})

	It("should not support 'z' before Python 3.11", func() {
		mustEvalToError(DialectPython3_10, "z.2f", 0)
		mustEvalToMatch(DialectPython3_10, ".2f", math.Copysign(0.0, -1), "-0.00")
	})
})

var _ = Describe("PastFailures", func() {
	It("should handle special cases correctly", func() {
		// Special cases we don't want to forget about.
		mustEvalToMatch(DefaultDialect, " #1b", 16789, " 0b100000110010101")
		mustEvalToMatch(DefaultDialect, "#1b", 16789, "0b100000110010101")

		mustEvalToMatch(DefaultDialect, "2.1", "foobar", "f ")

		mustEvalToMatch(DefaultDialect, "02.1", "foobar", "f0")

		mustEvalToMatch(DefaultDialect, "#10o", 16789, "   0o40625")

		mustEvalToMatch(DefaultDialect, "#030_x", 100000000, "0x000_0000_0000_0000_05f5_e100")

		mustEvalToMatch(DefaultDialect, "#_b", 16789, "0b100_0001_1001_0101")

		mustEvalToMatch(DefaultDialect, "#10,d", 16789, "    16,789")
		mustEvalToMatch(DefaultDialect, "#010_o", 0o0040625, "0o004_0625")

		mustEvalToMatch(DefaultDialect, "#0_b", 167, "0b1010_0111")
		mustEvalToMatch(DefaultDialect, "#0,d", 16789, "16,789")

		mustEvalToError(DefaultDialect, " 0", "foobar")
		mustEvalToMatch(DefaultDialect, "010", "foobar", "foobar0000")

		mustEvalToMatch(DefaultDialect, "#010d", 16789, "0000016789")
		mustEvalToMatch(DefaultDialect, ",=-10.5G", 77.11121111111112, ",,,,77.111")
		mustEvalToMatch(DefaultDialect, "A=#x", -53, "-0x35")
		mustEvalToMatch(DefaultDialect, "`=87", 10, "`````````````````````````````````````````````````````````````````````````````````````10")
		mustEvalToError(DefaultDialect, "}<1", "no_used")
	})
})

var _ = Describe("SimpleJSONPathSplit", func() {
	DescribeTable("splits JSONPath correctly",
		func(input string, expected []string) {
			result := simpleJSONPathSplit(input)
			Expect(result).To(Equal(expected))
		},
		Entry("empty string", "", []string{}),
		Entry("single property", "$.store", []string{"$", "store"}),
		Entry("nested properties", "$.store.book", []string{"$", "store", "book"}),
		Entry("array index", "$.store.book[0]", []string{"$", "store", "book", "0"}),
		Entry("array index with quotes", "$.store.book['0']", []string{"$", "store", "book", "0"}),
		Entry("nested properties with array", "$.store.book[0].title", []string{"$", "store", "book", "0", "title"}),
		Entry("properties with single quotes", "$['store']['book'][0]['title']", []string{"$", "store", "book", "0", "title"}),
		Entry("properties with mixed quotes", "$.store['book'][0].title", []string{"$", "store", "book", "0", "title"}),
	)
})

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

func mustFormatSpec(d Dialect, template string) FormatSpec {
	spec, err := d.NewFormatterSpecFromStr(template)
	Expect(err).NotTo(HaveOccurred(), "Failed to parse format spec: %s => %v", template, err)
	return spec
}

func mustEvalToMatch(d Dialect, template string, val any, expected string) {
	spec := mustFormatSpec(d, template)
	s, err := spec.Format(val)
	Expect(err).NotTo(HaveOccurred(), "Failed to format spec: %s", err)
	Expect(s).To(Equal(expected), "Expected: '%s', got: '%s'", expected, s)
}

func mustEvalToError(d Dialect, template string, val any) {
	spec, err := d.NewFormatterSpecFromStr(template)
	if err == nil {
		_, err2 := spec.Format(val)
		Expect(err2).To(HaveOccurred(), "should fail to parse but got %#v: %s", template, spec)
	}
}
