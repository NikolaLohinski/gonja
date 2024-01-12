package integration_test

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("lists", func() {
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
	Context("when getting an item by index", func() {
		Context("default", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					[]:    {{ value[]    }}
					[1]:   {{ value[1]   }}
					[-2]:  {{ value[-2]  }}
					[256]: {{ value[256] }}
					[-99]: {{ value[-99] }}
				`),
				})
				(*environment).Context.Set("value", []interface{}{"1", 2, 3, 4, "five"})
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				expected := heredoc.Doc(`
					[]:    
					[1]:   2
					[-2]:  4
					[256]: 
					[-99]: 
				`)
				AssertPrettyDiff(expected, *returnedResult)
			})
		})
	})

	Context("when getting a slice using the '[...]' syntax", func() {
		Context("default", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					[]:    {{ value[]    }}
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
				By("returning the expected result")
				expected := heredoc.Doc(`
					[]:    
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
