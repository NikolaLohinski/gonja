package methods

import (
	"github.com/nikolalohinski/gonja/v2/builtins/methods/pyint"
	. "github.com/nikolalohinski/gonja/v2/exec"
)

var intMethods = NewMethodSet[int](map[string]Method[int]{
	"is_integer": func(self int, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return true, nil
	},
	"bit_length": func(self int, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pyint.New(self).BitLength(), nil
	},
	"bit_count": func(self int, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pyint.New(self).BitCount(), nil
	},
	"as_integer_ratio": func(self int, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		a, b := pyint.New(self).AsIntegerRatio()
		return []int{a, b}, nil // lack of tuple type we reuse the list type
	},
})
