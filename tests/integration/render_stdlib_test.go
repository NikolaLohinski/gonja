package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/gomega"
)

func renderTemplate(source string, data map[string]any) string {
	template, err := gonja.FromString(source)
	ExpectWithOffset(1, err).To(BeNil(), "parse template")

	var context *exec.Context
	if data != nil {
		context = exec.NewContext(data)
	}

	rendered, err := template.ExecuteToString(context)
	ExpectWithOffset(1, err).To(BeNil(), "render template")

	return rendered
}
