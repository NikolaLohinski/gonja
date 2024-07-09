package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EndsWith", func() {
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
		tt := tt // capture range variable
		It(fmt.Sprintf("%q should end with %q to '%v'", tt.str, tt.subStr, tt.res), func() {
			pys := PyString(tt.str)
			Expect(pys.EndsWith(tt.subStr, nil, nil)).To(Equal(tt.res))
		})
	}
})
