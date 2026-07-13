package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Added filters", func() {
	Context("count", func() {
		It("should be an alias for length", func() {
			result := renderTemplate(`{{ [1,2,3,4] | count }}`, nil)
			Expect(result).To(Equal("4"))
		})
	})
})
