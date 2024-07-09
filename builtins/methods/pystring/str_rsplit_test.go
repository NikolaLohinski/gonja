package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RSplit", func() {
	tests := []struct {
		input    PyString
		delim    string
		maxSplit int
		expected []PyString
	}{
		{input: "hello", delim: "", maxSplit: 0, expected: []PyString{"hello"}},
		{input: "hello", delim: "", maxSplit: 1, expected: []PyString{"hello"}},
		{input: "", delim: "", maxSplit: 0, expected: []PyString{""}},
		{input: "hello", delim: ",", maxSplit: 0, expected: []PyString{"hello"}},
		{input: "hello,world", delim: ",", maxSplit: 0, expected: []PyString{"hello,world"}},
		{input: "hello,world", delim: ",", maxSplit: 1, expected: []PyString{"hello,world"}},
		{input: "hello,world", delim: ",", maxSplit: -1, expected: []PyString{"hello", "world"}},
		{input: "hello,world", delim: ",", maxSplit: 2, expected: []PyString{"hello", "world"}},
		{input: "hello,world", delim: ",", maxSplit: 3, expected: []PyString{"hello", "world"}},
		{input: "hello,world", delim: ",", maxSplit: 4, expected: []PyString{"hello", "world"}},
		{input: "hello,world", delim: ",", maxSplit: 5, expected: []PyString{"hello", "world"}},
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("For input '%s' with delimiter '%s' and max split %d, should return %v", test.input, test.delim, test.maxSplit, test.expected), func() {
			Expect(test.input.RSplit(test.delim, test.maxSplit)).To(Equal(test.expected))
		})
	}
})
