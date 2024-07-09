package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemovePrefix", func() {
	tests := []struct {
		input    PyString
		prefix   string
		expected PyString
	}{
		{input: "", prefix: "", expected: ""},                            // Empty string with empty prefix
		{input: "", prefix: "abc", expected: ""},                         // Empty string with non-empty prefix
		{input: "hello", prefix: "", expected: "hello"},                  // Non-empty string with empty prefix
		{input: "hello", prefix: "abc", expected: "hello"},               // Non-empty string with prefix not present
		{input: "hello", prefix: "he", expected: "llo"},                  // Non-empty string with prefix at beginning
		{input: "hello", prefix: "x", expected: "hello"},                 // Non-empty string with prefix not present
		{input: "hello", prefix: "hello", expected: ""},                  // Non-empty string with prefix being the entire string
		{input: "hello", prefix: "hel", expected: "lo"},                  // Non-empty string with prefix at beginning
		{input: "hello world", prefix: "hello", expected: " world"},      // Non-empty string with prefix at beginning
		{input: "hello world", prefix: "world", expected: "hello world"}, // Non-empty string with prefix not present
		{input: "hello world", prefix: "hello world", expected: ""},      // Non-empty string with prefix being the entire string
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input '%s' and prefix '%s', should return '%s'", test.input, test.prefix, test.expected), func() {
			Expect(test.input.RemovePrefix(test.prefix)).To(Equal(test.expected))
		})
	}
})
