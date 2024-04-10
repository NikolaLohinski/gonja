package pystring

import (
	"testing"
)

func TestRJust(t *testing.T) {
	tests := []struct {
		input    PyString
		width    int
		fillchar rune
		expected PyString
	}{
		{input: "", width: 5, fillchar: ' ', expected: "     "},            // Empty input
		{input: "hello", width: 3, fillchar: ' ', expected: "hello"},       // Width less than string length
		{input: "hello", width: 10, fillchar: ' ', expected: "     hello"}, // Normal case
		{input: "hello", width: 7, fillchar: '*', expected: "**hello"},     // Custom fill character
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.RJust(test.width, test.fillchar)
			if result != test.expected {
				t.Errorf("For input '%s' with width '%d' and fill character '%c', expected '%s' but got '%s'",
					test.input, test.width, test.fillchar, test.expected, result)
			}
		})
	}
}
