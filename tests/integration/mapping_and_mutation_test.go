package integration_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func decodeOrderedContext(t *testing.T, fixture string) map[string]any {
	t.Helper()
	value, err := gonja.DecodeOrdered([]byte(fixture))
	if err != nil {
		t.Fatal(err)
	}
	dict, ok := value.(*exec.Dict)
	if !ok {
		t.Fatalf("fixture root is %T, want *exec.Dict", value)
	}
	context := make(map[string]any, len(dict.Pairs))
	for _, pair := range dict.Pairs {
		context[pair.Key.String()] = pair.Value.Interface()
	}
	return context
}

func orderedJSONConfig() *config.Config {
	cfg := config.New()
	cfg.HuggingFaceToJSON = true
	return cfg
}

func renderDictItemKeys(t *testing.T, value any) string {
	t.Helper()
	out, err := render(t, config.New(),
		`{% for key, value in subject.items() %}{{ key }};{% endfor %}`,
		map[string]any{"subject": value})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestDictItemsSortsGoMapKeysCaseSensitively(t *testing.T) {
	got := renderDictItemKeys(t, map[string]any{"a": 1, "B": 2, "c": 3})
	if want := "B;a;c;"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestDictItemsPreservesNativeDictInsertionOrder(t *testing.T) {
	dict := exec.NewDict()
	for _, key := range []string{"z", "A", "m"} {
		dict.Pairs = append(dict.Pairs, &exec.Pair{Key: exec.AsValue(key), Value: exec.AsValue(key)})
	}
	if got, want := renderDictItemKeys(t, dict), "z;A;m;"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestItemsPreserveOrderedNestedValues(t *testing.T) {
	data := decodeOrderedContext(t, `{"arguments":{"city":"Paris","opts":{"z":1,"a":2}}}`)
	out, err := render(t, orderedJSONConfig(),
		`{% for k, v in arguments.items() %}{% if v is mapping %}{{ k }}={{ v|tojson }}{% endif %}{% endfor %}`,
		data)
	if err != nil {
		t.Fatal(err)
	}
	if want := `opts={"z": 1, "a": 2}`; out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestSplitStringsMatchDictionaryKeys(t *testing.T) {
	data := decodeOrderedContext(t, `{
		"tools":[{"function":{"name":"weather.get"}},{"name":"math.add"}],
		"tool_namespace_descriptions":{"weather":"Weather tools","math":"Math tools"}
	}`)
	const source = `{% macro metadata(tools) %}{% set ns = namespace(seen=[]) %}{% for tool in tools %}{% set fn = tool.function if tool.function is defined else tool %}{% set name = fn.name.split('.')[0] %}{% if name not in ns.seen %}{% set ns.seen = ns.seen + [name] %}{% endif %}{% endfor %}{% set descriptions = tool_namespace_descriptions %}{% for name in ns.seen %}{{ name }}={{ descriptions[name] if name in descriptions else '' }};{% endfor %}{% endmacro %}{{ metadata(tools) }}`
	out, err := render(t, orderedJSONConfig(), source, data)
	if err != nil {
		t.Fatal(err)
	}
	if want := "weather=Weather tools;math=Math tools;"; out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestListAppendMutationPersists(t *testing.T) {
	out, err := render(t, orderedJSONConfig(), `{% set values = [] %}{{ values.append('kept') }}{{ values|join }}`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out != "Nonekept" {
		t.Fatalf("got %q want %q", out, "Nonekept")
	}
}
