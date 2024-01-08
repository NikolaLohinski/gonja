package parser_test

import (
	"fmt"
	"strconv"

	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

var _ = Context("parser", func() {
	var (
		input = new(string)

		returnedTemplate = new(nodes.Template)
		returnedError    = new(error)
	)
	JustBeforeEach(func() {
		stream := tokens.Lex(*input, config.New())
		returnedTemplate, *returnedError = parser.NewParser("tests", stream, config.New(), loaders.MustNewFileSystemLoader(""), builtins.ControlStructures).Parse()
	})
	for _, testCase := range []struct {
		description  string
		inputs       []string
		nodeMatchers []types.GomegaMatcher
	}{
		{
			"is a comment",
			[]string{"{# My comment #}"},
			[]types.GomegaMatcher{
				And(
					BeAssignableToTypeOf(nodes.Comment{}),
					MatchFields(IgnoreExtras, Fields{"Text": Equal(" My comment ")}),
				),
			},
		},
		{
			"is a multiline comment",
			[]string{"{# My\nmultiline\ncomment #}"},
			[]types.GomegaMatcher{
				And(
					BeAssignableToTypeOf(nodes.Comment{}),
					MatchFields(IgnoreExtras, Fields{"Text": Equal(" My\nmultiline\ncomment ")}),
				),
			},
		},
		{
			"is raw text data",
			[]string{"raw text"},
			[]types.GomegaMatcher{
				And(
					BeAssignableToTypeOf(nodes.Data{}),
					MatchFields(IgnoreExtras, Fields{
						"Data": PointTo(MatchFields(IgnoreExtras, Fields{
							"Type": Equal(tokens.Data),
							"Val":  Equal("raw text"),
						})),
					}),
				),
			},
		},
		{
			"is a single quoted string",
			[]string{"{{ 'test' }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(MatchStringNode("test")),
			},
		},
		{
			"is a single quoted string with nested spaces, lines and tabs",
			[]string{"{{ '  \r\n\ttest' }}", `{{ '  \r\n\ttest' }}`},
			[]types.GomegaMatcher{
				MatchNodeOutput(MatchStringNode("  \r\n\ttest")),
			},
		},
		{
			"is an integer",
			[]string{`{{ 42 }}`},
			[]types.GomegaMatcher{
				MatchNodeOutput(MatchIntegerNode(42)),
			},
		},
		{
			"is a negative integer",
			[]string{`{{ -42 }}`},
			[]types.GomegaMatcher{
				MatchNodeOutput(MatchUnaryExpression(tokens.Subtraction, MatchIntegerNode(42))),
			},
		},
		{
			"is a float",
			[]string{`{{ 1.0 }}`},
			[]types.GomegaMatcher{
				MatchNodeOutput(And(
					BeAssignableToTypeOf(nodes.Float{}),
					MatchFields(IgnoreExtras, Fields{
						"Val": Equal(1.0),
					}),
				),
				),
			},
		},
		{
			"is a true boolean",
			[]string{"{{ true }}", "{{ True }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(MatchNodeBool(true)),
			},
		},
		{
			"is a false boolean",
			[]string{"{{ false }}", "{{ False }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(MatchNodeBool(false)),
			},
		},
		{
			"is a list",
			[]string{"{{ ['a', 'b'] }}", "{{ [\"a\", \"b\"] }}", "{{ ['a', 'b', ] }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					And(
						BeAssignableToTypeOf(nodes.List{}),
						MatchFields(IgnoreExtras, Fields{
							"Val": MatchAllElementsWithIndex(func(index int, _ interface{}) string {
								return strconv.Itoa(index)
							}, Elements{
								"0": PointTo(MatchStringNode("a")),
								"1": PointTo(MatchStringNode("b")),
							}),
						}),
					),
				),
			},
		},
		{
			"is an empty list",
			[]string{"{{ [] }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					And(
						BeAssignableToTypeOf(nodes.List{}),
						MatchFields(IgnoreExtras, Fields{
							"Val": BeEmpty(),
						}),
					),
				),
			},
		},
		{
			"is a tuple",
			[]string{"{{ ('a', 'b') }}", "{{ (\"a\", \"b\") }}", "{{ ('a', 'b', ) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					And(
						BeAssignableToTypeOf(nodes.Tuple{}),
						MatchFields(IgnoreExtras, Fields{
							"Val": MatchAllElementsWithIndex(func(index int, _ interface{}) string {
								return strconv.Itoa(index)
							}, Elements{
								"0": PointTo(MatchStringNode("a")),
								"1": PointTo(MatchStringNode("b")),
							}),
						}),
					),
				),
			},
		},
		{
			"is an empty dict",
			[]string{"{{ {} }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					And(
						BeAssignableToTypeOf(nodes.Dict{}),
						MatchFields(IgnoreExtras, Fields{
							"Pairs": BeEmpty(),
						}),
					),
				),
			},
		},
		{
			"is an inline dict with string keys",
			[]string{"{{ {'foo': 'bar'} }}", "{{ {\"foo\": \"bar\"} }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					And(
						BeAssignableToTypeOf(nodes.Dict{}),
						MatchFields(IgnoreExtras, Fields{
							"Pairs": MatchAllElementsWithIndex(func(index int, _ interface{}) string {
								return strconv.Itoa(index)
							}, Elements{
								"0": PointTo(And(
									BeAssignableToTypeOf(nodes.Pair{}),
									MatchFields(IgnoreExtras, Fields{
										"Key":   PointTo(MatchStringNode("foo")),
										"Value": PointTo(MatchStringNode("bar")),
									}),
								)),
							}),
						}),
					),
				),
			},
		},
		{
			"is an inline dict with integer keys",
			[]string{"{{ {1: 2} }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					And(
						BeAssignableToTypeOf(nodes.Dict{}),
						MatchFields(IgnoreExtras, Fields{
							"Pairs": MatchAllElementsWithIndex(func(index int, _ interface{}) string {
								return strconv.Itoa(index)
							}, Elements{
								"0": PointTo(And(
									BeAssignableToTypeOf(nodes.Pair{}),
									MatchFields(IgnoreExtras, Fields{
										"Key":   PointTo(MatchIntegerNode(1)),
										"Value": PointTo(MatchIntegerNode(2)),
									}),
								)),
							}),
						}),
					),
				),
			},
		},
		{
			"is an addition",
			[]string{"{{ 1 + 2 }}", "{{ (1 + 2) }}", "{{ 1 + (2) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchIntegerNode(1),
						tokens.Addition,
						MatchIntegerNode(2),
					),
				),
			},
		},
		{
			"is a multi-addition",
			[]string{"{{ 1 + 2 + 3 }}", "{{ (1 + 2) + 3 }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeBinaryExpression(
							MatchIntegerNode(1),
							tokens.Addition,
							MatchIntegerNode(2),
						),
						tokens.Addition,
						MatchIntegerNode(3),
					),
				),
			},
		},
		{
			"is a multi-addition and power operation",
			[]string{"{{ 1 + 2 ** 3 + 4 }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeBinaryExpression(
							MatchIntegerNode(1),
							tokens.Addition,
							MatchNodeBinaryExpression(
								MatchIntegerNode(2),
								tokens.Power,
								MatchIntegerNode(3),
							),
						),
						tokens.Addition,
						MatchIntegerNode(4),
					),
				),
			},
		},
		{
			"is a subtraction",
			[]string{"{{ 1 - 2 }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchIntegerNode(1),
						tokens.Subtraction,
						MatchIntegerNode(2),
					),
				),
			},
		},
		{
			"is a complex math expression",
			[]string{"{{ -1 * (-(-(10-100)) ** 2) ** 3 + 3 * (5 - 17) + 1 + 2 }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeBinaryExpression(
							MatchNodeBinaryExpression(
								MatchNodeBinaryExpression(
									MatchUnaryExpression(
										tokens.Subtraction,
										MatchIntegerNode(1),
									),
									tokens.Multiply,
									MatchNodeBinaryExpression(
										MatchUnaryExpression(
											tokens.Subtraction,
											MatchNodeBinaryExpression(
												MatchUnaryExpression(
													tokens.Subtraction,
													MatchNodeBinaryExpression(
														MatchIntegerNode(10),
														tokens.Subtraction,
														MatchIntegerNode(100),
													)),
												tokens.Power,
												MatchIntegerNode(2),
											)),
										tokens.Power,
										MatchIntegerNode(3),
									),
								),
								tokens.Addition,
								MatchNodeBinaryExpression(
									MatchIntegerNode(3),
									tokens.Multiply,
									MatchNodeBinaryExpression(
										MatchIntegerNode(5),
										tokens.Subtraction,
										MatchIntegerNode(17),
									)),
							),
							tokens.Addition,
							MatchIntegerNode(1),
						),
						tokens.Addition,
						MatchIntegerNode(2),
					),
				),
			},
		},
		{
			"is a negative expression",
			[]string{"{{ -(1 + 2) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchUnaryExpression(
						tokens.Subtraction,
						MatchNodeBinaryExpression(
							MatchIntegerNode(1),
							tokens.Addition,
							MatchIntegerNode(2),
						),
					),
				),
			},
		},
		{
			"is a combination of all operators",
			[]string{"{{ 2 * 3 + 4 % 2 + 1 - 2 }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeBinaryExpression(
							MatchNodeBinaryExpression(
								MatchNodeBinaryExpression(
									MatchIntegerNode(2),
									tokens.Multiply,
									MatchIntegerNode(3),
								),
								tokens.Addition,
								MatchNodeBinaryExpression(
									MatchIntegerNode(4),
									tokens.Modulo,
									MatchIntegerNode(2),
								),
							),
							tokens.Addition,
							MatchIntegerNode(1),
						),
						tokens.Subtraction,
						MatchIntegerNode(2),
					),
				),
			},
		},
		{
			"is a combination of all operators with parenthesis precedence",
			[]string{"{{ 2 * (3 + 4) % 2 + (1 - 2) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeBinaryExpression(
							MatchNodeBinaryExpression(
								MatchIntegerNode(2),
								tokens.Multiply,
								MatchNodeBinaryExpression(
									MatchIntegerNode(3),
									tokens.Addition,
									MatchIntegerNode(4),
								),
							),
							tokens.Modulo,
							MatchIntegerNode(2),
						),
						tokens.Addition,
						MatchNodeBinaryExpression(
							MatchIntegerNode(1),
							tokens.Subtraction,
							MatchIntegerNode(2),
						),
					),
				),
			},
		},
		{
			"is a variable",
			[]string{"{{ var }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNameNode("var"),
				),
			},
		},
		{
			"is a variable attribute access",
			[]string{"{{ var.attribute }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchGetAttributeNode(
						MatchNameNode("var"),
						"attribute",
					),
				),
			},
		},
		{
			"is a variable attribute integer access",
			[]string{"{{ var.1 }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchGetAttributeNode(
						MatchNameNode("var"),
						1,
					),
				),
			},
		},
		{
			"is a variable item access with a string",
			[]string{"{{ var['item'] }}", `{{ var["item"] }}`, `{{ var[("item")] }}`},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchGetItemNode(
						MatchNameNode("var"),
						MatchStringNode("item"),
					),
				),
			},
		},
		{
			"is a variable item access with an integer",
			[]string{"{{ var[123] }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchGetItemNode(
						MatchNameNode("var"),
						MatchIntegerNode(123),
					),
				),
			},
		},
		{
			"is a variable item access with a variable",
			[]string{"{{ var[other] }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchGetItemNode(
						MatchNameNode("var"),
						MatchNameNode("other"),
					),
				),
			},
		},
		{
			"is a variable passed through a filter",
			[]string{"{{ var | filter }}", "{{ var|filter }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeFilteredExpressionNode(
						MatchNameNode("var"),
						MatchFilterCallNode("filter", nil, nil),
					),
				),
			},
		},
		{
			"is an integer passed through two filters",
			[]string{"{{ 1 | first | second }}", "{{ 1|first|second }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeFilteredExpressionNode(
						MatchIntegerNode(1),
						MatchFilterCallNode("first", nil, nil),
						MatchFilterCallNode("second", nil, nil),
					),
				),
			},
		},
		{
			"is an integer passed through a filter with a positional argument",
			[]string{"{{ 1 | filter(arg) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeFilteredExpressionNode(
						MatchIntegerNode(1),
						MatchFilterCallNode(
							"filter",
							[]types.GomegaMatcher{
								PointTo(MatchNameNode("arg")),
							},
							nil,
						),
					),
				),
			},
		},
		{
			"is an integer passed through a filter with a keyword argument",
			[]string{"{{ 1 | filter(name=arg) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeFilteredExpressionNode(
						MatchIntegerNode(1),
						MatchFilterCallNode(
							"filter",
							nil,
							map[string]types.GomegaMatcher{
								"name": PointTo(MatchNameNode("arg")),
							},
						),
					),
				),
			},
		},
		{
			"is an integer passed through a filter with a keyword argument and a positional argument",
			[]string{"{{ 1 | filter(\"str\", name=arg) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeFilteredExpressionNode(
						MatchIntegerNode(1),
						MatchFilterCallNode(
							"filter",
							[]types.GomegaMatcher{
								PointTo(MatchStringNode("str")),
							},
							map[string]types.GomegaMatcher{
								"name": PointTo(MatchNameNode("arg")),
							},
						),
					),
				),
			},
		},
		{
			"is a logical expression",
			[]string{"{{ true and false }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeBool(true),
						tokens.And,
						MatchNodeBool(false),
					),
				),
			},
		},
		{
			"is a negated boolean",
			[]string{"{{ not false }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeNegation(
						MatchNodeBool(false),
					),
				),
			},
		},
		{
			"is a negation over a portion of a logical expression",
			[]string{"{{ not false and true }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeNegation(
							MatchNodeBool(false),
						),
						tokens.And,
						MatchNodeBool(true),
					),
				),
			},
		},
		{
			"is a negation over of a logical expression",
			[]string{"{{ not (false and true) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeNegation(
						MatchNodeBinaryExpression(
							MatchNodeBool(false),
							tokens.And,
							MatchNodeBool(true),
						),
					),
				),
			},
		},
		{
			"is an arithmetic expression with a comparison",
			[]string{"{{ 40 + 2 > 5 }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeBinaryExpression(
							MatchIntegerNode(40),
							tokens.Addition,
							MatchIntegerNode(2),
						),
						tokens.GreaterThan,
						MatchIntegerNode(5),
					),
				),
			},
		},
		{
			"is a logical expression with a filter",
			[]string{"{{ true and false | filter }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeBinaryExpression(
						MatchNodeBool(true),
						tokens.And,
						MatchNodeFilteredExpressionNode(
							MatchNodeBool(false),
							MatchFilterCallNode("filter", nil, nil),
						),
					),
				),
			},
		},
		{
			"is a function call with both positional and named arguments",
			[]string{"{{ func(101, name=arg) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchCallNode(
						MatchNameNode("func"),
						[]types.GomegaMatcher{
							PointTo(MatchIntegerNode(101)),
						},
						map[string]types.GomegaMatcher{
							"name": PointTo(MatchNameNode("arg")),
						},
					),
				),
			},
		},
		{
			"is a function call with an argument passed through a filter",
			[]string{"{{ func('filter me' | filter) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchCallNode(
						MatchNameNode("func"),
						[]types.GomegaMatcher{
							PointTo(MatchNodeFilteredExpressionNode(
								MatchStringNode("filter me"),
								MatchFilterCallNode("filter", nil, nil),
							)),
						},
						nil,
					),
				),
			},
		},
		{
			"is a function call from an object attribute",
			[]string{"{{ object.func() }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchCallNode(
						MatchGetAttributeNode(
							MatchNameNode("object"),
							"func",
						),
						nil,
						nil,
					),
				),
			},
		},
		{
			"is a function call with math expression as an argument",
			[]string{"{{ func(1 + 2) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchCallNode(
						MatchNameNode("func"),
						[]types.GomegaMatcher{
							PointTo(
								MatchNodeBinaryExpression(
									MatchIntegerNode(1),
									tokens.Addition,
									MatchIntegerNode(2),
								),
							),
						},
						nil,
					),
				),
			},
		},
		{
			"is a filtered variable call with a filtered argument",
			[]string{"{{ first | func(arg | second) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchNodeFilteredExpressionNode(
						MatchNameNode("first"),
						MatchFilterCallNode(
							"func",
							[]types.GomegaMatcher{
								PointTo(MatchNodeFilteredExpressionNode(
									MatchNameNode("arg"),
									MatchFilterCallNode("second", nil, nil),
								)),
							},
							nil,
						),
					),
				),
			},
		},
		{
			"is an equality test",
			[]string{"{{ 3 is equal 3 }}", "{{ 3 is equal(3) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchTestExpressionNode(
						MatchIntegerNode(3),
						MatchTestCall(
							"equal",
							[]types.GomegaMatcher{PointTo(MatchIntegerNode(3))},
							nil,
						),
					),
				),
			},
		},
		{
			"is an == test",
			[]string{"{{ 3 is == 3 }}", "{{ 3 is ==(3) }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchTestExpressionNode(
						MatchIntegerNode(3),
						MatchTestCall(
							"==",
							[]types.GomegaMatcher{PointTo(MatchIntegerNode(3))},
							nil,
						),
					),
				),
			},
		},
		{
			"is an inline if condition",
			[]string{"{{ 'foo' if 2 is odd }}"},
			[]types.GomegaMatcher{
				MatchNodeConditionalOutput(
					MatchStringNode("foo"),
					PointTo(
						MatchTestExpressionNode(
							MatchIntegerNode(2),
							MatchTestCall("odd", nil, nil),
						),
					),
					BeNil(),
				),
			},
		},
		{
			"is an inline if/else condition",
			[]string{"{{ 'foo' if 2 is odd else 'bar' }}"},
			[]types.GomegaMatcher{
				MatchNodeConditionalOutput(
					MatchStringNode("foo"),
					PointTo(
						MatchTestExpressionNode(
							MatchIntegerNode(2),
							MatchTestCall("odd", nil, nil),
						),
					),
					PointTo(MatchStringNode("bar")),
				),
			},
		},
		{
			"is a getter after an expression",
			[]string{"{{ (2 is odd).one }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchGetAttributeNode(
						MatchTestExpressionNode(
							MatchIntegerNode(2),
							MatchTestCall("odd", nil, nil),
						),
						"one",
					),
				),
			},
		},
		{
			"is a getter after an expression",
			[]string{"{{ (2 is odd)[\"one\"] }}"},
			[]types.GomegaMatcher{
				MatchNodeOutput(
					MatchGetItemNode(
						MatchTestExpressionNode(
							MatchIntegerNode(2),
							MatchTestCall("odd", nil, nil),
						),
						MatchStringNode("one"),
					),
				),
			},
		},
	} {
		t := testCase
		Context(fmt.Sprintf("when the input %s", t.description), func() {
			for _, testInput := range t.inputs {
				i := testInput
				Context(i, func() {
					BeforeEach(func() {
						*input = i
					})
					elements := make(Elements)
					for index, matcher := range t.nodeMatchers {
						elements[strconv.Itoa(index)] = PointTo(matcher)
					}
					It("should return the expected tree", func() {
						Expect(*returnedError).To(BeNil())
						Expect(returnedTemplate).To(PointTo(MatchFields(IgnoreExtras, Fields{
							"Identifier": Equal("tests"),
							"Nodes": MatchAllElementsWithIndex(
								func(index int, _ interface{}) string { return strconv.Itoa(index) },
								elements,
							),
						})))
					})
				})
			}

		})
	}
})
