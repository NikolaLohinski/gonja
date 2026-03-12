package exec_test

import (
	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testAttributeGetter struct{}

func (testAttributeGetter) GetAttribute(name string) (*exec.Value, bool) {
	if name == "name" {
		return exec.AsValue("alice"), true
	}
	return exec.AsValue(nil), false
}

type testInt64Value struct{}

func (testInt64Value) Int64() int64 {
	return 42
}

type testFloat64Value struct{}

func (testFloat64Value) Float64() float64 {
	return 3.5
}

var _ = Context("value compatibility helpers", func() {
	It("uses custom attribute getters", func() {
		value, ok := exec.AsValue(testAttributeGetter{}).GetAttribute("name")
		Expect(ok).To(BeTrue())
		Expect(value.String()).To(Equal("alice"))
	})

	It("parses integers from strings and custom int64 values", func() {
		Expect(exec.AsValue("12.9").Integer()).To(Equal(12))
		Expect(exec.AsValue(testInt64Value{}).Integer()).To(Equal(42))
	})

	It("parses floats from strings and custom float64 values", func() {
		Expect(exec.AsValue("12.5").Float()).To(Equal(12.5))
		Expect(exec.AsValue(testFloat64Value{}).Float()).To(Equal(3.5))
	})
})
