package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RStrip", func() {
	testCases := []struct {
		input    string
		cutset   string
		expected string
	}{
		{"   spacious   ", "", "   spacious"},
		{"mississippi", "ipz", "mississ"},
		{"Monty Python", " Python", "M"},
		{"Monty Python", "Python ", "M"},
		// Regression cases for the cutFrom-init-to-zero bug: when nothing
		// trailing matches the cutset, rstrip used to return "".
		{"hello", "", "hello"},
		{"Hello\nWorld", "", "Hello\nWorld"},
		{"Hello\nWorld   ", "", "Hello\nWorld"},
		{"abc", "xyz", "abc"},
		// Non-ASCII regression: rstrip used to slice into the byte string
		// using a rune index, which mangles multi-byte characters.
		{"日本語", "", "日本語"},
		{"日本語  ", "", "日本語"},
		{"日本語xyz", "xyz", "日本語"},
	}

	for _, tc := range testCases {
		It(fmt.Sprintf("RStrip(%q, %q) should return %q", tc.input, tc.cutset, tc.expected), func() {
			pys := PyString(tc.input)
			actual := pys.RStrip(tc.cutset)
			Expect(actual).To(Equal(PyString(tc.expected)))
		})
	}
})
