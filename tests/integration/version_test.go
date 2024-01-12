//go:build version
// +build version

package integration_test

import (
	"bytes"
	osExec "os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("version", func() {
	var (
		identifier = new(string)

		environment = new(*exec.Environment)
		loader      = new(loaders.Loader)

		context = new(*exec.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
	)
	BeforeEach(func() {
		*identifier = "/test"
		*environment = gonja.DefaultEnvironment
		*loader = loaders.MustNewMemoryLoader(nil)
	})
	JustBeforeEach(func() {
		var t *exec.Template
		t, *returnedErr = exec.NewTemplate(*identifier, gonja.DefaultConfig, *loader, *environment)
		if *returnedErr != nil {
			return
		}
		*returnedResult, *returnedErr = t.ExecuteToString(*context)
	})
	Context("when using 'gonja.version'", func() {
		var (
			version = new(string)
		)
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					gonja.version is '{{ gonja.version }}'
				`),
			})
			cmd := osExec.Command("git", "describe", "--tags", "--abbrev=0")
			buf := new(bytes.Buffer)
			cmd.Stdout = buf
			Must(cmd.Run())
			*version = strings.TrimSpace(buf.String())
		})
		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("returning the correct result")
			expected := heredoc.Doc(`
				gonja.version is '` + *version + `'
			`)
			AssertPrettyDiff(expected, *returnedResult)
		})
	})
})
