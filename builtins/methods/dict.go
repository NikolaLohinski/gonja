package methods

import (
	"sort"

	. "github.com/nikolalohinski/gonja/v2/exec"
)

var dictMethods = NewMethodSet[map[string]interface{}](map[string]Method[map[string]interface{}]{
	"keys": func(self map[string]interface{}, selfValue *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		keys := make([]string, 0)
		for key := range self {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys, nil
	},
})
