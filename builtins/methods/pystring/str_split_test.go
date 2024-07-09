package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Split", func() {
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
			description: "Split with double comma delimiter and no max split",
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
		tc := tc // capture range variable
		It(fmt.Sprintf("%s", tc.description), func() {
			actual := Split(tc.input, tc.delim, tc.maxSplit)
			Expect(actual).To(Equal(tc.expected))
		})
	}
})
