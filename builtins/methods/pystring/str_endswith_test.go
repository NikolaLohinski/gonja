package pystring

import (
	"testing"
)

func TestEndsWith(t *testing.T) {
	tests := []struct {
		str    string
		subStr PyString
		res    bool
	}{
		{str: "", subStr: "", res: true},
		{str: "12345", subStr: "", res: true},
		{str: "12345", subStr: "1", res: false},
		{str: "12345", subStr: "2", res: false},
		{str: "123455", subStr: "5", res: true},
		{str: "1234555", subStr: "55", res: true},
		{str: "12345555", subStr: "55", res: true},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			pys := PyString(tt.str)
			if got := pys.EndsWith(tt.subStr, nil, nil); got != tt.res {
				t.Errorf("%q.EndWith(%q) = '%v', want '%v'", tt.str, tt.subStr, got, tt.res)
			}
		})
	}
}
