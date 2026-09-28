package integration_test

import (
	"math"
	"testing"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

func render(t *testing.T, cfg *config.Config, source string, data map[string]any) (string, error) {
	t.Helper()
	const identifier = "/test"
	loader := loaders.MustNewMemoryLoader(map[string]string{identifier: source})
	tpl, err := exec.NewTemplate(identifier, cfg, loader, gonja.DefaultEnvironment)
	if err != nil {
		t.Fatal(err)
	}
	return tpl.ExecuteToString(exec.NewContext(data))
}

func TestHuggingFaceToJSONIsOptInAndInherited(t *testing.T) {
	cfg := config.New()
	if cfg.HuggingFaceToJSON || cfg.Inherit().HuggingFaceToJSON {
		t.Fatal("HF JSON must be disabled by default")
	}
	value, err := gonja.DecodeOrdered([]byte(`{"z":1,"a":"<>&'"}`))
	if err != nil {
		t.Fatal(err)
	}
	out, err := render(t, cfg, `{{ x|tojson }}`, map[string]any{"x": value})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"a":"\u003c\u003e\u0026\u0027","z":1}`; out != want {
		t.Fatalf("default tojson changed: got %q want %q", out, want)
	}

	cfg.HuggingFaceToJSON = true
	if !cfg.Inherit().HuggingFaceToJSON {
		t.Fatal("HF JSON configuration lost on inheritance")
	}
	out, err = render(t, cfg, `{{ x|tojson }}`, map[string]any{"x": value})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"z": 1, "a": "<>&'"}`; out != want {
		t.Fatalf("HF tojson: got %q want %q", out, want)
	}
	if _, err := render(t, cfg, `{{ x|tojson }}`, map[string]any{"x": map[string]any{"a": 1}}); err == nil {
		t.Fatal("HF tojson accepted an unordered Go map")
	}
}

func TestHuggingFaceToJSON(t *testing.T) {
	cfg := config.New()
	cfg.HuggingFaceToJSON = true
	cases := []struct{ source, input, want string }{
		{`{{ x|tojson }}`, `{"z":1,"a":{"b":[true,null,"<>&' café 😀\n\u2028"]}}`, `{"z": 1, "a": {"b": [true, null, "<>&' café 😀\n` + "\u2028" + `"]}}`},
		{`{{ x|tojson }}`, `{"z":1,"a":2,"z":3}`, `{"z": 3, "a": 2}`},
		{`{{ x|tojson }}`, `[{},[],"\"\\\b\f\r\t\u0000"]`, `[{}, [], "\"\\\b\f\r\t\u0000"]`},
		{`{{ x|tojson }}`, `["\\u2028","\u2028"]`, `["\\u2028", "` + "\u2028" + `"]`},
		{`{{ x|tojson }}`, `[1.0,-0.0,0.0001,0.00001,1000000000000000.0,1e16]`, `[1.0, -0.0, 0.0001, 1e-05, 1000000000000000.0, 1e+16]`},
		{`{% for k,v in x|items %}{{ k }}={{ v }};{% endfor %}`, `{"z":1,"a":2}`, `z=1;a=2;`},
		{`{{ {'z': 1, 'a': 2}|tojson }}`, `null`, `{"z": 1, "a": 2}`},
	}
	for _, c := range cases {
		v, err := gonja.DecodeOrdered([]byte(c.input))
		if err != nil {
			t.Fatal(err)
		}
		out, err := render(t, cfg, c.source, map[string]any{"x": v})
		if err != nil {
			t.Fatal(err)
		}
		if out != c.want {
			t.Fatalf("got %q want %q", out, c.want)
		}
	}
	for _, raw := range []string{`{} {}`, `{"x":}`, `[1,]`, `{1:2}`, `{`, `[`, `]`, `true false`, `9223372036854775808`, `1e999`} {
		if _, err := gonja.DecodeOrdered([]byte(raw)); err == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
	if _, err := render(t, cfg, `{{ x|tojson(indent=2) }}`, map[string]any{"x": 1}); err == nil {
		t.Fatal("silently accepted unsupported option")
	}
}

func TestHuggingFaceToJSONGoValues(t *testing.T) {
	cfg := config.New()
	cfg.HuggingFaceToJSON = true
	cases := []struct {
		name  string
		value any
		want  string
	}{
		{"empty_dict", exec.Dict{}, `{}`},
		{"empty_dict_pointer", &exec.Dict{}, `{}`},
		{"array", [2]int{1, 2}, `[1, 2]`},
		{"slice_pointer", &[]int{1, 2}, `[1, 2]`},
		{"empty_slice", []int{}, `[]`},
		{"values_list", exec.ValuesList{exec.AsValue(1), exec.AsValue("a")}, `[1, "a"]`},
		{"nonfinite", []float64{math.NaN(), math.Inf(1), math.Inf(-1)}, `[NaN, Infinity, -Infinity]`},
		{"unicode_separators", "\u2028\u2029", "\"\u2028\u2029\""},
		{"control_characters", "\a\v\x1f", `"\u0007\u000b\u001f"`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := render(t, cfg, `{{ x|tojson }}`, map[string]any{"x": c.value})
			if err != nil {
				t.Fatal(err)
			}
			if out != c.want {
				t.Fatalf("got %q want %q", out, c.want)
			}
		})
	}
	if _, err := render(t, cfg, `{{ x|tojson }}`, map[string]any{"x": []any{map[string]any{"a": 1}}}); err == nil {
		t.Fatal("accepted an unordered map nested in a list")
	}
}
