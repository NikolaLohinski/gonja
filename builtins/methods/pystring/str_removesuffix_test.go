package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemoveSuffix", func() {
	tests := []struct {
		input    PyString
		suffix   string
		expected PyString
	}{
		{input: "", suffix: "", expected: ""},                            // Empty string with empty suffix
		{input: "", suffix: "abc", expected: ""},                         // Empty string with non-empty suffix
		{input: "hello", suffix: "", expected: "hello"},                  // Non-empty string with empty suffix
		{input: "hello", suffix: "abc", expected: "hello"},               // Non-empty string with suffix not present
		{input: "hello", suffix: "lo", expected: "hel"},                  // Non-empty string with suffix at end
		{input: "hello", suffix: "x", expected: "hello"},                 // Non-empty string with suffix not present
		{input: "hello", suffix: "hello", expected: ""},                  // Non-empty string with suffix being the entire string
		{input: "hello", suffix: "lo", expected: "hel"},                  // Non-empty string with suffix at end
		{input: "hello world", suffix: "world", expected: "hello "},      // Non-empty string with suffix at end
		{input: "hello world", suffix: "hello", expected: "hello world"}, // Non-empty string with suffix not present
		{input: "hello world", suffix: "hello world", expected: ""},      // Non-empty string with suffix being the entire string
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input '%s' and suffix '%s', should return '%s'", test.input, test.suffix, test.expected), func() {
			Expect(test.input.RemoveSuffix(test.suffix)).To(Equal(test.expected))
		})
	}
})
