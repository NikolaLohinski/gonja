package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("containment over values supplied by the host program", func() {
	testCases := []struct {
		name     string
		template string
		context  map[string]any
		want     string
	}{
		{
			name:     "a mapping element is found in a slice of mappings",
			template: `{{ needle in rows }}`,
			context: map[string]any{
				"needle": map[string]any{"x": 2},
				"rows":   []any{map[string]any{"x": 1}, map[string]any{"x": 2}},
			},
			want: "True",
		},
		{
			name:     "a mapping element that is absent reports false",
			template: `{{ needle in rows }}`,
			context: map[string]any{
				"needle": map[string]any{"x": 9},
				"rows":   []any{map[string]any{"x": 1}, map[string]any{"x": 2}},
			},
			want: "False",
		},
		{
			name:     "the in test resolves mapping elements too",
			template: `{{ needle is in(rows) }}`,
			context: map[string]any{
				"needle": map[string]any{"x": 2},
				"rows":   []any{map[string]any{"x": 1}, map[string]any{"x": 2}},
			},
			want: "True",
		},
		{
			name:     "a slice element is found in a slice of slices",
			template: `{{ needle in rows }}`,
			context: map[string]any{
				"needle": []any{1},
				"rows":   []any{[]any{1}, []any{2}},
			},
			want: "True",
		},
		{
			name:     "a slice from the context answers as a list literal does",
			template: `{{ 1 in numbers }}|{{ 1 in [1.0, 2.0] }}`,
			context:  map[string]any{"numbers": []any{1.0, 2.0}},
			want:     "True|True",
		},
		{
			name:     "scalar and string containment are unchanged",
			template: `{{ 1 in [1, 2] }}|{{ "z" in words }}|{{ "ell" in "hello" }}|{{ "a" in mapping }}`,
			context: map[string]any{
				"words":   []any{"a", "b"},
				"mapping": map[string]any{"a": 1},
			},
			want: "True|False|True|True",
		},
	}

	for _, tc := range testCases {
		testCase := tc
		It(testCase.name, func() {
			Expect(renderTemplate(testCase.template, testCase.context)).To(Equal(testCase.want))
		})
	}
})
