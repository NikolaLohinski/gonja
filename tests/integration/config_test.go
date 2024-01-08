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

var _ = Context("config", func() {
	var (
		identifier = new(string)

		environment   = new(*exec.Environment)
		configuration = new(*config.Config)
		loader        = new(loaders.Loader)

		context = new(*exec.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
	)
	BeforeEach(func() {
		*identifier = "/test"
		*environment = gonja.DefaultEnvironment
		*configuration = config.New()
		*loader = loaders.MustNewMemoryLoader(nil)
	})
	JustBeforeEach(func() {
		var t *exec.Template
		t, *returnedErr = exec.NewTemplate(*identifier, *configuration, *loader, *environment)
		if *returnedErr != nil {
			return
		}
		*returnedResult, *returnedErr = t.Execute(*context)
	})
	Context("when toggling Config.StrictUndefined behavior", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: "Accessing data.nope: '{{ data.nope }}'",
			})
			(*environment).Context.Set("data", map[string]interface{}{})
		})
		Context("when Config.StrictUndefined = false", func() {
			BeforeEach(func() {
				(*configuration).StrictUndefined = false
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff("Accessing data.nope: ''", *returnedResult)
			})
		})
		Context("when Config.StrictUndefined = true", func() {
			BeforeEach(func() {
				(*configuration).StrictUndefined = true
			})
			It("should failed to render", func() {
				By("not returning an error")
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("when changing delimiters", func() {
		BeforeEach(func() {
			(*configuration).BlockStartString = "[%"
			(*configuration).BlockEndString = "%]"
			(*configuration).VariableStartString = "<<"
			(*configuration).VariableEndString = ">>"
			(*configuration).CommentStartString = "|#"
			(*configuration).CommentEndString = "#|"
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					[%- if "foo" in "foo bar" %]
					I am cornered
					[%- endif %]
					<< "but pointy" >>
					|# "and can be invisible!" #|
				`),
			})
		})

		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("returning the expected result")
			expected := heredoc.Doc(`

				I am cornered
				but pointy

			`)
			AssertPrettyDiff(expected, *returnedResult)
		})
	})
	Context("when toggling Config.AutoEscape behavior", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: "{{ data }}",
			})
			(*environment).Context.Set("data", "<a>test</a>")
		})
		Context("when Config.AutoEscape = false", func() {
			BeforeEach(func() {
				(*configuration).AutoEscape = false
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff("<a>test</a>", *returnedResult)
			})
		})
		Context("when Config.AutoEscape = true", func() {
			BeforeEach(func() {
				(*configuration).AutoEscape = true
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff("&lt;a&gt;test&lt;/a&gt;", *returnedResult)
			})
		})
	})
	Context("when toggling Config.TrimBlocks behavior", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					Some text
					{%- set block_example = "test" %}

					{{ "The empty line should have been removed" }}
					
					The empty line above should stay
				`),
			})
		})
		Context("when Config.TrimBlock = false", func() {
			BeforeEach(func() {
				(*configuration).TrimBlocks = false
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff(heredoc.Doc(`
					Some text

					The empty line should have been removed

					The empty line above should stay
				`), *returnedResult)
			})
		})
		Context("when Config.TrimBlocks = true", func() {
			BeforeEach(func() {
				(*configuration).TrimBlocks = true
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff(heredoc.Doc(`
					Some text
					The empty line should have been removed

					The empty line above should stay
				`), *returnedResult)
			})
		})
	})
	Context("when toggling Config.LeftStripBlocks behavior", func() {
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					  	{% set _ = "" %}block indented with spaces and tabs
					-
					  {{ "variable indented with spaces" }}
				`),
			})
		})
		Context("when Config.LeftStripBlocks = false", func() {
			BeforeEach(func() {
				(*configuration).LeftStripBlocks = false
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff(heredoc.Doc(`
					  	block indented with spaces and tabs
					-
					  variable indented with spaces
				`), *returnedResult)
			})
		})
		Context("when Config.LeftStripBlocks = true", func() {
			BeforeEach(func() {
				(*configuration).LeftStripBlocks = true
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff(heredoc.Doc(`
					block indented with spaces and tabs
					-
					  variable indented with spaces
				`), *returnedResult)
			})
		})
	})
})
