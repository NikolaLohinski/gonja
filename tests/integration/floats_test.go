package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("floats", func() {
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
	Context("when defining a float with multiple dots", func() {
		shouldFail("{{ 32.32.32 }}", "two dots in numeric token")
	})
	Context("when using native python methods", func() {
		Context("as_integer_ratio", func() {
			shouldRender("{{ 42.32.as_integer_ratio() }}", "[5956010507197481, 140737488355328]")
			shouldRender("{{ 42.00.as_integer_ratio() }}", "[42, 1]")
			shouldFail("{{ 42.35.as_integer_ratio('nope') }}", "received 1 unexpected positional argument")
		})
		Context("hex", func() {
			shouldRender("{{ 99999999999.00.hex() }}", "0x1.74876e7ff0000p+36")
			shouldRender("{{ 1.0.hex() }}", "0x1.0000000000000p+0")
			shouldRender("{{ 42.32.hex() }}", "0x1.528f5c28f5c29p+5")
			shouldRender("{{ 42.00.hex() }}", "0x1.5000000000000p+5")
			shouldFail("{{ 42.35.hex('nope') }}", "received 1 unexpected positional argument")
		})
		Context("is_integer", func() {
			shouldRender("{{ 42.32.is_integer() }}", "False")
			shouldRender("{{ 42.00.is_integer() }}", "True")
			shouldFail("{{ 42.35.is_integer('nope') }}", "received 1 unexpected positional argument")
		})
	})
})
