package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Macro parameters", func() {
	Context("*args (varargs)", func() {
		It("should collect overflow positional args into a list", func() {
			result := renderTemplate(
				`{% macro f(*args) %}{{ args }}{% endmacro %}{{ f(1,2,3) }}`, nil)
			Expect(result).To(Equal("[1, 2, 3]"))
		})
	})

	Context("**kwargs", func() {
		It("should collect unmatched keyword args into a dict", func() {
			result := renderTemplate(
				`{% macro f(**kw) %}{{ kw.a }}-{{ kw.b }}{% endmacro %}{{ f(a=1, b=2) }}`, nil)
			Expect(result).To(Equal("1-2"))
		})
	})

	Context("mixed: positional + defaults + *args + **kwargs", func() {
		It("should bind all argument types correctly", func() {
			result := renderTemplate(
				`{% macro f(x, y=10, *args, **kw) %}x={{x}} y={{y}} args={{args}} kw={{kw.k}}{% endmacro %}{{ f(1, 2, 3, 4, k=5) }}`,
				nil)
			Expect(result).To(Equal("x=1 y=2 args=[3, 4] kw=5"))
		})
	})
})
