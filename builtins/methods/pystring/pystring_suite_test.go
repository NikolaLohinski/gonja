package pystring_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPystring(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pystring-methods")
}
