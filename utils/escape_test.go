package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("escape", func() {
	It("uses HTML-compatible quote entities", func() {
		Expect(Escape(`<tag "quote" 'apostrophe'>`)).To(Equal("&lt;tag &#34;quote&#34; &#39;apostrophe&#39;&gt;"))
	})
})
