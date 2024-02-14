package loaders_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "loader")
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustReturn[T interface{}](returned T, err error) T {
	if err != nil {
		panic(err)
	}
	return returned
}
