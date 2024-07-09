package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Title", func() {
	tests := []struct {
		s        string
		expected string
	}{
		{"Hello world", "Hello World"},
		{"they're bill's friends from the UK", "They'Re Bill'S Friends From The Uk"},
		{"", ""},
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input %q should return %q", test.s, test.expected), func() {
			result := PyString(test.s).Title()
			Expect(string(result)).To(Equal(test.expected))
		})
	}
})
