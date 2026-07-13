package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("{% set %} tag", func() {
	Context("block form", func() {
		It("should capture rendered content", func() {
			result := renderTemplate(
				`{% set greeting %}Hello {{ name }}{% endset %}{{ greeting | upper }}`,
				map[string]any{"name": "World"})
			Expect(result).To(Equal("HELLO WORLD"))
		})
	})

	Context("expression form (regression)", func() {
		It("should still work with a simple assignment", func() {
			result := renderTemplate(`{% set x = 42 %}{{ x }}`, nil)
			Expect(result).To(Equal("42"))
		})
	})
})
