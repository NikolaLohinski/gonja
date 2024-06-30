package pystring

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	testCases := []struct {
		input       string
		delim       string
		maxSplit    int
		expected    []string
		description string
	}{
		{
			input:       "1,2,3",
			delim:       ",",
			maxSplit:    -1,
			expected:    []string{"1", "2", "3"},
			description: "Split with comma delimiter and no max split",
		},
		{
			input:       "1,2,3",
			delim:       ",",
			maxSplit:    1,
			expected:    []string{"1,2,3"},
			description: "Split with comma delimiter and max split of 1",
		},
		{
			input:       "1,,2,,3",
			delim:       ",,",
			maxSplit:    -1,
			expected:    []string{"1", "2", "3"},
			description: "Split with comma delimiter and max split of 1",
		},
		{
			input:       "1,2,,3,",
			delim:       ",",
			maxSplit:    -1,
			expected:    []string{"1", "2", "", "3", ""},
			description: "Split with comma delimiter and no max split, including empty strings",
		},
		{
			input:       "1    2   3",
			delim:       "", // This is actually undefined behavior in Python (error raised) but we use it as a proxy for None
			maxSplit:    -1,
			expected:    []string{"1", "2", "3"},
			description: "Split with no delimiter and no max split",
		},
		{
			input:       "1 2 3",
			delim:       "", // This is actually undefined behavior in Python (error raised) but we use it as a proxy for None
			maxSplit:    1,
			expected:    []string{"1 2 3"},
			description: "Split with no delimiter and max split of 1",
		},
		{
			input:       "   1   2   3   ",
			delim:       "", // This is actually undefined behavior in Python (error raised) but we use it as a proxy for None
			maxSplit:    -1,
			expected:    []string{"1", "2", "3"},
			description: "Split with no delimiter and no max split, including whitespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := Split(tc.input, tc.delim, tc.maxSplit)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, actual)
			}
		})
	}
}
