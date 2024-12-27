package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("tests", func() {
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
	Context("greaterthan", func() {
		shouldRender("{{ 42 is greaterthan(31) }}", "True")
		shouldRender("{{ 42 is greaterthan 31 }}", "True")
		shouldRender("{{ 42 is gt 31 }}", "True")
		shouldRender("{{ 42 is > 31 }}", "True")
		shouldFail("{{ 42 is greaterthan(True) }}", "True is not a number")
	})
	Context("https://github.com/NikolaLohinski/gonja/issues/19", func() {
		BeforeEach(func() {
			*context = exec.NewContext(map[string]interface{}{
				"var1": "1",
				"var2": "3",
			})
		})
		shouldRender(
			"var1 in ['1', '2'] or (var2 == '3'): {% if var1 in ['1', '2'] or (var2 == '3') %}ok{% endif %}",
			"var1 in ['1', '2'] or (var2 == '3'): ok",
		)
		shouldRender(
			"var1 in (['1'] + ['2']) or (var2 == '3'): {% if var1 in (['1'] + ['2']) or (var2 == '3') %}ok{% endif %}",
			"var1 in (['1'] + ['2']) or (var2 == '3'): ok",
		)
		shouldRender(
			"(var1 in ['1', '2']) or var2 == '3': {% if (var1 in ['1', '2']) or var2 == '3' %}ok{% endif %}",
			"(var1 in ['1', '2']) or var2 == '3': ok",
		)
		shouldRender(
			"(var1 in ['1', '2']) or (var2 == '3'): {% if (var1 in ['1', '2']) or (var2 == '3') %}ok{% endif %}",
			"(var1 in ['1', '2']) or (var2 == '3'): ok",
		)
		shouldRender(
			"var1 in ['1', '2'] or var2 == '3': {% if var1 in ['1', '2'] or var2 == '3' %}ok{% endif %}",
			"var1 in ['1', '2'] or var2 == '3': ok",
		)
	})
	Context("https://github.com/NikolaLohinski/gonja/issues/30", func() {
		shouldRender(
			`{{ "a" < "b" }}`,
			`True`,
		)
		shouldRender(
			`{{ "b" > "a" }}`,
			`True`,
		)
		shouldRender(
			`{{ "a" > "b" }}`,
			`False`,
		)
		shouldRender(
			`{{ "a" == "b" }}`,
			`False`,
		)
		shouldRender(
			`{{ "a" == "a" }}`,
			`True`,
		)
		shouldRender(
			`{{ "b" >= "a" }}`,
			`True`,
		)
		shouldRender(
			`{{ "a" >= "b" }}`,
			`False`,
		)
		shouldRender(
			`{{ "b" <= "a" }}`,
			`False`,
		)
		shouldRender(
			`{{ "a" <= "b" }}`,
			`True`,
		)
	})
})
