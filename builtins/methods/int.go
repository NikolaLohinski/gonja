package methods

import (
	. "github.com/nikolalohinski/gonja/v2/exec"
)

var intMethods = NewMethodSet[int](map[string]Method[int]{
	"is_integer": func(self int, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return true, nil
	},
})
