package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Strip", func() {
	tests := []struct {
		s        string
		cutset   string
		expected string
	}{
		{"   spacious   ", "", "spacious"},
		{"www.example.com", "cmowz.", "example"},
		{"#....... Section 3.2.1 Issue #32 .......", ".#! ", "Section 3.2.1 Issue #32"},
		{"    leading and trailing     ", "ing ", "leading and trail"},
		{"  space at the end     ", " ", "space at the end"},
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input %q with cutset %q should return %q", test.s, test.cutset, test.expected), func() {
			result := PyString(test.s).Strip(test.cutset)
			Expect(string(result)).To(Equal(test.expected))
		})
	}
})
