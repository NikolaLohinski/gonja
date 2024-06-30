package pystring

import (
	"testing"
)

func TestRemoveSuffix(t *testing.T) {
	tests := []struct {
		input    PyString
		suffix   string
		expected PyString
	}{
		{input: "", suffix: "", expected: ""},                            // Empty string with empty suffix
		{input: "", suffix: "abc", expected: ""},                         // Empty string with non-empty suffix
		{input: "hello", suffix: "", expected: "hello"},                  // Non-empty string with empty suffix
		{input: "hello", suffix: "abc", expected: "hello"},               // Non-empty string with suffix not present
		{input: "hello", suffix: "lo", expected: "hel"},                  // Non-empty string with suffix at end
		{input: "hello", suffix: "x", expected: "hello"},                 // Non-empty string with suffix not present
		{input: "hello", suffix: "hello", expected: ""},                  // Non-empty string with suffix being the entire string
		{input: "hello", suffix: "lo", expected: "hel"},                  // Non-empty string with suffix at end
		{input: "hello world", suffix: "world", expected: "hello "},      // Non-empty string with suffix at end
		{input: "hello world", suffix: "hello", expected: "hello world"}, // Non-empty string with suffix not present
		{input: "hello world", suffix: "hello world", expected: ""},      // Non-empty string with suffix being the entire string
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.RemoveSuffix(test.suffix)
			if result != test.expected {
				t.Errorf("For input '%s' and suffix '%s', expected '%s' but got '%s'",
					test.input, test.suffix, test.expected, result)
			}
		})
	}
}
