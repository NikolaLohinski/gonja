package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("Capitalize", func() {
	tests := []struct {
		in  string
		out PyString
	}{
		{in: "", out: ""},
		{in: "hello", out: "Hello"},
		{in: "HELLO", out: "Hello"},
		{in: "hELLO", out: "Hello"},
		{in: "hello world", out: "Hello world"},
		{in: "ätö", out: "Ätö"},
		{in: "işğüı", out: "Işğüı"},
	}

	for _, tt := range tests {
		in := tt.in
		out := tt.out

		It(fmt.Sprintf("should capitalize '%s' to '%s'", in, out), func() {
			pys := PyString(in)
			Expect(pys.Capitalize()).To(Equal(out))
		})
	}
})
