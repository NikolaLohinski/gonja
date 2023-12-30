package exec_test

import (
	"fmt"

	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("value", func() {
	Context("AsValue", func() {
		var (
			input = new(interface{})

			returnedValue = new(*exec.Value)
		)
		JustBeforeEach(func() {
			*returnedValue = exec.AsValue(*input)
		})
		for _, testCase := range []struct {
			golangObject interface{}
			description  string
			matchers     []func()
		}{
			{
				nil,
				"a nil value",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("")) },
					func() { Expect((*returnedValue).IsNil()).To(BeTrue()) },
				},
			},
			{
				"Hello World",
				"a string",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("Hello World"), ".String()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
					func() { Expect((*returnedValue).IsString()).To(BeTrue(), ".IsString()") },
					func() { Expect((*returnedValue).IsIterable()).To(BeTrue(), ".IsIterable()") },
				},
			},
			{
				42,
				"an integer",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("42"), ".String()") },
					func() { Expect((*returnedValue).IsInteger()).To(BeTrue(), ".IsInteger()") },
					func() { Expect((*returnedValue).IsNumber()).To(BeTrue(), ".IsNumber()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{
				0,
				"a zero value integer",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("0"), ".String()") },
					func() { Expect((*returnedValue).IsInteger()).To(BeTrue(), ".IsInteger()") },
					func() { Expect((*returnedValue).IsNumber()).To(BeTrue(), ".IsNumber()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeFalse(), ".IsTrue()") },
				},
			},
			{
				42.0,
				"a float",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("42.0"), ".String()") },
					func() { Expect((*returnedValue).IsInteger()).To(BeFalse(), ".IsInteger()") },
					func() { Expect((*returnedValue).IsFloat()).To(BeTrue(), ".IsFloat()") },
					func() { Expect((*returnedValue).IsNumber()).To(BeTrue(), ".IsNumber()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{
				42.5556700089099,
				"a float over the with maximal precision",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("42.55567000891"), ".String()") },
					func() { Expect((*returnedValue).IsInteger()).To(BeFalse(), ".IsInteger()") },
					func() { Expect((*returnedValue).IsFloat()).To(BeTrue(), ".IsFloat()") },
					func() { Expect((*returnedValue).IsNumber()).To(BeTrue(), ".IsNumber()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{
				0.0,
				"a zero value float",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("0.0"), ".String()") },
					func() { Expect((*returnedValue).IsInteger()).To(BeFalse(), ".IsInteger()") },
					func() { Expect((*returnedValue).IsFloat()).To(BeTrue(), ".IsFloat()") },
					func() { Expect((*returnedValue).IsNumber()).To(BeTrue(), ".IsNumber()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeFalse(), ".IsTrue()") },
				},
			},
			{
				true,
				"a true boolean",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("True"), ".String()") },
					func() { Expect((*returnedValue).IsBool()).To(BeTrue(), ".IsBool()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{
				false,
				"a falsy boolean",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("False"), ".String()") },
					func() { Expect((*returnedValue).IsBool()).To(BeTrue(), ".IsBool()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeFalse(), ".IsTrue()") },
				},
			},
			{
				[]int{1, 2, 3},
				"a slice of integers",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("[1, 2, 3]"), ".String()") },
					func() { Expect((*returnedValue).IsIterable()).To(BeTrue(), ".IsIterable()") },
					func() { Expect((*returnedValue).IsList()).To(BeTrue(), ".IsList()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{
				[]string{"a", "b", "c"},
				"a slice of strings",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("['a', 'b', 'c']"), ".String()") },
					func() { Expect((*returnedValue).IsIterable()).To(BeTrue(), ".IsIterable()") },
					func() { Expect((*returnedValue).IsList()).To(BeTrue(), ".IsList()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{
				[3]*exec.Value{exec.AsValue("a"), exec.AsValue("b"), exec.AsValue("c")},
				"a array of *exec.Values",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("['a', 'b', 'c']"), ".String()") },
					func() { Expect((*returnedValue).IsIterable()).To(BeTrue(), ".IsIterable()") },
					func() { Expect((*returnedValue).IsList()).To(BeTrue(), ".IsList()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{
				map[string]interface{}{
					"a": "a",
					"b": "b",
				},
				"a dictionary as a map[string]interface{}",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("{'a': 'a', 'b': 'b'}"), ".String()") },
					func() { Expect((*returnedValue).IsIterable()).To(BeTrue(), ".IsIterable()") },
					func() { Expect((*returnedValue).IsDict()).To(BeTrue(), ".IsDict()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{

				&exec.Dict{
					Pairs: []*exec.Pair{
						{Key: exec.AsValue("a"), Value: exec.AsValue("a")},
						{Key: exec.AsValue("b"), Value: exec.AsValue("b")},
					},
				},
				"a dictionary as an *exec.Dict of key/value pairs",
				[]func(){
					func() { Expect((*returnedValue).String()).To(Equal("{'a': 'a', 'b': 'b'}"), ".String()") },
					func() { Expect((*returnedValue).IsIterable()).To(BeTrue(), ".IsIterable()") },
					func() { Expect((*returnedValue).IsDict()).To(BeTrue(), ".IsDict()") },
					func() { Expect((*returnedValue).IsTrue()).To(BeTrue(), ".IsTrue()") },
				},
			},
			{
				func() {},
				"a function",
				[]func(){
					func() { Expect((*returnedValue).IsCallable()).To(BeTrue(), ".IsCallable()") },
				},
			},
		} {
			t := testCase
			Context(fmt.Sprintf("when the value is %s", t.description), func() {
				BeforeEach(func() {
					*input = t.golangObject
				})
				for _, matcher := range t.matchers {
					m := matcher
					It("should return the correct value", func() {
						By("being defined")
						Expect(*returnedValue).ToNot(BeNil())
						By("responding correctly to the expected *exec.Value methods")
						m()
					})
				}
			})
		}
	})

	Context("GetAttribute", func() {
		var (
			value     = new(*exec.Value)
			attribute = new(string)

			returnedAttribute = new(*exec.Value)
			returnedOk        = new(bool)
		)
		JustBeforeEach(func() {
			*returnedAttribute, *returnedOk = (*value).GetAttribute(*attribute)
		})
		Context("when the holding value is undefined", func() {
			BeforeEach(func() {
				*value = exec.AsValue(nil)
			})
			It("should fail", func() {
				By("returning a not ok flag")
				Expect(*returnedOk).To(BeFalse())
				By("returning an error")
				Expect((*returnedAttribute).IsError()).To(BeTrue(), ".isError()")
			})
		})
		Context("when the holding value is a struct", func() {
			BeforeEach(func() {
				type structure struct{ attribute string }
				*value = exec.AsValue(structure{
					attribute: "attribute",
				})
				*attribute = "attribute"
			})
			It("should return the expect content", func() {
				By("returning an ok flag")
				Expect(*returnedOk).To(BeTrue())
				By("not returning an error")
				Expect((*returnedAttribute).IsError()).To(BeFalse(), ".isError()")
				By("returning the expected attribute")
				Expect((*returnedAttribute).IsString()).To(BeTrue(), ".IsString()")
				Expect((*returnedAttribute).String()).To(Equal("attribute"), ".String()")
			})
		})
		Context("when the holding value is a struct but the attribute is not found", func() {
			BeforeEach(func() {
				type structure struct{}
				*value = exec.AsValue(structure{})
			})
			It("should return the expect content", func() {
				By("returning a not ok flag")
				Expect(*returnedOk).To(BeFalse())
				By("not returning an error")
				Expect((*returnedAttribute).IsError()).To(BeFalse(), ".isError()")
				By("returning a nil value")
				Expect((*returnedAttribute).IsNil()).To(BeTrue(), ".isNil()")
			})
		})
		Context("when the holding value is a map", func() {
			BeforeEach(func() {
				*value = exec.AsValue(map[string]interface{}{
					"attribute": "attribute",
				})
				*attribute = "attribute"
			})
			It("should return the expect content", func() {
				By("returning a not ok flag")
				Expect(*returnedOk).To(BeFalse())
				By("not returning an error")
				Expect((*returnedAttribute).IsError()).To(BeFalse(), ".isError()")
				By("returning a nil value")
				Expect((*returnedAttribute).IsNil()).To(BeTrue(), ".IsNil()")
			})
		})
	})

	Context("GetItem", func() {
		var (
			value = new(*exec.Value)
			item  = new(interface{})

			returnedItem = new(*exec.Value)
			returnedOk   = new(bool)
		)
		JustBeforeEach(func() {
			*returnedItem, *returnedOk = (*value).GetItem(*item)
		})
		Context("when the holding value is undefined", func() {
			BeforeEach(func() {
				*value = exec.AsValue(nil)
			})
			It("should fail", func() {
				By("returning a not ok flag")
				Expect(*returnedOk).To(BeFalse())
				By("returning an error")
				Expect((*returnedItem).IsError()).To(BeTrue(), ".isError()")
			})
		})
		Context("when the holding value is a map", func() {
			BeforeEach(func() {
				*value = exec.AsValue(map[string]interface{}{
					"item": "item",
				})
				*item = "item"
			})
			It("should return the expect content", func() {
				By("returning an ok flag")
				Expect(*returnedOk).To(BeTrue())
				By("not returning an error")
				Expect((*returnedItem).IsError()).To(BeFalse(), ".isError()")
				By("returning the expected attribute")
				Expect((*returnedItem).IsString()).To(BeTrue(), ".IsString()")
				Expect((*returnedItem).String()).To(Equal("item"), ".String()")
			})
		})
		Context("when the holding value is a map but the item is not found", func() {
			BeforeEach(func() {
				*value = exec.AsValue(map[string]interface{}{})
			})
			It("should return the expect content", func() {
				By("returning a not ok flag")
				Expect(*returnedOk).To(BeFalse())
				By("not returning an error")
				Expect((*returnedItem).IsError()).To(BeFalse(), ".isError()")
				By("returning a nil value")
				Expect((*returnedItem).IsNil()).To(BeTrue(), ".isNil()")
			})
		})
		Context("when the holding value is a struct", func() {
			BeforeEach(func() {
				type structure struct{}
				*value = exec.AsValue(structure{})
			})
			It("should return the expect content", func() {
				By("returning a not ok flag")
				Expect(*returnedOk).To(BeFalse())
				By("not returning an error")
				Expect((*returnedItem).IsError()).To(BeFalse(), ".isError()")
				By("returning a nil value")
				Expect((*returnedItem).IsNil()).To(BeTrue(), ".IsNil()")
			})
		})
		Context("when the holding value is map as a key/value pairs dictionary", func() {
			BeforeEach(func() {
				*value = exec.AsValue(&exec.Dict{[]*exec.Pair{
					{exec.AsValue("key"), exec.AsValue("value")},
				}})
				*item = "key"
			})
			It("should return the expect content", func() {
				By("returning an ok flag")
				Expect(*returnedOk).To(BeTrue())
				By("not returning an error")
				Expect((*returnedItem).IsError()).To(BeFalse(), ".isError()")
				By("returning a the correct item")
				Expect((*returnedItem).String()).To(Equal("value"), ".IsNil()")
			})
		})
	})

	Context("Set", func() {
		var (
			holder = new(*exec.Value)
			key    = new(*exec.Value)
			value  = new(interface{})

			returnedErr = new(error)
		)
		BeforeEach(func() {
			*holder = exec.AsValue(nil)
			*key = exec.AsValue("")
		})
		JustBeforeEach(func() {
			*returnedErr = (*holder).Set(*key, *value)
		})
		Context("when the holder value is nil", func() {
			It("should fail", func() {
				By("returning an non nil error")
				Expect(*returnedErr).To(MatchError("Can't set attribute or item on None"))
			})
		})
		Context("when overriding an existing attribute on a struct pointer", func() {
			BeforeEach(func() {
				type structure struct {
					Attribute string
				}
				*holder = exec.AsValue(&structure{
					Attribute: "attribute",
				})
				*key = exec.AsValue("Attribute")
				*value = "override"
			})
			It("should set the holder correctly", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("setting the correct value on the holder")
				item, ok := (*holder).GetAttribute((*key).String())
				Expect(ok).To(BeTrue(), "item should exist")
				Expect(item.String()).To(Equal("override"), "item should be correct")
			})
		})
		Context("when overriding an existing attribute on a direct struct", func() {
			BeforeEach(func() {
				type structure struct {
					Attribute string
				}
				*holder = exec.AsValue(structure{
					Attribute: "attribute",
				})
				*key = exec.AsValue("Attribute")
				*value = "override"
			})
			It("should fail to set the attribute", func() {
				By("returning an error")
				Expect(*returnedErr).To(MatchError("Can't write field \"Attribute\""))
			})
		})
		Context("when setting an unknown attribute on a struct", func() {
			BeforeEach(func() {
				type structure struct {
					Attribute string
				}
				*holder = exec.AsValue(&structure{
					Attribute: "attribute",
				})
				*key = exec.AsValue("Missing")
				*value = "override"
			})
			It("should fail to set the attribute", func() {
				By("returning an error")
				Expect(*returnedErr).To(MatchError("Can't write field \"Missing\""))
			})
		})
		Context("when setting a new key on a map", func() {
			BeforeEach(func() {
				*holder = exec.AsValue(map[string]interface{}{
					"existing": "item",
				})
				*key = exec.AsValue("new")
				*value = "new"
			})
			It("should set the holder correctly", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("setting the correct value on the holder")
				item, ok := (*holder).GetItem((*key).String())
				Expect(ok).To(BeTrue(), "item should exist")
				Expect(item.String()).To(Equal("new"), "item should be correct")
			})
		})
		Context("when setting an existing key on a map", func() {
			BeforeEach(func() {
				*holder = exec.AsValue(map[string]interface{}{
					"existing": "item",
				})
				*key = exec.AsValue("existing")
				*value = "new"
			})
			It("should set the holder correctly", func() {
				By("not returning an error")
				Expect(*returnedErr).To(BeNil())
				By("setting the correct value on the holder")
				item, ok := (*holder).GetItem((*key).String())
				Expect(ok).To(BeTrue(), "item should exist")
				Expect(item.String()).To(Equal("new"), "item should be correct")
			})
		})
	})

	Context("Keys", func() {
		var (
			value = new(*exec.Value)

			returnedKeys = new(exec.ValuesList)
		)
		JustBeforeEach(func() {
			*returnedKeys = (*value).Keys()
		})

		Context("when the value can not be tested for keys", func() {
			for _, kind := range []interface{}{
				nil,
				"string",
				42,
				33.0,
				true,
				false,
				[1]int{3},
				[]string{"a", "b"},
				func() {},
			} {
				v := exec.AsValue(kind)
				BeforeEach(func() {
					*value = v
				})
				Context(v.String(), func() {
					It("should return an empty list", func() {
						Expect(*returnedKeys).To(BeEmpty())
					})
				})
			}
		})
		Context("when the value is a map", func() {
			BeforeEach(func() {
				*value = exec.AsValue(map[string]interface{}{
					"c": "c",
					"a": "a",
					"B": "B",
				})
			})
			It("should return the correct list", func() {
				Expect((*returnedKeys).String()).To(Equal("['a', 'B', 'c']"))
			})
		})
		Context("when the value is a key/value pairs map", func() {
			BeforeEach(func() {
				*value = exec.AsValue(&exec.Dict{[]*exec.Pair{
					{exec.AsValue("c"), exec.AsValue("c")},
					{exec.AsValue("A"), exec.AsValue("A")},
					{exec.AsValue("b"), exec.AsValue("b")},
				}})
			})
			It("should return the correct list", func() {
				Expect((*returnedKeys).String()).To(Equal("['c', 'A', 'b']"))
			})
		})
	})
})
