package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsDecimal", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString("123"), expected: true},
		{input: PyString("123.45"), expected: false},
		{input: PyString("١٢٣"), expected: true},  // Arabic-Indic digits
		{input: PyString(""), expected: false},    // Empty string
		{input: PyString("abc"), expected: false}, // Non-decimal characters
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is decimal as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsDecimal()).To(Equal(test.expected))
		})
	}
})
