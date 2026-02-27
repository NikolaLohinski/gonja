package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upper", func() {
	tests := []struct {
		s        string
		expected string
	}{
		{"Hello world", "HELLO WORLD"},
		{"they're bill's friends from the UK", "THEY'RE BILL'S FRIENDS FROM THE UK"},
		{"unicode is 😊", "UNICODE IS 😊"},
		{"こんにちは、世界", "こんにちは、世界"},
		{"Привет, мир", "ПРИВЕТ, МИР"},
		{"مرحبا العالم", "مرحبا العالم"},
	}

	for _, test := range tests {
		It(fmt.Sprintf("For input %q should return %q", test.s, test.expected), func() {
			result := PyString(test.s).Upper()
			Expect(string(result)).To(Equal(test.expected))
		})
	}
})
