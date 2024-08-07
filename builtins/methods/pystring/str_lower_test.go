package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lower", func() {
	tests := []struct {
		input    PyString
		expected PyString
	}{
		{input: "", expected: ""},                     // Empty string
		{input: "HELLO", expected: "hello"},           // All uppercase
		{input: "hello", expected: "hello"},           // All lowercase
		{input: "HeLLo", expected: "hello"},           // Mixed case
		{input: "123", expected: "123"},               // Non-letter characters
		{input: "He12LLo", expected: "he12llo"},       // Mixed case with non-letter characters
		{input: "hElLo", expected: "hello"},           // Mixed case with some lowercase
		{input: "hElLo123", expected: "hello123"},     // Mixed case with some lowercase and non-letter characters
		{input: " ", expected: " "},                   // Space character
		{input: "\t", expected: "\t"},                 // Tab character
		{input: "\n", expected: "\n"},                 // Newline character
		{input: "\uFFFD", expected: "\ufffd"},         // Replacement character
		{input: "\U0001F609", expected: "\U0001f609"}, // Emoji character
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should lower '%s' to '%s'", test.input, test.expected), func() {
			Expect(test.input.Lower()).To(Equal(test.expected))
		})
	}
})
