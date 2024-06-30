package pystring

import "testing"

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "", expected: true},                     // Empty string
		{input: "Hello, World!", expected: true},        // Printable ASCII characters
		{input: "\t", expected: false},                  // Nonprintable ASCII character (TAB)
		{input: "こんにちは", expected: true},                // Printable non-ASCII characters (Japanese)
		{input: "Hello, \x07World!", expected: false},   // Nonprintable ASCII character (BELL)
		{input: "\u200b", expected: false},              // Nonprintable Unicode character (ZERO WIDTH SPACE)
		{input: " ", expected: true},                    // Printable ASCII character (SPACE)
		{input: "Hello, \x1bWorld!", expected: false},   // Nonprintable ASCII character (ESCAPE)
		{input: "foo\x00bar", expected: false},          // Nonprintable ASCII character (NULL)
		{input: "Hello, \uFEFFWorld!", expected: false}, // Nonprintable Unicode character (ZERO WIDTH NO-BREAK SPACE)
		{input: "Hello, \rWorld!", expected: false},     // Nonprintable ASCII character (CARRIAGE RETURN)
		{input: "Hello, \x1fWorld!", expected: false},   // Nonprintable ASCII character (UNIT SEPARATOR)
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := IsPrintable(test.input)
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
