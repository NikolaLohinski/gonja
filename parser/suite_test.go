package parser_test

import (
	"strconv"
	"testing"

	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/tokens"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "parser")
}

func MatchNodeOutput(expression types.GomegaMatcher) types.GomegaMatcher {
	return MatchNodeConditionalOutput(expression, BeNil(), BeNil())
}

func MatchNodeConditionalOutput(expression, condition, alternative types.GomegaMatcher) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.Output{}),
		MatchFields(IgnoreExtras, Fields{
			"Expression":  PointTo(expression),
			"Condition":   condition,
			"Alternative": alternative,
		}),
	)
}

func MatchNodeNegation(term types.GomegaMatcher) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.Negation{}),
		MatchFields(IgnoreExtras, Fields{
			"Term": PointTo(term),
			"Operator": PointTo(And(
				BeAssignableToTypeOf(tokens.Token{}),
				MatchFields(IgnoreExtras, Fields{
					"Type": Equal(tokens.Not),
				}),
			)),
		}),
	)
}

func MatchNodeBinOperator(operator tokens.Type) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.BinOperator{}),
		MatchFields(IgnoreExtras, Fields{
			"Token": PointTo(And(
				BeAssignableToTypeOf(tokens.Token{}),
				MatchFields(IgnoreExtras, Fields{
					"Type": Equal(operator),
				}),
			)),
		}),
	)
}

func MatchNodeBinaryExpression(left types.GomegaMatcher, operator tokens.Type, right types.GomegaMatcher) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.BinaryExpression{}),
		MatchFields(IgnoreExtras, Fields{
			"Left":     PointTo(left),
			"Operator": PointTo(MatchNodeBinOperator(operator)),
			"Right":    PointTo(right),
		}),
	)
}

func MatchUnaryExpression(operator tokens.Type, term types.GomegaMatcher) types.GomegaMatcher {
	negative := false
	if operator == tokens.Subtraction {
		negative = true
	}
	return And(
		BeAssignableToTypeOf(nodes.UnaryExpression{}),
		MatchFields(IgnoreExtras, Fields{
			"Negative": Equal(negative),
			"Operator": PointTo(And(
				BeAssignableToTypeOf(tokens.Token{}),
				MatchFields(IgnoreExtras, Fields{
					"Type": Equal(operator),
				}),
			)),
			"Term": PointTo(term),
		}),
	)
}

func MatchNodeBool(boolean bool) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.Bool{}),
		MatchFields(IgnoreExtras, Fields{
			"Val": Equal(boolean),
		}),
	)
}

func MatchIntegerNode(integer int) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.Integer{}),
		MatchFields(IgnoreExtras, Fields{
			"Val": Equal(integer),
		}),
	)
}

func MatchStringNode(str string) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.String{}),
		MatchFields(IgnoreExtras, Fields{
			"Val": Equal(str),
		}),
	)
}

func MatchNameNode(name string) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.Name{}),
		MatchFields(IgnoreExtras, Fields{
			"Name": PointTo(And(
				BeAssignableToTypeOf(tokens.Token{}),
				MatchFields(IgnoreExtras, Fields{
					"Val": Equal(name),
				}),
			)),
		}),
	)
}

func MatchFilterCallNode(name string, args []types.GomegaMatcher, kwargs Fields) types.GomegaMatcher {
	filterArgs := make(Elements)
	for index, arg := range args {
		filterArgs[strconv.Itoa(index)] = arg
	}
	filterKwargs := make(Keys)
	for key, kwarg := range kwargs {
		filterKwargs[key] = kwarg
	}
	return And(
		BeAssignableToTypeOf(nodes.FilterCall{}),
		MatchFields(IgnoreExtras, Fields{
			"Name": Equal(name),
			"Args": MatchAllElementsWithIndex(
				func(index int, element interface{}) string { return strconv.Itoa(index) },
				filterArgs,
			),
			"Kwargs": MatchAllKeys(filterKwargs),
		}),
	)
}

func MatchNodeFilteredExpressionNode(expression types.GomegaMatcher, filters ...types.GomegaMatcher) types.GomegaMatcher {
	filterMatchers := make(Elements)
	for index, filter := range filters {
		filterMatchers[strconv.Itoa(index)] = PointTo(filter)
	}
	return And(
		BeAssignableToTypeOf(nodes.FilteredExpression{}),
		MatchFields(IgnoreExtras, Fields{
			"Expression": PointTo(expression),
			"Filters": MatchAllElementsWithIndex(
				func(index int, _ interface{}) string { return strconv.Itoa(index) },
				filterMatchers,
			),
		}),
	)
}

func MatchGetItemNode(node, arg types.GomegaMatcher) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.GetItem{}),
		MatchFields(IgnoreExtras, Fields{
			"Node": PointTo(node),
			"Arg":  PointTo(arg),
		}),
	)
}

func MatchGetAttributeNode(node types.GomegaMatcher, attribute interface{}) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.GetAttribute{}),
		Or(
			MatchFields(IgnoreExtras, Fields{
				"Node":      PointTo(node),
				"Attribute": Equal(attribute),
			}),
			MatchFields(IgnoreExtras, Fields{
				"Node":  PointTo(node),
				"Index": Equal(attribute),
			}),
		),
	)
}

func MatchCallNode(node types.GomegaMatcher, args []types.GomegaMatcher, kwargs Fields) types.GomegaMatcher {
	filterArgs := make(Elements)
	for index, arg := range args {
		filterArgs[strconv.Itoa(index)] = arg
	}
	filterKwargs := make(Keys)
	for key, kwarg := range kwargs {
		filterKwargs[key] = kwarg
	}
	return And(
		BeAssignableToTypeOf(nodes.Call{}),
		MatchFields(IgnoreExtras, Fields{
			"Func": PointTo(node),
			"Args": MatchAllElementsWithIndex(
				func(index int, element interface{}) string { return strconv.Itoa(index) },
				filterArgs,
			),
			"Kwargs": MatchAllKeys(filterKwargs),
		}),
	)
}

func MatchTestExpressionNode(expression, test types.GomegaMatcher) types.GomegaMatcher {
	return And(
		BeAssignableToTypeOf(nodes.TestExpression{}),
		MatchFields(IgnoreExtras, Fields{
			"Expression": PointTo(expression),
			"Test":       PointTo(test),
		}),
	)
}

func MatchTestCall(name string, args []types.GomegaMatcher, kwargs Fields) types.GomegaMatcher {
	filterArgs := make(Elements)
	for index, arg := range args {
		filterArgs[strconv.Itoa(index)] = arg
	}
	filterKwargs := make(Keys)
	for key, kwarg := range kwargs {
		filterKwargs[key] = kwarg
	}
	return And(
		BeAssignableToTypeOf(nodes.TestCall{}),
		MatchFields(IgnoreExtras, Fields{
			"Name": Equal(name),
			"Args": MatchAllElementsWithIndex(
				func(index int, element interface{}) string { return strconv.Itoa(index) },
				filterArgs,
			),
			"Kwargs": MatchAllKeys(filterKwargs),
		}),
	)
}
