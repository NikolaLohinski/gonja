package integration_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

// Models an agate-style adapter: a struct with positional items and optional
// named columns, rather than a Go slice or an exec.Dict.
type tableSequence struct {
	items []any
	names []string
}

func (s *tableSequence) Len() int { return len(s.items) }

func (s *tableSequence) GetItem(key any) (*exec.Value, bool) {
	if name, ok := key.(string); ok {
		for i, column := range s.names {
			if column == name {
				return exec.ToValue(s.items[i]), true
			}
		}
	}
	if i, ok := key.(int); ok && i >= 0 && i < len(s.items) {
		return exec.ToValue(s.items[i]), true
	}
	return exec.AsValue(nil), false
}

func TestTableSequence(t *testing.T) {
	row := &tableSequence{items: []any{"alice", 42}, names: []string{"name", "age"}}
	table := &tableSequence{items: []any{row, &tableSequence{items: []any{"bob", 17}, names: row.names}}}
	data := exec.NewContext(map[string]any{"table": table, "empty": &tableSequence{}})
	for _, tc := range []struct{ name, source, want string }{
		{"iterate", `{% for row in table %}{{ loop.index }}/{{ loop.length }}:{{ row.name }}={{ row[1] }};{% endfor %}`, "1/2:alice=42;2/2:bob=17;"},
		{"row iteration", `{% for row in table %}{% for cell in row %}{{ cell }};{% endfor %}{% endfor %}`, "alice;42;bob;17;"},
		{"unpack rows", `{% for name, age in table %}{{ name }}={{ age }};{% endfor %}`, "alice=42;bob=17;"},
		{"classification", `{{ table is sequence }}|{{ table is iterable }}|{{ table is mapping }}|{{ table[0] is sequence }}`, "True|True|False|True"},
		{"length", `{{ table|length }}|{{ table[0]|length }}|{{ empty|length }}`, "2|2|0"},
		{"truthiness", `{% if table %}full{% endif %}|{% if empty %}bad{% else %}empty{% endif %}|{{ not empty }}|{{ not table }}`, "full|empty|True|False"},
		{"empty iteration", `{% for row in empty %}bad{% else %}empty{% endfor %}`, "empty"},
		{"negative indices", `{{ table[-2][-2] }}|{{ table[-1]['age'] }}|{{ table[-3] }}|{{ table[2] }}`, "alice|17||"},
		{"slice", `{% for row in table[1:] %}{{ row.name }}{% endfor %}|{{ table[0][::-1]|join(',') }}`, "bob|42,alice"},
		{"filters", `{{ (table|first).name }}|{{ (table|last).name }}|{{ table|map(attribute='name')|join(',') }}|{{ table|reverse|map(attribute='name')|join(',') }}`, "alice|bob|alice,bob|bob,alice"},
		{"membership", `{{ 'alice' in table[0] }}|{{ 'name' in table[0] }}|{{ ['alice', 42] in table }}`, "True|False|True"},
		{"list conversion", `{{ table|list }}|{{ table[0] }}`, "[['alice', 42], ['bob', 17]]|['alice', 42]"},
		{"concatenate", `{{ table[0] + [9] }}|{{ [9] + table[0] }}`, "['alice', 42, 9]|[9, 'alice', 42]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			template, err := gonja.FromString(tc.source)
			if err != nil {
				t.Fatal(err)
			}
			got, err := template.ExecuteToString(data)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestIterationOfUnsupportedStruct(t *testing.T) {
	template, err := gonja.FromString(`{% for item in value %}bad{% else %}empty{% endfor %}`)
	if err != nil {
		t.Fatal(err)
	}
	got, err := template.ExecuteToString(exec.NewContext(map[string]any{"value": struct{ Name string }{"alice"}}))
	if err != nil {
		t.Fatal(err)
	}
	if got != "empty" {
		t.Fatalf("got %q, want empty", got)
	}
}
