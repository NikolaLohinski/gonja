package builtins

import (
	stdjson "encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ardanlabs/gonja/exec"
	"github.com/ardanlabs/gonja/utils"
)

type tupleValue []any

func (t tupleValue) String() string {
	if len(t) == 0 {
		return "()"
	}

	var out strings.Builder
	out.WriteByte('(')
	for i, item := range t {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(pythonRepr(item))
	}
	if len(t) == 1 {
		out.WriteByte(',')
	}
	out.WriteByte(')')
	return out.String()
}

type groupTupleValue []any

func (g groupTupleValue) String() string {
	return tupleValue(g).String()
}

func (g groupTupleValue) GetAttribute(name string) (*exec.Value, bool) {
	if len(g) != 2 {
		return exec.AsValue(nil), false
	}
	switch name {
	case "grouper":
		return exec.AsValue(g[0]), true
	case "list":
		return exec.AsValue(g[1]), true
	default:
		return exec.AsValue(nil), false
	}
}

func pythonRepr(value any) string {
	rendered := exec.AsValue(value)
	if rendered.IsNil() {
		return "None"
	}
	if rendered.IsString() {
		return fmt.Sprintf(`'%s'`, rendered.String())
	}
	return rendered.String()
}

func stringifyFilterValue(value *exec.Value) string {
	if value == nil || value.IsNil() {
		return "None"
	}
	return value.String()
}

func escapeFilterValue(value *exec.Value) string {
	if value == nil || value.IsNil() {
		return "None"
	}
	return utils.Escape(value.String())
}

func resolveAttributeValue(value *exec.Value, attribute *exec.Value, defaultValue *exec.Value) (*exec.Value, bool) {
	if attribute == nil || attribute.IsNil() {
		return value, true
	}
	if attribute.IsInteger() {
		return resolveAttributeIndex(value, attribute.Integer(), defaultValue)
	}
	return resolveAttributePath(value, attribute.String(), defaultValue)
}

func resolveAttributePath(value *exec.Value, path string, defaultValue *exec.Value) (*exec.Value, bool) {
	current := value
	if path == "" {
		return current, true
	}

	for part := range strings.SplitSeq(path, ".") {
		if current == nil || current.IsNil() {
			if defaultValue != nil {
				return defaultValue, true
			}
			return exec.AsValue(nil), false
		}

		if index, err := strconv.Atoi(part); err == nil {
			next, found := current.GetItem(index)
			if !found {
				if defaultValue != nil {
					return defaultValue, true
				}
				return exec.AsValue(nil), false
			}
			current = next
			continue
		}

		next, found := current.Get(part)
		if !found {
			if defaultValue != nil {
				return defaultValue, true
			}
			return exec.AsValue(nil), false
		}
		current = next
	}

	return current, true
}

func resolveAttributeIndex(value *exec.Value, index int, defaultValue *exec.Value) (*exec.Value, bool) {
	if value == nil || value.IsNil() {
		if defaultValue != nil {
			return defaultValue, true
		}
		return exec.AsValue(nil), false
	}
	next, found := value.GetItem(index)
	if !found {
		if defaultValue != nil {
			return defaultValue, true
		}
		return exec.AsValue(nil), false
	}
	return next, true
}

func compareValues(left, right *exec.Value, caseSensitive bool) int {
	switch {
	case left == nil || left.IsNil():
		if right == nil || right.IsNil() {
			return 0
		}
		return -1
	case right == nil || right.IsNil():
		return 1
	case left.IsNumber() && right.IsNumber():
		lf := left.Float()
		rf := right.Float()
		switch {
		case lf < rf:
			return -1
		case lf > rf:
			return 1
		default:
			return 0
		}
	default:
		ls := stringifyFilterValue(left)
		rs := stringifyFilterValue(right)
		if !caseSensitive {
			ls = strings.ToLower(ls)
			rs = strings.ToLower(rs)
		}
		return strings.Compare(ls, rs)
	}
}

func takeValueArgument(output **exec.Value) exec.ArgumentTransmuter {
	return func(v *exec.Value) error {
		if output == nil {
			return fmt.Errorf("received nil pointer to value output")
		}
		*output = v
		return nil
	}
}

func takeStringArgument(output *string) exec.ArgumentTransmuter {
	return func(v *exec.Value) error {
		if output == nil {
			return fmt.Errorf("received nil pointer to string output")
		}
		*output = v.String()
		return nil
	}
}

func takeBoolArgument(output *bool) exec.ArgumentTransmuter {
	return func(v *exec.Value) error {
		if output == nil {
			return fmt.Errorf("received nil pointer to bool output")
		}
		*output = v.Bool()
		return nil
	}
}

func takeIntArgument(output *int) exec.ArgumentTransmuter {
	return func(v *exec.Value) error {
		if output == nil {
			return fmt.Errorf("received nil pointer to int output")
		}
		*output = v.Integer()
		return nil
	}
}

func takeFloatArgument(output *float64) exec.ArgumentTransmuter {
	return func(v *exec.Value) error {
		if output == nil {
			return fmt.Errorf("received nil pointer to float output")
		}
		*output = v.Float()
		return nil
	}
}

func normalizeJSONValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case stdjson.Marshaler:
		return typed
	case *exec.Dict:
		object := make(map[string]any, len(typed.Pairs))
		for _, pair := range typed.Pairs {
			object[pair.Key.String()] = normalizeJSONValue(pair.Value.Interface())
		}
		return object
	case exec.Dict:
		object := make(map[string]any, len(typed.Pairs))
		for _, pair := range typed.Pairs {
			object[pair.Key.String()] = normalizeJSONValue(pair.Value.Interface())
		}
		return object
	case *exec.Value:
		return normalizeJSONValue(typed.Interface())
	}

	resolved := reflect.ValueOf(value)
	if !resolved.IsValid() {
		return nil
	}

	switch resolved.Kind() {
	case reflect.Slice, reflect.Array:
		if resolved.Type().Elem().Kind() == reflect.Uint8 {
			return value
		}
		out := make([]any, resolved.Len())
		for i := 0; i < resolved.Len(); i++ {
			out[i] = normalizeJSONValue(resolved.Index(i).Interface())
		}
		return out
	case reflect.Map:
		if resolved.Type().Key().Kind() != reflect.String {
			return value
		}
		out := make(map[string]any, resolved.Len())
		iter := resolved.MapRange()
		for iter.Next() {
			out[iter.Key().String()] = normalizeJSONValue(iter.Value().Interface())
		}
		return out
	default:
		return value
	}
}

func urlEncodePair(item *exec.Value) (*exec.Value, *exec.Value, bool) {
	first, ok := item.GetItem(0)
	if !ok {
		return nil, nil, false
	}
	second, ok := item.GetItem(1)
	if !ok {
		return nil, nil, false
	}
	return first, second, true
}

func urlizeToken(token string, trimURLLimit int, rel string, target string, extraSchemes []string) (string, bool) {
	lower := strings.ToLower(token)
	switch {
	case strings.HasPrefix(lower, "mailto:"):
		email := token[len("mailto:"):]
		return buildURLAnchor("mailto:"+email, trimURLTitle(email, trimURLLimit), "", target), true
	case filterURLizeEmailRegexp.MatchString(token):
		return buildURLAnchor("mailto:"+token, trimURLTitle(token, trimURLLimit), "", target), true
	case strings.HasPrefix(lower, "http://"), strings.HasPrefix(lower, "https://"), hasURLScheme(lower, extraSchemes):
		return buildURLAnchor(token, trimURLTitle(token, trimURLLimit), rel, target), true
	case filterURLizeDomainRegexp.MatchString(token):
		return buildURLAnchor("https://"+token, trimURLTitle(token, trimURLLimit), rel, target), true
	default:
		return "", false
	}
}

func hasURLScheme(token string, schemes []string) bool {
	for _, scheme := range schemes {
		if strings.HasPrefix(token, strings.ToLower(scheme)) {
			return true
		}
	}
	return false
}

func trimURLTitle(title string, trimURLLimit int) string {
	if trimURLLimit > 3 && len(title) > trimURLLimit {
		return title[:trimURLLimit-3] + "..."
	}
	return title
}

func buildURLAnchor(href string, title string, rel string, target string) string {
	var attrs strings.Builder
	if rel != "" {
		attrs.WriteString(fmt.Sprintf(` rel="%s"`, rel))
	}
	if target != "" {
		attrs.WriteString(fmt.Sprintf(` target="%s"`, target))
	}
	return fmt.Sprintf(`<a href="%s"%s>%s</a>`, utils.IRIEncode(href), attrs.String(), utils.Escape(title))
}
