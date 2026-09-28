package integration_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func TestSequenceIncludesStringsListsArraysAndMappings(t *testing.T) {
	ordered := exec.NewDict()
	ordered.Pairs = append(ordered.Pairs,
		&exec.Pair{Key: exec.AsValue("z"), Value: exec.AsValue(1)},
		&exec.Pair{Key: exec.AsValue("a"), Value: exec.AsValue(2)},
	)

	cases := []struct {
		name  string
		value any
		want  string
	}{
		{name: "string", value: "text", want: "true"},
		{name: "list", value: []string{"text"}, want: "true"},
		{name: "array", value: [1]string{"text"}, want: "true"},
		{name: "map", value: map[string]any{"text": "value"}, want: "true"},
		{name: "ordered_exec_dict", value: ordered, want: "true"},
		{name: "none", value: nil, want: "false"},
		{name: "number", value: 42, want: "false"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := render(t, config.New(),
				`{% if value is sequence %}true{% else %}false{% endif %}`,
				map[string]any{"value": tc.value})
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("sequence test returned %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMappingToolBodyIsRejectedAsMalformedContentParts(t *testing.T) {
	content := exec.NewDict()
	content.Pairs = append(content.Pairs,
		&exec.Pair{Key: exec.AsValue("z"), Value: exec.AsValue(1)},
		&exec.Pair{Key: exec.AsValue("a"), Value: exec.AsValue("x")},
	)

	// Mapping keys are strings, not content-part objects with a get method.
	const source = `{%- set tool_body = follow.get('content') -%}
{%- if tool_body is string -%}
string
{%- elif tool_body is sequence and tool_body is not string -%}
{%- for part in tool_body -%}{{ part.get('type') }}{%- endfor -%}
{%- else -%}
fallback
{%- endif -%}`

	got, err := render(t, config.New(), source, map[string]any{
		"follow": map[string]any{"content": content},
	})
	if err == nil {
		t.Fatalf("mapping tool body was accepted with output %q; want an error from calling .get on a mapping key", got)
	}
}

func TestMappingMessageContentIsRejectedAsMalformedContentParts(t *testing.T) {
	// Mappings enter the sequence branch but cannot supply content-part objects.
	const source = `{%- if message.get('content') is string -%}
string
{%- elif message.get('content') is sequence -%}
{%- for item in message['content'] -%}{{ item.get('type') }}{%- endfor -%}
{%- else -%}
fallback
{%- endif -%}`

	got, err := render(t, config.New(), source, map[string]any{
		"message": map[string]any{
			"content": map[string]any{"text": "not a parts list"},
		},
	})
	if err == nil {
		t.Fatalf("mapping message content was accepted with output %q; want an error from calling .get on a mapping key", got)
	}
}
