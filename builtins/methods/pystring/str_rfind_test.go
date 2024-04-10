package pystring

import (
	"testing"
)

func TestRFind(t *testing.T) {
	tests := []struct {
		input    PyString
		substr   string
		start    *int
		end      *int
		expected int
	}{
		{input: "hello", substr: "", start: nil, end: nil, expected: 5},               // Empty substring
		{input: "", substr: "", start: nil, end: nil, expected: 0},                    // Empty input and substring
		{input: "hello", substr: "hello", start: nil, end: nil, expected: 0},          // Substring equals input
		{input: "hello", substr: "llo", start: nil, end: nil, expected: 2},            // Normal case
		{input: "hello", substr: "llo", start: nil, end: nil, expected: 2},            // Normal case
		{input: "hello", substr: "l", start: nil, end: nil, expected: 3},              // Substring at the end
		{input: "hello", substr: "lo", start: nil, end: nil, expected: 3},             // Overlapping substrings
		{input: "hello", substr: "x", start: nil, end: nil, expected: -1},             // Substring not found
		{input: "hello", substr: "he", start: nil, end: nil, expected: 0},             // Substring at the beginning
		{input: "hello", substr: "hello", start: intP(1), end: intP(4), expected: -1}, // Substring not found within the specified range
		{input: "hello", substr: "ll", start: intP(1), end: intP(4), expected: 1},     // Substring found within the specified range
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.RFind(test.substr, test.start, test.end)
			if result != test.expected {
				t.Errorf("For input '%s', substring '%s', start '%v', end '%v', expected '%d' but got '%d'",
					test.input, test.substr, test.start, test.end, test.expected, result)
			}
		})
	}
}
