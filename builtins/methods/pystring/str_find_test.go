package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Find", func() {
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
		tt := tt // capture range variable
		It(fmt.Sprintf("%q should find %q to '%v'", tt.str, tt.subStr, tt.res), func() {
			pys := PyString(tt.str)
			Expect(pys.Find(tt.subStr, nil, nil)).To(Equal(tt.res))
		})
	}
})
