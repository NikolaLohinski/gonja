package legacy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLegacy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "legacy")
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustReturn(returned any, err error) any {
	if err != nil {
		panic(err)
	}
	return returned
}
