package methods

import (
	"fmt"
	"sort"

	. "github.com/nikolalohinski/gonja/v2/exec"
)

var dictMethods = MethodSet[map[string]interface{}]{
	"keys": func(self map[string]interface{}, selfValue *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, fmt.Errorf("wrong signature for '%s.keys': %s", selfValue.String(), err)
		}
		keys := make([]string, 0)
		for key := range self {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys, nil
	},
}
