package exec

import (
	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/nodes"
)

// TestFunction is the type test functions must fulfil
type TestFunction func(*Context, *Value, *VarArgs) (bool, error)

// TestSet maps test names to their TestFunction handler
type TestSet map[string]TestFunction

// Exists returns true if the given test is already registered
func (ts TestSet) Exists(name string) bool {
	_, existing := ts[name]
	return existing
}

// Register registers a new test. If there's already a test with the same
// name, RegisterTest will panic. You usually want to call this
// function in the test's init() function:
// http://golang.org/doc/effective_go.html#init
func (ts *TestSet) Register(name string, fn TestFunction) error {
	if ts.Exists(name) {
		return errors.Errorf("test with name '%s' is already registered", name)
	}
	(*ts)[name] = fn
	return nil
}

// Replace replaces an already registered test with a new implementation. Use this
// function with caution since it allows you to change existing test behaviour.
func (ts *TestSet) Replace(name string, fn TestFunction) error {
	if !ts.Exists(name) {
		return errors.Errorf("test with name '%s' does not exist (therefore cannot be overridden)", name)
	}
	(*ts)[name] = fn
	return nil
}

func (ts *TestSet) Update(other TestSet) TestSet {
	for name, test := range other {
		(*ts)[name] = test
	}
	return *ts
}

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
	if !e.Tests.Exists(name) {
		return AsValue(errors.Errorf(`Test "%s" not found`, name))
	}
	test, _ := (*e.Tests)[name]

	result, err := test(e.Ctx, in, params)
	if err != nil {
		return AsValue(errors.Wrapf(err, `Unable to execute test %s`, name))
	} else {
		return AsValue(result)
	}
}
