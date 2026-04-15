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
		{s: "abcdef", prefix: "", start: new(0), end: new(6), expected: true},
		{s: "abcdef", prefix: "abc", start: nil, end: nil, expected: true},
		{"abcdef", "def", nil, nil, false},
		{"abcdef", "bcd", nil, nil, false},
		{"abcdef", "a", nil, nil, true},
		{"abcdef", "cde", nil, nil, false},
		{"abcdef", "abc", nil, new(2), false},
		{"abcdef", "cd", new(2), nil, true},
		{s: "test123", prefix: "", start: new(3), end: new(1), expected: false},
		{"test123", "st", new(2), new(-1), true},
		{"abcdef", "abc", new(1), new(3), false},
		{"abcdef", "abc", new(0), new(2), false},
		{"abcdef", "abc", new(0), new(4), true},
	}

	for _, test := range tests {
		It(fmt.Sprintf("For input %q.StartsWith(%q, %v, %v) should return %v", test.s, test.prefix, test.start, test.end, test.expected), func() {
			result := PyString(test.s).StartsWith(test.prefix, test.start, test.end)
			Expect(result).To(Equal(test.expected))
		})
	}
})
