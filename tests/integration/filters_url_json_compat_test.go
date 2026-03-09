package integration_test

import "testing"

type compatJSONMarshaler struct{}

func (compatJSONMarshaler) MarshalJSON() ([]byte, error) {
	return []byte("42"), nil
}

func TestURLAndJSONFilterCompatibility(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		context  map[string]any
		want     string
	}{
		{
			name:     "tojson preserves compact spacing and escapes apostrophes",
			template: `{{ value|tojson }}`,
			context: map[string]any{
				"value": map[string]any{"a": "b's"},
			},
			want: `{"a": "b\u0027s"}`,
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
		t.Run(tc.name, func(t *testing.T) {
			if got := renderTemplate(t, tc.template, tc.context); got != tc.want {
				t.Fatalf("rendered output mismatch\nwant: %q\ngot:  %q", tc.want, got)
			}
		})
	}
}
