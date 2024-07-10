package methods

import (
	"errors"

	"github.com/nikolalohinski/gonja/v2/builtins/methods/pyerrors"
	"github.com/nikolalohinski/gonja/v2/builtins/methods/pystring"
	. "github.com/nikolalohinski/gonja/v2/exec"
	"golang.org/x/exp/utf8string"
)

var strMethods = NewMethodSet[string](map[string]Method[string]{
	"capitalize": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Capitalize(), nil
	},
	"capwords": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).CapWords(), nil
	},
	"casefold": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Casefold(), nil
	},
	"center": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			width    int
			fillchar string
		)
		if err := arguments.Take(
			PositionalArgument("width", nil, IntArgument(&width)),
			PositionalArgument("fillchar", AsValue(' '), StringArgument(&fillchar)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Center(width, utf8string.NewString(fillchar).At(0)), nil
	},
	"count": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			sub   string
			start int
			end   int
		)
		if err := arguments.Take(
			PositionalArgument("sub", nil, StringArgument(&sub)),
			PositionalArgument("start", AsValue(0), IntArgument(&start)),
			PositionalArgument("end", AsValue(len(self)), IntArgument(&end)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Count(pystring.New(self), &start, &end), nil
	},
	"encode": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			encoding  string
			errorsArg string
		)
		if err := arguments.Take(
			KeywordArgument("encoding", AsValue("utf-8"), StringEnumArgument(&encoding, []string{
				// https://docs.python.org/3/library/codecs.html#standard-encodings
				// Supporting just the most common on that are easy to implement in Go
				"utf-8", "utf_8", "U8", "UTF", "utf8", "cp65001", // UTF-8 and its aliases
				"latin_1", "iso-8859-1", "iso8859-1", "8859", "cp819", "latin", "latin1", "L1", // ISO-8859-1 and its aliases
			})),
			KeywordArgument("errors", AsValue("strict"), StringEnumArgument(&errorsArg, []string{
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

		res, err := pystring.PyString(self).Encode(encoding, errorsArg)
		if err != nil {
			if errors.Is(err, pyerrors.ErrArguments) {
				return nil, ErrInvalidCall(err)
			}
			return nil, err
		}

		return res, nil
	},
	"endswith": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			suffix string
			start  int
			end    int
		)
		if err := arguments.Take(
			PositionalArgument("suffix", nil, StringArgument(&suffix)),
			PositionalArgument("start", AsValue(0), IntArgument(&start)),
			PositionalArgument("end", AsValue(len(self)), IntArgument(&end)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).EndsWith(pystring.New(suffix), &start, &end), nil
	},
	"expandtabs": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			substr  string
			tabsize int
			end     int
		)
		if err := arguments.Take(
			PositionalArgument("substr", nil, StringArgument(&substr)),
			PositionalArgument("tabsize", AsValue(0), IntArgument(&tabsize)),
			PositionalArgument("end", AsValue(len(self)), IntArgument(&end)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).ExpandTabs(pystring.New(substr), &tabsize), nil
	},
	"find": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			sub   string
			start int
			end   int
		)
		if err := arguments.Take(
			PositionalArgument("sub", nil, StringArgument(&sub)),
			PositionalArgument("start", AsValue(0), IntArgument(&start)),
			PositionalArgument("end", AsValue(len(self)), IntArgument(&end)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Find(pystring.New(sub), &start, &end), nil
	},
	"format": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		args := make([]any, 0, len(arguments.Args))
		for _, arg := range arguments.Args {
			args = append(args, arg.Interface())
		}
		kwargs := make(map[string]any)
		for key, value := range arguments.KwArgs {
			kwargs[key] = value.Interface()
		}
		return pystring.PyString(self).Format(args, kwargs)
	},
	"format_map": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		args := make([]any, 0, len(arguments.Args))
		for _, arg := range arguments.Args {
			args = append(args, arg.Interface())
		}
		kwargs := make(map[string]any)
		for key, value := range arguments.KwArgs {
			kwargs[key] = value.Interface()
		}
		return pystring.PyString(self).FormatMap(args, kwargs)
	},
	"isalnum": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsAlnum(), nil
	},
	"isalpha": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsAlpha(), nil
	},
	"isascii": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsASCII(), nil
	},
	"isdecimal": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsDecimal(), nil
	},
	"isdigit": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsDigit(), nil
	},
	"islower": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsLower(), nil
	},
	"isnumeric": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsNumeric(), nil
	},
	"isprintable": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsPrintable(), nil
	},
	"isspace": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsSpace(), nil
	},
	"istitle": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsTitle(), nil
	},
	"isupper": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).IsUpper(), nil
	},
	"join": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			strs []string
		)
		if err := arguments.Take(
			PositionalArgument("iterable", nil, StringListArgument(&strs)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.JoinString(self, strs), nil
	},
	"ljust": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			width    int
			fillchar string
		)
		if err := arguments.Take(
			PositionalArgument("width", nil, IntArgument(&width)),
			PositionalArgument("fillchar", AsValue(' '), StringArgument(&fillchar)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).LJust(width, utf8string.NewString(fillchar).At(0)), nil
	},
	"lower": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Lower(), nil
	},
	"lstrip": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			cutset string
		)
		if err := arguments.Take(
			PositionalArgument("cutset", AsValue(' '), StringArgument(&cutset)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).LStrip(cutset), nil
	},
	"partition": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			sep string
		)
		if err := arguments.Take(
			PositionalArgument("sep", AsValue(' '), StringArgument(&sep)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		p1, p2, p3 := pystring.PyString(self).Partition(sep)
		return []string{p1.String(), p2.String(), p3.String()}, nil
	},
	"removeprefix": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			prefix string
		)
		if err := arguments.Take(
			PositionalArgument("prefix", nil, StringArgument(&prefix)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).RemovePrefix(prefix), nil
	},
	"removesuffix": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			suffix string
		)
		if err := arguments.Take(
			PositionalArgument("suffix", nil, StringArgument(&suffix)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).RemoveSuffix(suffix), nil
	},
	"replace": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			old   string
			new   string
			count int
		)
		if err := arguments.Take(
			PositionalArgument("old", AsValue(' '), StringArgument(&old)),
			PositionalArgument("new", AsValue(' '), StringArgument(&new)),
			PositionalArgument("count", nil, IntArgument(&count)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Replace(old, new, count), nil
	},
	"rfind": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			sub   string
			start int
			end   int
		)
		if err := arguments.Take(
			PositionalArgument("sub", nil, StringArgument(&sub)),
			PositionalArgument("start", AsValue(0), IntArgument(&start)),
			PositionalArgument("end", AsValue(len(self)), IntArgument(&end)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).RFind(sub, &start, &end), nil
	},
	"rjust": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			width    int
			fillchar string
		)
		if err := arguments.Take(
			PositionalArgument("width", nil, IntArgument(&width)),
			PositionalArgument("fillchar", AsValue(' '), StringArgument(&fillchar)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).RJust(width, utf8string.NewString(fillchar).At(0)), nil
	},
	"rpartition": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			sep string
		)
		if err := arguments.Take(
			PositionalArgument("sep", AsValue(' '), StringArgument(&sep)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		p1, p2, p3 := pystring.PyString(self).RPartition(sep)
		return []string{p1.String(), p2.String(), p3.String()}, nil
	},
	"rsplit": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			sep      string
			maxsplit int
		)
		if err := arguments.Take(
			PositionalArgument("sep", AsValue(' '), StringArgument(&sep)),
			PositionalArgument("maxsplit", AsValue(-1), IntArgument(&maxsplit)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).RSplit(sep, maxsplit), nil
	},
	"rstrip": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			cutset string
		)
		if err := arguments.Take(
			PositionalArgument("cutset", AsValue(' '), StringArgument(&cutset)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).RStrip(cutset), nil
	},
	"split": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			sep      string
			maxsplit int
		)
		if err := arguments.Take(
			PositionalArgument("sep", AsValue(' '), StringArgument(&sep)),
			PositionalArgument("maxsplit", AsValue(-1), IntArgument(&maxsplit)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Split(sep, maxsplit), nil
	},
	"splitlines": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			keepends bool
		)
		if err := arguments.Take(
			PositionalArgument("keepends", AsValue(false), BoolArgument(&keepends)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).SplitLines(keepends), nil
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

		// Handle array case.
		if prefixes != nil {
			for _, prefix := range prefixes {
				if pystring.PyString(self).StartsWith(prefix, &start, &end) {
					return true, nil
				}
			}
			return false, nil
		}

		// Handle single prefix case.
		return pystring.PyString(self).StartsWith(prefix, &start, &end), nil
	},
	"strip": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			cutset string
		)
		if err := arguments.Take(
			PositionalArgument("cutset", AsValue(' '), StringArgument(&cutset)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Strip(cutset), nil
	},
	"swapcase": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).SwapCase(), nil
	},
	"title": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Title(), nil
	},
	"upper": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		if err := arguments.Take(); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).Upper(), nil
	},
	"zfill": func(self string, _ *Value, arguments *VarArgs) (interface{}, error) {
		var (
			width int
		)
		if err := arguments.Take(
			PositionalArgument("width", nil, IntArgument(&width)),
		); err != nil {
			return nil, ErrInvalidCall(err)
		}
		return pystring.PyString(self).ZFill(width), nil
	},
})
