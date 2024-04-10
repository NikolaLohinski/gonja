package pystring

import (
	"testing"
)

func TestRPartition(t *testing.T) {
	tests := []struct {
		input    PyString
		substr   string
		expected struct{ Before, Separator, After PyString }
	}{
		{input: "hello", substr: "", expected: struct{ Before, Separator, After PyString }{Before: "", Separator: "", After: "hello"}},      // Empty substring
		{input: "", substr: "", expected: struct{ Before, Separator, After PyString }{Before: "", Separator: "", After: ""}},                // Empty input and substring
		{input: "hello", substr: "hello", expected: struct{ Before, Separator, After PyString }{Before: "", Separator: "hello", After: ""}}, // Substring equals input
		{input: "hello", substr: "llo", expected: struct{ Before, Separator, After PyString }{Before: "he", Separator: "llo", After: ""}},   // Normal case
		{input: "hello", substr: "x", expected: struct{ Before, Separator, After PyString }{Before: "", Separator: "", After: "hello"}},     // Substring not found
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			before, separator, after := test.input.RPartition(test.substr)
			if before != test.expected.Before || separator != test.expected.Separator || after != test.expected.After {
				t.Errorf("For input '%s' with substring '%s', expected (%s, %s, %s) but got (%s, %s, %s)",
					test.input, test.substr, test.expected.Before, test.expected.Separator, test.expected.After, before, separator, after)
			}
		})
	}
}
