package pystring

import (
	"testing"
)

func TestStrip(t *testing.T) {
	tests := []struct {
		s        string
		cutset   string
		expected string
	}{
		{"   spacious   ", "", "spacious"},
		{"www.example.com", "cmowz.", "example"},
		{"#....... Section 3.2.1 Issue #32 .......", ".#! ", "Section 3.2.1 Issue #32"},
		{"    leading and trailing     ", "ing ", "leading and trail"},
		{"  space at the end     ", " ", "space at the end"},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			result := PyString(test.s).Strip(test.cutset)
			if string(result) != test.expected {
				t.Errorf("Expected %s but got %s", test.expected, result)
			}
		})
	}
}
