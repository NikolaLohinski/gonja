package methods

import (
	. "github.com/nikolalohinski/gonja/v2/exec"
)

var boolMethods = NewMethodSet[bool](map[string]Method[bool]{
	"bit_length": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		if self {
			return 1, nil
		}
		return 0, nil
	},
	"bit_count": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		if self {
			return 1, nil
		}
		return 0, nil
	},
})
