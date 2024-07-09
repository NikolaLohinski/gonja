package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SwapCase", func() {
	tests := []struct {
		s        string
		expected string
	}{
		{"", ""},
		{"Hello World", "hELLO wORLD"},
		{"Spam and EGGS", "sPAM AND eggs"},
		{"12345", "12345"},
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input %q should return %q", test.s, test.expected), func() {
			result := PyString(test.s).SwapCase()
			Expect(string(result)).To(Equal(test.expected))
		})
	}
})
