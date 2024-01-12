package integration_test

import (
	"os"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("helpers", func() {
	var (
		returnedTemplate **exec.Template
		returnedErr      *error

		executeResult string
	)
	BeforeEach(func() {
		returnedTemplate = new(*exec.Template)
		returnedErr = new(error)
	})
	Context("FromString", func() {
		const source = "Hello {{ 'bob' | capitalize  }}!"
		JustBeforeEach(func() {
			*returnedTemplate, *returnedErr = gonja.FromString(source)
		})
		It("should return the expected template object", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(*returnedTemplate).ToNot(BeNil())
			executeResult, *returnedErr = (*returnedTemplate).ExecuteToString(nil)
			Expect(*returnedErr).To(BeNil())
			Expect(executeResult).To(Equal("Hello Bob!"))
		})
	})
	Context("FromBytes", func() {
		var source = []byte("Hello {{ 'bob' | capitalize  }}!")
		JustBeforeEach(func() {
			*returnedTemplate, *returnedErr = gonja.FromBytes(source)
		})
		It("should return the expected template object", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(*returnedTemplate).ToNot(BeNil())
			executeResult, *returnedErr = (*returnedTemplate).ExecuteToString(nil)
			Expect(*returnedErr).To(BeNil())
			Expect(executeResult).To(Equal("Hello Bob!"))
		})
	})
	Context("FromFile", func() {
		var (
			filepath = new(string)
		)
		BeforeEach(func() {
			file := MustReturn(os.CreateTemp("", "helpers.*.tpl")).(*os.File)
			MustReturn(file.WriteString("Hello {{ 'bob' | capitalize  }}!"))
			*filepath = file.Name()
		})
		AfterEach(func() {
			os.Remove(*filepath)
		})
		JustBeforeEach(func() {
			*returnedTemplate, *returnedErr = gonja.FromFile(*filepath)
		})
		It("should return the expected template object", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(*returnedTemplate).ToNot(BeNil())
			executeResult, *returnedErr = (*returnedTemplate).ExecuteToString(nil)
			Expect(*returnedErr).To(BeNil())
			Expect(executeResult).To(Equal("Hello Bob!"))
		})
	})
})
