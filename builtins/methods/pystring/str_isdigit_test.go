package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsDigit", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString("123"), expected: true},
		{input: PyString("123.45"), expected: false},
		{input: PyString("१२३"), expected: true},  // Devanagari digits
		{input: PyString("๑๒๓"), expected: true},  // Thai digits
		{input: PyString(""), expected: false},    // Empty string
		{input: PyString("abc"), expected: false}, // Non-digit characters
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is digit as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsDigit()).To(Equal(test.expected))
		})
	}
})
