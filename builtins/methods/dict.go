package methods

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/nikolalohinski/gonja/v2/exec"
)

var dictMethods = exec.NewMethodSet[map[string]any](map[string]exec.Method[map[string]any]{
	"keys": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		keys := make([]string, 0, len(self))
		for key := range self {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys, nil
	},

	"values": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		keys := make([]string, 0, len(self))
		for key := range self {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		values := make([]any, 0, len(keys))
		for _, key := range keys {
			values = append(values, self[key])
		}
		return values, nil
	},

	"items": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		keys := make([]string, 0, len(self))
		for key := range self {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		items := make([]any, 0, len(keys))
		for _, key := range keys {
			items = append(items, []any{key, self[key]})
		}
		return items, nil
	},

	"get": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if len(arguments.Args) < 1 || len(arguments.Args) > 2 {
			return nil, exec.ErrInvalidCall(errors.New("get() takes 1 or 2 positional arguments"))
		}
		key := arguments.Args[0].String()
		if val, ok := self[key]; ok {
			return val, nil
		}
		if len(arguments.Args) == 2 {
			return arguments.Args[1].Interface(), nil
		}
		return nil, nil
	},

	"pop": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if len(arguments.Args) < 1 || len(arguments.Args) > 2 {
			return nil, exec.ErrInvalidCall(errors.New("pop() takes 1 or 2 positional arguments"))
		}
		key := arguments.Args[0].String()
		if val, ok := self[key]; ok {
			delete(self, key)
			return val, nil
		}
		if len(arguments.Args) == 2 {
			return arguments.Args[1].Interface(), nil
		}
		return nil, nil
	},

	"setdefault": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if len(arguments.Args) < 1 || len(arguments.Args) > 2 {
			return nil, exec.ErrInvalidCall(errors.New("setdefault() takes 1 or 2 positional arguments"))
		}
		key := arguments.Args[0].String()
		if val, ok := self[key]; ok {
			return val, nil
		}
		var def any
		if len(arguments.Args) == 2 {
			def = arguments.Args[1].Interface()
		}
		self[key] = def
		return def, nil
	},

	"update": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if len(arguments.Args) > 1 {
			return nil, exec.ErrInvalidCall(errors.New("update() takes at most 1 positional argument"))
		}
		if len(arguments.Args) == 1 {
			arg := arguments.Args[0]
			arg.Iterate(func(_, _ int, k, v *exec.Value) bool {
				if v != nil {
					self[k.String()] = v.Interface()
				}
				return true
			}, func() {})
		}
		for k, v := range arguments.KwArgs {
			self[k] = v.Interface()
		}
		return nil, nil
	},

	"copy": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		clone := make(map[string]any, len(self))
		for k, v := range self {
			clone[k] = v
		}
		return clone, nil
	},

	"clear": func(self map[string]any, selfValue *exec.Value, arguments *exec.VarArgs) (any, error) {
		if err := arguments.Take(); err != nil {
			return nil, exec.ErrInvalidCall(err)
		}
		for k := range self {
			delete(self, k)
		}
		return nil, nil
	},
})
