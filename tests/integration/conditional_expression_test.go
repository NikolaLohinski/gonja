package integration_test

import (
	"fmt"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"testing"
)

func TestGroupedConditional(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"true", `{{ 'a' + ('b' if true else '') }}`, "ab"},
		{"false", `{{ 'a' + ('b' if false else '') }}`, "a"},
		{"lazy_true", `{{ ('ok' if true else fail()) }}`, "ok"},
		{"lazy_false", `{{ (fail() if false else 'ok') }}`, "ok"},
		{"chain", `{{ ('a' if false else 'b' if false else 'c') }}`, "c"},
		{"nested", `{{ (('a' if false else 'b') if true else 'c') }}`, "b"},
		{"boolean", `{{ ('a' if true and not false else 'b') }}`, "a"},
		{"filter", `{{ ('a' if false else 'b') | upper }}`, "B"},
		{"set", `{% set x = ('a' if false else 'b') %}{{ x }}`, "b"},
		{"tuple", `{{ ('a' if false else 'b', 'c' if true else 'd')[1] }}`, "c"},
		{"missing_else", `{{ ('a' if false) }}`, ""},
		{"for_filter", `{% for x in [1,2,3] if x > 1 %}{{ x }}{% endfor %}`, "23"},
		{"output_conditional", `{{ 'a' if false else 'b' }}`, "b"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tpl, err := gonja.FromString(c.source)
			if err != nil {
				t.Fatal(err)
			}
			out, err := tpl.ExecuteToString(exec.NewContext(map[string]interface{}{"fail": func() (string, error) { return "", fmt.Errorf("selected failure") }}))
			if err != nil {
				t.Fatal(err)
			}
			if out != c.want {
				t.Fatalf("got %q want %q", out, c.want)
			}
		})
	}
	for _, source := range []string{`{{ ('a' if ) }}`, `{{ ('a' if true else ) }}`} {
		if _, err := gonja.FromString(source); err == nil {
			t.Fatalf("accepted malformed %s", source)
		}
	}
	for _, source := range []string{`{{ ('a' if fail() else 'b') }}`, `{{ (fail() if true else 'b') }}`} {
		tpl, err := gonja.FromString(source)
		if err != nil {
			t.Fatal(err)
		}
		_, err = tpl.ExecuteToString(exec.NewContext(map[string]interface{}{"fail": func() (string, error) { return "", fmt.Errorf("selected failure") }}))
		if err == nil {
			t.Fatalf("swallowed error: %s", source)
		}
	}
}
