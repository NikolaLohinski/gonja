package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("bools", func() {
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
		Context("string", func() {
			shouldRender("{{ True.string() }}", "True")
			shouldRender("{{ False.string() }}", "False")
			shouldFail("{{ False.string('nope') }}", "received 1 unexpected positional argument")
		})
		Context("int", func() {
			shouldRender("{{ True.int() }}", "1")
			shouldRender("{{ False.int() }}", "0")
			shouldFail("{{ False.int('nope') }}", "received 1 unexpected positional argument")
		})
		Context("bit_count", func() {
			shouldRender("{{ True.bit_count() }}", "1")
			shouldRender("{{ False.bit_count() }}", "0")
			shouldFail("{{ True.bit_count('nope') }}", "received 1 unexpected positional argument")
		})
		Context("bit_length", func() {
			shouldRender("{{ True.bit_length() }}", "1")
			shouldRender("{{ False.bit_length() }}", "0")
			shouldFail("{{ False.bit_length('nope') }}", "received 1 unexpected positional argument")
		})
		Context("as_integer_ratio", func() {
			shouldRender("{{ True.as_integer_ratio() }}", "[1, 1]")
			shouldRender("{{ False.as_integer_ratio() }}", "[0, 1]")
			shouldFail("{{ False.as_integer_ratio('nope') }}", "received 1 unexpected positional argument")
		})
	})
})
