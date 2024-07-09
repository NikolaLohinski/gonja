package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Casefold", func() {
	tests := []struct {
		in  string
		out PyString
	}{
		{in: "", out: ""},
		{in: "grüßen", out: "grüssen"},
		{in: "ÄÖÜ", out: "äöü"},
		{in: "ǅABCDǄ", out: "ǆabcdǆ"},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		It(fmt.Sprintf("should casefold '%s' to '%s'", tt.in, tt.out), func() {
			pys := PyString(tt.in)
			Expect(pys.Casefold()).To(Equal(tt.out))
		})
	}
})
