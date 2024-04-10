package pystring

import "testing"

func TestSwapCase(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{"", ""},
		{"Hello World", "hELLO wORLD"},
		{"Spam and EGGS", "sPAM AND eggs"},
		{"12345", "12345"},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			result := PyString(test.s).SwapCase()
			if string(result) != test.expected {
				t.Errorf("Expected %s but got %s", test.expected, result)
			}
		})
	}
}
