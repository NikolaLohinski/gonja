package legacy_test

import (
	"math/rand"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(42)
}

func TestLegacy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "legacy")
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustReturn(returned interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return returned
}
