package pystring

import "testing"

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: false},          // Empty string
		{input: PyString("Hello"), expected: true},      // All alphabetic characters
		{input: PyString("123"), expected: false},       // Non-alphabetic characters
		{input: PyString("Hello, 世界"), expected: false}, // Non-alphabetic characters
		{input: PyString("١٢٣"), expected: false},       // Non-alphabetic characters
		{input: PyString("Hello123"), expected: false},  // Non-alphabetic characters
		{input: PyString("Hello!"), expected: false},    // Non-alphabetic characters
		{input: PyString("123 456"), expected: false},   // Non-alphabetic characters
		{input: PyString("     "), expected: false},     // Non-alphabetic characters
		{input: PyString("ᠠᡠᠰᠱᠲ"), expected: true},      // Mongolian characters
		{input: PyString("ᚠᛁᚻ"), expected: true},        // Runic characters
		{input: PyString("ᚠᛁᚻ123"), expected: false},    // Mixed characters
		{input: PyString("ᚠᛁᚻ 123"), expected: false},   // Mixed characters with space
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsAlpha()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
