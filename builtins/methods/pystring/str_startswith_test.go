package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StartsWith", func() {
	tests := []struct {
		s        string
		prefix   string
		start    *int
		end      *int
		expected bool
	}{
		{"abcdef", "", nil, nil, true},
		{s: "abcdef", prefix: "", start: intP(0), end: intP(6), expected: true},
		{s: "abcdef", prefix: "abc", start: nil, end: nil, expected: true},
		{"abcdef", "def", nil, nil, false},
		{"abcdef", "bcd", nil, nil, false},
		{"abcdef", "a", nil, nil, true},
		{"abcdef", "cde", nil, nil, false},
		{"abcdef", "abc", nil, intP(2), false},
		{"abcdef", "cd", intP(2), nil, true},
		{s: "test123", prefix: "", start: intP(3), end: intP(1), expected: false},
		{"test123", "st", intP(2), intP(-1), true},
		{"abcdef", "abc", intP(1), intP(3), false},
		{"abcdef", "abc", intP(0), intP(2), false},
		{"abcdef", "abc", intP(0), intP(4), true},
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input %q.StartsWith(%q, %v, %v) should return %v", test.s, test.prefix, test.start, test.end, test.expected), func() {
			result := PyString(test.s).StartsWith(test.prefix, test.start, test.end)
			Expect(result).To(Equal(test.expected))
		})
	}
})
