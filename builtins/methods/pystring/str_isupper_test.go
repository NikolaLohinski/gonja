package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsUpper", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: "", expected: false},        // Empty string
		{input: "BANANA", expected: true},   // All characters are uppercase
		{input: "banana", expected: false},  // All characters are lowercase
		{input: "baNana", expected: false},  // Mixed case
		{input: " ", expected: false},       // Space character
		{input: "BananA", expected: false},  // Mixed case with one uppercase
		{input: "bananas", expected: false}, // All characters are lowercase with one non-letter
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is upper as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsUpper()).To(Equal(test.expected))
		})
	}
})
