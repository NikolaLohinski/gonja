package methods

import (
	"github.com/nikolalohinski/gonja/v2/builtins/methods/pybool"
	. "github.com/nikolalohinski/gonja/v2/exec"
)

var boolMethods = NewMethodSet[bool](map[string]Method[bool]{
	"string": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pybool.New(self).String(), nil
	},
	"int": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pybool.New(self).Int(), nil
	},
	"bit_length": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pybool.New(self).BitLength(), nil
	},
	"bit_count": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pybool.New(self).BitCount(), nil
	},
	"as_integer_ratio": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		a, b := pybool.New(self).AsIntegerRatio()
		return []int{a, b}, nil // lack of tuple type we reuse the list type
	},
})
