package integration_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestTemplates(t *testing.T) {
	// Add a global to the default set
	root := "./testdata"
	env := testEnv(root)
	env.Globals.Set("this_is_a_global_variable", "this is a global text")
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	env.Globals.Set("CI_COMMIT_TAG", strings.TrimSpace(buf.String()))
	GlobTemplateTests(t, root, env)
}

func TestExpressions(t *testing.T) {
	root := "./testdata/expressions"
	env := testEnv(root)
	GlobTemplateTests(t, root, env)
}

func TestFilters(t *testing.T) {
	root := "./testdata/filters"
	env := testEnv(root)
	GlobTemplateTests(t, root, env)
}

func TestFunctions(t *testing.T) {
	root := "./testdata/functions"
	env := testEnv(root)
	GlobTemplateTests(t, root, env)
}

func TestTests(t *testing.T) {
	root := "./testdata/tests"
	env := testEnv(root)
	GlobTemplateTests(t, root, env)
}

func TestStatements(t *testing.T) {
	root := "./testdata/statements"
	env := testEnv(root)
	GlobTemplateTests(t, root, env)
}
