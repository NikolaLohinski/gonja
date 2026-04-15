package integration_test

import (
	"bytes"

	"github.com/ardanlabs/gonja"
	"github.com/ardanlabs/gonja/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("examples", func() {
	Context("when running the main README example", func() {
		It("should render the expected content", func() {
			template, err := gonja.FromString("Hello {{ name | capitalize }}!")
			Expect(err).To(BeNil())

			out := bytes.NewBufferString("")
			err = template.Execute(out, exec.NewContext(map[string]any{"name": "bob"}))
			Expect(err).To(BeNil())

			Expect(out.String()).To(Equal("Hello Bob!"))
		})
	})
})
