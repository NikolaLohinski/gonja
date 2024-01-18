package methods

import (
	"fmt"

	. "github.com/nikolalohinski/gonja/v2/exec"
)

var intMethods = MethodSet[int]{
	"is_integer": func(self int, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, fmt.Errorf("wrong signature for '%d.is_integer': %s", self, err)
		}
		return true, nil
	},
}
