package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Count", func() {
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
		tt := tt // capture range variable
		It(fmt.Sprintf("%q should count %q to '%v'", tt.str, tt.subStr, tt.res), func() {
			pys := PyString(tt.str)
			Expect(pys.Count(tt.subStr, nil, nil)).To(Equal(tt.res))
		})
	}
})
