package integration_test

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("expressions", func() {
	var (
		identifier = new(string)

		environment   = new(*exec.Environment)
		loader        = new(loaders.Loader)
		configuration = new(*config.Config)

		context = new(*exec.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
	)
	BeforeEach(func() {
		*identifier = "/test"
		*environment = gonja.DefaultEnvironment
		*loader = loaders.MustNewMemoryLoader(nil)
		*configuration = config.New()
	})
	JustBeforeEach(func() {
		var t *exec.Template
		t, *returnedErr = exec.NewTemplate(*identifier, *configuration, *loader, *environment)
		if *returnedErr != nil {
			return
		}
		*returnedResult, *returnedErr = t.ExecuteToString(*context)
	})
	Context("https://github.com/NikolaLohinski/gonja/issues/49", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					Output 1:
					{%- if not(variable) %}
					    Negated
					{%- else %}
					    Original
					{%- endif %}

					Output 2:
					{%- if not (variable) %}
					    Negated
					{%- else %}
					    Original
					{%- endif %}
				`),
			})
			(*environment).Context.Set("variable", false)
		})

		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("returning the expected result")
			expected := heredoc.Doc(`
				Output 1:
				    Negated

				Output 2:
				    Negated`)
			AssertPrettyDiff(expected, *returnedResult)
		})
	})
	Context("https://github.com/NikolaLohinski/gonja/issues/86", func() {
		// Python/Jinja2 `or` and `and` are value-preserving short-circuit
		// operators: they return one of the operands, not a coerced bool.
		Context("`or` returns the first truthy operand", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: `{{ 'first' or 'fallback' }}`,
				})
			})
			It("renders 'first'", func() {
				Expect(*returnedErr).To(BeNil())
				AssertPrettyDiff("first", *returnedResult)
			})
		})
		Context("`or` falls through when the left side is empty", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: `{{ '' or 'fallback' }}`,
				})
			})
			It("renders 'fallback'", func() {
				Expect(*returnedErr).To(BeNil())
				AssertPrettyDiff("fallback", *returnedResult)
			})
		})
		Context("`or` falls through when the left side is undefined", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: `{{ missing or 'fallback' }}`,
				})
			})
			It("renders 'fallback'", func() {
				Expect(*returnedErr).To(BeNil())
				AssertPrettyDiff("fallback", *returnedResult)
			})
		})
		Context("`and` returns the last operand when both are truthy", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: `{{ 'first' and 'last' }}`,
				})
			})
			It("renders 'last'", func() {
				Expect(*returnedErr).To(BeNil())
				AssertPrettyDiff("last", *returnedResult)
			})
		})
		Context("`and` short-circuits on the first falsy operand", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: `{{ '' and 'never' }}`,
				})
			})
			It("renders an empty string", func() {
				Expect(*returnedErr).To(BeNil())
				AssertPrettyDiff("", *returnedResult)
			})
		})
		Context("`or` chains preserve the first truthy value", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: `{{ '' or 0 or 'winner' or 'never' }}`,
				})
			})
			It("renders 'winner'", func() {
				Expect(*returnedErr).To(BeNil())
				AssertPrettyDiff("winner", *returnedResult)
			})
		})
	})
	Context("https://github.com/NikolaLohinski/gonja/issues/40", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					{%- set output = "foo" if variable else bar -%}
					{{- output -}}
				`),
			})
			(*environment).Context.Set("variable", true)
		})

		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("returning the expected result")
			expected := "foo"
			AssertPrettyDiff(expected, *returnedResult)
		})
	})
	Context("`{% set v = X if cond else Y %}` honors Pythonic truthiness on cond", func() {
		// Regression: SetControlStructure.Execute previously used
		// Value.Bool() which returns true only for reflect.Bool kinds,
		// silently routing every non-bool-typed condition (int 1,
		// non-empty string, populated list/map, struct) to the else
		// branch. Now uses IsTrue() to match the {{ X if c else Y }}
		// renderer behavior.
		DescribeTable(
			"set + ternary with truthy non-bool conditions resolves the X branch",
			func(template, expected string, ctxKV map[string]interface{}) {
				*loader = loaders.MustNewMemoryLoader(map[string]string{*identifier: template})
				for k, v := range ctxKV {
					(*environment).Context.Set(k, v)
				}
				t, err := exec.NewTemplate(*identifier, *configuration, *loader, *environment)
				Expect(err).To(BeNil())
				out, err := t.ExecuteToString(*context)
				Expect(err).To(BeNil())
				AssertPrettyDiff(expected, out)
			},
			Entry("int 1 condition", `{% set v = 'X' if 1 else 'Y' %}{{ v }}`, "X", nil),
			Entry("string condition", `{% set v = 'X' if 'truthy' else 'Y' %}{{ v }}`, "X", nil),
			Entry("list condition", `{% set v = 'X' if [1] else 'Y' %}{{ v }}`, "X", nil),
			Entry("`and` chain returning string", `{% set v = 'X' if ('hello' and 'world') else 'Y' %}{{ v }}`, "X", nil),
			Entry("attribute access X branch", `{% set v = a.b if 1 else 'Y' %}{{ v }}`, "Z",
				map[string]interface{}{"a": map[string]interface{}{"b": "Z"}}),
			Entry("falsy int 0 still goes else", `{% set v = 'X' if 0 else 'Y' %}{{ v }}`, "Y", nil),
			Entry("empty string still goes else", `{% set v = 'X' if '' else 'Y' %}{{ v }}`, "Y", nil),
			Entry("empty list still goes else", `{% set v = 'X' if [] else 'Y' %}{{ v }}`, "Y", nil),
		)
	})
})
