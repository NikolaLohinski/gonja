package loaders_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("shifted", func() {
	var (
		loader loaders.Loader

		rootID      = new(string)
		rootContent = new(io.Reader)
		subLoader   = new(loaders.Loader)

		returnedErr = new(error)
	)

	BeforeEach(func() {
		*rootID = "rootID"
		*rootContent = bytes.NewBufferString("root content")
		*subLoader = loaders.MustNewMemoryLoader(map[string]string{
			"/foo": "bar",
		})
	})

	JustBeforeEach(func() {
		loader = loaders.MustNewShiftedLoader(*rootID, *rootContent, *subLoader)
	})

	Context("Resolve", func() {
		var (
			path = new(string)

			returnedPath = new(string)
		)
		JustBeforeEach(func() {
			*returnedPath, *returnedErr = loader.Resolve(*path)
		})
		Context("when reaching out to root", func() {
			BeforeEach(func() {
				*path = *rootID
			})
			It("should retrieve the correct path", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning the correct path")
				Expect(string(*returnedPath)).To(Equal(*rootID))
			})
		})
		Context("when the path is valid in the sub-loader", func() {
			BeforeEach(func() {
				*path = "/foo"
			})
			It("should retrieve the expected path", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning the correct path")
				Expect(string(*returnedPath)).To(Equal("/foo"))
			})
		})
		Context("when the path is not valid in the sub-loader", func() {
			BeforeEach(func() {
				By("returning an error")
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when using in conjunction with the file system loader", Ordered, func() {
			var (
				directory = new(string)
			)
			BeforeAll(func() {
				*directory = MustReturn(os.MkdirTemp("", ""))
				Must(os.MkdirAll(filepath.Join(*directory, "nested"), os.ModePerm))
				file := MustReturn(os.CreateTemp(filepath.Join(*directory, "nested"), ""))
				*path = "./nested/" + filepath.Base(file.Name())
				MustReturn(file.WriteString("content"))
			})
			AfterAll(func() {
				Must(os.RemoveAll(*directory))
			})
			BeforeEach(func() {
				*rootContent = bytes.NewBufferString("root content")
				*subLoader = loaders.MustNewFileSystemLoader(*directory)
			})
			It("should retrieve the expected path", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning the correct path")
				Expect(string(*returnedPath)).To(Equal(filepath.Join(*directory, *path)))
			})
		})
	})
	Context("Read", func() {
		var (
			path = new(string)

			returnedContent = new(io.Reader)
		)
		JustBeforeEach(func() {
			*returnedContent, *returnedErr = loader.Read(*path)
		})
		Context("when reaching out to the root", func() {
			BeforeEach(func() {
				*path = *rootID
			})
			It("should retrieve the root content", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning the correct content")
				Expect(*returnedContent).ToNot(BeNil())
				Expect(string(MustReturn(io.ReadAll(*returnedContent)))).To(Equal("root content"))
			})
		})
		Context("when the path is valid in the sub-loader", func() {
			BeforeEach(func() {
				*path = "/foo"
			})
			It("should retrieve the expected content", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning the correct content")
				Expect(*returnedContent).ToNot(BeNil())
				Expect(string(MustReturn(io.ReadAll(*returnedContent)))).To(Equal("bar"))
			})
		})
		Context("when the path is not valid in the sub-loader", func() {
			BeforeEach(func() {
				By("returning an error")
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
})
