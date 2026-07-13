package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("namespace() global", func() {
	It("should allow mutation across a for-loop", func() {
		// Classic Jinja2 idiom: use namespace() so that {% set ns.x = ... %}
		// inside a loop persists after the loop ends.
		result := renderTemplate(
			`{% set ns = namespace(x=0) %}{% for i in [1,2,3,4] %}{% set ns.x = ns.x + i %}{% endfor %}{{ ns.x }}`,
			nil)
		Expect(result).To(Equal("10"))
	})
})
