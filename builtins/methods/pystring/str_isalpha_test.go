package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsAlpha", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: false},          // Empty string
		{input: PyString("Hello"), expected: true},      // All alphabetic characters
		{input: PyString("123"), expected: false},       // Non-alphabetic characters
		{input: PyString("Hello, 世界"), expected: false}, // Non-alphabetic characters
		{input: PyString("١٢٣"), expected: false},       // Non-alphabetic characters
		{input: PyString("Hello123"), expected: false},  // Non-alphabetic characters
		{input: PyString("Hello!"), expected: false},    // Non-alphabetic characters
		{input: PyString("123 456"), expected: false},   // Non-alphabetic characters
		{input: PyString("     "), expected: false},     // Non-alphabetic characters
		{input: PyString("ᠠᡠᠰᠱᠲ"), expected: true},      // Mongolian characters
		{input: PyString("ᚠᛁᚻ"), expected: true},        // Runic characters
		{input: PyString("ᚠᛁᚻ123"), expected: false},    // Mixed characters
		{input: PyString("ᚠᛁᚻ 123"), expected: false},   // Mixed characters with space
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is alphabetic as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsAlpha()).To(Equal(test.expected))
		})
	}
})
