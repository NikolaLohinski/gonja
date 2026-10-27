package integration_test

import (
	"errors"
	"io/fs"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("error identity", func() {
	var (
		errSentinel = errors.New("sentinel failure")
		environment = new(*exec.Environment)
	)
	BeforeEach(func() {
		*environment = &exec.Environment{
			Context:           exec.EmptyContext().Update(builtins.GlobalFunctions).Update(builtins.GlobalVariables),
			Filters:           exec.NewFilterSet(map[string]exec.FilterFunction{}).Update(builtins.Filters),
			Tests:             exec.NewTestSet(map[string]exec.TestFunction{}).Update(builtins.Tests),
			ControlStructures: builtins.ControlStructures,
			Methods:           builtins.Methods,
		}
		(*environment).Context.Set("fail", func() (any, error) { return nil, errSentinel })
		Expect((*environment).Filters.Register("failing", func(_ *exec.Evaluator, _ *exec.Value, _ *exec.VarArgs) *exec.Value {
			return exec.AsValue(errSentinel)
		})).To(Succeed())
		Expect((*environment).Tests.Register("failing", func(_ *exec.Evaluator, _ *exec.Value, _ *exec.VarArgs) (bool, error) {
			return false, errSentinel
		})).To(Succeed())
	})
	render := func(source string) error {
		loader := loaders.MustNewMemoryLoader(map[string]string{"/test": source})
		template, err := exec.NewTemplate("/test", gonja.DefaultConfig, loader, *environment)
		if err != nil {
			return err
		}
		_, err = template.ExecuteToString(exec.EmptyContext())
		return err
	}
	for _, testCase := range []struct {
		description, source, message string
	}{
		{"function", "{{ fail() }}", "invalid call to function 'fail': sentinel failure"},
		{"filter", "{{ 1 | failing }}", "invalid call to filter 'failing': sentinel failure"},
		{"test", "{{ 1 is failing }}", "invalid call to test 'failing': sentinel failure"},
	} {
		testCase := testCase
		Context(testCase.description, func() {
			It("should keep the message and the original error", func() {
				err := render(testCase.source)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(testCase.message))
				Expect(errors.Is(err, errSentinel)).To(BeTrue())
			})
		})
	}
	Context("missing template", func() {
		It("should keep the loader error", func() {
			loader := loaders.MustNewFileSystemLoader(GinkgoT().TempDir())
			_, err := exec.NewTemplate("missing.j2", gonja.DefaultConfig, loader, *environment)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, fs.ErrNotExist)).To(BeTrue())
		})
	})
})
