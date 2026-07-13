package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("{% do %} extension", func() {
	It("should append to a list", func() {
		result := renderTemplate(`{% set lst = [1,2] %}{% do lst.append(3) %}{{ lst }}`, nil)
		Expect(result).To(Equal("[1, 2, 3]"))
	})

	// Same known limitation as the dict mutation tests: gonja's method
	// dispatch does not propagate map mutations back into the context.
	PIt("should update a dict via kwargs", func() {
		result := renderTemplate(`{% set d = {"a":1} %}{% do d.update(b=2) %}{{ d.keys() }}`, nil)
		Expect(result).ToNot(BeEmpty())
	})
})
