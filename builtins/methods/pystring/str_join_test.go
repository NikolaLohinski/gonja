package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Join", func() {
	tests := []struct {
		s    string
		it   []string
		want string
	}{
		{",", []string{"a", "b", "c"}, "a,b,c"},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		It(fmt.Sprintf("should join %v with '%s' as '%s'", tt.it, tt.s, tt.want), func() {
			Expect(JoinString(tt.s, tt.it)).To(Equal(tt.want))
		})
	}
})
