package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bug fixes", func() {
	Context("testNotEqual with uncomparable types", func() {
		It("should return False for two equal lists", func() {
			result := renderTemplate(`{{ [1,2,3] is ne [1,2,3] }}`, nil)
			Expect(result).To(Equal("False"))
		})

		It("should return True for two different lists", func() {
			result := renderTemplate(`{{ [1,2,3] is ne [4,5,6] }}`, nil)
			Expect(result).To(Equal("True"))
		})
	})
})
