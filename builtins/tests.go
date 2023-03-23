package builtins

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nikolalohinski/gonja/exec"
)

var Tests = exec.TestSet{
	"callable":    testCallable,
	"defined":     testDefined,
	"divisibleby": testDivisibleby,
	"eq":          testEqual,
	"equalto":     testEqual,
	"==":          testEqual,
	"even":        testEven,
	"ge":          testGreaterEqual,
	">=":          testGreaterEqual,
	"gt":          testGreaterThan,
	"greaterthan": testGreaterThan,
	">":           testGreaterThan,
	"in":          testIn,
	"iterable":    testIterable,
	"le":          testLessEqual,
	"<=":          testLessEqual,
	"lower":       testLower,
	"lt":          testLessThan,
	"lessthan":    testLessThan,
	"<":           testLessThan,
	"mapping":     testMapping,
	"ne":          testNotEqual,
	"!=":          testNotEqual,
	"none":        testNone,
	"number":      testNumber,
	"odd":         testOdd,
	"sameas":      testSameas,
	"sequence":    testSequence,
	"string":      testString,
	"undefined":   testUndefined,
	"upper":       testUpper,
	"empty":       testEmpty,
}

func testCallable(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	return in.IsCallable(), nil
}

func testDefined(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	return !(in.IsError() || in.IsNil()), nil
}

func testDivisibleby(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	param := params.First()
	if param.Integer() == 0 {
		return false, nil
	}
	return in.Integer()%param.Integer() == 0, nil
}

func testEqual(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	param := params.First()
	return in.Interface() == param.Interface(), nil
}

func testEven(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	if !in.IsInteger() {
		return false, nil
	}
	return in.Integer()%2 == 0, nil
}

func testGreaterEqual(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	param := params.Args[0]
	if !in.IsNumber() || !param.IsNumber() {
		return false, nil
	}
	return in.Float() >= param.Float(), nil
}

func testGreaterThan(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	param := params.Args[0]
	if !in.IsNumber() || !param.IsNumber() {
		return false, nil
	}
	return in.Float() > param.Float(), nil
}

func testIn(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	seq := params.First()
	return seq.Contains(in), nil
}

func testIterable(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	return in.IsDict() || in.IsList() || in.IsString(), nil
}

func testSequence(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	return in.IsList(), nil
}

func testLessEqual(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	param := params.Args[0]
	if !in.IsNumber() || !param.IsNumber() {
		return false, nil
	}
	return in.Float() <= param.Float(), nil
}

func testLower(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	if !in.IsString() {
		return false, nil
	}
	return strings.ToLower(in.String()) == in.String(), nil
}

func testLessThan(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	param := params.Args[0]
	if !in.IsNumber() || !param.IsNumber() {
		return false, nil
	}
	return in.Float() < param.Float(), nil
}

func testMapping(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	return in.IsDict(), nil
}

func testNotEqual(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	param := params.Args[0]
	return in.Interface() != param.Interface(), nil
}

func testNone(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	return in.IsNil(), nil
}

func testNumber(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	return in.IsNumber(), nil
}

func testOdd(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	if !in.IsInteger() {
		return false, nil
	}
	return in.Integer()%2 == 1, nil
}

func testSameas(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	param := params.Args[0]
	if in.IsNil() && param.IsNil() {
		return true, nil
	} else if param.Val.CanAddr() && in.Val.CanAddr() {
		return param.Val.Addr() == in.Val.Addr(), nil
	}
	return reflect.Indirect(param.Val) == reflect.Indirect(in.Val), nil
}

func testString(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	return in.IsString(), nil
}

func testUndefined(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	defined, err := testDefined(ctx, in, params)
	return !defined, err
}

func testUpper(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	if !in.IsString() {
		return false, nil
	}
	return strings.ToUpper(in.String()) == in.String(), nil
}

func testEmpty(ctx *exec.Context, in *exec.Value, params *exec.VarArgs) (bool, error) {
	if !in.IsList() && !in.IsDict() && !in.IsString() {
		return false, exec.AsValue(fmt.Errorf("test 'empty' can only be called for list, map or string"))
	} else {
		return in.Len() == 0, nil
	}
}
