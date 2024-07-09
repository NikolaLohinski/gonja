package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsPrintable", func() {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "", expected: true},                     // Empty string
		{input: "Hello, World!", expected: true},        // Printable ASCII characters
		{input: "\t", expected: false},                  // Nonprintable ASCII character (TAB)
		{input: "こんにちは", expected: true},                // Printable non-ASCII characters (Japanese)
		{input: "Hello, \x07World!", expected: false},   // Nonprintable ASCII character (BELL)
		{input: "\u200b", expected: false},              // Nonprintable Unicode character (ZERO WIDTH SPACE)
		{input: " ", expected: true},                    // Printable ASCII character (SPACE)
		{input: "Hello, \x1bWorld!", expected: false},   // Nonprintable ASCII character (ESCAPE)
		{input: "foo\x00bar", expected: false},          // Nonprintable ASCII character (NULL)
		{input: "Hello, \uFEFFWorld!", expected: false}, // Nonprintable Unicode character (ZERO WIDTH NO-BREAK SPACE)
		{input: "Hello, \rWorld!", expected: false},     // Nonprintable ASCII character (CARRIAGE RETURN)
		{input: "Hello, \x1fWorld!", expected: false},   // Nonprintable ASCII character (UNIT SEPARATOR)
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is printable as '%t'", test.input, test.expected), func() {
			Expect(IsPrintable(test.input)).To(Equal(test.expected))
		})
	}
})
