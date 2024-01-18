package methods

import (
	"fmt"

	. "github.com/nikolalohinski/gonja/v2/exec"
)

var boolMethods = MethodSet[bool]{
	"bit_length": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, fmt.Errorf("wrong signature for '%t.bit_length': %s", self, err)
		}
		if self {
			return 1, nil
		}
		return 0, nil
	},
	"bit_count": func(self bool, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, fmt.Errorf("wrong signature for '%t.bit_count': %s", self, err)
		}
		if self {
			return 1, nil
		}
		return 0, nil
	},
}
