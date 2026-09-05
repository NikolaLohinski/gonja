package integration_test

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TaggedUser struct {
	UserID string `jinja:"user_id"`
}

var _ = Context("struct fields", func() {
	It("should resolve jinja tags through all shared attribute access paths", func() {
		template, err := gonja.FromString(`{{ user.UserID }}|{{ user.user_id }}|{{ user["user_id"] }}|{{ user|attr("user_id") }}|{{ [user]|map(attribute="user_id")|first }}`)
		Expect(err).ToNot(HaveOccurred())

		result, err := template.ExecuteToString(exec.NewContext(map[string]any{
			"user": TaggedUser{UserID: "123456"},
		}))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal("123456|123456|123456|123456|123456"))
	})

	It("should preserve missing-field behavior in non-strict mode", func() {
		template, err := gonja.FromString(`{{ user.missing }}`)
		Expect(err).ToNot(HaveOccurred())

		result, err := template.ExecuteToString(exec.NewContext(map[string]any{
			"user": TaggedUser{UserID: "123456"},
		}))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(""))
	})

	It("should preserve missing-field behavior in strict mode", func() {
		configuration := config.New()
		configuration.StrictUndefined = true
		loader := loaders.MustNewMemoryLoader(map[string]string{
			"/test": `{{ user.missing }}`,
		})
		template, err := exec.NewTemplate("/test", configuration, loader, gonja.DefaultEnvironment)
		Expect(err).ToNot(HaveOccurred())

		result, err := template.ExecuteToString(exec.NewContext(map[string]any{
			"user": TaggedUser{UserID: "123456"},
		}))
		Expect(result).To(BeEmpty())
		Expect(err).To(HaveOccurred())
	})
})
