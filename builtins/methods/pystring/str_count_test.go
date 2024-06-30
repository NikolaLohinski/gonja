package pystring

import "testing"

func TestCount(t *testing.T) {
	tests := []struct {
		str    string
		subStr PyString
		res    int
	}{
		{str: "", subStr: "", res: 1},
		{str: "12345", subStr: "", res: 6},
		{str: "12345", subStr: "1", res: 1},
		{str: "12345", subStr: "2", res: 1},
		{str: "123455", subStr: "5", res: 2},
		{str: "1234555", subStr: "55", res: 1},
		{str: "12345555", subStr: "55", res: 2},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			pys := PyString(tt.str)
			if got := pys.Count(tt.subStr, nil, nil); got != tt.res {
				t.Errorf("%q.Count(%q) = '%v', want '%v'", tt.str, tt.subStr, got, tt.res)
			}
		})
	}
}
