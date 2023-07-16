//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/nikolalohinski/gonja"
	"github.com/nikolalohinski/gonja/loaders"

	u "github.com/nikolalohinski/gonja/utils"
)

func testEnv(root string) *gonja.Environment {
	cfg := gonja.NewConfig()
	loader := loaders.MustNewFileSystemLoader(root)
	env := gonja.NewEnvironment(cfg, loader)
	env.Autoescape = true
	env.Globals.Set("lorem", u.Lorem) // Predictable random content
	return env
}

func GlobTemplateTests(t *testing.T, root string, env *gonja.Environment) {
	pattern := filepath.Join(root, `*.tpl`)
	matches, err := filepath.Glob(pattern)
	// env := TestEnv(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, match := range matches {
		filename, err := filepath.Rel(root, match)
		if err != nil {
			t.Fatalf("Unable to compute path from `%s`:\n%s", match, err.Error())
		}
		testName := strings.Replace(path.Base(match), ".tpl", "", 1)
		t.Run(testName, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			rand.Seed(42) // Make tests deterministics

			tpl, err := env.FromFile(filename)
			if err != nil {
				t.Fatalf("Error on FromFile('%s'):\n%s", filename, err.Error())
			}
			testFilename := fmt.Sprintf("%s.out", match)
			expected, rerr := ioutil.ReadFile(testFilename)
			if rerr != nil {
				t.Fatalf("Error on ReadFile('%s'):\n%s", testFilename, rerr.Error())
			}
			rendered, err := tpl.ExecuteBytes(Fixtures)
			if err != nil {
				t.Fatalf("Error on Execute('%s'):\n%s", filename, err.Error())
			}
			// rendered = testTemplateFixes.fixIfNeeded(filename, rendered)
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
				t.Errorf("%s rendered with diff:\n%v", testFilename, result)
			}
		})
	}
}

func GlobErrorTests(t *testing.T, root string) {
	pattern := filepath.Join(root, `*.err`)
	matches, err := filepath.Glob(pattern)
	env := testEnv(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, match := range matches {
		testName := strings.Replace(path.Base(match), ".err", "", 1)
		t.Run(testName, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			testData, err := ioutil.ReadFile(match)
			tests := strings.Split(string(testData), "\n")

			checkFilename := fmt.Sprintf("%s.out", match)
			checkData, err := ioutil.ReadFile(checkFilename)
			if err != nil {
				t.Fatalf("Error on ReadFile('%s'):\n%s", checkFilename, err.Error())
			}
			checks := strings.Split(string(checkData), "\n")

			if len(checks) != len(tests) {
				t.Fatal("Template lines != Checks lines")
			}

			for idx, test := range tests {
				if strings.TrimSpace(test) == "" {
					continue
				}
				if strings.TrimSpace(checks[idx]) == "" {
					t.Fatalf("[%s Line %d] Check is empty (must contain an regular expression).",
						match, idx+1)
				}

				tpl, err := env.FromString(test)
				if err != nil {
					t.Fatalf("Error on FromString('%s'):\n%s", test, err.Error())
				}

				tpl, err = env.FromBytes([]byte(test))
				if err != nil {
					t.Fatalf("Error on FromBytes('%s'):\n%s", test, err.Error())
				}

				_, err = tpl.ExecuteBytes(Fixtures)
				if err == nil {
					t.Fatalf("[%s Line %d] Expected error for (got none): %s",
						match, idx+1, tests[idx])
				}

				re := regexp.MustCompile(fmt.Sprintf("^%s$", checks[idx]))
				if !re.MatchString(err.Error()) {
					t.Fatalf("[%s Line %d] Error for '%s' (err = '%s') does not match the (regexp-)check: %s",
						match, idx+1, test, err.Error(), checks[idx])
				}
			}
		})
	}
}
