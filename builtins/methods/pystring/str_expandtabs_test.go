package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExpandTabs", func() {
	tests := []struct {
		s       string
		tabSize *int
		want    string
	}{
		{
			s:       "01\t012\t0123\t01234",
			tabSize: nil,
			want:    "01      012     0123    01234",
		},
		{
			s:       "01\t012\t0123\t01234",
			tabSize: intP(4),
			want:    "01  012 0123    01234",
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		It(fmt.Sprintf("should expand tabs in %q to '%v'", tt.s, tt.want), func() {
			Expect(ExpandTabs(tt.s, tt.tabSize)).To(Equal(tt.want))
		})
	}
})

// Mock implementation for intP
func intP(i int) *int {
	return &i
}
