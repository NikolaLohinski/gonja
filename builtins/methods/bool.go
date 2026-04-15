package methods

import (
	"github.com/ardanlabs/gonja/builtins/methods/pybool"
	"github.com/ardanlabs/gonja/exec"
)

var boolMethods = exec.NewMethodSet[bool](map[string]exec.Method[bool]{
	"string": func(self bool, _ *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		return pybool.New(self).String(), nil
	},
	"int": func(self bool, _ *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		return pybool.New(self).Int(), nil
	},
	"bit_length": func(self bool, _ *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		return pybool.New(self).BitLength(), nil
	},
	"bit_count": func(self bool, _ *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		return pybool.New(self).BitCount(), nil
	},
	"as_integer_ratio": func(self bool, _ *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		a, b := pybool.New(self).AsIntegerRatio()
		return []int{a, b}, nil // lack of tuple type we reuse the list type
	},
})
