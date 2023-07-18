package integration_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "integration")
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
