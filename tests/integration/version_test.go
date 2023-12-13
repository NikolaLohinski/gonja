//go:build version
// +build version

package integration_test

import (
	"bytes"
	osExec "os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja"
	"github.com/nikolalohinski/gonja/config"
	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("version", func() {
	var (
		identifier = new(string)

		environment = new(*exec.Environment)
		config      = new(*config.Config)
		loader      = new(loaders.Loader)

		context = new(*exec.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
	)
	BeforeEach(func() {
		*identifier = "/test"
		*environment = gonja.DefaultEnvironment
		*config = gonja.DefaultConfig
		*loader = loaders.MustNewMemoryLoader(nil)
	})
	JustBeforeEach(func() {
		var t *exec.Template
		t, *returnedErr = exec.NewTemplate(*identifier, *config, *loader, *environment)
		if *returnedErr != nil {
			return
		}
		*returnedResult, *returnedErr = t.Execute(*context)
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
			By("not returning the correct result")
			expected := heredoc.Doc(`
				gonja.version is '` + *version + `'
			`)
			AssertPrettyDiff(expected, *returnedResult)
		})
	})
})
