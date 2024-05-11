package integration_test

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("control structure 'block'", func() {
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
	Context("when defining two blocks and accessing each one using the 'self' object", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					{% block first -%}
					first block
					{%- endblock first %}
					{% block second -%}
					second block
					{%- endblock second %}
					self {{ self.first() }}
					self {{ self.second() }}
				`),
			})
			(*environment).Context.Set("value", map[string]interface{}{"exists": "content"})
		})

		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("returning the expected result")
			expected := heredoc.Doc(`
				first block
				second block
				self first block
				self second block
			`)
			AssertPrettyDiff(expected, *returnedResult)
		})
	})
	Context("when reusing the self block multiple times", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					{% block block -%}
					reused content
					{%- endblock block %}

					some content in between

					{{ self.block() }}
					{{ self.block() }}
				`),
			})
		})

		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("returning the expected result")
			expected := heredoc.Doc(`
				reused content

				some content in between

				reused content
				reused content
			`)
			AssertPrettyDiff(expected, *returnedResult)
		})
	})

})
