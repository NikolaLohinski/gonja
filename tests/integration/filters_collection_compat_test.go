package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("collection filter compatibility", func() {
	testCases := []struct {
		name     string
		template string
		context  map[string]any
		want     string
	}{
		{
			name:     "items filter returns tuple values",
			template: `{% for item in values|items|sort(attribute=0) %}{{ item[0] }}={{ item[1] }};{% endfor %}`,
			context: map[string]any{
				"values": map[string]any{"b": "two", "a": "one"},
			},
			want: "a=one;b=two;",
		},
		{
			name:     "dict items method returns ordered tuples",
			template: `{% for item in values.items() %}{{ item[0] }}={{ item[1] }};{% endfor %}`,
			context: map[string]any{
				"values": map[string]any{"b": 2, "a": 1},
			},
			want: "a=1;b=2;",
		},
		{
			name:     "join resolves attributes and escapes separator",
			template: `{% autoescape true %}{{ users|join("<hr>", attribute="name") }}{% endautoescape %}`,
			context: map[string]any{
				"users": []any{
					map[string]any{"name": "<b>Alice</b>"},
					map[string]any{"name": "Bob"},
				},
			},
			want: "&lt;b&gt;Alice&lt;/b&gt;&lt;hr&gt;Bob",
		},
		{
			name:     "groupby supports tuple indexes",
			template: `{% for group in entries|groupby(0) %}{{ group.grouper }}={{ group.list|sum(attribute=1) }};{% endfor %}`,
			context: map[string]any{
				"entries": []any{
					[]any{"a", 1},
					[]any{"a", 2},
					[]any{"b", 3},
				},
			},
			want: "a=3;b=3;",
		},
		{
			name:     "groupby supports defaults and case insensitive grouping",
			template: `{% for group in users|groupby("city", default="Unknown", case_sensitive=false) %}{{ group.grouper }}={{ group.list|length }};{% endfor %}`,
			context: map[string]any{
				"users": []any{
					map[string]any{"name": "Alice", "city": "ca"},
					map[string]any{"name": "Bob", "city": "CA"},
					map[string]any{"name": "Cara"},
				},
			},
			want: "ca=2;Unknown=1;",
		},
		{
			name:     "map resolves nested attributes with defaults",
			template: `{{ users|map(attribute="profile.name", default="Anonymous")|join(", ") }}`,
			context: map[string]any{
				"users": []any{
					map[string]any{"profile": map[string]any{"name": "Alice"}},
					map[string]any{},
					map[string]any{"profile": map[string]any{"name": "Bob"}},
				},
			},
			want: "Alice, Anonymous, Bob",
		},
		{
			name:     "map forwards filter arguments",
			template: `{{ names|map("replace", "a", "o")|join(", ") }}`,
			context: map[string]any{
				"names": []any{"bar", "baz"},
			},
			want: "bor, boz",
		},
		{
			name:     "selectattr and rejectattr handle missing nested attributes",
			template: `{{ users|selectattr("profile.active")|map(attribute="name")|join(",") }}|{{ users|rejectattr("profile.active")|map(attribute="name")|join(",") }}`,
			context: map[string]any{
				"users": []any{
					map[string]any{"name": "Alice", "profile": map[string]any{"active": true}},
					map[string]any{"name": "Bob"},
					map[string]any{"name": "Cara", "profile": map[string]any{"active": false}},
				},
			},
			want: "Alice|Bob,Cara",
		},
		{
			name:     "sort supports multiple attributes",
			template: `{{ users|sort(attribute="meta.rank,name")|map(attribute="name")|join(",") }}`,
			context: map[string]any{
				"users": []any{
					map[string]any{"name": "Bob", "meta": map[string]any{"rank": 2}},
					map[string]any{"name": "Alice", "meta": map[string]any{"rank": 1}},
					map[string]any{"name": "Aaron", "meta": map[string]any{"rank": 1}},
				},
			},
			want: "Aaron,Alice,Bob",
		},
		{
			name:     "sum reads tuple indexes",
			template: `{{ entries|sum(attribute=1) }}`,
			context: map[string]any{
				"entries": []any{
					[]any{"a", 1},
					[]any{"b", 2},
					[]any{"c", 3},
				},
			},
			want: "6",
		},
		{
			name:     "unique skips items missing the selected attribute",
			template: `{{ users|unique(attribute="meta.city")|map(attribute="name")|join(",") }}`,
			context: map[string]any{
				"users": []any{
					map[string]any{"name": "Alice", "meta": map[string]any{"city": "SF"}},
					map[string]any{"name": "Bob"},
					map[string]any{"name": "Cara", "meta": map[string]any{"city": "SF"}},
					map[string]any{"name": "Dan", "meta": map[string]any{"city": "NY"}},
				},
			},
			want: "Alice,Dan",
		},
		{
			name:     "slice follows jinja column distribution",
			template: `{{ values|slice(3, fill_with="x") }}`,
			context: map[string]any{
				"values": []any{1, 2, 3, 4, 5, 6, 7},
			},
			want: "[[1, 2, 3], [4, 5, 'x'], [6, 7, 'x']]",
		},
		{
			name:     "float parses strings and respects defaults",
			template: `{{ "32.32"|float }}|{{ "abc"|float(default=1.5) }}`,
			want:     "32.32|1.5",
		},
	}

	for _, tc := range testCases {
		testCase := tc
		It(testCase.name, func() {
			Expect(renderTemplate(testCase.template, testCase.context)).To(Equal(testCase.want))
		})
	}
})
