package integration_test

import (
	"os"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("control structure 'include'", func() {
	var (
		identifier = new(string)

		environment = new(*exec.Environment)
		loader      = new(loaders.Loader)

		context = new(*exec.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
	)
	BeforeEach(func() {
		*identifier = "/test"
		*environment = gonja.DefaultEnvironment
		*loader = loaders.MustNewMemoryLoader(nil)
	})
	JustBeforeEach(func() {
		var t *exec.Template
		t, *returnedErr = exec.NewTemplate(*identifier, gonja.DefaultConfig, *loader, *environment)
		if *returnedErr != nil {
			return
		}
		*returnedResult, *returnedErr = t.ExecuteToString(*context)
	})
	BeforeEach(func() {
		*loader = loaders.MustNewMemoryLoader(map[string]string{
			*identifier:          "{% include '/included/template' %}",
			"/included/template": "included content",
		})
	})

	It("should return the expected rendered content", func() {
		By("not returning any error")
		Expect(*returnedErr).To(BeNil())
		By("returning the expected result")
		AssertPrettyDiff("included content", *returnedResult)
	})

	Context("when the include statement has `ignore missing` defined", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier:          "{% include '/included/template' ignore missing %}",
				"/included/template": "included content",
			})
		})

		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("returning the expected result")
			AssertPrettyDiff("included content", *returnedResult)
		})

		Context("and included template is missing", func() {
			var (
				file *os.File
			)
			BeforeEach(func() {
				directory := os.TempDir()
				file = MustReturn(os.CreateTemp(directory, ""))
				MustReturn(file.WriteString("root content{% include '/tmp/does/not/exist' ignore missing %}"))
				*identifier = file.Name()
				*loader = loaders.MustNewFileSystemLoader(directory)
			})
			AfterEach(func() {
				os.RemoveAll(file.Name())
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff("root content", *returnedResult)
			})
		})
	})
})
