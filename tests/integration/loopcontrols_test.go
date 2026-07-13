package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("{% break %} and {% continue %}", func() {
	It("break should stop iteration", func() {
		result := renderTemplate(
			`{% for i in [1,2,3,4,5] %}{% if i == 3 %}{% break %}{% endif %}{{ i }}{% endfor %}`, nil)
		Expect(result).To(Equal("12"))
	})

	It("continue should skip the current iteration", func() {
		result := renderTemplate(
			`{% for i in [1,2,3,4,5] %}{% if i == 3 %}{% continue %}{% endif %}{{ i }}{% endfor %}`, nil)
		Expect(result).To(Equal("1245"))
	})

	It("break inside nested loop should only break inner loop", func() {
		result := renderTemplate(
			`{% for i in [1,2] %}{% for j in [1,2,3] %}{% if j == 2 %}{% break %}{% endif %}{{ i }}{{ j }} {% endfor %}{% endfor %}`, nil)
		Expect(result).To(Equal("11 21 "))
	})
})
