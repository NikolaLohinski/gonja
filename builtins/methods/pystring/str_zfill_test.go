package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ZFill", func() {
	tests := []struct {
		s        string
		width    int
		expected string
	}{
		{"42", 5, "00042"},
		{"-42", 5, "-0042"},
		{"hello", 8, "000hello"},
		{"world", 3, "world"},
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input %q with width %d should return %q", test.s, test.width, test.expected), func() {
			result := PyString(test.s).ZFill(test.width)
			Expect(string(result)).To(Equal(test.expected))
		})
	}
})
