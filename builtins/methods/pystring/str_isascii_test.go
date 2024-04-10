package pystring

import "testing"

func TestIsASCII(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: true},           // Empty string
		{input: PyString("Hello"), expected: true},      // All ASCII characters
		{input: PyString("Hello, 世界"), expected: false}, // Non-ASCII characters
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsASCII()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
