package pystring

import (
	"testing"
)

func TestFind(t *testing.T) {
	tests := []struct {
		str    string
		subStr PyString
		res    int
	}{
		{str: "", subStr: "", res: 0},
		{str: "12345", subStr: "", res: 0},
		{str: "12345", subStr: "1", res: 0},
		{str: "12345", subStr: "2", res: 1},
		{str: "123455", subStr: "5", res: 4},
		{str: "1234555", subStr: "55", res: 4},
		{str: "12345555", subStr: "A", res: -1},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			pys := PyString(tt.str)
			if got := pys.Find(tt.subStr, nil, nil); got != tt.res {
				t.Errorf("%q.Find(%q) = '%v', want '%v'", tt.str, tt.subStr, got, tt.res)
			}
		})
	}
}
