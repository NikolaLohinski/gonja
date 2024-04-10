package pystring

import (
	"reflect"
	"testing"
)

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input    string
		keepends bool
		expected []string
	}{
		{
			input:    "",
			keepends: false,
			expected: []string{},
		},
		{
			input:    "\n",
			keepends: false,
			expected: []string{""},
		},
		{
			input:    "\na\n",
			keepends: false,
			expected: []string{"", "a"},
		},
		{
			input:    "Line 1\nLine 2\nLine 3",
			keepends: false,
			expected: []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			input:    "Line 1\r\nLine 2\r\nLine 3",
			keepends: false,
			expected: []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			input:    "Line 1\r\nLine 2\r\nLine 3\r\n",
			keepends: true,
			expected: []string{"Line 1\r\n", "Line 2\r\n", "Line 3\r\n"},
		},
		{
			input:    "Line 1\n\rLine 2\n\rLine 3",
			keepends: false,
			expected: []string{"Line 1", "", "Line 2", "", "Line 3"},
		},
		{
			input:    "Line 1\rLine 2\rLine 3",
			keepends: true,
			expected: []string{"Line 1\r", "Line 2\r", "Line 3"},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := SplitLines(test.input, test.keepends)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %#v but got %#v", test.expected, result)
			}
		})
	}
}
