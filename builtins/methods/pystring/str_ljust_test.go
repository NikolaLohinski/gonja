package pystring

import "testing"

func TestLJust(t *testing.T) {
	tests := []struct {
		input    PyString
		width    int
		fillchar rune
		expected PyString
	}{
		{input: "", width: 5, fillchar: ' ', expected: "     "},                    // Empty string with width > 0
		{input: "hello", width: 10, fillchar: '*', expected: "hello*****"},         // String shorter than width
		{input: "world", width: 5, fillchar: '-', expected: "world"},               // String length equals width
		{input: "foo", width: 2, fillchar: ' ', expected: "foo"},                   // Width less than string length
		{input: "bar", width: 6, fillchar: '!', expected: "bar!!!"},                // Width greater than string length
		{input: "baz", width: 0, fillchar: '*', expected: "baz"},                   // Width equals 0
		{input: "qux", width: -5, fillchar: '.', expected: "qux"},                  // Negative width
		{input: "quux", width: 6, fillchar: '*', expected: "quux**"},               // Width equals string length
		{input: "grault", width: 7, fillchar: '#', expected: "grault#"},            // Width greater than string length
		{input: "garply", width: 8, fillchar: ' ', expected: "garply  "},           // Width greater than string length
		{input: "waldo", width: 5, fillchar: '_', expected: "waldo"},               // Width equals string length
		{input: "fred", width: 4, fillchar: '+', expected: "fred"},                 // Width equals string length
		{input: "xyzzy", width: 10, fillchar: ' ', expected: "xyzzy     "},         // Width greater than string length
		{input: "thud", width: 5, fillchar: '*', expected: "thud*"},                // Width equals string length
		{input: " ", width: 5, fillchar: '-', expected: " ----"},                   // Single space character
		{input: " ", width: 0, fillchar: '*', expected: " "},                       // Single space character with width = 0
		{input: " ", width: -5, fillchar: '+', expected: " "},                      // Single space character with negative width
		{input: "\t", width: 5, fillchar: '*', expected: "\t****"},                 // Single tab character
		{input: "\n", width: 5, fillchar: '*', expected: "\n****"},                 // Single newline character
		{input: "hello", width: 0, fillchar: '*', expected: "hello"},               // Width equals 0
		{input: "hello", width: -5, fillchar: '*', expected: "hello"},              // Negative width
		{input: "\x00", width: 5, fillchar: '*', expected: "\x00****"},             // Single null character
		{input: "\uFFFD", width: 5, fillchar: '*', expected: "\uFFFD****"},         // Single replacement character
		{input: "\U0001F609", width: 5, fillchar: '*', expected: "\U0001F609****"}, // Single emoji character
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.LJust(test.width, test.fillchar)
			if result != test.expected {
				t.Errorf("For input '%s', expected '%s' but got '%s'", test.input, test.expected, result)
			}
		})
	}
}
