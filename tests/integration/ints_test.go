package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("ints", func() {
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
	Context("when using native python methods", func() {
		var (
			shouldRender = func(template, result string) {
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
		Context("is_integer", func() {
			shouldRender("{{ 42.is_integer() }}", "True")
			shouldFail("{{ 42.is_integer('nope') }}", "wrong signature for '42.is_integer': received 1 unexpected positional argument")
		})
	})
})
