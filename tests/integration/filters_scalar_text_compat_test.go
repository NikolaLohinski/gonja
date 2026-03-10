package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type compatIntValue struct{}

func (compatIntValue) Int() int {
	return 42
}

var _ = Context("scalar and text filter compatibility", func() {
	testCases := []struct {
		name     string
		template string
		context  map[string]any
		want     string
	}{
		{
			name:     "int supports large values and custom int methods",
			template: `{{ "12345678901234567890"|int }}|{{ value|int }}`,
			context: map[string]any{
				"value": compatIntValue{},
			},
			want: "12345678901234567890|42",
		},
		{
			name:     "pprint renders simple lists as indented json",
			template: `{{ values|pprint }}`,
			context: map[string]any{
				"values": []any{1, 2, 3},
			},
			want: "[\n  1,\n  2,\n  3\n]",
		},
		{
			name:     "round supports negative precision",
			template: `{{ 1234.567|round(-2, "floor") }}`,
			want:     "1200.0",
		},
		{
			name:     "striptags removes comments and collapses whitespace",
			template: `{{ value|striptags }}`,
			context: map[string]any{
				"value": `Hello <!-- hidden --> <b>world</b>`,
			},
			want: "Hello world",
		},
		{
			name:     "title stringifies non string values",
			template: `{{ 5|title }}`,
			want:     "5",
		},
		{
			name:     "trim with none uses trimspace",
			template: `{{ "  hello  "|trim(chars=None) }}`,
			want:     "hello",
		},
		{
			name:     "truncate keeps values inside default leeway",
			template: `{{ "foo bar baz qux"|truncate(11, False) }}`,
			want:     "foo bar baz qux",
		},
		{
			name:     "wordwrap wraps by character width",
			template: `{{ "foo bar baz"|wordwrap(7) }}`,
			want:     "foo bar\nbaz",
		},
	}

	for _, tc := range testCases {
		testCase := tc
		It(testCase.name, func() {
			Expect(renderTemplate(testCase.template, testCase.context)).To(Equal(testCase.want))
		})
	}
})
