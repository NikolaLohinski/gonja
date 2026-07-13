package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("loop.* variables", func() {
	Context("depth tracking in nested loops", func() {
		It("loop.depth should count from 1", func() {
			result := renderTemplate(
				`{% for i in [1,2] %}{% for j in [3,4] %}{{ loop.depth }} {% endfor %}{% endfor %}`, nil)
			Expect(result).To(Equal("2 2 2 2 "))
		})

		It("loop.depth0 should count from 0", func() {
			result := renderTemplate(
				`{% for i in [1,2] %}{% for j in [3,4] %}{{ loop.depth0 }} {% endfor %}{% endfor %}`, nil)
			Expect(result).To(Equal("1 1 1 1 "))
		})
	})

	Context("loop.cycle", func() {
		It("should cycle values with lowercase name", func() {
			result := renderTemplate(
				`{% for i in [1,2,3,4] %}{{ loop.cycle("odd","even") }} {% endfor %}`, nil)
			Expect(result).To(Equal("odd even odd even "))
		})

		It("should cycle values with capitalized Cycle for backward compat", func() {
			result := renderTemplate(
				`{% for i in [1,2,3,4] %}{{ loop.Cycle("odd","even") }} {% endfor %}`, nil)
			Expect(result).To(Equal("odd even odd even "))
		})
	})

	Context("loop.changed", func() {
		It("should detect value changes between iterations", func() {
			result := renderTemplate(
				`{% for i in [1,1,2,2,3] %}{{ loop.changed(i) }} {% endfor %}`, nil)
			Expect(result).To(Equal("True False True False True "))
		})
	})
})
