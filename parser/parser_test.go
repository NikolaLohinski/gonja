package parser_test

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/parser"
	"github.com/nikolalohinski/gonja/tokens"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logLevel = flag.String("log.level", "", "Log Level")

func TestMain(m *testing.M) {
	flag.Parse()

	log.SetFormatter(&prefixed.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
		ForceFormatting:  true,
	})

	switch *logLevel {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warning", "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.PanicLevel)
	}
	os.Exit(m.Run())
}

var testCases = []struct {
	name     string
	text     string
	expected specs
}{
	{"comment", "{# My comment #}", specs{nodes.Comment{}, attrs{
		"Text": val{" My comment "},
	}}},
	{"multiline comment", "{# My\nmultiline\ncomment #}", specs{nodes.Comment{}, attrs{
		"Text": val{" My\nmultiline\ncomment "},
	}}},
	{"empty comment", "{##}", specs{nodes.Comment{}, attrs{
		"Text": val{""},
	}}},
	{"raw text", "raw text", specs{nodes.Data{}, attrs{
		"Data": _token("raw text"),
	}}},
	// Literals
	{"single quotes string", "{{ 'test' }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.String{}, "test"),
	}}},
	{"single quotes string with whitespace chars", "{{ '  \n\ttest' }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.String{}, "  \n\ttest"),
	}}},
	{"single quotes string with raw whitespace chars", `{{ '  \n\ttest' }}`, specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.String{}, "  \n\ttest"),
	}}},
	{"double quotes string", `{{ "test" }}`, specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.String{}, "test"),
	}}},
	{"double quotes string with whitespace chars", "{{ \"  \n\ttest\" }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.String{}, "  \n\ttest"),
	}}},
	{"double quotes string with raw whitespace chars", `{{ "  \n\ttest" }}`, specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.String{}, "  \n\ttest"),
	}}},
	{"single quotes inside double quotes string", `{{ "'quoted' test" }}`, specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.String{}, "'quoted' test"),
	}}},
	{"integer", "{{ 42 }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Integer{}, int64(42)),
	}}},
	{"negative-integer", "{{ -42 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.UnaryExpression{}, attrs{
			"Negative": val{true},
			"Term":     _literal(nodes.Integer{}, int64(42)),
		}},
	}}},
	{"float", "{{ 42.0 }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Float{}, float64(42)),
	}}},
	{"negative-float", "{{ -42.0 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.UnaryExpression{}, attrs{
			"Negative": val{true},
			"Term":     _literal(nodes.Float{}, float64(42)),
		}},
	}}},
	{"bool-true", "{{ true }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Bool{}, true),
	}}},
	{"bool-True", "{{ True }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Bool{}, true),
	}}},
	{"bool-false", "{{ false }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Bool{}, false),
	}}},
	{"bool-False", "{{ False }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Bool{}, false),
	}}},
	{"list", "{{ ['list', \"of\", 'objects'] }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.List{}, slice{
			_literal(nodes.String{}, "list"),
			_literal(nodes.String{}, "of"),
			_literal(nodes.String{}, "objects"),
		}),
	}}},
	{"list with trailing coma", "{{ ['list', \"of\", 'objects',] }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.List{}, slice{
			_literal(nodes.String{}, "list"),
			_literal(nodes.String{}, "of"),
			_literal(nodes.String{}, "objects"),
		}),
	}}},
	{"single entry list", "{{ ['list'] }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.List{}, slice{
			_literal(nodes.String{}, "list"),
		}),
	}}},
	{"empty list", "{{ [] }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.List{}, slice{}),
	}}},
	{"tuple", "{{ ('tuple', \"of\", 'objects') }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Tuple{}, slice{
			_literal(nodes.String{}, "tuple"),
			_literal(nodes.String{}, "of"),
			_literal(nodes.String{}, "objects"),
		}),
	}}},
	{"tuple with trailing coma", "{{ ('tuple', \"of\", 'objects',) }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Tuple{}, slice{
			_literal(nodes.String{}, "tuple"),
			_literal(nodes.String{}, "of"),
			_literal(nodes.String{}, "objects"),
		}),
	}}},
	{"single entry tuple", "{{ ('tuple',) }}", specs{nodes.Output{}, attrs{
		"Expression": _literal(nodes.Tuple{}, slice{
			_literal(nodes.String{}, "tuple"),
		}),
	}}},
	{"empty dict", "{{ {} }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Dict{}, attrs{}},
	}}},
	{"dict string", "{{ {'dict': 'of', 'key': 'and', 'value': 'pairs'} }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Dict{}, attrs{
			"Pairs": slice{
				specs{nodes.Pair{}, attrs{
					"Key":   _literal(nodes.String{}, "dict"),
					"Value": _literal(nodes.String{}, "of"),
				}},
				specs{nodes.Pair{}, attrs{
					"Key":   _literal(nodes.String{}, "key"),
					"Value": _literal(nodes.String{}, "and"),
				}},
				specs{nodes.Pair{}, attrs{
					"Key":   _literal(nodes.String{}, "value"),
					"Value": _literal(nodes.String{}, "pairs"),
				}},
			},
		}},
	}}},
	{"dict int", "{{ {1: 'one', 2: 'two', 3: 'three'} }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Dict{}, attrs{
			"Pairs": slice{
				specs{nodes.Pair{}, attrs{
					"Key":   _literal(nodes.Integer{}, int64(1)),
					"Value": _literal(nodes.String{}, "one"),
				}},
				specs{nodes.Pair{}, attrs{
					"Key":   _literal(nodes.Integer{}, int64(2)),
					"Value": _literal(nodes.String{}, "two"),
				}},
				specs{nodes.Pair{}, attrs{
					"Key":   _literal(nodes.Integer{}, int64(3)),
					"Value": _literal(nodes.String{}, "three"),
				}},
			},
		}},
	}}},
	{"addition", "{{ 40 + 2 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left":     _literal(nodes.Integer{}, int64(40)),
			"Right":    _literal(nodes.Integer{}, int64(2)),
			"Operator": _binOp("+"),
		}},
	}}},
	{"multiple additions", "{{ 40 + 1 + 1 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left": specs{nodes.BinaryExpression{}, attrs{
				"Left":     _literal(nodes.Integer{}, int64(40)),
				"Right":    _literal(nodes.Integer{}, int64(1)),
				"Operator": _binOp("+"),
			}},
			"Right":    _literal(nodes.Integer{}, int64(1)),
			"Operator": _binOp("+"),
		}},
	}}},
	{"multiple additions with power", "{{ 40 + 2 ** 1 + 0 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left": specs{nodes.BinaryExpression{}, attrs{
				"Left": _literal(nodes.Integer{}, int64(40)),
				"Right": specs{nodes.BinaryExpression{}, attrs{
					"Left":     _literal(nodes.Integer{}, int64(2)),
					"Right":    _literal(nodes.Integer{}, int64(1)),
					"Operator": _binOp("**"),
				}},
				"Operator": _binOp("+"),
			}},
			"Right":    _literal(nodes.Integer{}, int64(0)),
			"Operator": _binOp("+"),
		}},
	}}},
	{"substract", "{{ 40 - 2 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left":     _literal(nodes.Integer{}, int64(40)),
			"Right":    _literal(nodes.Integer{}, int64(2)),
			"Operator": _binOp("-"),
		}},
	}}},
	{"complex math", "{{ -1 * (-(-(10-100)) ** 2) ** 3 + 3 * (5 - 17) + 1 + 2 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left": specs{nodes.BinaryExpression{}, attrs{
				"Left": specs{nodes.BinaryExpression{}, attrs{
					"Left": specs{nodes.BinaryExpression{}, attrs{
						"Left": specs{nodes.UnaryExpression{}, attrs{
							"Negative": val{true},
							"Term":     _literal(nodes.Integer{}, int64(1)),
						}},
						"Right": specs{nodes.BinaryExpression{}, attrs{
							"Left": specs{nodes.UnaryExpression{}, attrs{
								"Negative": val{true},
								"Term": specs{nodes.BinaryExpression{}, attrs{
									"Left": specs{nodes.UnaryExpression{}, attrs{
										"Negative": val{true},
										"Term": specs{nodes.BinaryExpression{}, attrs{
											"Left":     _literal(nodes.Integer{}, int64(10)),
											"Right":    _literal(nodes.Integer{}, int64(100)),
											"Operator": _binOp("-"),
										}},
									}},
									"Right":    _literal(nodes.Integer{}, int64(2)),
									"Operator": _binOp("**"),
								}},
							}},
							"Right":    _literal(nodes.Integer{}, int64(3)),
							"Operator": _binOp("**"),
						}},
						"Operator": _binOp("*"),
					}},
					"Right": specs{nodes.BinaryExpression{}, attrs{
						"Left": _literal(nodes.Integer{}, int64(3)),
						"Right": specs{nodes.BinaryExpression{}, attrs{
							"Left":     _literal(nodes.Integer{}, int64(5)),
							"Right":    _literal(nodes.Integer{}, int64(17)),
							"Operator": _binOp("-"),
						}},
						"Operator": _binOp("*"),
					}},
					"Operator": _binOp("+"),
				}},
				"Right":    _literal(nodes.Integer{}, int64(1)),
				"Operator": _binOp("+"),
			}},
			"Right":    _literal(nodes.Integer{}, int64(2)),
			"Operator": _binOp("+"),
		}},
	}}},
	{"negative-expression", "{{ -(40 + 2) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.UnaryExpression{}, attrs{
			"Negative": val{true},
			"Term": specs{nodes.BinaryExpression{}, attrs{
				"Left":     _literal(nodes.Integer{}, int64(40)),
				"Right":    _literal(nodes.Integer{}, int64(2)),
				"Operator": _binOp("+"),
			}},
		}},
	}}},
	{"Operators precedence", "{{ 2 * 3 + 4 % 2 + 1 - 2 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left": specs{nodes.BinaryExpression{}, attrs{
				"Left": specs{nodes.BinaryExpression{}, attrs{
					"Left": specs{nodes.BinaryExpression{}, attrs{
						"Left":     _literal(nodes.Integer{}, int64(2)),
						"Right":    _literal(nodes.Integer{}, int64(3)),
						"Operator": _binOp("*"),
					}},
					"Right": specs{nodes.BinaryExpression{}, attrs{
						"Left":     _literal(nodes.Integer{}, int64(4)),
						"Right":    _literal(nodes.Integer{}, int64(2)),
						"Operator": _binOp("%"),
					}},
					"Operator": _binOp("+"),
				}},
				"Right":    _literal(nodes.Integer{}, int64(1)),
				"Operator": _binOp("+"),
			}},
			"Right":    _literal(nodes.Integer{}, int64(2)),
			"Operator": _binOp("-"),
		}},
	}}},
	{"Operators precedence with parenthesis", "{{ 2 * (3 + 4) % 2 + (1 - 2) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left": specs{nodes.BinaryExpression{}, attrs{
				"Left": specs{nodes.BinaryExpression{}, attrs{
					"Left": _literal(nodes.Integer{}, int64(2)),
					"Right": specs{nodes.BinaryExpression{}, attrs{
						"Left":     _literal(nodes.Integer{}, int64(3)),
						"Right":    _literal(nodes.Integer{}, int64(4)),
						"Operator": _binOp("+"),
					}},
					"Operator": _binOp("*"),
				}},
				"Right":    _literal(nodes.Integer{}, int64(2)),
				"Operator": _binOp("%"),
			}},
			"Right": specs{nodes.BinaryExpression{}, attrs{
				"Left":     _literal(nodes.Integer{}, int64(1)),
				"Right":    _literal(nodes.Integer{}, int64(2)),
				"Operator": _binOp("-"),
			}},
			"Operator": _binOp("+"),
		}},
	}}},
	{"variable", "{{ a_var }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Name{}, attrs{
			"Name": _token("a_var"),
		}},
	}}},
	{"variable attribute", "{{ a_var.attr }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Getattr{}, attrs{
			"Node": specs{nodes.Name{}, attrs{
				"Name": _token("a_var"),
			}},
			"Attr": val{"attr"},
		}},
	}}},
	{"variable item", "{{ a_var[\"item\"] }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Getitem{}, attrs{
			"Node": specs{nodes.Name{}, attrs{
				"Name": _token("a_var"),
			}},
			"Arg": val{"item"},
		}},
	}}},
	{"variable and filter", "{{ a_var|safe }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.FilteredExpression{}, attrs{
			"Expression": specs{nodes.Name{}, attrs{
				"Name": _token("a_var"),
			}},
			"Filters": slice{
				filter{"safe", slice{}, attrs{}},
			},
		}},
	}}},
	{"integer and filter", "{{ 42|safe }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.FilteredExpression{}, attrs{
			"Expression": _literal(nodes.Integer{}, int64(42)),
			"Filters": slice{
				filter{"safe", slice{}, attrs{}},
			},
		}},
	}}},
	{"negative integer and filter", "{{ -42|safe }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.FilteredExpression{}, attrs{
			"Expression": specs{nodes.UnaryExpression{}, attrs{
				"Negative": val{true},
				"Term":     _literal(nodes.Integer{}, int64(42)),
			}},
			"Filters": slice{
				filter{"safe", slice{}, attrs{}},
			},
		}},
	}}},
	{"logical expressions", "{{ true and false }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left":     _literal(nodes.Bool{}, true),
			"Right":    _literal(nodes.Bool{}, false),
			"Operator": _binOp("and"),
		}},
	}}},
	{"negated boolean", "{{ not true }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Negation{}, attrs{
			"Term": _literal(nodes.Bool{}, true),
		}},
	}}},
	{"negated logical expression", "{{ not false and true }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left": specs{nodes.Negation{}, attrs{
				"Term": _literal(nodes.Bool{}, false),
			}},
			"Right":    _literal(nodes.Bool{}, true),
			"Operator": _binOp("and"),
		}},
	}}},
	{"negated logical expression with parenthesis", "{{ not (false and true) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Negation{}, attrs{
			"Term": specs{nodes.BinaryExpression{}, attrs{
				"Left":     _literal(nodes.Bool{}, false),
				"Right":    _literal(nodes.Bool{}, true),
				"Operator": _binOp("and"),
			}},
		}},
	}}},
	{"logical expression with math comparison", "{{ 40 + 2 > 5 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left": specs{nodes.BinaryExpression{}, attrs{
				"Left":     _literal(nodes.Integer{}, int64(40)),
				"Right":    _literal(nodes.Integer{}, int64(2)),
				"Operator": _binOp("+"),
			}},
			"Right":    _literal(nodes.Integer{}, int64(5)),
			"Operator": _binOp(">"),
		}},
	}}},
	{"logical expression with filter", "{{ false and true|safe }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.BinaryExpression{}, attrs{
			"Left": _literal(nodes.Bool{}, false),
			"Right": specs{nodes.FilteredExpression{}, attrs{
				"Expression": _literal(nodes.Bool{}, true),
				"Filters": slice{
					filter{"safe", slice{}, attrs{}},
				},
			}},
			"Operator": _binOp("and"),
		}},
	}}},
	{"logical expression with parenthesis and filter", "{{ (false and true)|safe }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.FilteredExpression{}, attrs{
			"Expression": specs{nodes.BinaryExpression{}, attrs{
				"Left":     _literal(nodes.Bool{}, false),
				"Right":    _literal(nodes.Bool{}, true),
				"Operator": _binOp("and"),
			}},
			"Filters": slice{
				filter{"safe", slice{}, attrs{}},
			},
		}},
	}}},
	{"function", "{{ a_func(42) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Call{}, attrs{
			"Func": specs{nodes.Name{}, attrs{"Name": _token("a_func")}},
			"Args": slice{_literal(nodes.Integer{}, int64(42))},
		}},
	}}},
	{"method", "{{ an_obj.a_method(42) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Call{}, attrs{
			"Func": specs{nodes.Getattr{}, attrs{
				"Node": specs{nodes.Name{}, attrs{"Name": _token("an_obj")}},
				"Attr": val{"a_method"},
			}},
			"Args": slice{_literal(nodes.Integer{}, int64(42))},
		}},
	}}},
	{"function with filtered args", "{{ a_func(42|safe) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.Call{}, attrs{
			"Func": specs{nodes.Name{}, attrs{"Name": _token("a_func")}},
			"Args": slice{
				specs{nodes.FilteredExpression{}, attrs{
					"Expression": _literal(nodes.Integer{}, int64(42)),
					"Filters": slice{
						filter{"safe", slice{}, attrs{}},
					},
				}},
			},
		}},
	}}},
	{"variable and multiple filters", "{{ a_var|add(42)|safe }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.FilteredExpression{}, attrs{
			"Expression": specs{nodes.Name{}, attrs{"Name": _token("a_var")}},
			"Filters": slice{
				filter{"add", slice{_literal(nodes.Integer{}, int64(42))}, attrs{}},
				filter{"safe", slice{}, attrs{}},
			},
		}},
	}}},
	{"variable and expression filters", "{{ a_var|add(40 + 2) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.FilteredExpression{}, attrs{
			"Expression": specs{nodes.Name{}, attrs{"Name": _token("a_var")}},
			"Filters": slice{
				filter{"add", slice{
					specs{nodes.BinaryExpression{}, attrs{
						"Left":     _literal(nodes.Integer{}, int64(40)),
						"Right":    _literal(nodes.Integer{}, int64(2)),
						"Operator": _binOp("+"),
					}},
				}, attrs{}},
			},
		}},
	}}},
	{"variable and nested filters", "{{ a_var|add( 42|add(2) ) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.FilteredExpression{}, attrs{
			"Expression": specs{nodes.Name{}, attrs{"Name": _token("a_var")}},
			"Filters": slice{
				filter{"add", slice{
					specs{nodes.FilteredExpression{}, attrs{
						"Expression": _literal(nodes.Integer{}, int64(42)),
						"Filters": slice{
							filter{"add", slice{_literal(nodes.Integer{}, int64(2))}, attrs{}},
						},
					}},
				}, attrs{}},
			},
		}},
	}}},
	{"Test equal", "{{ 3 is equal 3 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.TestExpression{}, attrs{
			"Expression": _literal(nodes.Integer{}, int64(3)),
			"Test": specs{nodes.TestCall{}, attrs{
				"Name": val{"equal"},
				"Args": slice{_literal(nodes.Integer{}, int64(3))},
			}},
		}},
	}}},
	{"Test equal parenthesis", "{{ 3 is equal(3) }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.TestExpression{}, attrs{
			"Expression": _literal(nodes.Integer{}, int64(3)),
			"Test": specs{nodes.TestCall{}, attrs{
				"Name": val{"equal"},
				"Args": slice{_literal(nodes.Integer{}, int64(3))},
			}},
		}},
	}}},
	{"Test ==", "{{ 3 is == 3 }}", specs{nodes.Output{}, attrs{
		"Expression": specs{nodes.TestExpression{}, attrs{
			"Expression": _literal(nodes.Integer{}, int64(3)),
			"Test": specs{nodes.TestCall{}, attrs{
				"Name": val{"=="},
				"Args": slice{_literal(nodes.Integer{}, int64(3))},
			}},
		}},
	}}},
}

// func parseText(text string) (*nodeDocument, *Error) {
// 	tokens, err := lex("test", text)
// 	if err != nil {
// 		return nil, err
// 	}
// 	parser := newParser("test", tokens, &Template{
// 		set: &TemplateSet{},
// 	})
// 	return parser.parseDocument()
// }

func _deref(value reflect.Value) reflect.Value {
	for (value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr) && !value.IsNil() {
		value = value.Elem()
	}
	return value
}

type asserter interface {
	assert(t *testing.T, value reflect.Value)
}

type specs struct {
	typ   interface{}
	attrs attrs
}

func (specs specs) assert(t *testing.T, value reflect.Value) {
	assert := assert.New(t)
	value = _deref(value)
	// t.Logf("type(expected %+v, actual %+v)", reflect.TypeOf(specs.typ), value.Type())
	if !assert.Equal(reflect.TypeOf(specs.typ), value.Type()) {
		return
	}
	if specs.attrs != nil {
		specs.attrs.assert(t, value)
	}
}

type val struct {
	value interface{}
}

func (val val) assert(t *testing.T, value reflect.Value) {
	assert := assert.New(t)
	value = _deref(value)
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		assert.Equal(val.value, value.Int())
	case reflect.Float32, reflect.Float64:
		assert.Equal(val.value, value.Float())
	case reflect.String:
		assert.Equal(val.value, value.String())
	case reflect.Bool:
		assert.Equal(val.value, value.Bool())
	case reflect.Slice:
		current, ok := val.value.(asserter)
		if assert.True(ok) {
			current.assert(t, value)
		}
	case reflect.Map:
		assert.Len(val.value, value.Len())
		v2 := reflect.ValueOf(val.value)

		iter := value.MapRange()
		for iter.Next() {
			assert.Equal(iter.Value(), v2.MapIndex(iter.Key()))
		}
	case reflect.Func:
		assert.Equal(value, reflect.ValueOf(val.value))
	default:
		assert.Failf("Unknown value", "Unknown value kind '%s'", value.Kind())
	}
}

func _literal(typ interface{}, value interface{}) asserter {
	return specs{typ, attrs{
		"Val": val{value},
	}}
}

func _token(value string) asserter {
	return specs{tokens.Token{}, attrs{
		"Val": val{value},
	}}
}

func _binOp(value string) asserter {
	return specs{nodes.BinOperator{}, attrs{
		"Token": _token(value),
	}}
}

type attrs map[string]asserter

func (attrs attrs) assert(t *testing.T, value reflect.Value) {
	assert := assert.New(t)
	for attr, specs := range attrs {
		field := value.FieldByName(attr)
		if assert.True(field.IsValid(), fmt.Sprintf("No field named '%s' found", attr)) {
			specs.assert(t, field)
		}
	}
}

type slice []asserter

func (slice slice) assert(t *testing.T, value reflect.Value) {
	if assert.Equal(t, reflect.Slice, value.Kind()) {
		if assert.Equal(t, len(slice), value.Len()) {
			for idx, specs := range slice {
				specs.assert(t, value.Index(idx))
			}
		}
	}
}

type filter struct {
	name   string
	args   slice
	kwargs attrs
}

func (filter filter) assert(t *testing.T, value reflect.Value) {
	value = _deref(value)
	assert := assert.New(t)
	assert.Equal(reflect.TypeOf(nodes.FilterCall{}), value.Type())
	assert.Equal(filter.name, value.FieldByName("Name").String())
	args := value.FieldByName("Args")
	kwargs := value.FieldByName("Kwargs")
	if assert.Equal(len(filter.args), args.Len()) {
		for idx, specs := range filter.args {
			specs.assert(t, args.Index(idx))
		}
	}
	if assert.Equal(len(filter.kwargs), kwargs.Len()) {
		for key, specs := range filter.kwargs {
			specs.assert(t, args.MapIndex(reflect.ValueOf(key)))
		}
	}
}

func TestParser(t *testing.T) {
	for _, tc := range testCases {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			// t.Parallel()
			assert := assert.New(t)
			tpl, err := parser.Parse(test.text)
			if assert.Nil(err, "Unable to parse template: %s", err) {
				if assert.Equal(1, len(tpl.Nodes), "Expected one node") {
					test.expected.assert(t, reflect.ValueOf(tpl.Nodes[0]))
				} else {
					t.Logf("Nodes %+v", tpl.Nodes)
				}
			}
		})
	}
}
