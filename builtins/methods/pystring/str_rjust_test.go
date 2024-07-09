package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RJust", func() {
	tests := []struct {
		input    PyString
		width    int
		fillchar rune
		expected PyString
	}{
		{input: "", width: 5, fillchar: ' ', expected: "     "},            // Empty input
		{input: "hello", width: 3, fillchar: ' ', expected: "hello"},       // Width less than string length
		{input: "hello", width: 10, fillchar: ' ', expected: "     hello"}, // Normal case
		{input: "hello", width: 7, fillchar: '*', expected: "**hello"},     // Custom fill character
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input '%s' with width '%d' and fill character '%c', should return '%s'", test.input, test.width, test.fillchar, test.expected), func() {
			Expect(test.input.RJust(test.width, test.fillchar)).To(Equal(test.expected))
		})
	}
})
