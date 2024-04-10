package pystring

import "testing"


func TestIsUpper(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: "", expected: false},       // Empty string
		{input: "BANANA", expected: true},  // All characters are uppercase
		{input: "banana", expected: false}, // All characters are lowercase
		{input: "baNana", expected: false}, // Mixed case
		{input: " ", expected: false},      // Space character
		{input: "BananA", expected: false}, // Mixed case with one uppercase
		{input: "bananas", expected: false}, // All characters are lowercase with one non-letter
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsUpper()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
	
}
