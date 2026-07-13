package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recursive for-loop", func() {
	// Full recursive loop support requires refactoring gonja's callability
	// model to allow `loop` to be both an attribute-bearing value and a
	// callable in the same expression. Marked pending until addressed.
	PIt("should render a tree recursively", func() {
		tree := []map[string]any{
			{"name": "a", "children": []map[string]any{
				{"name": "b", "children": []map[string]any{}},
				{"name": "c", "children": []map[string]any{
					{"name": "d", "children": []map[string]any{}},
				}},
			}},
		}
		result := renderTemplate(
			`{% for node in tree recursive %}{{ node.name }}({{ loop(node.children) }}){% endfor %}`,
			map[string]any{"tree": tree})
		Expect(result).To(Equal("a(b()c(d()))"))
	})
})
