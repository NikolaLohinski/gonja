package pystring

import "testing"

func TestZFill(t *testing.T) {
	tests := []struct {
		s        string
		width    int
		expected string
	}{
		{"42", 5, "00042"},
		{"-42", 5, "-0042"},
		{"hello", 8, "000hello"},
		{"world", 3, "world"},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			result := PyString(test.s).ZFill(test.width)
			if string(result) != test.expected {
				t.Errorf("Expected %s but got %s", test.expected, result)
			}
		})
	}
}
