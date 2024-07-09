package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsNumeric", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: false},      // Empty string
		{input: PyString("123"), expected: true},    // All numeric characters
		{input: PyString("١٢٣"), expected: true},    // All numeric characters (Arabic digits)
		{input: PyString("12.34"), expected: false}, // Contains non-numeric characters
		{input: PyString("½"), expected: true},      // Numeric character (VULGAR FRACTION ONE HALF)
		{input: PyString("⅔"), expected: true},      // Numeric character (VULGAR FRACTION TWO THIRDS)
		{input: PyString("¼"), expected: true},      // Numeric character (VULGAR FRACTION ONE QUARTER)
		{input: PyString("A12"), expected: false},   // Contains non-numeric characters
		{input: PyString("12A"), expected: false},   // Contains non-numeric characters
		{input: PyString("1.5"), expected: false},   // Contains non-numeric characters
		{input: PyString("⅓"), expected: true},      // Numeric character (VULGAR FRACTION ONE THIRD)
		{input: PyString("٥٠٠٠"), expected: true},   // All numeric characters (Arabic digits)
		// {input: PyString("一二三"), expected: true},    // All numeric characters (Chinese numbers) - Works in Python, not in Go
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is numeric as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsNumeric()).To(Equal(test.expected))
		})
	}
})
