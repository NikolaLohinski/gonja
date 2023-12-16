package integration_test

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja"
	"github.com/nikolalohinski/gonja/config"
	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("lists", func() {
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
	Context("when getting a slice using the '[...]' syntax", func() {
		Context("default", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					[:]:   {{ value[:]   }}
					[2:]:  {{ value[2:]  }}
					[:3]:  {{ value[:3]  }}
					[:-2]: {{ value[:-2] }}
					[-4:]: {{ value[-4:] }}
				`),
				})
				(*environment).Context.Set("value", []interface{}{"1", 2, 3, 4, "five"})
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("not returning the expected result")
				expected := heredoc.Doc(`
					[:]:   ['1', 2, 3, 4, 'five']
					[2:]:  [3, 4, 'five']
					[:3]:  ['1', 2, 3]
					[:-2]: ['1', 2, 3]
					[-4:]: [2, 3, 4, 'five']
				`)
				AssertPrettyDiff(expected, *returnedResult)
			})
		})

	})
})
