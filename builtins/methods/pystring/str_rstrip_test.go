package pystring

import (
	"fmt"
	"testing"
)

func TestRStrip(t *testing.T) {
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
		t.Run(fmt.Sprintf("RStrip(%q, %q)", tc.input, tc.cutset), func(t *testing.T) {
			pys := PyString(tc.input)
			actual := pys.RStrip(tc.cutset)
			if actual != PyString(tc.expected) {
				t.Errorf("Expected RStrip(%q, %q) to be %q, got %q", tc.input, tc.cutset, tc.expected, actual)
			}
		})
	}
}
