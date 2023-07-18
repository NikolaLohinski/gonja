package integration_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nikolalohinski/gonja"
	"github.com/nikolalohinski/gonja/config"
	"github.com/pmezard/go-difflib/difflib"
)

func TestWhiteSpace(t *testing.T) {
	files, err := filepath.Glob("testdata/whitespaces/*.tpl")
	if err != nil {
		panic(err)
	}
	for _, path := range files {
		source := path
		output := path + ".out"
		name := strings.TrimSuffix(filepath.Base(source), ".tpl")
		t.Run(name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			cfg := config.NewConfig()
			env := gonja.NewEnvironment(cfg, gonja.DefaultLoader)

			tpl, err := env.FromFile(source)
			if err != nil {
				t.Fatalf("Error on FromFile('%s'): %s", source, err.Error())
			}
			expected, rerr := ioutil.ReadFile(output)
			if rerr != nil {
				t.Fatalf("Error on ReadFile('%s'): %s", output, rerr.Error())
			}
			rendered, err := tpl.ExecuteBytes(Fixtures)
			if err != nil {
				t.Fatalf("Error on Execute('%s'): %s", source, err.Error())
			}
			// rendered = testTemplateFixes.fixIfNeeded(match, rendered)
			if bytes.Compare(expected, rendered) != 0 {
				diff := difflib.UnifiedDiff{
					A:        difflib.SplitLines(string(expected)),
					B:        difflib.SplitLines(string(rendered)),
					FromFile: "Expected",
					ToFile:   "Rendered",
					Context:  2,
					Eol:      "\n",
				}
				result, _ := difflib.GetUnifiedDiffString(diff)
				t.Errorf("%s rendered with diff:\n%v", source, result)
			}
		})
	}
}
