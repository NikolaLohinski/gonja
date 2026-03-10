package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type compatJSONMarshaler struct{}

func (compatJSONMarshaler) MarshalJSON() ([]byte, error) {
	return []byte("42"), nil
}

var _ = Context("url and json filter compatibility", func() {
	testCases := []struct {
		name     string
		template string
		context  map[string]any
		want     string
	}{
		{
			name:     "tojson uses stdlib output and escapes apostrophes",
			template: `{{ value|tojson }}`,
			context: map[string]any{
				"value": map[string]any{"a": "b's"},
			},
			want: `{"a":"b\u0027s"}`,
		},
		{
			name:     "tojson honors custom marshalers",
			template: `{{ value|tojson }}`,
			context: map[string]any{
				"value": compatJSONMarshaler{},
			},
			want: "42",
		},
		{
			name:     "urlencode preserves path separators in strings",
			template: `{{ "a b/c"|urlencode }}`,
			want:     "a%20b/c",
		},
		{
			name:     "urlencode supports mappings and pair lists",
			template: `{{ query|urlencode }}|{{ pairs|urlencode }}`,
			context: map[string]any{
				"query": map[string]any{"f": 1, "z": 2},
				"pairs": []any{
					[]any{"a b/c", "a b/c"},
				},
			},
			want: "f=1&z=2|a+b%2Fc=a+b%2Fc",
		},
		{
			name:     "urlize uses https for bare domains",
			template: `{{ "example.com"|urlize|safe }}`,
			want:     `<a href="https://example.com" rel="noopener">example.com</a>`,
		},
		{
			name:     "urlize supports custom schemes",
			template: `{{ "xmpp:foo@example.com"|urlize(extra_schemes=["xmpp:"])|safe }}`,
			want:     `<a href="xmpp:foo@example.com" rel="noopener">xmpp:foo@example.com</a>`,
		},
		{
			name:     "urlize keeps explicit rel values intact",
			template: `{{ "http://example.com"|urlize(rel="nofollow", target="_blank")|safe }}`,
			want:     `<a href="http://example.com" rel="nofollow" target="_blank">http://example.com</a>`,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		It(testCase.name, func() {
			Expect(renderTemplate(testCase.template, testCase.context)).To(Equal(testCase.want))
		})
	}
})
