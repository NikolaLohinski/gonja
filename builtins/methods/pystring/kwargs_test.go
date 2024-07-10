package pystring

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetNestedKwArgs", func() {
	type Nested struct {
		Field string
	}

	type Root struct {
		Nested Nested
	}

	DescribeTable("retrieves the correct value",
		func(keys []string, kwarg AttributeGetter, expected any, expectError bool) {
			result, err := getNestedKwArgs(keys, kwarg)
			if expectError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expected))
			}
		},
		Entry("simple key",
			[]string{"root", "Nested", "Field"},
			KwArgs{
				"root": Root{
					Nested: Nested{
						Field: "value",
					},
				},
			},
			"value",
			false,
		),
		Entry("missing key",
			[]string{"root", "Missing"},
			KwArgs{
				"root": Root{
					Nested: Nested{
						Field: "value",
					},
				},
			},
			"",
			true,
		),
		Entry("map of any",
			[]string{"map", "key"},
			KwArgs{
				"map": map[string]any{
					"key": "value",
				},
			},
			"value",
			false,
		),
		Entry("list of string",
			[]string{"map", "1"},
			KwArgs{
				"map": []any{
					"foo", "bar",
				},
			},
			"bar",
			false,
		),
		Entry("map of string",
			[]string{"map", "key"},
			KwArgs{
				"map": map[string]string{
					"key": "value",
				},
			},
			"value",
			false,
		),
		Entry("map of int",
			[]string{"map", "key"},
			KwArgs{
				"map": map[string]int{
					"key": 42,
				},
			},
			42,
			false,
		),
		Entry("struct without AttributeGetter",
			[]string{"root", "Nested", "Field"},
			KwArgs{
				"root": struct {
					Nested struct {
						Field string
					}
				}{
					Nested: struct {
						Field string
					}{
						Field: "value",
					},
				},
			},
			"value",
			false,
		),
		Entry("struct with AttributeGetter",
			[]string{"root", "Field"},
			KwArgs{
				"root": testStructWithGetter{
					Field: "value",
				},
			},
			"value",
			false,
		),
	)
})

type testStructWithGetter struct {
	Field string
}

func (t testStructWithGetter) Get(key string) (any, bool) {
	if key == "Field" {
		return t.Field, true
	}
	return nil, false
}
