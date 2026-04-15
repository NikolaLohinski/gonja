package exec_test

import (
	"math/big"

	"github.com/ardanlabs/gonja/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testItemGetter struct {
	items []string
}

func (t testItemGetter) GetItem(key any) (*exec.Value, bool) {
	index, ok := key.(int)
	if !ok || index < 0 || index >= len(t.items) {
		return exec.AsValue(nil), false
	}
	return exec.AsValue(t.items[index]), true
}

var _ = Context("value scalar compatibility helpers", func() {
	It("uses pointer stringers before dereferencing", func() {
		Expect(exec.AsValue(big.NewInt(42)).String()).To(Equal("42"))
	})

	It("returns embedded exec values directly from index", func() {
		value := exec.AsValue([]*exec.Value{exec.AsValue("foo")})
		item := value.Index(0)
		Expect(item.IsString()).To(BeTrue())
		Expect(item.String()).To(Equal("foo"))
	})

	It("unwraps map entries into plain values", func() {
		items := exec.AsValue(map[string]*exec.Value{
			"a": exec.AsValue(1),
		}).Items()
		Expect(items).To(HaveLen(1))
		Expect(items[0].Value.IsInteger()).To(BeTrue())
	})

	It("uses custom item getters", func() {
		value, ok := exec.AsValue(testItemGetter{items: []string{"foo", "bar"}}).GetItem(1)
		Expect(ok).To(BeTrue())
		Expect(value.String()).To(Equal("bar"))
	})
})
