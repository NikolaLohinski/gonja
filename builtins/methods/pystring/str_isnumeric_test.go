package pystring

import "testing"

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: PyString(""), expected: false},      // Empty string
		{input: PyString("123"), expected: true},    // All numeric characters
		{input: PyString("١٢٣"), expected: true},    // All numeric characters (Arabic digits)
		{input: PyString("12.34"), expected: false}, // Contains non-numeric characters
		{input: PyString("½"), expected: true},      // Numeric character (VULGAR FRACTION ONE HALF)
		{input: PyString("⅔"), expected: true},      // Numeric character (VULGAR FRACTION TWO THIRDS)
		{input: PyString("¼"), expected: true},      // Numeric character (VULGAR FRACTION ONE QUARTER)
		{input: PyString("A12"), expected: false},   // Contains non-numeric characters
		{input: PyString("12A"), expected: false},   // Contains non-numeric characters
		{input: PyString("1.5"), expected: false},   // Contains non-numeric characters
		{input: PyString("⅓"), expected: true},      // Numeric character (VULGAR FRACTION ONE THIRD)
		{input: PyString("٥٠٠٠"), expected: true},   // All numeric characters (Arabic digits)

		// TODO: works in python; not in golang.
		// {input: PyString("一二三"), expected: true},    // All numeric characters (Chinese numbers)
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsNumeric()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
