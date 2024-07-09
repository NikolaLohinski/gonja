package pystring

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CapWords", func() {
	tests := []struct {
		name     string
		input    PyString
		expected PyString
	}{
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Single word",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "Multiple words",
			input:    "hello world",
			expected: "Hello World",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  hello   world  ",
			expected: "Hello World",
		},
		{
			name:     "They're bill's friends from the uk",
			input:    "They're bill's friends from the uk",
			expected: "They're Bill's Friends From The Uk",
		},
		{
			name:     "Empty string after splitting",
			input:    " , , ",
			expected: ", ,",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		It(tt.name, func() {
			result := PyString(tt.input).CapWords()
			Expect(result).To(Equal(tt.expected))
		})
	}
})
