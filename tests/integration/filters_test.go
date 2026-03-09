package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("filters", func() {
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
	Context("indent", func() {
		shouldRender("{{ 'jinja' | indent }}", "jinja")
		shouldRender("{{ '\nfoo\nbar' | indent }}", "\n    foo\n    bar")
		shouldRender(`{{ '\nfoo bar\n"baz"\n' | indent(2) }}`, "\n  foo bar\n  \"baz\"\n")
		shouldRender(`{{ '\nfoo bar\n"baz"\n' | indent(2, true) }}`, "  \n  foo bar\n  \"baz\"\n")
		shouldRender(`{{ 'jinja\nflask' | indent(width='>>> ', first=True) }}`, ">>> jinja\n>>> flask")
		shouldRender(`{% autoescape true %}{{ '\n<b>foo</b>\n<i>bar</i>\n' | safe | indent(2, true) }}{% endautoescape %}`, "  \n  <b>foo</b>\n  <i>bar</i>\n")
		shouldFail("{{ True | indent }}", "invalid call to filter 'indent': True is not a string")
		shouldFail("{{ 'jinja' | indent(width=True) }}", "invalid call to filter 'indent': failed to validate argument 'width': True is neither a string nor an integer")
	})
	Context("slice", func() {
		shouldRender("{{ [1, 2, 3, 4, 5, 6] | slice(2) }}", "[[1, 2, 3], [4, 5, 6]]")
		shouldRender("{{ [1, 2, 3, 4, 5] | slice(3) }}", "[[1, 2], [3, 4], [5]]")
		shouldRender("{{ [1, 2, 3, 4, 5] | slice(3, 42) }}", "[[1, 2], [3, 4], [5, 42]]")
		shouldRender("{{ [1, 2, 3, 4, 5, 6, 7] | slice(3, fill_with='this') }}", "[[1, 2, 3], [4, 5, 'this'], [6, 7, 'this']]")
		shouldFail("{{ True | slice(42) }}", "invalid call to filter 'slice': True is not a list")
		shouldFail("{{ True | slice('yolo') }}", "invalid call to filter 'slice': failed to validate argument 'slices': yolo is not an integer")
		shouldFail("{{ True | slice(-32) }}", "invalid call to filter 'slice': slices argument -32 must be > 0")
	})
	Context("default", func() {
		shouldRender(`{{ undefined_var | default("default_value") }}`, "default_value")
		shouldRender(`{{ "" | default("default_value", true) }}`, "default_value")
		shouldRender(`{{ "is_true" | default("default_value", true) }}`, "is_true")
	})
})
