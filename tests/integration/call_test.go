package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("{% call %} blocks with caller()", func() {
	It("should render a basic call block", func() {
		result := renderTemplate(
			`{% macro wrap() %}<div>{{ caller() }}</div>{% endmacro %}{% call wrap() %}hello{% endcall %}`,
			nil)
		Expect(result).To(Equal("<div>hello</div>"))
	})

	It("should support arguments to the wrapping macro", func() {
		result := renderTemplate(
			`{% macro wrap(cls) %}<div class="{{ cls }}">{{ caller() }}</div>{% endmacro %}{% call wrap("box") %}body{% endcall %}`,
			nil)
		Expect(result).To(Equal(`<div class="box">body</div>`))
	})
})
