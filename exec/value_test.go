package exec_test

import (
	// "fmt"

	"reflect"
	"testing"

	"github.com/nikolalohinski/gonja/exec"
	"github.com/stretchr/testify/assert"
)

type flags struct {
	IsString   bool
	IsCallable bool
	IsBool     bool
	IsFloat    bool
	IsInteger  bool
	IsNumber   bool
	IsList     bool
	IsDict     bool
	IsIterable bool
	IsNil      bool
	IsTrue     bool
	IsError    bool
}

func (f *flags) assert(t *testing.T, value *exec.Value) {
	assert := assert.New(t)

	val := reflect.ValueOf(value)
	fval := reflect.ValueOf(f).Elem()

	for i := 0; i < fval.NumField(); i++ {
		name := fval.Type().Field(i).Name
		method := val.MethodByName(name)
		bVal := fval.Field(i).Interface().(bool)
		result := method.Call([]reflect.Value{})
		bResult := result[0].Interface().(bool)
		if bVal {
			assert.Truef(bResult, `%s() should be true`, name)
		} else {
			assert.Falsef(bResult, `%s() should be false`, name)
		}
	}
}

var valueCases = []struct {
	name     string
	value    interface{}
	asString string
	flags    flags
}{
	{"nil", nil, "", flags{IsNil: true}},
	{"string", "Hello World", "Hello World", flags{IsString: true, IsTrue: true, IsIterable: true}},
	{"int", 42, "42", flags{IsInteger: true, IsNumber: true, IsTrue: true}},
	{"int 0", 0, "0", flags{IsInteger: true, IsNumber: true}},
	{"float", 42., "42.0", flags{IsFloat: true, IsNumber: true, IsTrue: true}},
	{"float with trailing zeros", 42.04200, "42.042", flags{IsFloat: true, IsNumber: true, IsTrue: true}},
	{"float max precision", 42.5556700089099, "42.55567000891", flags{IsFloat: true, IsNumber: true, IsTrue: true}},
	{"float max precision rounded up", 42.555670008999999, "42.555670009", flags{IsFloat: true, IsNumber: true, IsTrue: true}},
	{"float 0.0", 0., "0.0", flags{IsFloat: true, IsNumber: true}},
	{"true", true, "True", flags{IsBool: true, IsTrue: true}},
	{"false", false, "False", flags{IsBool: true}},
	{"slice", []int{1, 2, 3}, "[1, 2, 3]", flags{IsTrue: true, IsIterable: true, IsList: true}},
	{"strings slice", []string{"a", "b", "c"}, "['a', 'b', 'c']", flags{IsTrue: true, IsIterable: true, IsList: true}},
	{
		"values slice",
		[]*exec.Value{exec.AsValue(1), exec.AsValue(2), exec.AsValue(3)},
		"[1, 2, 3]",
		flags{IsTrue: true, IsIterable: true, IsList: true},
	},
	{"string values slice",
		[]*exec.Value{exec.AsValue("a"), exec.AsValue("b"), exec.AsValue("c")},
		"['a', 'b', 'c']",
		flags{IsTrue: true, IsIterable: true, IsList: true},
	},
	{"array", [3]int{1, 2, 3}, "[1, 2, 3]", flags{IsTrue: true, IsIterable: true, IsList: true}},
	{"strings array", [3]string{"a", "b", "c"}, "['a', 'b', 'c']", flags{IsTrue: true, IsIterable: true, IsList: true}},
	{
		"dict as map",
		map[string]string{"a": "a", "b": "b"},
		"{'a': 'a', 'b': 'b'}",
		flags{IsTrue: true, IsIterable: true, IsDict: true},
	},
	{
		"dict as Dict/Pairs",
		&exec.Dict{[]*exec.Pair{
			{exec.AsValue("a"), exec.AsValue("a")},
			{exec.AsValue("b"), exec.AsValue("b")},
		}},
		"{'a': 'a', 'b': 'b'}",
		flags{IsTrue: true, IsIterable: true, IsDict: true},
	},
	{"func", func() {}, "<func() Value>", flags{IsCallable: true}},
}

func TestValue(t *testing.T) {
	for _, lc := range valueCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			assert := assert.New(t)

			value := exec.AsValue(test.value)

			assert.Equal(test.asString, value.String())
			test.flags.assert(t, value)
		})
	}
}

func TestValueFromMap(t *testing.T) {
	for _, lc := range valueCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			assert := assert.New(t)

			data := map[string]interface{}{"value": test.value}
			value := exec.AsValue(data["value"])

			assert.Equal(test.asString, value.String())
			test.flags.assert(t, value)
		})
	}
}

type testStruct struct {
	Attr string
}

func (t testStruct) String() string {
	return t.Attr
}

var getattrCases = []struct {
	name     string
	value    interface{}
	attr     string
	found    bool
	asString string
	flags    flags
}{
	{"nil", nil, "missing", false, "", flags{IsError: true}},
	{"attr found", testStruct{"test"}, "Attr", true, "test", flags{IsString: true, IsTrue: true, IsIterable: true}},
	{"attr not found", testStruct{"test"}, "Missing", false, "", flags{IsNil: true}},
	{"item", map[string]interface{}{"Attr": "test"}, "Attr", false, "", flags{IsNil: true}},
}

func TestValueGetAttr(t *testing.T) {
	for _, lc := range getattrCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			assert := assert.New(t)

			value := exec.AsValue(test.value)
			out, found := value.Getattr(test.attr)

			if !test.flags.IsError && out.IsError() {
				t.Fatalf(`Unexpected error: %s`, out.Error())
			}

			if test.found {
				assert.Truef(found, `Attribute '%s' should be found on %s`, test.attr, value)
				assert.Equal(test.asString, out.String())
			} else {
				assert.Falsef(found, `Attribute '%s' should not be found on %s`, test.attr, value)
			}

			test.flags.assert(t, out)
		})
	}
}

var getitemCases = []struct {
	name     string
	value    interface{}
	key      interface{}
	found    bool
	asString string
	flags    flags
}{
	{"nil", nil, "missing", false, "", flags{IsError: true}},
	{"item found", map[string]interface{}{"Attr": "test"}, "Attr", true, "test", flags{IsString: true, IsTrue: true, IsIterable: true}},
	{"item not found", map[string]interface{}{"Attr": "test"}, "Missing", false, "test", flags{IsNil: true}},
	{"attr", testStruct{"test"}, "Attr", false, "", flags{IsNil: true}},
	{"dict found", &exec.Dict{[]*exec.Pair{
		{exec.AsValue("key"), exec.AsValue("value")},
		{exec.AsValue("otherKey"), exec.AsValue("otherValue")},
	}}, "key", true, "value", flags{IsTrue: true, IsString: true, IsIterable: true}},
}

func TestValueGetitem(t *testing.T) {
	for _, lc := range getitemCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			assert := assert.New(t)

			value := exec.AsValue(test.value)
			out, found := value.Getitem(test.key)

			if !test.flags.IsError && out.IsError() {
				t.Fatalf(`Unexpected error: %s`, out.Error())
			}

			if test.found {
				assert.Truef(found, `Key '%s' should be found on %s`, test.key, value)
				assert.Equal(test.asString, out.String())
			} else {
				assert.Falsef(found, `Key '%s' should not be found on %s`, test.key, value)
			}

			test.flags.assert(t, out)
		})
	}
}

var setCases = []struct {
	name     string
	value    interface{}
	attr     string
	set      interface{}
	error    bool
	asString string
}{
	{"nil", nil, "missing", "whatever", true, ""},
	{"existing attr on struct by ref", &testStruct{"test"}, "Attr", "value", false, "value"},
	{"existing attr on struct by value", testStruct{"test"}, "Attr", "value", true, `Can't write field "Attr"`},
	{"missing attr on struct by ref", &testStruct{"test"}, "Missing", "value", true, "test"},
	{"missing attr on struct by value", testStruct{"test"}, "Missing", "value", true, "test"},
	{
		"existing key on map",
		map[string]interface{}{"Attr": "test"},
		"Attr",
		"value",
		false,
		"{'Attr': 'value'}",
	},
	{
		"new key on map",
		map[string]interface{}{"Attr": "test"},
		"New",
		"value",
		false,
		"{'Attr': 'test', 'New': 'value'}",
	},
}

func TestValueSet(t *testing.T) {
	for _, lc := range setCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			assert := assert.New(t)

			value := exec.AsValue(test.value)
			err := value.Set(exec.AsValue(test.attr), test.set)

			if test.error {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Equal(test.asString, value.String())
			}
		})
	}
}

var valueKeysCases = []struct {
	name     string
	value    interface{}
	asString string
	isError  bool
}{
	{"nil", nil, "", true},
	{"string", "Hello World", "", true},
	{"int", 42, "", true},
	{"float", 42., "", true},
	{"true", true, "", true},
	{"false", false, "", true},
	{"slice", []int{1, 2, 3}, "", true},
	// Map keys are sorted alphabetically, case insensitive
	{"dict as map", map[string]string{"c": "c", "a": "a", "B": "B"}, "['a', 'B', 'c']", false},
	// Dict as Pairs keys are kept in order
	{
		"dict as Dict/Pairs",
		&exec.Dict{[]*exec.Pair{
			{exec.AsValue("c"), exec.AsValue("c")},
			{exec.AsValue("A"), exec.AsValue("A")},
			{exec.AsValue("b"), exec.AsValue("b")},
		}},
		"['c', 'A', 'b']",
		false,
	},
	{"func", func() {}, "", true},
}

func TestValueKeys(t *testing.T) {
	for _, lc := range valueKeysCases {
		test := lc
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			assert := assert.New(t)

			value := exec.AsValue(test.value)
			keys := value.Keys()
			if test.isError {
				assert.Len(keys, 0)
			} else {
				assert.Equal(test.asString, keys.String())
			}
		})
	}
}
