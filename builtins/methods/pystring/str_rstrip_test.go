package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RStrip", func() {
	testCases := []struct {
		input    string
		cutset   string
		expected string
	}{
		{"   spacious   ", "", "   spacious"},
		{"mississippi", "ipz", "mississ"},
		{"Monty Python", " Python", "M"},
		{"Monty Python", "Python ", "M"},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		It(fmt.Sprintf("RStrip(%q, %q) should return %q", tc.input, tc.cutset, tc.expected), func() {
			pys := PyString(tc.input)
			actual := pys.RStrip(tc.cutset)
			Expect(actual).To(Equal(PyString(tc.expected)))
		})
	}
})
