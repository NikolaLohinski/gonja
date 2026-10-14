package integration_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func TestAdditionRejectsUndefinedAndNoneOperands(t *testing.T) {
	context := exec.NewContext(map[string]interface{}{
		"x": map[string]interface{}{},
	})

	for _, source := range []string{
		`{{ 'prefix' + x.missing }}`,
		`{{ x.missing + 'suffix' }}`,
		`{{ 'prefix' + None }}`,
		`{{ None + 'suffix' }}`,
		`{{ 1 + None }}`,
	} {
		t.Run(source, func(t *testing.T) {
			tpl, err := gonja.FromString(source)
			if err != nil {
				t.Fatal(err)
			}
			if output, err := tpl.ExecuteToString(context); err == nil {
				t.Fatalf("rendered %q instead of rejecting nil operand", output)
			}
		})
	}
}

func TestUndefinedRemainsPermissiveOutsideAddition(t *testing.T) {
	context := exec.NewContext(map[string]interface{}{
		"x": map[string]interface{}{},
	})

	for _, test := range []struct {
		source string
		want   string
	}{
		{`{{ x.missing }}`, ""},
		{`{{ 'prefix' ~ x.missing }}`, "prefix"},
	} {
		tpl, err := gonja.FromString(test.source)
		if err != nil {
			t.Fatal(err)
		}
		output, err := tpl.ExecuteToString(context)
		if err != nil {
			t.Fatal(err)
		}
		if output != test.want {
			t.Fatalf("%s rendered %q, want %q", test.source, output, test.want)
		}
	}
}
