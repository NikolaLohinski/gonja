package exec_test

import (
	"reflect"

	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("context", func() {
	var (
		ctx = new(*exec.Context)

		name = new(string)

		returnedOk    = new(bool)
		returnedValue = new(interface{})
	)
	BeforeEach(func() {
		*ctx = exec.EmptyContext()
	})
	JustBeforeEach(func() {
		*returnedValue, *returnedOk = (*ctx).Get(*name)
	})
	Context("root", func() {
		BeforeEach(func() {
			*name = "key"
		})
		Context("when the value is nil", func() {
			BeforeEach(func() {
				(*ctx).Set(*name, nil)
			})
			assert := func() {
				It("should return the expected value", func() {
					By("returning an OK flag")
					Expect(*returnedOk).To(BeTrue())
					By("returning the expected value")
					Expect(*returnedValue).To(BeNil())
				})
			}
			assert()
			Context("when using a sub context", func() {
				BeforeEach(func() {
					*ctx = (*ctx).Inherit()
				})
				assert()
			})
		})
		Context("when the value is a string", func() {
			BeforeEach(func() {
				(*ctx).Set(*name, "string")
			})
			assert := func() {
				It("should return the expected value", func() {
					By("returning an OK flag")
					Expect(*returnedOk).To(BeTrue())
					By("returning the expected value")
					Expect(*returnedValue).To(Equal("string"))
				})
			}
			assert()
			Context("when using a sub context", func() {
				BeforeEach(func() {
					*ctx = (*ctx).Inherit()
				})
				assert()
			})
		})
		Context("when the value is an integer", func() {
			BeforeEach(func() {
				(*ctx).Set(*name, 42)
			})
			assert := func() {
				It("should return the expected value", func() {
					By("returning an OK flag")
					Expect(*returnedOk).To(BeTrue())
					By("returning the expected value")
					Expect(*returnedValue).To(Equal(42))
				})
			}
			Context("when using a sub context", func() {
				BeforeEach(func() {
					*ctx = (*ctx).Inherit()
				})
				assert()
			})
		})
		Context("when the value is a float", func() {
			BeforeEach(func() {
				(*ctx).Set(*name, 1.2)
			})
			assert := func() {
				It("should return the expected value", func() {
					By("returning an OK flag")
					Expect(*returnedOk).To(BeTrue())
					By("returning the expected value")
					Expect(*returnedValue).To(Equal(1.2))
				})
			}
			Context("when using a sub context", func() {
				BeforeEach(func() {
					*ctx = (*ctx).Inherit()
				})
				assert()
			})
		})
		Context("when the value is a boolean", func() {
			BeforeEach(func() {
				(*ctx).Set(*name, true)
			})
			assert := func() {
				It("should return the expected value", func() {
					By("returning an OK flag")
					Expect(*returnedOk).To(BeTrue())
					By("returning the expected value")
					Expect(*returnedValue).To(BeTrue())
				})
			}
			Context("when using a sub context", func() {
				BeforeEach(func() {
					*ctx = (*ctx).Inherit()
				})
				assert()
			})
		})
		Context("when the value is a func", func() {
			var (
				fun = func() {}
			)
			BeforeEach(func() {
				(*ctx).Set(*name, fun)
			})
			assert := func() {
				It("should return the expected value", func() {
					By("returning an OK flag")
					Expect(*returnedOk).To(BeTrue())
					By("returning the expected value")
					Expect(reflect.ValueOf(*returnedValue)).To(Equal(reflect.ValueOf(fun)))
				})
			}
			Context("when using a sub context", func() {
				BeforeEach(func() {
					*ctx = (*ctx).Inherit()
				})
				assert()
			})
		})
	})
})
