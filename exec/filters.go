package exec

import (
	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/nodes"
)

// FilterFunction is the type filter functions must fulfil
type FilterFunction func(e *Evaluator, in *Value, params *VarArgs) *Value

// EvaluateFiltered evaluate a filtered expression
func (e *Evaluator) EvaluateFiltered(expr *nodes.FilteredExpression) *Value {
	value := e.Eval(expr.Expression)

	for _, filter := range expr.Filters {
		value = e.ExecuteFilter(filter, value)
		if value.IsError() {
			return AsValue(errors.Wrapf(value, `Unable to evaluate filter %s`, filter))
		}
	}

	return value
}

// ExecuteFilter execute a filter node
func (e *Evaluator) ExecuteFilter(fc *nodes.FilterCall, v *Value) *Value {
	params := NewVarArgs()

	for _, param := range fc.Args {
		value := e.Eval(param)
		if value.IsError() {
			return AsValue(errors.Wrapf(value, `Unable to evaluate parameter %s`, param))
		}
		params.Args = append(params.Args, value)
	}

	for key, param := range fc.Kwargs {
		value := e.Eval(param)
		if value.IsError() {
			return AsValue(errors.Wrapf(value, `Unable to evaluate parameter %s=%s`, key, param))
		}
		params.KwArgs[key] = value
	}
	return e.ExecuteFilterByName(fc.Name, v, params)
}

// ExecuteFilterByName execute a filter given its name
func (e *Evaluator) ExecuteFilterByName(name string, in *Value, params *VarArgs) *Value {
	if !e.Environment.Filters.Exists(name) {
		return AsValue(errors.Errorf(`Filter "%s" not found`, name))
	}
	fn, _ := (e.Environment.Filters)[name]

	return fn(e, in, params)
}
