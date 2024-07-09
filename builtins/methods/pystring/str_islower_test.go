package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsLower", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: false},             // Empty string
		{input: PyString("hello"), expected: true},         // All lowercase
		{input: PyString("Hello"), expected: false},        // Mixed case
		{input: PyString("HELLO"), expected: false},        // All uppercase
		{input: PyString("123"), expected: false},          // Non-letter characters
		{input: PyString("Hello123"), expected: false},     // Alphanumeric characters with uppercase
		{input: PyString("hello, world!"), expected: true}, // Lowercase with symbols
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is lower as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsLower()).To(Equal(test.expected))
		})
	}
})
