package pystring

import (
	"testing"
)

func TestCapWords(t *testing.T) {
	tests := []struct {
		name     string
		input    PyString
		expected PyString
	}{
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Single word",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "Multiple words",
			input:    "hello world",
			expected: "Hello World",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  hello   world  ",
			expected: "Hello World",
		},
		{
			name:     "They're bill's friends from the uk",
			input:    "They're bill's friends from the uk",
			expected: "They're Bill's Friends From The Uk",
		},
		{
			name:     "Empty string after splitting",
			input:    " , , ",
			expected: ", ,",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := PyString(test.input).CapWords()
			if result != test.expected {
				t.Errorf("Expected %q, but got %q", test.expected, result)
			}
		})
	}
}
