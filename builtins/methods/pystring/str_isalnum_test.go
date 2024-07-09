package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsAlnum", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: false},           // Empty string
		{input: PyString("Hello"), expected: true},       // All alphabetic characters
		{input: PyString("123"), expected: true},         // All numeric characters
		{input: PyString("Hello123"), expected: true},    // Alphanumeric characters
		{input: PyString("Hello, 123"), expected: false}, // Non-alphanumeric characters
		{input: PyString("123!"), expected: false},       // Non-alphanumeric characters
		{input: PyString("    "), expected: false},       // Non-alphanumeric characters
		{input: PyString("١٢٣"), expected: true},         // Numeric characters
		{input: PyString("ᠠᡠᠰᠱᠲ"), expected: true},       // Alphanumeric characters
		{input: PyString("ᚠᛁᚻ"), expected: true},         // Alphanumeric characters
		{input: PyString("Hello, 世界"), expected: false},  // Non-alphanumeric characters
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is alphanumeric as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsAlnum()).To(Equal(test.expected))
		})
	}
})
