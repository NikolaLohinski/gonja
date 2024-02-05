package methods

import (
	"fmt"
	"strings"

	"golang.org/x/text/encoding/charmap"

	. "github.com/nikolalohinski/gonja/v2/exec"
)

var strMethods = MethodSet[string]{
	"upper": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return strings.ToUpper(self), nil
	},
	"startswith": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			prefix   string
			prefixes []string
			start    int
			end      int
		)
		if err := arguments.Take(
			PositionalArgument("prefix", nil, OrArgument(
				StringArgument(&prefix),
				StringListArgument(&prefixes),
			)),
			PositionalArgument("start", AsValue(0), IntArgument(&start)),
			PositionalArgument("end", AsValue(len(self)), IntArgument(&end)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		if start >= len(self) {
			return false, nil
		}
		if prefixes == nil {
			prefixes = []string{prefix}
		}
		if end >= len(self) {
			end = len(self)
		}
		for start < 0 {
			start += len(self)
		}
		for end < 0 {
			end += len(self)
		}
		if start > end {
			return false, nil
		}
		for _, p := range prefixes {
			if strings.HasPrefix(self[start:end], p) {
				return true, nil
			}
		}
		return false, nil
	},
	"encode": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			encoding string
			errors   string
		)
		if err := arguments.Take(
			KeywordArgument("encoding", AsValue("utf-8"), StringEnumArgument(&encoding, []string{
				// https://docs.python.org/3/library/codecs.html#standard-encodings
				// Supporting just the most common on that are easy to implement in Go
				"utf-8", "utf_8", "U8", "UTF", "utf8", "cp65001", // UTF-8 and its aliases
				"latin_1", "iso-8859-1", "iso8859-1", "8859", "cp819", "latin", "latin1", "L1", // ISO-8859-1 and its aliases
			})),
			KeywordArgument("errors", AsValue("strict"), StringEnumArgument(&errors, []string{
				// See https://docs.python.org/3/library/codecs.html#error-handlers
				"strict",
				"ignore", // only makes sense for ISO-8859-1
				// not implementing the arguments below as it is too much work
				// and it is very likely no one is using them in Jinja templates
				// "replace",
				// "xmlcharrefreplace",
				// "backslashreplace",
			})),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}

		switch encoding {
		case "utf-8", "utf_8", "U8", "UTF", "utf8", "cp65001": // UTF-8 and its aliases
			return []byte(self), nil
		case "latin_1", "iso-8859-1", "iso8859-1", "8859", "cp819", "latin", "latin1", "L1": // ISO-8859-1 and its aliases
			encoder := charmap.ISO8859_1.NewEncoder()
			result, err := encoder.Bytes([]byte(self))
			if err != nil && errors == "strict" {
				return nil, fmt.Errorf("failed to encode %s to ISO-8859-1: %s", self, err)
			}
			return result, nil
		default:
			return nil, ErrInvalidCall(fmt.Errorf("unsupported encoding '%s'", encoding))
		}
	},
}
