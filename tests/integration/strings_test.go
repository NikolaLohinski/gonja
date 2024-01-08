package integration_test

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("strings", func() {
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
		*returnedResult, *returnedErr = t.Execute(*context)
	})
	Context("when getting a character by index", func() {
		Context("default", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					[]:    {{ value[]    }}
					[4]:   {{ value[4]   }}
					[-3]:  {{ value[-3]  }}
					[256]: {{ value[256] }}
					[-99]: {{ value[-99] }}
				`),
				})
				(*environment).Context.Set("value", "testing is fun!")
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				expected := heredoc.Doc(`
					[]:    
					[4]:   i
					[-3]:  u
					[256]: 
					[-99]: 
				`)
				AssertPrettyDiff(expected, *returnedResult)
			})
		})
	})
	Context("when getting a substring with '[...]' notation", func() {
		Context("default", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					[]:    {{ value[]    }}
					[:]:   {{ value[:]   }}
					[2:]:  {{ value[2:]  }}
					[:3]:  {{ value[:3]  }}
					[:-2]: {{ value[:-2] }}
					[-5:]: {{ value[-5:] }}
				`),
				})
				(*environment).Context.Set("value", "testing is fun!")
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				expected := heredoc.Doc(`
					[]:    
					[:]:   testing is fun!
					[2:]:  sting is fun!
					[:3]:  tes
					[:-2]: testing is fu
					[-5:]:  fun!
				`)
				AssertPrettyDiff(expected, *returnedResult)
			})
		})

	})
})
