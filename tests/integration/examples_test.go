package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("examples", func() {
	Context("when running the main README example", func() {
		It("should render the expected content", func() {
			loader, err := loaders.NewMemoryLoader(map[string]string{"/_": "Hello {{ name | capitalize }}!"})
			Expect(err).To(BeNil())

			template, err := exec.NewTemplate("/_", gonja.DefaultConfig, loader, gonja.DefaultEnvironment)
			Expect(err).To(BeNil())

			out, err := template.Execute(exec.NewContext(map[string]interface{}{"name": "bob"}))
			Expect(err).To(BeNil())

			Expect(out).To(Equal("Hello Bob!"))
		})
	})
})
