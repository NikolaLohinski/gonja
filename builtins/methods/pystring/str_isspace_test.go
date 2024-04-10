package pystring

import "testing"

func TestIsSpace(t *testing.T) {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: "", expected: false},                  // Empty string
		{input: " ", expected: true},                  // ASCII space
		{input: "\t", expected: true},                 // ASCII horizontal tab
		{input: "\n", expected: true},                 // ASCII newline
		{input: "\r", expected: true},                 // ASCII carriage return
		{input: "\v", expected: true},                 // ASCII vertical tab
		{input: "\f", expected: true},                 // ASCII form feed
		{input: "\u00A0", expected: true},             // Non-breaking space
		{input: "\u2000", expected: true},             // En quad
		{input: "\u2001", expected: true},             // Em quad
		{input: "\u2002", expected: true},             // En space
		{input: "\u2003", expected: true},             // Em space
		{input: "\u2004", expected: true},             // Three-per-em space
		{input: "\u2005", expected: true},             // Four-per-em space
		{input: "\u2006", expected: true},             // Six-per-em space
		{input: "\u2007", expected: true},             // Figure space
		{input: "\u2008", expected: true},             // Punctuation space
		{input: "\u2009", expected: true},             // Thin space
		{input: "\u200A", expected: true},             // Hair space
		{input: "\u2028", expected: true},             // Line separator
		{input: "\u2029", expected: true},             // Paragraph separator
		{input: "\u202F", expected: true},             // Narrow no-break space
		{input: "\u205F", expected: true},             // Medium mathematical space
		{input: "\u3000", expected: true},             // Ideographic space
		{input: "Hello, World!", expected: false},     // Non-whitespace string
		{input: "  Hello, World!  ", expected: false}, // Non-whitespace string with leading and trailing spaces
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.IsSpace()
			if result != test.expected {
				t.Errorf("For input '%s', expected %t but got %t", test.input, test.expected, result)
			}
		})
	}
}
