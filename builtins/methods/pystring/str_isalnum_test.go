package pystring

import (
	"testing"
)

func TestIsAlnum(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: false},           // Empty string
		{input: PyString("Hello"), expected: true},       // All alphabetic characters
		{input: PyString("123"), expected: true},         // All numeric characters
		{input: PyString("Hello123"), expected: true},    // Alphanumeric characters
		{input: PyString("Hello, 123"), expected: false}, // Non-alphanumeric characters
		{input: PyString("123!"), expected: false},       // Non-alphanumeric characters
		{input: PyString("    "), expected: false},       // Non-alphanumeric characters
		{input: PyString("١٢٣"), expected: true},         // Numeric characters
		{input: PyString("ᠠᡠᠰᠱᠲ"), expected: true},       // Alphanumeric characters
		{input: PyString("ᚠᛁᚻ"), expected: true},         // Alphanumeric characters
		{input: PyString("Hello, 世界"), expected: false},  // Non-alphanumeric characters
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsAlnum()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
