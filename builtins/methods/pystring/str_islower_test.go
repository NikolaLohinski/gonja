package pystring

import "testing"

func TestIsLower(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: false},             // Empty string
		{input: PyString("hello"), expected: true},         // All lowercase
		{input: PyString("Hello"), expected: false},        // Mixed case
		{input: PyString("HELLO"), expected: false},        // All uppercase
		{input: PyString("123"), expected: false},          // Non-letter characters
		{input: PyString("Hello123"), expected: false},     // Alphanumeric characters with uppercase
		{input: PyString("hello, world!"), expected: true}, // Lowercase with symbols
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsLower()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
