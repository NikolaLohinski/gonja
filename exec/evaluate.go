package exec

import (
	"math"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/nodes"
	"github.com/nikolalohinski/gonja/tokens"
)

var (
	typeOfValuePtr   = reflect.TypeOf(new(Value))
	typeOfExecCtxPtr = reflect.TypeOf(new(Context))
)

type Evaluator struct {
	*EvalConfig
	Ctx *Context
}

func (r *Renderer) Evaluator() *Evaluator {
	return &Evaluator{
		EvalConfig: r.EvalConfig,
		Ctx:        r.Ctx,
	}
}

func (r *Renderer) Eval(node nodes.Expression) *Value {
	e := r.Evaluator()
	return e.Eval(node)
}

func (e *Evaluator) Eval(node nodes.Expression) *Value {
	switch n := node.(type) {
	case *nodes.None:
		return AsValue(nil)
	case *nodes.String:
		return AsValue(n.Val)
	case *nodes.Integer:
		return AsValue(n.Val)
	case *nodes.Float:
		return AsValue(n.Val)
	case *nodes.Bool:
		return AsValue(n.Val)
	case *nodes.List:
		return e.evalList(n)
	case *nodes.Tuple:
		return e.evalTuple(n)
	case *nodes.Dict:
		return e.evalDict(n)
	case *nodes.Pair:
		return e.evalPair(n)
	case *nodes.Name:
		return e.evalName(n)
	case *nodes.Call:
		return e.evalCall(n)
	case *nodes.Getitem:
		return e.evalGetitem(n)
	case *nodes.Getattr:
		return e.evalGetattr(n)
	case *nodes.Error:
		return AsValue(n.Error)
	case *nodes.Negation:
		result := e.Eval(n.Term)
		if result.IsError() {
			return result
		}
		return result.Negate()
	case *nodes.BinaryExpression:
		return e.evalBinaryExpression(n)
	case *nodes.UnaryExpression:
		return e.evalUnaryExpression(n)
	case *nodes.FilteredExpression:
		return e.EvaluateFiltered(n)
	case *nodes.TestExpression:
		return e.EvalTest(n)
	default:
		return AsValue(errors.Errorf(`Unknown expression type "%T"`, n))
	}
}

func (e *Evaluator) evalBinaryExpression(node *nodes.BinaryExpression) *Value {
	var (
		left  *Value
		right *Value
	)
	left = e.Eval(node.Left)
	if left.IsError() {
		return AsValue(errors.Wrapf(left, `Unable to evaluate left parameter %s`, node.Left))
	}

	switch node.Operator.Token.Val {
	// These operators allow lazy right expression evluation
	case "and", "or":
	default:
		right = e.Eval(node.Right)
		if right.IsError() {
			return AsValue(errors.Wrapf(right, `Unable to evaluate right parameter %s`, node.Right))
		}
	}

	switch node.Operator.Token.Type {
	case tokens.Add:
		if left.IsList() {
			if !right.IsList() {
				return AsValue(errors.Wrapf(right, `Unable to concatenate list to %s`, node.Right))
			}

			v := &Value{Val: reflect.ValueOf([]interface{}{})}

			for ix := 0; ix < left.getResolvedValue().Len(); ix++ {
				v.Val = reflect.Append(v.Val, left.getResolvedValue().Index(ix))
			}

			for ix := 0; ix < right.getResolvedValue().Len(); ix++ {
				v.Val = reflect.Append(v.Val, right.getResolvedValue().Index(ix))
			}

			return v
		}
		if left.IsFloat() || right.IsFloat() {
			// Result will be a float
			return AsValue(left.Float() + right.Float())
		}
		// Result will be an integer
		return AsValue(left.Integer() + right.Integer())
	case tokens.Sub:
		if left.IsFloat() || right.IsFloat() {
			// Result will be a float
			return AsValue(left.Float() - right.Float())
		}
		// Result will be an integer
		return AsValue(left.Integer() - right.Integer())
	case tokens.Mul:
		if left.IsFloat() || right.IsFloat() {
			// Result will be float
			return AsValue(left.Float() * right.Float())
		}
		if left.IsString() {
			return AsValue(strings.Repeat(left.String(), right.Integer()))
		}
		// Result will be int
		return AsValue(left.Integer() * right.Integer())
	case tokens.Div:
		// Float division
		return AsValue(left.Float() / right.Float())
	case tokens.Floordiv:
		// Int division
		return AsValue(int(left.Float() / right.Float()))
	case tokens.Mod:
		// Result will be int
		return AsValue(left.Integer() % right.Integer())
	case tokens.Pow:
		return AsValue(math.Pow(left.Float(), right.Float()))
	case tokens.Tilde:
		return AsValue(strings.Join([]string{left.String(), right.String()}, ""))
	case tokens.And:
		if !left.IsTrue() {
			return AsValue(false)
		}
		right = e.Eval(node.Right)
		if right.IsError() {
			return AsValue(errors.Wrapf(right, `Unable to evaluate right parameter %s`, node.Right))
		}
		return AsValue(right.IsTrue())
	case tokens.Or:
		if left.IsTrue() {
			return AsValue(true)
		}
		right = e.Eval(node.Right)
		if right.IsError() {
			return AsValue(errors.Wrapf(right, `Unable to evaluate right parameter %s`, node.Right))
		}
		return AsValue(right.IsTrue())
	case tokens.Lteq:
		if left.IsFloat() || right.IsFloat() {
			return AsValue(left.Float() <= right.Float())
		}
		return AsValue(left.Integer() <= right.Integer())
	case tokens.Gteq:
		if left.IsFloat() || right.IsFloat() {
			return AsValue(left.Float() >= right.Float())
		}
		return AsValue(left.Integer() >= right.Integer())
	case tokens.Eq:
		return AsValue(left.EqualValueTo(right))
	case tokens.Gt:
		if left.IsFloat() || right.IsFloat() {
			return AsValue(left.Float() > right.Float())
		}
		return AsValue(left.Integer() > right.Integer())
	case tokens.Lt:
		if left.IsFloat() || right.IsFloat() {
			return AsValue(left.Float() < right.Float())
		}
		return AsValue(left.Integer() < right.Integer())
	case tokens.Ne:
		return AsValue(!left.EqualValueTo(right))
	case tokens.In:
		return AsValue(right.Contains(left))
	default:
		return AsValue(errors.Errorf(`Unknown operator "%s"`, node.Operator.Token))
	}
}

func (e *Evaluator) evalUnaryExpression(expr *nodes.UnaryExpression) *Value {
	result := e.Eval(expr.Term)
	if result.IsError() {
		return AsValue(errors.Wrapf(result, `Unable to evaluate term %s`, expr.Term))
	}
	if expr.Negative {
		if result.IsNumber() {
			switch {
			case result.IsFloat():
				return AsValue(-1 * result.Float())
			case result.IsInteger():
				return AsValue(-1 * result.Integer())
			default:
				return AsValue(errors.New("Operation between a number and a non-(float/integer) is not possible"))
			}
		} else {
			return AsValue(errors.Errorf("Negative sign on a non-number expression %s", expr.Position()))
		}
	}
	return result
}

func (e *Evaluator) evalList(node *nodes.List) *Value {
	values := ValuesList{}
	for _, val := range node.Val {
		value := e.Eval(val)
		values = append(values, value)
	}
	return AsValue(values)
}

func (e *Evaluator) evalTuple(node *nodes.Tuple) *Value {
	values := ValuesList{}
	for _, val := range node.Val {
		value := e.Eval(val)
		values = append(values, value)
	}
	return AsValue(values)
}

func (e *Evaluator) evalDict(node *nodes.Dict) *Value {
	pairs := []*Pair{}
	for _, pair := range node.Pairs {
		p := e.evalPair(pair)
		if p.IsError() {
			return AsValue(errors.Wrapf(p, `Unable to evaluate pair "%s"`, pair))
		}
		pairs = append(pairs, p.Interface().(*Pair))
	}
	return AsValue(&Dict{pairs})
}

func (e *Evaluator) evalPair(node *nodes.Pair) *Value {
	key := e.Eval(node.Key)
	if key.IsError() {
		return AsValue(errors.Wrapf(key, `Unable to evaluate key "%s"`, node.Key))
	}
	value := e.Eval(node.Value)
	if value.IsError() {
		return AsValue(errors.Wrapf(value, `Unable to evaluate value "%s"`, node.Value))
	}
	return AsValue(&Pair{key, value})
}

func (e *Evaluator) evalName(node *nodes.Name) *Value {
	val, ok := e.Ctx.Get(node.Name.Val)
	if !ok && e.Config.StrictUndefined {
		return AsValue(errors.Errorf(`Unable to evaluate name "%s"`, node.Name.Val))
	}
	return ToValue(val)
}

func (e *Evaluator) evalGetitem(node *nodes.Getitem) *Value {
	value := e.Eval(node.Node)
	if value.IsError() {
		return AsValue(errors.Wrapf(value, `Unable to evaluate target %s`, node.Node))
	}
	if node.Arg == nil {
		return AsValue(errors.Wrapf(value, `Argument not provided %s`, node.Node))
	}

	argument := e.Eval(node.Arg)
	var key interface{}
	switch {
	case argument != nil && argument.IsString():
		key = argument.String()
	case argument != nil && argument.IsInteger():
		key = argument.Integer()
	default:
		return AsValue(errors.Wrapf(value, `Argument %s does not evaluate to string or integer in: %s`, node.Arg, node.Node))
	}

	item, found := value.Getitem(key)
	if !found && argument.IsString() {
		item, found = value.Getattr(argument.String())
	}
	if !found {
		if item.IsError() || argument.IsInteger() /* always fail when accessing array indexes */ {
			return AsValue(errors.Wrapf(item, `Unable to evaluate %s`, node))
		}
		if e.Config.StrictUndefined {
			return AsValue(errors.Errorf(`Unable to evaluate %s: item '%s' not found`, node, node.Arg))
		}
		return AsValue(nil)
	}
	return item
}

func (e *Evaluator) evalGetattr(node *nodes.Getattr) *Value {
	value := e.Eval(node.Node)
	if value.IsError() {
		return AsValue(errors.Wrapf(value, `Unable to evaluate target %s`, node.Node))
	}

	if node.Attr != "" {
		attr, found := value.Getattr(node.Attr)
		if !found {
			attr, found = value.Getitem(node.Attr)
		}
		if !found {
			if attr.IsError() {
				return AsValue(errors.Wrapf(attr, `Unable to evaluate %s`, node))
			}
			if e.Config.StrictUndefined {
				return AsValue(errors.Errorf(`Unable to evaluate %s: attribute '%s' not found`, node, node.Attr))
			}
			return AsValue(nil)
		}
		return attr
	} else {
		item, found := value.Getitem(node.Index)
		if !found {
			if item.IsError() {
				return AsValue(errors.Wrapf(item, `Unable to evaluate %s`, node))
			}
			if e.Config.StrictUndefined {
				return AsValue(errors.Errorf(`Unable to evaluate %s: item %d not found`, node, node.Index))
			}
			return AsValue(nil)
		}
		return item
	}
}

func (e *Evaluator) evalCall(node *nodes.Call) *Value {
	fn := e.Eval(node.Func)
	if fn.IsError() {
		return AsValue(errors.Wrapf(fn, `Unable to evaluate function "%s"`, node.Func))
	}

	if !fn.IsCallable() {
		return AsValue(errors.Errorf(`%s is not callable`, node.Func))
	}

	// current := reflect.ValueOf(fn) // Get the initial value

	var current reflect.Value
	var isSafe bool

	var params []reflect.Value
	var err error
	t := fn.Val.Type()

	if t.NumIn() == 1 && t.In(0) == reflect.TypeOf(&VarArgs{}) {
		params, err = e.evalVarArgs(node)
	} else {
		params, err = e.evalParams(node, fn)
	}
	if err != nil {
		return AsValue(errors.Wrapf(err, `Unable to evaluate parameters`))
	}

	// Call it and get first return parameter back
	values := fn.Val.Call(params)
	rv := values[0]
	if t.NumOut() == 2 {
		e := values[1].Interface()
		if e != nil {
			err, ok := e.(error)
			if !ok {
				return AsValue(errors.Errorf("The second return value is not an error"))
			}
			if err != nil {
				return AsValue(err)
			}
		}
	}

	if rv.Type() != typeOfValuePtr {
		current = reflect.ValueOf(rv.Interface())
	} else {
		// Return the function call value
		current = rv.Interface().(*Value).Val
		isSafe = rv.Interface().(*Value).Safe
	}

	if !current.IsValid() {
		// Value is not valid (e. g. NIL value)
		return AsValue(nil)
	}

	return &Value{Val: current, Safe: isSafe}
}

func (e *Evaluator) evalVariable(node *nodes.Variable) (*Value, error) {
	var current reflect.Value
	var isSafe bool

	for idx, part := range node.Parts {
		if idx == 0 {
			val, ok := e.Ctx.Get(node.Parts[0].S)
			if !ok && e.Config.StrictUndefined {
				return nil, errors.Errorf(`Unable to evaluate name "%s"`, node.Parts[0].S)
			}
			current = reflect.ValueOf(val) // Get the initial value
		} else {
			// Next parts, resolve it from current

			// Before resolving the pointer, let's see if we have a method to call
			// Problem with resolving the pointer is we're changing the receiver
			isFunc := false
			if part.Type == nodes.VarTypeIdent {
				funcValue := current.MethodByName(part.S)
				if funcValue.IsValid() {
					current = funcValue
					isFunc = true
				}
			}

			if !isFunc {
				// If current a pointer, resolve it
				if current.Kind() == reflect.Ptr {
					current = current.Elem()
					if !current.IsValid() {
						// Value is not valid (anymore)
						return AsValue(nil), nil
					}
				}

				// Look up which part must be called now
				switch part.Type {
				case nodes.VarTypeInt:
					// Calling an index is only possible for:
					// * slices/arrays/strings
					switch current.Kind() {
					case reflect.String, reflect.Array, reflect.Slice:
						if part.I >= 0 && current.Len() > part.I {
							current = current.Index(part.I)
						} else {
							// In Django, exceeding the length of a list is just empty.
							return AsValue(nil), nil
						}
					default:
						return nil, errors.Errorf("Can't access an index on type %s (variable %s)",
							current.Kind().String(), node.String())
					}
				case nodes.VarTypeIdent:
					// debugging:
					// fmt.Printf("now = %s (kind: %s)\n", part.s, current.Kind().String())

					// Calling a field or key
					switch current.Kind() {
					case reflect.Struct:
						current = current.FieldByName(part.S)
					case reflect.Map:
						current = current.MapIndex(reflect.ValueOf(part.S))
					default:
						return nil, errors.Errorf("Can't access a field by name on type %s (variable %s)",
							current.Kind().String(), node.String())
					}
				default:
					panic("unimplemented")
				}
			}
		}

		if !current.IsValid() {
			// Value is not valid (anymore)
			return AsValue(nil), nil
		}

		// If current is a reflect.ValueOf(gonja.Value), then unpack it
		// Happens in function calls (as a return value) or by injecting
		// into the execution context (e.g. in a for-loop)
		if current.Type() == typeOfValuePtr {
			tmpValue := current.Interface().(*Value)
			current = tmpValue.Val
			isSafe = tmpValue.Safe
		}

		// Check whether this is an interface and resolve it where required
		if current.Kind() == reflect.Interface {
			current = reflect.ValueOf(current.Interface())
		}

		// Check if the part is a function call
		if part.IsFunctionCall {

			var params []reflect.Value
			var err error
			t := current.Type()

			if t.NumIn() == 1 && t.In(0) == reflect.TypeOf(&VarArgs{}) {
				// params, err = e.evalVarArgs(node, t, part)
			} else {
				// params, err = e.evalParams(node, t, part)
			}
			if err != nil {
				return nil, err
			}

			// Call it and get first return parameter back
			values := current.Call(params)
			rv := values[0]
			if t.NumOut() == 2 {
				e := values[1].Interface()
				if e != nil {
					err, ok := e.(error)
					if !ok {
						return nil, errors.Errorf("The second return value is not an error")
					}
					if err != nil {
						return nil, err
					}
				}
			}

			if rv.Type() != typeOfValuePtr {
				current = reflect.ValueOf(rv.Interface())
			} else {
				// Return the function call value
				current = rv.Interface().(*Value).Val
				isSafe = rv.Interface().(*Value).Safe
			}
		}

		if !current.IsValid() {
			// Value is not valid (e. g. NIL value)
			return AsValue(nil), nil
		}
	}

	return &Value{Val: current, Safe: isSafe}, nil
}

func (e *Evaluator) evalVarArgs(node *nodes.Call) ([]reflect.Value, error) {
	params := &VarArgs{
		Args:   []*Value{},
		KwArgs: map[string]*Value{},
	}
	for _, param := range node.Args {
		value := e.Eval(param)
		if value.IsError() {
			return nil, value
		}
		params.Args = append(params.Args, value)
	}

	for key, param := range node.Kwargs {
		value := e.Eval(param)
		if value.IsError() {
			return nil, value
		}
		params.KwArgs[key] = value
	}
	// va := AsValue(VarArgs{})
	return []reflect.Value{reflect.ValueOf(params)}, nil
}

func (e *Evaluator) evalParams(node *nodes.Call, fn *Value) ([]reflect.Value, error) {
	args := node.Args
	t := fn.Val.Type()

	if len(args) != t.NumIn() && !(len(args) >= t.NumIn()-1 && t.IsVariadic()) {
		msg := "Function input argument count (%d) of '%s' must be equal to the calling argument count (%d)."
		return nil, errors.Errorf(msg, t.NumIn(), node.String(), len(args))
	}

	// Output arguments
	if t.NumOut() != 1 && t.NumOut() != 2 {
		msg := "'%s' must have exactly 1 or 2 output arguments, the second argument must be of type error"
		return nil, errors.Errorf(msg, node.String())
	}

	// Evaluate all parameters
	var parameters []reflect.Value

	numArgs := t.NumIn()
	isVariadic := t.IsVariadic()
	var fnArg reflect.Type

	for idx, arg := range args {
		pv := e.Eval(arg)
		if pv.IsError() {
			return nil, pv
		}

		if isVariadic {
			if idx >= numArgs-1 {
				fnArg = t.In(numArgs - 1).Elem()
			} else {
				fnArg = t.In(idx)
			}
		} else {
			fnArg = t.In(idx)
		}

		if fnArg != typeOfValuePtr {
			// Function's argument is not a *gonja.Value, then we have to check whether input argument is of the same type as the function's argument
			if !isVariadic {
				if fnArg != reflect.TypeOf(pv.Interface()) && fnArg.Kind() != reflect.Interface {
					msg := "Function input argument %d of '%s' must be of type %s or *gonja.Value (not %T)."
					return nil, errors.Errorf(msg, idx, node.String(), fnArg.String(), pv.Interface())
				}
				// Function's argument has another type, using the interface-value
				parameters = append(parameters, reflect.ValueOf(pv.Interface()))
			} else {
				if fnArg != reflect.TypeOf(pv.Interface()) && fnArg.Kind() != reflect.Interface {
					msg := "Function variadic input argument of '%s' must be of type %s or *gonja.Value (not %T)."
					return nil, errors.Errorf(msg, node.String(), fnArg.String(), pv.Interface())
				}
				// Function's argument has another type, using the interface-value
				parameters = append(parameters, reflect.ValueOf(pv.Interface()))
			}
		} else {
			// Function's argument is a *gonja.Value
			parameters = append(parameters, reflect.ValueOf(pv))
		}
	}

	// Check if any of the values are invalid
	for _, p := range parameters {
		if p.Kind() == reflect.Invalid {
			return nil, errors.Errorf("Calling a function using an invalid parameter")
		}
	}

	return parameters, nil
}
