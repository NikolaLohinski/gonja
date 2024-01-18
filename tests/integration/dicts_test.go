package integration_test

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("dicts", func() {
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
	Context("when getting an item", func() {
		Context("default", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					[]: {{ value[] }}
					["exists"]: {{ value["exists"] }}
					['exists']: {{ value['exists'] }}
					['nope']: {{ value['nope'] }}
					["nope"]: {{ value["nope"] }}
					["exi" ~ "sts"]: {{ value["exi" ~ "sts"] }}
				`),
				})
				(*environment).Context.Set("value", map[string]interface{}{"exists": "content"})
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				expected := heredoc.Doc(`
					[]: 
					["exists"]: content
					['exists']: content
					['nope']: 
					["nope"]: 
					["exi" ~ "sts"]: content
				`)
				AssertPrettyDiff(expected, *returnedResult)
			})
		})
	})
	Context("when accessing a raw dict", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: `{{ {"foo": "bar"}["foo"] }}`,
			})
		})
		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("returning the expected result")
			AssertPrettyDiff("bar", *returnedResult)
		})
	})

	Context("when doing an invalid access with ..", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: "{{ data..test }}",
			})
			(*environment).Context.Set("data", map[string]interface{}{})
		})

		It("should return the expected error", func() {
			By("returning an error")
			Expect(*returnedErr).ToNot(BeNil())
		})
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
		Context("keys", func() {
			shouldRender("{{ {'foo': 'bar', 'yolo': 1}.keys() }}", "['foo', 'yolo']")
			shouldFail("{{ {}.keys('nope') }}", "wrong signature for '{}.keys': received 1 unexpected positional argument")
		})
	})

})
