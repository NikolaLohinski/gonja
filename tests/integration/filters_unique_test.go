package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// renderTemplateError renders a template and returns its execution error, if any.
func renderTemplateError(source string, data map[string]any) (string, error) {
	template, err := gonja.FromString(source)
	ExpectWithOffset(1, err).To(BeNil(), "parse template")

	var context *exec.Context
	if data != nil {
		context = exec.NewContext(data)
	}

	return template.ExecuteToString(context)
}

var _ = Context("the unique filter over unhashable elements", func() {
	hashable := []struct {
		name     string
		template string
		context  map[string]any
		want     string
	}{
		{
			name:     "scalars are still de-duplicated",
			template: `{{ values|unique }}`,
			context:  map[string]any{"values": []any{3, 1, 2, 1}},
			want:     "[3, 1, 2]",
		},
		{
			name:     "strings are still de-duplicated case-insensitively by default",
			template: `{{ values|unique }}`,
			context:  map[string]any{"values": []any{"B", "b", "a"}},
			want:     "['B', 'a']",
		},
		{
			name:     "booleans are still de-duplicated",
			template: `{{ [true, false, true]|unique }}`,
			want:     "[True, False]",
		},
		{
			name:     "an attribute that resolves to a scalar still de-duplicates the mappings",
			template: `{{ rows|unique(attribute='x') }}`,
			context: map[string]any{
				"rows": []any{map[string]any{"x": 1}, map[string]any{"x": 2}},
			},
			want: "[{'x': 1}, {'x': 2}]",
		},
		{
			name:     "an empty sequence stays empty",
			template: `{{ values|unique }}`,
			context:  map[string]any{"values": []any{}},
			want:     "[]",
		},
	}

	for _, tc := range hashable {
		testCase := tc
		It(testCase.name, func() {
			Expect(renderTemplate(testCase.template, testCase.context)).To(Equal(testCase.want))
		})
	}

	unhashable := []struct {
		name     string
		template string
		context  map[string]any
	}{
		{
			name:     "a sequence of mappings reports an error instead of panicking",
			template: `{{ values|unique }}`,
			context: map[string]any{
				"values": []any{map[string]any{"x": 1}, map[string]any{"x": 2}},
			},
		},
		{
			name:     "a sequence of sequences reports an error instead of panicking",
			template: `{{ values|unique }}`,
			context:  map[string]any{"values": []any{[]any{1}, []any{2}}},
		},
		{
			name:     "an attribute resolving to a mapping reports an error instead of panicking",
			template: `{{ values|unique(attribute='x') }}`,
			context: map[string]any{
				"values": []any{map[string]any{"x": map[string]any{"deep": 1}}},
			},
		},
	}

	for _, tc := range unhashable {
		testCase := tc
		It(testCase.name, func() {
			_, err := renderTemplateError(testCase.template, testCase.context)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unhashable type"))
		})
	}
})
