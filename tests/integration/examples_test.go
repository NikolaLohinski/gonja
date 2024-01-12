package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("examples", func() {
	Context("when running the main README example", func() {
		It("should render the expected content", func() {
			template, err := gonja.FromString("Hello {{ name | capitalize }}!")
			Expect(err).To(BeNil())

			out, err := template.ExecuteToString(exec.NewContext(map[string]interface{}{"name": "bob"}))
			Expect(err).To(BeNil())

			Expect(out).To(Equal("Hello Bob!"))
		})
	})
})
