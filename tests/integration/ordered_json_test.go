package integration_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func TestDecodeOrderedNestedAndEmpty(t *testing.T) {
	got, err := gonja.DecodeOrdered([]byte(`{"z":[{"b":true,"a":null},[],{}],"a":{"y":"café","x":[1,1.5]}}`))
	if err != nil {
		t.Fatal(err)
	}
	root, ok := got.(*exec.Dict)
	if !ok {
		t.Fatalf("root type %T", got)
	}
	if len(root.Pairs) != 2 || root.Pairs[0].Key.String() != "z" || root.Pairs[1].Key.String() != "a" {
		t.Fatalf("root order: %v", root)
	}
	array, ok := root.Pairs[0].Value.Interface().([]any)
	if !ok || len(array) != 3 {
		t.Fatalf("array: %#v", root.Pairs[0].Value.Interface())
	}
	first, ok := array[0].(*exec.Dict)
	if !ok || len(first.Pairs) != 2 || first.Pairs[0].Key.String() != "b" || first.Pairs[1].Key.String() != "a" || first.Pairs[0].Value.Interface() != true || first.Pairs[1].Value.Interface() != nil {
		t.Fatalf("nested object: %#v", array[0])
	}
	if empty, ok := array[1].([]any); !ok || empty == nil || len(empty) != 0 {
		t.Fatalf("empty array: %#v", array[1])
	}
	if empty, ok := array[2].(*exec.Dict); !ok || len(empty.Pairs) != 0 {
		t.Fatalf("empty object: %#v", array[2])
	}
	nested, ok := root.Pairs[1].Value.Interface().(*exec.Dict)
	if !ok || len(nested.Pairs) != 2 || nested.Pairs[0].Key.String() != "y" || nested.Pairs[0].Value.Interface() != "café" || nested.Pairs[1].Key.String() != "x" || !reflect.DeepEqual(nested.Pairs[1].Value.Interface(), []any{int64(1), float64(1.5)}) {
		t.Fatalf("nested values: %#v", root.Pairs[1].Value.Interface())
	}
}

func TestDecodeOrderedDuplicateKeys(t *testing.T) {
	got, err := gonja.DecodeOrdered([]byte(`{"z":1,"a":2,"z":3}`))
	if err != nil {
		t.Fatal(err)
	}
	dict := got.(*exec.Dict)
	if len(dict.Pairs) != 2 || dict.Pairs[0].Key.String() != "z" || dict.Pairs[0].Value.Interface() != int64(3) || dict.Pairs[1].Key.String() != "a" || dict.Pairs[1].Value.Interface() != int64(2) {
		t.Fatalf("duplicate key position/value: %v", dict)
	}
}

func TestDecodeOrderedNumbers(t *testing.T) {
	got, err := gonja.DecodeOrdered([]byte(`[1,-2,1.0,-0.0,1e2]`))
	if err != nil {
		t.Fatal(err)
	}
	if want := []any{int64(1), int64(-2), float64(1), float64(0), float64(100)}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
	for i, value := range got.([]any) {
		if i < 2 {
			if _, ok := value.(int64); !ok {
				t.Fatalf("number %d has type %T, want int64", i, value)
			}
		} else if _, ok := value.(float64); !ok {
			t.Fatalf("number %d has type %T, want float64", i, value)
		}
	}
}

func TestDecodeOrderedRejectsInvalidJSON(t *testing.T) {
	for _, raw := range []string{``, `{`, `[`, `{"x":}`, `[1,]`, `{1:2}`, `]`, `{} {}`, `true false`, `9223372036854775808`, `1e999`} {
		t.Run(raw, func(t *testing.T) {
			if _, err := gonja.DecodeOrdered([]byte(raw)); err == nil {
				t.Fatalf("accepted %q", raw)
			}
		})
	}
}

func TestToJSONSortedDictionaryKeys(t *testing.T) {
	decoded, err := gonja.DecodeOrdered([]byte(`{"z":[{"y":2,"x":1}],"a":{"d":4,"c":3}}`))
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, source string
		data         map[string]any
		want         string
	}{
		{"literal", `{{ {"z": {"y": 2, "x": 1}, "a": 3}|tojson }}`, nil, `{"a":3,"z":{"x":1,"y":2}}`},
		{"decoded", `{{ x|tojson }}`, map[string]any{"x": decoded}, `{"a":{"c":3,"d":4},"z":[{"x":1,"y":2}]}`},
		{"map", `{{ x|tojson }}`, map[string]any{"x": map[string]any{"z": decoded, "a": 1}}, `{"a":1,"z":{"a":{"c":3,"d":4},"z":[{"x":1,"y":2}]}}`},
		{"escaping", `{{ x|tojson }}`, map[string]any{"x": &exec.Dict{Pairs: []*exec.Pair{{Key: exec.AsValue("z<"), Value: exec.AsValue("<>&'\u2028\u2029é")}}}}, `{"z\u003c":"\u003c\u003e\u0026\u0027\u2028\u2029é"}`},
		{"indent", `{{ x|tojson(indent=2) }}`, map[string]any{"x": decoded}, "{\n  \"a\": {\n    \"c\": 3,\n    \"d\": 4\n  },\n  \"z\": [\n    {\n      \"x\": 1,\n      \"y\": 2\n    }\n  ]\n}"},
		{"nil and empty", `{{ x|tojson }}`, map[string]any{"x": []any{(*exec.Dict)(nil), &exec.Dict{}, (*exec.Value)(nil)}}, `[null,{},null]`},
		{"ordinary values", `{{ x|tojson }}`, map[string]any{"x": []any{true, 1.5, []byte("hi"), compatJSONMarshaler{}}}, `[true,1.5,"aGk=",42]`},
		{"duplicates", `{{ x|tojson }}`, map[string]any{"x": &exec.Dict{Pairs: []*exec.Pair{{Key: exec.AsValue("z"), Value: exec.AsValue(1)}, {Key: exec.AsValue("a"), Value: exec.AsValue(2)}, {Key: exec.AsValue("z"), Value: exec.AsValue(3)}}}}, `{"a":2,"z":3}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := render(t, config.New(), tc.source, tc.data)
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.want {
				t.Fatalf("got %q want %q", out, tc.want)
			}
		})
	}
}

func TestToJSONNonfiniteError(t *testing.T) {
	value := &exec.Dict{Pairs: []*exec.Pair{{Key: exec.AsValue("z"), Value: exec.AsValue(math.NaN())}}}
	_, err := render(t, config.New(), `{{ x|tojson }}`, map[string]any{"x": value})
	if err == nil {
		t.Fatal("expected JSON encoding error")
	}
}
