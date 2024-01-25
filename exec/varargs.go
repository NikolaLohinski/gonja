package exec

import (
	"fmt"
	"sort"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

// VarArgs represents pythonic variadic args/kwargs
type VarArgs struct {
	Args   []*Value
	KwArgs map[string]*Value
}

func NewVarArgs() *VarArgs {
	return &VarArgs{
		Args:   []*Value{},
		KwArgs: map[string]*Value{},
	}
}

// First returns the first argument or nil AsValue
func (va *VarArgs) First() *Value {
	if len(va.Args) > 0 {
		return va.Args[0]
	}
	return AsValue(nil)
}

// GetKeywordArgument gets a keyword arguments with fallback on default value
func (va *VarArgs) GetKeywordArgument(key string, fallback interface{}) *Value {
	value, ok := va.KwArgs[key]
	if ok {
		return value
	}
	return AsValue(fallback)
}

type KwArg struct {
	Name    string
	Default interface{}
}

// Expect validates VarArgs against an expected signature
func (v *VarArgs) Expect(arguments int, keywordArguments []*KwArg) *ReducedVarArgs {
	result := &ReducedVarArgs{VarArgs: v}
	copiedVariableArguments := &VarArgs{
		Args:   v.Args,
		KwArgs: map[string]*Value{},
	}
	reduceIndex := -1
	unexpectedArgs := []string{}
	if len(v.Args) < arguments {
		// Priority on missing arguments
		if arguments > 1 {
			result.error = errors.Errorf(`expected %d arguments, got %d`, arguments, len(v.Args))
		} else {
			result.error = errors.Errorf(`expected an argument, got %d`, len(v.Args))
		}
		return result
	} else if len(v.Args) > arguments {
		copiedVariableArguments.Args = v.Args[:arguments]
		for index, arg := range v.Args[arguments:] {
			if len(keywordArguments) > index {
				copiedVariableArguments.KwArgs[keywordArguments[index].Name] = arg
				reduceIndex = index + 1
			} else {
				unexpectedArgs = append(unexpectedArgs, arg.String())
			}
		}
	}

	unexpectedKwArgs := []string{}
Loop:
	for key, value := range v.KwArgs {
		for index, keywordArgument := range keywordArguments {
			if key == keywordArgument.Name {
				if reduceIndex < 0 || index >= reduceIndex {
					copiedVariableArguments.KwArgs[key] = value
					continue Loop
				} else {
					result.error = errors.Errorf(`keyword '%s' has been submitted twice`, key)
					break Loop
				}
			}
		}
		kv := strings.Join([]string{key, value.String()}, "=")
		unexpectedKwArgs = append(unexpectedKwArgs, kv)
	}
	sort.Strings(unexpectedKwArgs)

	if result.error != nil {
		return result
	}

	switch {
	case len(unexpectedArgs) == 0 && len(unexpectedKwArgs) == 0:
	case len(unexpectedArgs) == 1 && len(unexpectedKwArgs) == 0:
		result.error = errors.Errorf(`unexpected argument '%s'`, unexpectedArgs[0])
	case len(unexpectedArgs) > 1 && len(unexpectedKwArgs) == 0:
		result.error = errors.Errorf(`unexpected arguments '%s'`, strings.Join(unexpectedArgs, ", "))
	case len(unexpectedArgs) == 0 && len(unexpectedKwArgs) == 1:
		result.error = errors.Errorf(`unexpected keyword argument '%s'`, unexpectedKwArgs[0])
	case len(unexpectedArgs) == 0 && len(unexpectedKwArgs) > 0:
		result.error = errors.Errorf(`unexpected keyword arguments '%s'`, strings.Join(unexpectedKwArgs, ", "))
	default:
		result.error = errors.Errorf(`unexpected arguments '%s, %s'`,
			strings.Join(unexpectedArgs, ", "),
			strings.Join(unexpectedKwArgs, ", "),
		)
	}

	if result.error != nil {
		return result
	}
	// fill defaults
	for _, kwarg := range keywordArguments {
		_, exists := copiedVariableArguments.KwArgs[kwarg.Name]
		if !exists {
			copiedVariableArguments.KwArgs[kwarg.Name] = AsValue(kwarg.Default)
		}
	}
	result.VarArgs = copiedVariableArguments
	return result
}

// ExpectArgs ensures VarArgs receive only arguments
func (va *VarArgs) ExpectArgs(args int) *ReducedVarArgs {
	return va.Expect(args, []*KwArg{})
}

// ExpectNothing ensures VarArgs does not receive any argument
func (va *VarArgs) ExpectNothing() *ReducedVarArgs {
	return va.ExpectArgs(0)
}

// ExpectKwArgs allow to specify optionnaly expected KwArgs
func (va *VarArgs) ExpectKwArgs(kwargs []*KwArg) *ReducedVarArgs {
	return va.Expect(0, kwargs)
}

// ReducedVarArgs represents python variadic arguments / keyword arguments
// but values are reduced (ie. keyword arguments given as arguments are accessible by name)
type ReducedVarArgs struct {
	*VarArgs
	error error
}

// IsError returns true if there was an error on Expect call
func (r *ReducedVarArgs) IsError() bool {
	return r.error != nil
}

func (r *ReducedVarArgs) Error() string {
	if r.IsError() {
		return r.error.Error()
	}
	return ""
}

type ArgumentTransmuter func(*Value) error

func StringArgument(output *string) func(*Value) error {
	return func(v *Value) error {
		if !v.IsString() {
			return fmt.Errorf("%s is not a string", v.String())
		}
		if output == nil {
			return errors.New("received nil pointer to string in StringArgument transposer")
		}
		*output = v.String()
		return nil
	}
}

func OrArgument(transmuters ...ArgumentTransmuter) func(*Value) error {
	return func(v *Value) error {
		var errors []string
		for _, transmuter := range transmuters {
			err := transmuter(v)
			if err == nil {
				return nil
			}
			errors = append(errors, err.Error())
		}
		return fmt.Errorf("failed to validate argument '%s' against any alternative: %s", v.String(), strings.Join(errors, " | "))
	}
}

func StringEnumArgument(output *string, options []string) func(v *Value) error {
	return func(v *Value) error {
		if !v.IsString() {
			return fmt.Errorf("%s is not a string", v.String())
		}
		if output == nil {
			return errors.New("received nil pointer to string in StringArgument transposer")
		}
		value := v.String()
		for _, option := range options {
			if option == value {
				*output = v.String()
				return nil
			}
		}
		return fmt.Errorf("unexpected value '%s' is not in: ['%s']", value, strings.Join(options, "','"))
	}
}

func IntArgument(output *int) func(v *Value) error {
	return func(v *Value) error {
		if !v.IsInteger() {
			return fmt.Errorf("%s is not an integer", v.String())
		}
		if output == nil {
			return errors.New("received nil pointer to int in IntArgument transposer")
		}
		*output = v.Integer()
		return nil
	}
}

func AnyArgument(output *interface{}) func(*Value) error {
	return func(v *Value) error {
		if output == nil {
			return errors.New("received nil pointer to string in AnyArgument transposer")
		}
		*output = v.Interface()
		return nil
	}
}

func StringListArgument(output *[]string) func(*Value) error {
	return func(v *Value) error {
		if !v.IsList() {
			return fmt.Errorf("%s is not a list", v.String())
		}
		if output == nil {
			return errors.New("received nil pointer to string in StringListArgument transposer")
		}
		*output = []string{}
		for i := 0; i < v.Len(); i++ {
			item := ToValue(v.Val.Index(i))
			if !item.IsString() {
				return fmt.Errorf("%s item '%s' of the argument list is not a string: %s", humanize.Ordinal(i+1), item.String(), v.String())
			}
			*output = append(*output, item.String())
		}
		return nil
	}
}

type argument struct {
	name        string
	positional  bool
	fallback    *Value
	transmuters []ArgumentTransmuter
}

func PositionalArgument(name string, fallback *Value, transmuters ...ArgumentTransmuter) *argument {
	return &argument{
		name:        name,
		transmuters: transmuters,
		positional:  true,
		fallback:    fallback,
	}
}

func KeywordArgument(name string, defaultValue *Value, transmuters ...ArgumentTransmuter) *argument {
	return &argument{
		name:        name,
		positional:  false,
		transmuters: transmuters,
		fallback:    defaultValue,
	}
}

func (v *VarArgs) Take(arguments ...*argument) error {
	unexpectedArgs := len(v.Args)
	unexpectedKwArgs := v.KwArgs
	for index, argument := range arguments {
		var value *Value
		if argument.positional {
			if index >= len(v.Args) {
				if argument.fallback != nil {
					value = argument.fallback
				} else {
					return fmt.Errorf("missing required %s positional argument '%s'", humanize.Ordinal(index+1), argument.name)
				}
			} else {
				value = v.Args[index]
				unexpectedArgs -= 1
			}
		} else {
			if unexpectedArgs > 0 && index < len(v.Args) {
				value = v.Args[index]
				unexpectedArgs -= 1
			} else if v, ok := unexpectedKwArgs[argument.name]; ok {
				value = v
				delete(unexpectedKwArgs, argument.name)
			} else {
				value = argument.fallback
			}
		}
		for _, transmute := range argument.transmuters {
			if err := transmute(value); err != nil {
				return errors.Errorf("failed to validate argument '%s': %s", argument.name, err.Error())
			}
		}
	}
	if unexpectedArgs != 0 {
		message := fmt.Sprintf("received %d unexpected positional argument", unexpectedArgs)
		if unexpectedArgs > 1 {
			message = message + "s"
		}
		return errors.New(message)
	}
	if len(unexpectedKwArgs) != 0 {
		message := fmt.Sprintf("received %d unexpected keyword argument", len(unexpectedKwArgs))
		if len(unexpectedKwArgs) > 1 {
			message = message + "s"
		}
		unexpectedKwArgNames := []string{}
		for name := range unexpectedKwArgs {
			unexpectedKwArgNames = append(unexpectedKwArgNames, name)
		}
		return fmt.Errorf("%s: '%s'", message, strings.Join(unexpectedKwArgNames, "','"))
	}
	return nil
}
