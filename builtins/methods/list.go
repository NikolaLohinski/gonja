package methods

import (
	"reflect"

	. "github.com/nikolalohinski/gonja/v2/exec"
)

var listMethods = NewMethodSet[[]interface{}](map[string]Method[[]interface{}]{
	"append": func(_ []interface{}, selfValue *Value, arguments *VarArgs) (interface{}, error) {
		var (
			x interface{}
		)
		if err := arguments.Take(
			PositionalArgument("x", nil, AnyArgument(&x)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}

		*selfValue = *ToValue(reflect.Append(selfValue.Val, reflect.ValueOf(ToValue(x))))

		return nil, nil
	},
	"reverse": func(_ []interface{}, selfValue *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		reversed := reflect.MakeSlice(selfValue.Val.Type(), 0, 0)
		for i := selfValue.Val.Len() - 1; i >= 0; i-- {
			reversed = reflect.Append(reversed, selfValue.Val.Index(i))
		}
		for i := 0; i < selfValue.Val.Len(); i++ {
			selfValue.Val.Index(i).Set(reversed.Index(i))
		}

		return nil, nil
	},
	"copy": func(self []interface{}, selfValue *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return self, nil
	},
})
