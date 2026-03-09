package exec_test

import (
	"math/big"
	"testing"

	"github.com/nikolalohinski/gonja/v2/exec"
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

func TestValueScalarCompatibilityHelpers(t *testing.T) {
	t.Run("string uses pointer stringers before dereferencing", func(t *testing.T) {
		if got := exec.AsValue(big.NewInt(42)).String(); got != "42" {
			t.Fatalf("expected big.Int stringer output, got %q", got)
		}
	})

	t.Run("index returns embedded exec values directly", func(t *testing.T) {
		value := exec.AsValue([]*exec.Value{exec.AsValue("foo")})
		item := value.Index(0)
		if !item.IsString() {
			t.Fatal("expected indexed value to behave like the wrapped string")
		}
		if got := item.String(); got != "foo" {
			t.Fatalf("unexpected indexed value: %q", got)
		}
	})

	t.Run("items unwrap map entries into plain values", func(t *testing.T) {
		items := exec.AsValue(map[string]*exec.Value{
			"a": exec.AsValue(1),
		}).Items()
		if len(items) != 1 {
			t.Fatalf("expected exactly one item, got %d", len(items))
		}
		if !items[0].Value.IsInteger() {
			t.Fatal("expected map value to be unwrapped as an integer")
		}
	})

	t.Run("custom item getters are used", func(t *testing.T) {
		value, ok := exec.AsValue(testItemGetter{items: []string{"foo", "bar"}}).GetItem(1)
		if !ok {
			t.Fatal("expected custom item getter to resolve the item")
		}
		if got := value.String(); got != "bar" {
			t.Fatalf("unexpected item value: %q", got)
		}
	})
}
