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

var _ = Context("dicts", func() {
	var (
		identifier = new(string)

		environment = new(*exec.Environment)
		config      = new(*config.Config)
		loader      = new(loaders.Loader)

		context = new(*exec.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
	)
	BeforeEach(func() {
		*identifier = "/test"
		*environment = gonja.DefaultEnvironment
		*config = gonja.DefaultConfig
		*loader = loaders.MustNewMemoryLoader(nil)
	})
	JustBeforeEach(func() {
		var t *exec.Template
		t, *returnedErr = exec.NewTemplate(*identifier, *config, *loader, *environment)
		if *returnedErr != nil {
			return
		}
		*returnedResult, *returnedErr = t.Execute(*context)
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
				By("not returning the expected result")
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
})
