package exec

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/nodes"
)

type InvalidFilterCallError error

// FilterFunction is the type filter functions must fulfil
type FilterFunction func(e *Evaluator, in *Value, params *VarArgs) *Value

// EvaluateFiltered evaluate a filtered expression
func (e *Evaluator) EvaluateFiltered(expr *nodes.FilteredExpression) *Value {
	value := e.Eval(expr.Expression)

	for _, filter := range expr.Filters {
		value = e.ExecuteFilter(filter, value)
		if value.IsError() {
			return AsValue(errors.Wrapf(value, "unable to evaluate filter %s", filter))
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
			return AsValue(errors.Wrapf(value, "unable to evaluate parameter %s", param))
		}
		params.Args = append(params.Args, value)
	}

	for key, param := range fc.Kwargs {
		value := e.Eval(param)
		if value.IsError() {
			return AsValue(errors.Wrapf(value, "unable to evaluate parameter %s=%s", key, param))
		}
		params.KwArgs[key] = value
	}
	return e.ExecuteFilterByName(fc.Name, v, params)
}

// ExecuteFilterByName execute a filter given its name
func (e *Evaluator) ExecuteFilterByName(name string, in *Value, params *VarArgs) *Value {
	filter, ok := (e.Environment.Filters)[name]
	if !e.Environment.Filters.Exists(name) || !ok {
		return AsValue(errors.Errorf("filter '%s' not found", name))
	}
	returnedValue := filter(e, in, params)
	if returnedValue.IsError() {
		err, ok := returnedValue.Interface().(InvalidFilterCallError)
		if ok {
			return AsValue(fmt.Errorf("invalid call to filter '%s': %s", name, err.Error()))
		}
	}

	return returnedValue
}
