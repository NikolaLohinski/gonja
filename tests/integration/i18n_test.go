package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("i18n support", func() {
	Context("{% trans %}", func() {
		It("should render a plain block", func() {
			result := renderTemplate(`{% trans %}Hello World{% endtrans %}`, nil)
			Expect(result).To(Equal("Hello World"))
		})

		It("should interpolate variables", func() {
			result := renderTemplate(`{% trans name="Alice" %}Hi {{ name }}!{% endtrans %}`, nil)
			Expect(result).To(Equal("Hi Alice!"))
		})

		It("should pick singular when count=1", func() {
			result := renderTemplate(
				`{% trans count=1 %}1 apple{% pluralize %}{{ count }} apples{% endtrans %}`, nil)
			Expect(result).To(Equal("1 apple"))
		})

		It("should pick plural when count>1", func() {
			result := renderTemplate(
				`{% trans count=5 %}1 apple{% pluralize %}{{ count }} apples{% endtrans %}`, nil)
			Expect(result).To(Equal("5 apples"))
		})
	})

	Context("gettext family", func() {
		It("_() should pass string through unchanged", func() {
			result := renderTemplate(`{{ _("hello") }}`, nil)
			Expect(result).To(Equal("hello"))
		})

		It("ngettext should return singular when n==1", func() {
			result := renderTemplate(`{{ ngettext("apple", "apples", 1) }}`, nil)
			Expect(result).To(Equal("apple"))
		})

		It("ngettext should return plural when n!=1", func() {
			result := renderTemplate(`{{ ngettext("apple", "apples", 5) }}`, nil)
			Expect(result).To(Equal("apples"))
		})
	})
})
