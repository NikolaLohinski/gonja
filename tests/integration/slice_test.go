package integration_test

import (
	"strings"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Slicing with step", func() {
	DescribeTable("slice expressions",
		func(tpl, expected string) {
			Expect(renderTemplate(tpl, nil)).To(Equal(expected))
		},
		Entry("every second element", `{{ [1,2,3,4,5,6][::2] }}`, "[1, 3, 5]"),
		Entry("start:stop:step", `{{ [1,2,3,4,5,6][1:5:2] }}`, "[2, 4]"),
		Entry("reverse list via negative step", `{{ [1,2,3,4,5,6][::-1] }}`, "[6, 5, 4, 3, 2, 1]"),
		Entry("string every second char", `{{ "abcdef"[::2] }}`, "ace"),
		Entry("string reverse", `{{ "abcdef"[::-1] }}`, "fedcba"),
		Entry("regression: no step", `{{ [1,2,3,4,5,6][1:5] }}`, "[2, 3, 4, 5]"),
		Entry("regression: open start", `{{ [1,2,3,4,5,6][:3] }}`, "[1, 2, 3]"),
		Entry("regression: string basic", `{{ "abcdef"[1:4] }}`, "bcd"),
	)

	It("should error gracefully on zero step", func() {
		// This test intentionally does NOT use renderTemplate() because that
		// helper asserts no error occurred. Here we expect an error, so we
		// call the underlying gonja API directly.
		template, err := gonja.FromString(`{{ [1,2,3][::0] }}`)
		Expect(err).To(BeNil(), "template with [::0] should parse OK")

		_, err = template.ExecuteToString(exec.NewContext(nil))
		Expect(err).ToNot(BeNil(), "step of zero should cause an execution error")
		Expect(strings.Contains(err.Error(), "zero")).To(BeTrue(),
			"error should mention zero, got: %s", err.Error())
	})
})
