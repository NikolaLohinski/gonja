package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Center", func() {
	tests := []struct {
		in    string
		width int
		out   PyString
	}{
		{in: "", width: 15, out: "               "},
		{in: "grüßen", width: 15, out: "     grüßen    "},
		{in: "ÄÖÜ", width: 15, out: "      ÄÖÜ      "},
		{in: "ǅABCDǄ", width: 15, out: "     ǅABCDǄ    "},
		{in: "ǅABCDǄ", width: 14, out: "    ǅABCDǄ    "},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		It(fmt.Sprintf("should center '%s' to '%s' with width %d", tt.in, tt.out, tt.width), func() {
			pys := PyString(tt.in)
			Expect(pys.Center(tt.width, 0)).To(Equal(tt.out))
		})
	}
})
