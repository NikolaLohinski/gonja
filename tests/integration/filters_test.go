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
		shouldRender("{{ '\nfoo\nbar' | indent }}", "\n    foo\n    bar\n")
		shouldFail("{{ True | indent }}", "invalid call to filter 'indent': True is not a string")
		shouldFail("{{ True | indent(width='yolo') }}", "invalid call to filter 'indent': failed to validate argument 'width': yolo is not an integer")
	})
	Context("slice", func() {
		shouldRender("{{ [1, 2, 3, 4, 5, 6] | slice(2) }}", "[[1, 2, 3], [4, 5, 6]]")
		shouldRender("{{ [1, 2, 3, 4, 5] | slice(3) }}", "[[1, 2], [3, 4], [5]]")
		shouldRender("{{ [1, 2, 3, 4, 5] | slice(3, 42) }}", "[[1, 2], [3, 4], [5, 42]]")
		shouldRender("{{ [1, 2, 3, 4, 5, 6, 7] | slice(3, fill_with='this') }}", "[[1, 2, 3], [4, 5, 6], [7, 'this', 'this']]")
		shouldFail("{{ True | slice(42) }}", "invalid call to filter 'slice': True is not a list")
		shouldFail("{{ True | slice('yolo') }}", "invalid call to filter 'slice': failed to validate argument 'slices': yolo is not an integer")
		shouldFail("{{ True | slice(-32) }}", "invalid call to filter 'slice': slices argument -32 must be > 0")
	})
})
