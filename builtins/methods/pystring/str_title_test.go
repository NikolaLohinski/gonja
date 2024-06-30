package pystring

import (
	"testing"
)

func TestTitle(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{"Hello world", "Hello World"},
		{"they're bill's friends from the UK", "They'Re Bill'S Friends From The Uk"},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			result := PyString(test.s).Title()
			if string(result) != test.expected {
				t.Errorf("Expected %s but got %s", test.expected, result)
			}
		})
	}
}
