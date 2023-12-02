//go:build version
// +build version

package integration_test

import (
	"bytes"
	"fmt"
	osExec "os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
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
		BeforeEach(func() {
			*loader = loaders.MustNewMemoryLoader(map[string]string{
				*identifier: heredoc.Doc(`
					{%- if "v" ~ gonja.version != CI_COMMIT_TAG -%}
					v{{- gonja.version }} != {{ CI_COMMIT_TAG }}
					{%- else -%}
					"v" ~ gonja.version == CI_COMMIT_TAG
					{%- endif %}
				`),
			})
			cmd := osExec.Command("git", "describe", "--tags", "--abbrev=0")
			buf := new(bytes.Buffer)
			cmd.Stdout = buf
			Must(cmd.Run())
			(*environment).Context.Set("CI_COMMIT_TAG", strings.TrimSpace(buf.String()))
		})
		It("should return the expected rendered content", func() {
			By("not returning any error")
			Expect(*returnedErr).To(BeNil())
			By("not returning the correct result")
			expected := heredoc.Doc(`
				"v" ~ gonja.version == CI_COMMIT_TAG
			`)
			edits := myers.ComputeEdits("expected", expected, *returnedResult)
			diffs := gotextdiff.ToUnified("expected", "got", expected, edits)
			Expect(diffs.Hunks).To(BeEmpty(), "\n"+fmt.Sprint(diffs))
		})
	})
})
