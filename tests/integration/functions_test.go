package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("functions", func() {
	var (
		identifier = new(string)

		environment = new(*exec.Environment)
		loader      = new(loaders.Loader)

		context = new(*exec.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
		shouldRender   = func(template, result string) {
			Context(template, func() {
				BeforeEach(func() {
					*loader = loaders.MustNewMemoryLoader(map[string]string{
						*identifier: template,
					})
				})
				It("should return the expected rendered content", func() {
					By("not returning any error")
					Expect(*returnedErr).To(BeNil())
					By("returning the expected result")
					AssertPrettyDiff(result, *returnedResult)
				})
			})
		}
		shouldFail = func(template, err string) {
			Context(template, func() {
				BeforeEach(func() {
					*loader = loaders.MustNewMemoryLoader(map[string]string{
						*identifier: template,
					})
				})
				It("should return the expected error", func() {
					Expect(*returnedErr).ToNot(BeNil())
					Expect((*returnedErr).Error()).To(MatchRegexp(err))
				})
			})
		}
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
	Context("joiner", func() {
		shouldRender(`{% set pipe = joiner("|") -%}{% for i in [0, 1, 2] %}{{ pipe() }}{{ i }}{% endfor %}`, "0|1|2")
		shouldFail("{% set pipe = joiner(True) -%}", "invalid call to function 'joiner': failed to validate argument 'sep': True is not a string")
	})
	Context("range", func() {
		shouldRender(`{% for i in range(10) %}{{ i }}{% endfor %}`, "0123456789")
		shouldRender(`{% for i in range(1, 10, 2) %}{{ i }}{% endfor %}`, "13579")
		shouldRender(`{% for i in range(10, 1, -1) %}{{ i }}{% endfor %}`, "1098765432")
		shouldRender(`{% for i in range(10, 1, -2) %}{{ i }}{% endfor %}`, "108642")
		shouldFail("{% set invalid = range(True) -%}", "invalid call to function 'range': expected signature is \\[start, ]stop\\[, step] where all arguments are integers")
	})
})
