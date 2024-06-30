package pystring

import "testing"

func TestIsDecimal(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString("123"), expected: true},
		{input: PyString("123.45"), expected: false},
		{input: PyString("١٢٣"), expected: true},  // Arabic-Indic digits
		{input: PyString(""), expected: false},    // Empty string
		{input: PyString("abc"), expected: false}, // Non-decimal characters
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsDecimal()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
