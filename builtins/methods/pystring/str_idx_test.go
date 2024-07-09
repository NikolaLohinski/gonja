package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Index", func() {
	tests := []struct {
		s     string
		start *int
		end   *int
		want  string
	}{
		{
			s:     "hello",
			start: nil,
			end:   nil,
			want:  "hello",
		},
		{
			s:     "",
			start: nil,
			end:   nil,
			want:  "",
		},
		{
			s:     "hello",
			start: intP(0),
			end:   intP(-100),
			want:  "",
		},
		{
			s:     "hello",
			start: intP(-100),
			end:   nil,
			want:  "hello",
		},
		{
			s:     "hello",
			start: intP(2),
			end:   nil,
			want:  "llo",
		},
		{
			s:     "hello",
			start: intP(-1),
			end:   nil,
			want:  "o",
		},
		{
			s:     "hello",
			start: intP(-5),
			end:   nil,
			want:  "hello",
		},
		{
			s:     "hello",
			start: intP(-2),
			end:   nil,
			want:  "lo",
		},
		{
			s:     "hello",
			start: intP(-2),
			end:   intP(100),
			want:  "lo",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		It(fmt.Sprintf("%q should index from %v to %v as '%v'", tt.s, tt.start, tt.end, tt.want), func() {
			got, _ := PyString(tt.s).Idx(tt.start, tt.end)
			Expect(string(got)).To(Equal(tt.want))
		})
	}
})
