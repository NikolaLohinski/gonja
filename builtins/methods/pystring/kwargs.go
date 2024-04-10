package pystring

import (
	"fmt"
)

// Used to resolve nested kw-args
type AttributeGetter interface {
	Get(string) (any, bool)
}

// KwArgs adds AttributeGetter interface to map[string]any
type KwArgs map[string]any

func (k KwArgs) Get(key string) (any, bool) {
	v, ok := k[key]
	return v, ok
}

func getNestedKwArgs(keys []string, kwarg AttributeGetter) (any, error) {
	if len(keys) == 0 {
		// Shouldn't happen if all other logic is correct
		return "", fmt.Errorf("%w: empty key", ErrInternal)
	}

	key := keys[0]
	tail := keys[1:]

	// Recursion stop case
	if len(tail) == 0 {
		if val, ok := kwarg.Get(key); ok {
			return val, nil
		}
		return "", fmt.Errorf("%w: '%s'", ErrKey, key)
	}

	// Fetch key
	maybeVal, ok := kwarg.Get(key)
	if !ok {
		return "", fmt.Errorf("%w: '%s'", ErrKey, key)
	}

	// See if we can recurse down to sub-keys
	if subKwarg, ok := maybeVal.(AttributeGetter); ok {
		return getNestedKwArgs(tail, subKwarg)
	}
	if subMap, ok := maybeVal.(map[string]any); ok {
		return getNestedKwArgs(tail, KwArgs(subMap))
	}
	// TODO: support more type assertions such as map[string]string,
	// TODO: support more type assertions such as map[string]AttributeGetter,
	// if len(tail) == 0  then we can also check basic types?? no... prob not... might need reflection
	// TODO: support more type assertions such as map[string]bool,
	// TODO: support index string

	return "", fmt.Errorf("%w: '%#v' is not a sub-gettable for key %s", ErrValue, maybeVal, key)
}
