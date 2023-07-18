package parser

import (
	// "fmt"

	// "github.com/juju/errors"
	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/tokens"
)

// FilterFunction is the type filter functions must fulfil
// type FilterFunction func(in *Value, params *VarArgs) (out *Value, err *Error)

// MustApplyFilter behaves like ApplyFilter, but panics on an error.
// func MustApplyFilter(name string, value *Value, params *VarArgs) *Value {
// 	val, err := ApplyFilter(name, value, params)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return val
// }

// // ApplyFilter applies a filter to a given value using the given parameters.
// // Returns a *gonja.Value or an error.
// func ApplyFilter(name string, value *Value, params *VarArgs) (*Value, *Error) {
// 	fn, existing := filters[name]
// 	if !existing {
// 		return nil, &Error{
// 			Sender:    "applyfilter",
// 			OrigError: errors.Errorf("Filter with name '%s' not found.", name),
// 		}
// 	}

// 	// // Make sure param is a *Value
// 	// if param == nil {
// 	// 	param = AsValue(nil)
// 	// }

// 	return fn(value, params)
// }

// type filterCall struct {
// 	token *Token

// 	name   string
// 	args   []Expression
// 	kwargs map[string]Expression

// 	filterFunc FilterFunction
// }

// func (fc *filterCall) Execute(v *Value, ctx *ExecutionContext) (*Value, *Error) {
// 	params := &VarArgs{
// 		Args:   []*Value{},
// 		KwArgs: map[string]*Value{},
// 	}
// 	var err *Error

// 	for _, param := range fc.args {
// 		value, err := param.Evaluate(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		params.Args = append(params.Args, value)
// 	}

// 	for key, param := range fc.kwargs {
// 		value, err := param.Evaluate(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		params.KwArgs[key] = value
// 	}

// 	filteredValue, err := fc.filterFunc(v, params)
// 	if err != nil {
// 		return nil, err.updateFromTokenIfNeeded(ctx.template, fc.token)
// 	}
// 	return filteredValue, nil
// }

// Filter = IDENT | IDENT ":" FilterArg | IDENT "|" Filter
func (p *Parser) ParseFilter() (*nodes.FilterCall, error) {
	identToken := p.Match(tokens.Name)

	// Check filter ident
	if identToken == nil {
		return nil, p.Error("Filter name must be an identifier.", p.Current())
	}

	filter := &nodes.FilterCall{
		Token:  identToken,
		Name:   identToken.Val,
		Args:   []nodes.Expression{},
		Kwargs: map[string]nodes.Expression{},
	}

	// // Get the appropriate filter function and bind it
	// filterFn, exists := filters[identToken.Val]
	// if !exists {
	// 	return nil, p.Error(fmt.Sprintf("Filter '%s' does not exist.", identToken.Val), identToken)
	// }

	// filter.filterFunc = filterFn

	// Check for filter-argument (2 tokens needed: ':' ARG)
	if p.Match(tokens.LeftParenthesis) != nil {
		if p.Current(tokens.VariableEnd) != nil {
			return nil, p.Error("Filter parameter required after '('.", nil)
		}

		for p.Match(tokens.Comma) != nil || p.Match(tokens.RightParenthesis) == nil {
			// TODO: Handle multiple args and kwargs
			v, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}

			if p.Match(tokens.Assign) != nil {
				key := v.Position().Val
				value, errValue := p.ParseExpression()
				if errValue != nil {
					return nil, errValue
				}
				filter.Kwargs[key] = value
			} else {
				filter.Args = append(filter.Args, v)
			}
		}
	}

	return filter, nil
}
