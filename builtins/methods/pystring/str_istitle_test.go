package pystring

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsTitle", func() {
	tests := []struct {
		input    PyString
		expected bool
	}{
		{input: "", expected: false},              // Empty string
		{input: "Hello", expected: true},          // Titlecased string
		{input: "Hello World", expected: true},    // Titlecased string with spaces
		{input: "hello", expected: false},         // Lowercase string
		{input: "hello world", expected: false},   // Lowercase string with spaces
		{input: "Hello World ", expected: true},   // Titlecased string with trailing space
		{input: "Hello World!", expected: true},   // Titlecased string with punctuation
		{input: "Hello123", expected: true},       // Titlecased string with numbers
		{input: "Hello  World", expected: true},   // Titlecased string with double space
		{input: "123Hello World", expected: true}, // Titlecased string with numbers at the beginning
		{input: "Hello_world", expected: false},   // Underscore not considered titlecased
		{input: "HellO", expected: false},         // Mixture of uppercase and lowercase letters
	}

	for _, test := range tests {
		test := test // capture range variable
		It(fmt.Sprintf("should determine if '%s' is title as '%t'", test.input, test.expected), func() {
			Expect(test.input.IsTitle()).To(Equal(test.expected))
		})
	}
})
