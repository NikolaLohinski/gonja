package pystring

import (
	"reflect"
	"testing"
)

func TestRSplit(t *testing.T) {
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
		t.Run(string(test.input), func(t *testing.T) {
			result := test.input.RSplit(test.delim, test.maxSplit)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("For input '%s' with delimiter '%s' and max split %d, expected %v but got %v",
					test.input, test.delim, test.maxSplit, test.expected, result)
			}
		})
	}
}
