package pystring

import "testing"

func TestLStrip(t *testing.T) {
	tests := []struct {
		input    PyString
		cutset   string
		expected PyString
	}{
		{input: "", cutset: "", expected: ""},                                              // Empty string with empty cutset
		{input: "", cutset: "abc", expected: ""},                                           // Empty string with non-empty cutset
		{input: "hello", cutset: "", expected: "hello"},                                    // Non-empty string with empty cutset
		{input: "hello", cutset: "abc", expected: "hello"},                                 // Non-empty string with cutset not present
		{input: "  	space and tab prefix", cutset: " ", expected: "	space and tab prefix"}, // Non-empty string with leading spaces and tab and single space cutset
		{input: "hello", cutset: "h", expected: "ello"},                                    // Non-empty string with single-character cutset
		{input: "hello", cutset: "he", expected: "llo"},                                    // Non-empty string with two-character cutset
		{input: "   spacious   ", cutset: "", expected: "spacious   "},                     // String with leading and trailing spaces, empty cutset
		{input: "   spacious   ", cutset: " ", expected: "spacious   "},                    // String with leading and trailing spaces, single-space cutset
		{input: "   spacious   ", cutset: "sp ", expected: "acious   "},                    // String with leading and trailing spaces, multi-character cutset
		{input: "www.example.com", cutset: "cmowz.", expected: "example.com"},              // String with leading characters from cutset
		{input: "   \t\n", cutset: "", expected: ""},                                       // String with only whitespace characters, empty cutset
		{input: "   \t\n", cutset: " \t\n", expected: ""},                                  // String with only whitespace characters, whitespace cutset
		{input: "   \t\n", cutset: " \t", expected: "\n"},                                  // String with only whitespace characters, single whitespace cutset
		{input: "   \t\n", cutset: "\t", expected: "   \t\n"},                              // String with only whitespace characters, single tab cutset
		{input: "\t\t\t", cutset: "\t", expected: ""},                                      // String with only tab characters, single tab cutset
		{input: "\n\n\n", cutset: "\n", expected: ""},                                      // String with only newline characters, single newline cutset
	}

	for n, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.LStrip(test.cutset)
			if result != test.expected {
				t.Errorf("%d: For input '%s' and cutset '%s', expected '%s' but got '%s'", n, test.input, test.cutset, test.expected, result)
			}
		})
	}
}
