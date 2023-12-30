package exec

import (
	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/nodes"
)

// TestFunction is the type test functions must fulfil
type TestFunction func(*Context, *Value, *VarArgs) (bool, error)

func (e *Evaluator) EvalTest(expr *nodes.TestExpression) *Value {
	value := e.Eval(expr.Expression)
	// if value.IsError() {
	// 	return AsValue(errors.Wrapf(value, `Unable to evaluate expresion %s`, expr.Expression))
	// }

	return e.ExecuteTest(expr.Test, value)
}

func (e *Evaluator) ExecuteTest(tc *nodes.TestCall, v *Value) *Value {
	params := &VarArgs{
		Args:   []*Value{},
		KwArgs: map[string]*Value{},
	}

	for _, param := range tc.Args {
		value := e.Eval(param)
		if value.IsError() {
			return AsValue(errors.Wrapf(value, `Unable to evaluate parameter %s`, param))
		}
		params.Args = append(params.Args, value)
	}

	for key, param := range tc.Kwargs {
		value := e.Eval(param)
		if value.IsError() {
			return AsValue(errors.Wrapf(value, `Unable to evaluate parameter %s`, param))
		}
		params.KwArgs[key] = value
	}

	return e.ExecuteTestByName(tc.Name, v, params)
}

func (e *Evaluator) ExecuteTestByName(name string, in *Value, params *VarArgs) *Value {
	if !e.Environment.Tests.Exists(name) {
		return AsValue(errors.Errorf(`Test "%s" not found`, name))
	}
	test, ok := e.Environment.Tests[name]
	if !ok {
		return AsValue(errors.Errorf(`Test "%s" not found`, name))
	}

	result, err := test(e.Environment.Context, in, params)
	if err != nil {
		return AsValue(errors.Wrapf(err, `Unable to execute test %s`, name))
	} else {
		return AsValue(result)
	}
}
