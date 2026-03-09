package integration_test

import (
	"testing"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func renderTemplate(t *testing.T, source string, data map[string]any) string {
	t.Helper()

	template, err := gonja.FromString(source)
	if err != nil {
		t.Fatalf("parse template: %v", err)
	}

	var context *exec.Context
	if data != nil {
		context = exec.NewContext(data)
	}

	rendered, err := template.ExecuteToString(context)
	if err != nil {
		t.Fatalf("render template: %v", err)
	}

	return rendered
}
