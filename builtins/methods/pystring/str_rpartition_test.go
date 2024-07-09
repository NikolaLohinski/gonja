package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RPartition", func() {
	tests := []struct {
		input    PyString
		substr   string
		expected struct{ Before, Separator, After PyString }
	}{
		{input: "hello", substr: "", expected: struct{ Before, Separator, After PyString }{Before: "", Separator: "", After: "hello"}},      // Empty substring
		{input: "", substr: "", expected: struct{ Before, Separator, After PyString }{Before: "", Separator: "", After: ""}},                // Empty input and substring
		{input: "hello", substr: "hello", expected: struct{ Before, Separator, After PyString }{Before: "", Separator: "hello", After: ""}}, // Substring equals input
		{input: "hello", substr: "llo", expected: struct{ Before, Separator, After PyString }{Before: "he", Separator: "llo", After: ""}},   // Normal case
		{input: "hello", substr: "x", expected: struct{ Before, Separator, After PyString }{Before: "", Separator: "", After: "hello"}},     // Substring not found
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input '%s' with substring '%s', should return (%s, %s, %s)", test.input, test.substr, test.expected.Before, test.expected.Separator, test.expected.After), func() {
			before, separator, after := test.input.RPartition(test.substr)
			Expect(before).To(Equal(test.expected.Before))
			Expect(separator).To(Equal(test.expected.Separator))
			Expect(after).To(Equal(test.expected.After))
		})
	}
})
