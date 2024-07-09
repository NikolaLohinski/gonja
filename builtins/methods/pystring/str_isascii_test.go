package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsASCII", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: true},           // Empty string
		{input: PyString("Hello"), expected: true},      // All ASCII characters
		{input: PyString("Hello, 世界"), expected: false}, // Non-ASCII characters
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is ASCII as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsASCII()).To(Equal(test.expected))
		})
	}
})
