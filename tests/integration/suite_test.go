package integration_test

import (
	"fmt"
	"testing"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLegacy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "integration")
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustReturn[I interface{}](returned I, err error) I {
	if err != nil {
		panic(err)
	}
	return returned
}

func AssertPrettyDiff(expected, got string) bool {
	edits := myers.ComputeEdits("expected", expected, got)
	diffs := gotextdiff.ToUnified("expected", "got", expected, edits)
	return Expect(diffs.Hunks).To(BeEmpty(), "\n"+fmt.Sprint(diffs))
}
