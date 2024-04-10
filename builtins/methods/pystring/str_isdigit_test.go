package pystring

import "testing"

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString("123"), expected: true},
		{input: PyString("123.45"), expected: false},
		{input: PyString("१२३"), expected: true},  // Devanagari digits
		{input: PyString("๑๒๓"), expected: true},  // Thai digits
		{input: PyString(""), expected: false},    // Empty string
		{input: PyString("abc"), expected: false}, // Non-digit characters
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsDigit()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
