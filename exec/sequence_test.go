package exec_test

import (
	"reflect"
	"testing"

	"github.com/nikolalohinski/gonja/v2/exec"
)

type valueSequence struct{ items []*exec.Value }

func (s valueSequence) Len() int { return len(s.items) }
func (s valueSequence) GetItem(key any) (*exec.Value, bool) {
	i, ok := key.(int)
	if !ok || i < 0 || i >= len(s.items) {
		return exec.AsValue(nil), false
	}
	return s.items[i], true
}

var _ exec.Sequence = valueSequence{}

func TestSequenceValues(t *testing.T) {
	safe := exec.AsSafeValue("<b>safe</b>")
	sequence := valueSequence{items: []*exec.Value{safe, exec.AsValue(nil)}}
	for _, input := range []any{sequence, &sequence} {
		value := exec.AsValue(input)
		if !value.IsList() || !value.IsIterable() || !value.CanSlice() || value.IsDict() {
			t.Fatalf("incorrect sequence classification for %T", input)
		}
		if item := value.Slice(0, 1).Index(0); value.Index(0) != safe || !item.Safe || !item.EqualValueTo(safe) {
			t.Fatal("sequence access discarded the item's safety flag or value")
		}
		if !value.Index(-1).IsNil() || !value.Index(2).IsNil() {
			t.Fatal("out-of-bounds index did not return None")
		}
		item, found := value.GetItem(-2)
		if !found || item != safe {
			t.Fatal("negative index did not resolve the first item")
		}
		want := []any{"<b>safe</b>", nil}
		if got := value.ToGoSimpleType(false); !reflect.DeepEqual(got, want) {
			t.Fatalf("conversion got %#v, want %#v", got, want)
		}
	}
	var output []string
	if err := exec.StringListArgument(&output)(exec.AsValue(valueSequence{items: []*exec.Value{safe}})); err != nil || !reflect.DeepEqual(output, []string{"<b>safe</b>"}) {
		t.Fatalf("string list argument: %v, %v", output, err)
	}
}

func TestSequenceIterationOrder(t *testing.T) {
	sequence := exec.AsValue(valueSequence{items: []*exec.Value{exec.AsValue("B"), exec.AsValue("a"), exec.AsValue("C")}})
	for _, tc := range []struct {
		reverse, sorted, caseSensitive bool
		want                           []string
	}{
		{false, false, false, []string{"B", "a", "C"}},
		{true, false, false, []string{"C", "a", "B"}},
		{false, true, false, []string{"a", "B", "C"}},
		{true, true, false, []string{"C", "B", "a"}},
		{false, true, true, []string{"B", "C", "a"}},
	} {
		var got []string
		sequence.IterateOrder(func(index, count int, item, value *exec.Value) bool {
			if index != len(got) || count != 3 || value != nil {
				t.Fatal("incorrect iteration metadata")
			}
			got = append(got, item.String())
			return true
		}, func() { t.Fatal("nonempty sequence reported empty") }, tc.reverse, tc.sorted, tc.caseSensitive)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("got %v, want %v", got, tc.want)
		}
	}
	if sequence.Index(0).String() != "B" {
		t.Fatal("sorting mutated the input")
	}
	calls := 0
	sequence.Iterate(func(_, _ int, _, _ *exec.Value) bool { calls++; return false }, func() { t.Fatal("stopped sequence reported empty") })
	if calls != 1 {
		t.Fatalf("early stop called callback %d times", calls)
	}
}

func TestSequenceMembershipWithMissingValues(t *testing.T) {
	for _, tc := range []struct {
		name  string
		items []*exec.Value
		other *exec.Value
		want  bool
		err   bool
	}{
		{"None equals None", []*exec.Value{exec.AsValue(nil)}, exec.AsValue(nil), true, false},
		{"None differs from undefined", []*exec.Value{exec.AsValue(nil)}, exec.UndefinedValue(false, "missing"), false, false},
		{"undefined equals undefined", []*exec.Value{exec.UndefinedValue(false, "missing")}, exec.UndefinedValue(false, "other"), true, false},
		{"empty sequence tolerates strict operand", nil, exec.UndefinedValue(true, "missing"), false, false},
		{"strict operand", []*exec.Value{exec.AsValue(1)}, exec.UndefinedValue(true, "missing"), false, true},
		{"strict item", []*exec.Value{exec.UndefinedValue(true, "missing")}, exec.AsValue(1), false, true},
		{"named string", []*exec.Value{exec.AsValue(namedSequenceString("alice"))}, exec.AsValue("alice"), true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := exec.AsValue(valueSequence{items: tc.items}).ContainsChecked(tc.other)
			if got != tc.want || (err != nil) != tc.err {
				t.Fatalf("got (%v, %v), want (%v, error=%v)", got, err, tc.want, tc.err)
			}
		})
	}
}

type namedSequenceString string

func TestIterationEmptyCallback(t *testing.T) {
	var nilSequence *valueSequence
	for _, tc := range []struct {
		name                  string
		input                 any
		emptyCalls, itemCalls int
	}{
		{"empty sequence", valueSequence{}, 1, 0},
		{"nil sequence", nilSequence, 1, 0},
		{"unsupported struct", struct{}{}, 1, 0},
		{"empty dict", exec.NewDict(), 1, 0},
		{"dict", &exec.Dict{Pairs: []*exec.Pair{{Key: exec.AsValue("k"), Value: exec.AsValue("v")}}}, 0, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			emptyCalls, itemCalls := 0, 0
			exec.AsValue(tc.input).Iterate(func(_, _ int, _, _ *exec.Value) bool { itemCalls++; return true }, func() { emptyCalls++ })
			if emptyCalls != tc.emptyCalls || itemCalls != tc.itemCalls {
				t.Fatalf("got %d empty callbacks and %d items", emptyCalls, itemCalls)
			}
		})
	}
}
