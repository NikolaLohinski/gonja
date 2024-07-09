package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Partition", func() {
	tests := []struct {
		input    PyString
		delim    string
		expected [3]PyString
	}{
		{input: "", delim: "", expected: [3]PyString{"", "", ""}},                         // Empty string with empty delimiter
		{input: "", delim: "abc", expected: [3]PyString{"", "", ""}},                      // Empty string with non-empty delimiter
		{input: "hello", delim: "", expected: [3]PyString{"hello", "", ""}},               // Non-empty string with empty delimiter
		{input: "hello", delim: "abc", expected: [3]PyString{"hello", "", ""}},            // Non-empty string with delimiter not present
		{input: "hello", delim: "e", expected: [3]PyString{"h", "e", "llo"}},              // Non-empty string with single-character delimiter
		{input: "hello", delim: "el", expected: [3]PyString{"h", "el", "lo"}},             // Non-empty string with two-character delimiter
		{input: "hello", delim: "o", expected: [3]PyString{"hell", "o", ""}},              // Non-empty string with delimiter at end
		{input: "hello", delim: "h", expected: [3]PyString{"", "h", "ello"}},              // Non-empty string with delimiter at beginning
		{input: "hello", delim: "x", expected: [3]PyString{"hello", "", ""}},              // Non-empty string with delimiter not present
		{input: "hello", delim: "hello", expected: [3]PyString{"", "hello", ""}},          // Non-empty string with delimiter being the entire string
		{input: "hello", delim: "hellx", expected: [3]PyString{"hello", "", ""}},          // Non-empty string with delimiter not present
		{input: "hello", delim: "lo", expected: [3]PyString{"hel", "lo", ""}},             // Non-empty string with delimiter in middle
		{input: "hello world", delim: " ", expected: [3]PyString{"hello", " ", "world"}},  // Non-empty string with space delimiter
		{input: "hello world", delim: "o ", expected: [3]PyString{"hell", "o ", "world"}}, // Non-empty string with two-character delimiter
		{input: "hello world", delim: "x", expected: [3]PyString{"hello world", "", ""}},  // Non-empty string with delimiter not present
	}

	for n, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("%d: For input '%s' and delimiter '%s', should return '%s', '%s', '%s'", n, test.input, test.delim, test.expected[0], test.expected[1], test.expected[2]), func() {
			r1, r2, r3 := test.input.Partition(test.delim)
			Expect(r1).To(Equal(test.expected[0]))
			Expect(r2).To(Equal(test.expected[1]))
			Expect(r3).To(Equal(test.expected[2]))
		})
	}
})
