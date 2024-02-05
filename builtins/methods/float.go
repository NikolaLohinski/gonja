package methods

import (
	. "github.com/nikolalohinski/gonja/v2/exec"
)

var floatMethods = MethodSet[float64]{
	"is_integer": func(self float64, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return self == float64(int(self)), nil
	},
}
