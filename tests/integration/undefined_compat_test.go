package integration_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
)

func TestUndefinedAndNoneAreDistinct(t *testing.T) {
	const source = `{{ missing }}|{{ missing is defined }}|{{ missing is undefined }}|{{ missing is none }}|{{ missing|default('D') }}|{{ None }}|{{ None is defined }}|{{ None is none }}|{{ None|default('D') }}|{{ [None] }}|{{ {}.get('x') }}`
	got, err := render(t, config.New(), source, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := "|False|True|False|D|None|True|True|None|[None]|None"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestDefaultUndefinedBehavior(t *testing.T) {
	cases := []struct {
		name, source, want string
	}{
		{"false", `{% if missing %}yes{% else %}no{% endif %}`, "no"},
		{"iteration", `{% for x in missing %}x{% else %}empty{% endfor %}`, "empty"},
		{"length", `{{ missing|length }}`, "0"},
		{"list", `{{ missing|list }}`, "[]"},
		{"map", `{{ missing|map|list }}`, "[]"},
		{"upper", `{{ missing|upper }}`, ""},
		{"concat", `{{ missing ~ 'x' }}`, "x"},
		{"equal undefined", `{{ missing == other_missing }}`, "True"},
		{"not equal none", `{{ missing != None }}`, "True"},
		{"conditional", `{{ (1 if false) is undefined }}`, "True"},
		{"macro omitted", `{% macro m(value) %}{{ value is undefined }}:{{ value|default('D') }}{% endmacro %}{{ m() }}`, "True:D"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := render(t, config.New(), tc.source, nil)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestUndefinedInvalidOperationsFail(t *testing.T) {
	for _, source := range []string{
		`{{ missing + 1 }}`,
		`{{ +missing }}`,
		`{{ missing < 1 }}`,
		`{{ missing.value }}`,
		`{{ missing[0] }}`,
		`{{ missing() }}`,
		`{{ missing|int }}`,
	} {
		t.Run(source, func(t *testing.T) {
			if got, err := render(t, config.New(), source, nil); err == nil {
				t.Fatalf("rendered %q; want an error", got)
			}
		})
	}
}

func TestStrictUndefinedOnlyAllowsInspectionAndDefault(t *testing.T) {
	cfg := config.New()
	cfg.StrictUndefined = true
	got, err := render(t, cfg, `{{ missing is defined }}|{{ missing is undefined }}|{{ missing|default('D') }}`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "False|True|D" {
		t.Fatalf("got %q", got)
	}

	for _, source := range []string{
		`{{ missing }}`,
		`{% if missing %}yes{% endif %}`,
		`{% for x in missing %}{{ x }}{% endfor %}`,
		`{{ missing|length }}`,
		`{{ missing == missing }}`,
		`{{ missing is equalto(missing) }}`,
		`{{ missing ~ 'x' }}`,
		`{{ not missing }}`,
	} {
		t.Run(source, func(t *testing.T) {
			if got, err := render(t, cfg, source, nil); err == nil {
				t.Fatalf("rendered %q; want an error", got)
			}
		})
	}
}

func TestUndefinedCollectionAndAttributeSemantics(t *testing.T) {
	ctx := map[string]any{
		"d":     map[string]any{},
		"users": []any{map[string]any{"x": 1}, map[string]any{}},
	}
	cases := []struct {
		name, source, want string
	}{
		{"nested repr", `{{ [missing] }}|{{ {'a': missing} }}`, `[Undefined]|{'a': Undefined}`},
		{"selectattr defined", `{{ users|selectattr('x', 'defined')|list }}`, `[{'x': 1}]`},
		{"selectattr undefined", `{{ users|selectattr('x', 'undefined')|list }}`, `[{}]`},
		{"map missing", `{{ users|map(attribute='x')|list }}`, `[1, Undefined]`},
		{"items missing", `{{ missing|items|list }}`, `[]`},
		{"empty ends", `{{ []|first is undefined }}|{{ []|last is undefined }}|{{ []|first|default('D') }}`, `True|True|D`},
		{"lowercase none", `{{ none }}|{{ none is none }}|{{ none is defined }}`, `None|True|True`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, strict := range []bool{false, true} {
				cfg := config.New()
				cfg.StrictUndefined = strict
				got, err := render(t, cfg, tc.source, ctx)
				if err != nil {
					t.Fatalf("strict=%v: %v", strict, err)
				}
				if got != tc.want {
					t.Fatalf("strict=%v: got %q want %q", strict, got, tc.want)
				}
			}
		})
	}

	got, err := render(t, config.New(), `{{ d[missing] is undefined }}`, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got != "True" {
		t.Fatalf("got %q", got)
	}
}

func TestNoneInvalidOperations(t *testing.T) {
	for _, source := range []string{
		`{{ None|length }}`,
		`{{ None|list }}`,
		`{% for item in None %}{{ item }}{% endfor %}`,
		`{{ None|items|list }}`,
		`{{ None < 1 }}`,
		`{{ None > 1 }}`,
		`{{ None <= 1 }}`,
		`{{ None >= 1 }}`,
		`{{ None - 1 }}`,
		`{{ 1 - None }}`,
		`{{ None * 1 }}`,
		`{{ 1 / None }}`,
		`{{ None // 1 }}`,
		`{{ None % 1 }}`,
		`{{ None ** 1 }}`,
		`{{ 1 in None }}`,
		`{{ None in 'text' }}`,
		`{{ missing in 'text' }}`,
	} {
		t.Run(source, func(t *testing.T) {
			if got, err := render(t, config.New(), source, nil); err == nil {
				t.Fatalf("rendered %q; want an error", got)
			}
		})
	}
}

func TestStrictUndefinedConsumers(t *testing.T) {
	cfg := config.New()
	cfg.StrictUndefined = true

	got, err := render(t, cfg, `{{ missing in [] }}|{{ missing|items|list }}`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != `False|[]` {
		t.Fatalf("got %q", got)
	}

	for _, source := range []string{
		`{{ missing in [None, 1] }}`,
		`{{ 1 in [missing] }}`,
		`{{ missing is iterable }}`,
		`{{ missing in {} }}`,
		`{{ 1 in missing }}`,
		`{{ missing|first is undefined }}`,
		`{{ missing|last is undefined }}`,
		`{% set d = {} %}{{ d[missing] }}`,
		`{{ [{}]|selectattr('x')|list }}`,
		`{{ +missing }}`,
	} {
		t.Run(source, func(t *testing.T) {
			if got, err := render(t, cfg, source, nil); err == nil {
				t.Fatalf("rendered %q; want an error", got)
			}
		})
	}
}

func TestUndefinedSequenceAndConditionalSemantics(t *testing.T) {
	for _, strict := range []bool{false, true} {
		cfg := config.New()
		cfg.StrictUndefined = strict
		source, want := `{{ missing is sequence }}|{{ missing is iterable }}`, `True|True`
		if strict {
			source, want = `{{ missing is sequence }}|{{ (1 if false) }}`, `False|`
		}
		got, err := render(t, cfg, source, nil)
		if err != nil {
			t.Fatalf("strict=%v: %v", strict, err)
		}
		if got != want {
			t.Fatalf("strict=%v: got %q want %q", strict, got, want)
		}
	}
}

func TestMalformedExpressionsReturnErrorsWithoutPanicking(t *testing.T) {
	for _, source := range []string{`{{ missing[ }}`, `{{ missing. }}`, `{{ missing( }}`} {
		t.Run(source, func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Fatalf("panicked: %v", recovered)
				}
			}()
			if _, err := gonja.FromString(source); err == nil {
				t.Fatal("malformed expression was accepted")
			}
		})
	}
}
