package loaders_test

import (
	"io"

	"github.com/nikolalohinski/gonja/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("memory", func() {
	var (
		loader loaders.Loader

		content = new(map[string]string)

		returnedErr = new(error)
	)

	BeforeEach(func() {
		*content = map[string]string{
			"/home/sweet": "home",
			"/home/of":    "content",
		}
	})

	JustBeforeEach(func() {
		loader = loaders.MustNewMemoryLoader(*content)
	})

	Context("Read", func() {
		var (
			path = new(string)

			reader = new(io.Reader)
		)
		BeforeEach(func() {
			*path = "/home/of"
		})
		JustBeforeEach(func() {
			*reader, *returnedErr = loader.Read(*path)
		})
		Context("when path is absolute", func() {
			It("should retrieve the expected file", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning a reader with the correct content")
				content, err := io.ReadAll(*reader)
				Expect(err).To(BeNil())
				Expect(string(content)).To(Equal("content"))
			})
		})
		Context("when path is relative", func() {
			BeforeEach(func() {
				*path = "of"
			})
			It("should retrieve the expected file", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning a reader with the correct content")
				content, err := io.ReadAll(*reader)
				Expect(err).To(BeNil())
				Expect(string(content)).To(Equal("content"))
			})
		})
	})
	Context("Resolve", func() {
		var (
			path = new(string)

			returnedPath = new(string)
		)
		BeforeEach(func() {
			*path = "/home/sweet"
		})
		JustBeforeEach(func() {
			*returnedPath, *returnedErr = loader.Resolve(*path)
		})
		Context("when path is absolute", func() {
			It("should retrieve the expected file", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning the correct path")
				Expect(string(*returnedPath)).To(Equal(*path))
			})
		})
		Context("when path is relative", func() {
			BeforeEach(func() {
				*path = "sweet"
			})
			It("should retrieve the expected file", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("returning the correct path")
				Expect(string(*returnedPath)).To(Equal("/home/sweet"))
			})
		})
	})
	Context("Inherit", func() {
		var (
			newIdentifier = new(string)

			returnedLoader = new(loaders.Loader)
		)
		BeforeEach(func() {
			*newIdentifier = ""
		})
		JustBeforeEach(func() {
			*returnedLoader, *returnedErr = loader.Inherit(*newIdentifier)
		})
		Context("when no root is given", func() {
			It("should create a new Loader without errors", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("having the loader operate relatively to the inherited root")
				_, err := (*returnedLoader).Read("sweet")
				Expect(err).To(BeNil())
			})
		})
		Context("when a new root is defined", func() {
			BeforeEach(func() {
				*newIdentifier = "/home/of"
			})
			It("should create a new Loader without errors", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("having the loader operate relatively to the new root")
				_, err := (*returnedLoader).Read("of")
				Expect(err).To(BeNil())
			})
		})
	})
})
