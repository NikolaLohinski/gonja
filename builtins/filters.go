// Package builtins provides built-in filters, tests, control structures, and global functions.
package builtins

import (
	stdjson "encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/nikolalohinski/gonja/v2/exec"
)

// Filters export all builtin filters
var Filters = exec.NewFilterSet(map[string]exec.FilterFunction{
	"abs":            filterAbs,
	"attr":           filterAttr,
	"batch":          filterBatch,
	"capitalize":     filterCapitalize,
	"center":         filterCenter,
	"default":        filterDefault,
	"d":              filterDefault,
	"dictsort":       filterDictSort,
	"e":              filterEscape,
	"escape":         filterEscape,
	"filesizeformat": filterFileSize,
	"first":          filterFirst,
	"float":          filterFloat,
	"forceescape":    filterForceEscape,
	"format":         filterFormat,
	"groupby":        filterGroupBy,
	"indent":         filterIndent,
	"int":            filterInteger,
	"join":           filterJoin,
	"items":          filterItems,
	"last":           filterLast,
	"length":         filterLength,
	"list":           filterList,
	"lower":          filterLower,
	"map":            filterMap,
	"max":            filterMax,
	"min":            filterMin,
	"pprint":         filterPPrint,
	"random":         filterRandom,
	"rejectattr":     filterRejectAttr,
	"reject":         filterReject,
	"replace":        filterReplace,
	"reverse":        filterReverse,
	"round":          filterRound,
	"safe":           filterSafe,
	"selectattr":     filterSelectAttr,
	"select":         filterSelect,
	"slice":          filterSlice,
	"sort":           filterSort,
	"string":         filterString,
	"striptags":      filterStriptags,
	"sum":            filterSum,
	"title":          filterTitle,
	"tojson":         filterToJSON,
	"trim":           filterTrim,
	"truncate":       filterTruncate,
	"unique":         filterUnique,
	"upper":          filterUpper,
	"urlencode":      filterUrlencode,
	"urlize":         filterUrlize,
	"wordcount":      filterWordcount,
	"wordwrap":       filterWordwrap,
	"xmlattr":        filterXMLAttr,
})

func filterAbs(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.IsInteger() {
		asInt := in.Integer()
		if asInt < 0 {
			return exec.AsValue(-asInt)
		}
		return in
	} else if in.IsFloat() {
		return exec.AsValue(math.Abs(in.Float()))
	}
	return exec.AsValue(math.Abs(in.Float())) // nothing to do here, just to keep track of the safe application
}

func filterAttr(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var name string
	if err := params.Take(
		exec.PositionalArgument("name", nil, takeStringArgument(&name)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	value, _ := in.GetAttribute(name)
	return value
}

func filterBatch(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		lineCount int
		fillWith  *exec.Value
	)
	if err := params.Take(
		exec.PositionalArgument("linecount", nil, takeIntArgument(&lineCount)),
		exec.KeywordArgument("fill_with", exec.AsValue(nil), takeValueArgument(&fillWith)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	out := make([]any, 0)
	var row []any
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		if math.Mod(float64(idx), float64(lineCount)) == 0 {
			if row != nil {
				out = append(out, exec.AsValue(row).Interface())
			}
			row = make([]any, 0)
		}
		row = append(row, key.Interface())
		return true
	}, func() {})
	if len(row) > 0 {
		if !fillWith.IsNil() {
			for len(row) < lineCount {
				row = append(row, fillWith.Interface())
			}
		}
		out = append(out, exec.AsValue(row).Interface())
	}
	return exec.AsValue(out)
}

func filterCapitalize(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}

	if !in.IsString() {
		return in
	}

	if in.Len() <= 0 {
		return exec.AsValue("")
	}
	t := in.String()
	r, size := utf8.DecodeRuneInString(t)
	return exec.AsValue(strings.ToUpper(string(r)) + strings.ToLower(t[size:]))
}

func filterCenter(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var width int
	if err := params.Take(
		exec.PositionalArgument("width", nil, takeIntArgument(&width)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	slen := in.Len()
	if width <= slen {
		return in
	}

	spaces := width - slen
	left := spaces/2 + spaces%2
	right := spaces / 2

	return exec.AsValue(fmt.Sprintf("%s%s%s", strings.Repeat(" ", left),
		in.String(), strings.Repeat(" ", right)))
}

func sortByKey(in *exec.Value, caseSensitive bool, reverse bool) []tupleValue {
	items := in.Items()
	sort.SliceStable(items, func(i, j int) bool {
		comparison := compareValues(items[i].Key, items[j].Key, caseSensitive)
		if reverse {
			return comparison > 0
		}
		return comparison < 0
	})

	out := make([]tupleValue, 0, len(items))
	for _, item := range items {
		out = append(out, tupleValue{item.Key.Interface(), item.Value.Interface()})
	}
	return out
}

func sortByValue(in *exec.Value, caseSensitive, reverse bool) []tupleValue {
	items := in.Items()
	sort.SliceStable(items, func(i, j int) bool {
		comparison := compareValues(items[i].Value, items[j].Value, caseSensitive)
		if reverse {
			return comparison > 0
		}
		return comparison < 0
	})

	out := make([]tupleValue, 0, len(items))
	for _, item := range items {
		out = append(out, tupleValue{item.Key.Interface(), item.Value.Interface()})
	}
	return out
}

func filterDictSort(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		caseSensitive bool
		by            string
		reverse       bool
	)
	if err := params.Take(
		exec.KeywordArgument("case_sensitive", exec.AsValue(false), takeBoolArgument(&caseSensitive)),
		exec.KeywordArgument("by", exec.AsValue("key"), takeStringArgument(&by)),
		exec.KeywordArgument("reverse", exec.AsValue(false), takeBoolArgument(&reverse)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}

	switch by {
	case "key":
		return exec.AsValue(sortByKey(in, caseSensitive, reverse))
	case "value":
		return exec.AsValue(sortByValue(in, caseSensitive, reverse))
	default:
		return exec.AsValue(errors.New(`by should be either 'key' or 'value`))
	}
}

func filterEscape(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.Safe {
		return in
	}
	return exec.AsSafeValue(in.Escaped())
}

var (
	bytesPrefixes  = []string{"kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	binaryPrefixes = []string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
)

func filterFileSize(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var binary bool
	if err := params.Take(
		exec.KeywordArgument("binary", exec.AsValue(false), takeBoolArgument(&binary)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	bytes := in.Float()
	var base float64
	var prefixes []string
	if binary {
		base = 1024.0
		prefixes = binaryPrefixes
	} else {
		base = 1000.0
		prefixes = bytesPrefixes
	}
	if bytes == 1.0 {
		return exec.AsValue("1 Byte")
	} else if bytes < base {
		return exec.AsValue(fmt.Sprintf("%.0f Bytes", bytes))
	} else {
		var i int
		var unit float64
		var prefix string
		for i, prefix = range prefixes {
			unit = math.Pow(base, float64(i+2))
			if bytes < unit {
				return exec.AsValue(fmt.Sprintf("%.1f %s", (base * bytes / unit), prefix))
			}
		}
		return exec.AsValue(fmt.Sprintf("%.1f %s", (base * bytes / unit), prefix))
	}
}

func filterFirst(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.CanSlice() && in.Len() > 0 {
		return in.Index(0)
	}
	return exec.AsValue("")
}

func filterFloat(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var defaultValue *exec.Value
	if err := params.Take(
		exec.KeywordArgument("default", exec.AsValue(0.0), takeValueArgument(&defaultValue)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.IsNil() {
		return defaultValue
	}
	if in.IsFloat() || in.IsInteger() {
		return exec.AsValue(in.Float())
	}
	if parsed, err := strconv.ParseFloat(in.String(), 64); err == nil {
		return exec.AsValue(parsed)
	}
	return defaultValue
}

func filterForceEscape(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	return exec.AsSafeValue(in.Escaped())
}

func filterFormat(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	args := []any{}
	for _, arg := range params.Args {
		args = append(args, arg.Interface())
	}
	return exec.AsValue(fmt.Sprintf(in.String(), args...))
}

func filterGroupBy(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		attribute     *exec.Value
		defaultValue  *exec.Value
		caseSensitive bool
	)
	if err := params.Take(
		exec.PositionalArgument("attribute", nil, takeValueArgument(&attribute)),
		exec.KeywordArgument("default", exec.AsValue(nil), takeValueArgument(&defaultValue)),
		exec.KeywordArgument("case_sensitive", exec.AsValue(false), takeBoolArgument(&caseSensitive)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}

	items := make([]*exec.Value, 0)
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		items = append(items, key)
		return true
	}, func() {})

	sort.SliceStable(items, func(i, j int) bool {
		left, _ := resolveAttributeValue(items[i], attribute, defaultValue)
		right, _ := resolveAttributeValue(items[j], attribute, defaultValue)
		return compareValues(left, right, caseSensitive) < 0
	})

	out := make([]groupTupleValue, 0)
	for _, item := range items {
		key, found := resolveAttributeValue(item, attribute, defaultValue)
		if !found {
			continue
		}
		if len(out) == 0 {
			out = append(out, groupTupleValue{key.Interface(), []any{item.Interface()}})
			continue
		}

		last := out[len(out)-1]
		lastKey := exec.AsValue(last[0])
		if compareValues(lastKey, key, caseSensitive) == 0 {
			last[1] = append(last[1].([]any), item.Interface())
			out[len(out)-1] = last
			continue
		}
		out = append(out, groupTupleValue{key.Interface(), []any{item.Interface()}})
	}
	return exec.AsValue(out)
}

func filterIndent(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		indent string
		first  bool
		blank  bool
	)
	if err := params.Take(
		exec.KeywordArgument("width", exec.AsValue(4), func(v *exec.Value) error {
			switch {
			case v.IsInteger():
				width := v.Integer()
				if width < 0 {
					indent = ""
				} else {
					indent = strings.Repeat(" ", width)
				}
			case v.IsString():
				indent = v.String()
			default:
				return fmt.Errorf("%s is neither a string nor an integer", v.String())
			}
			return nil
		}),
		exec.KeywordArgument("first", exec.AsValue(false), exec.BoolArgument(&first)),
		exec.KeywordArgument("blank", exec.AsValue(false), exec.BoolArgument(&blank)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if !in.IsString() {
		return exec.AsValue(exec.ErrInvalidCall(fmt.Errorf("%s is not a string", in.String())))
	}
	lines := strings.Split(in.String(), "\n")
	for idx, line := range lines {
		if idx == 0 {
			if first {
				lines[idx] = indent + line
			}
			continue
		}
		if line == "" && !blank {
			continue
		}
		lines[idx] = indent + line
	}
	out := strings.Join(lines, "\n")
	if in.Safe {
		return exec.AsSafeValue(out)
	}
	return exec.AsValue(out)
}

func filterInteger(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		defaultValue *exec.Value
		base         int
	)
	if err := params.Take(
		exec.KeywordArgument("default", exec.AsValue(0), takeValueArgument(&defaultValue)),
		exec.KeywordArgument("base", exec.AsValue(10), takeIntArgument(&base)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}

	switch {
	case in.IsNil():
		return defaultValue
	case in.IsInteger():
		return exec.AsValue(in.Interface())
	case in.IsFloat():
		return exec.AsValue(int(in.Float()))
	}

	if converted, ok := in.Interface().(interface{ Int() int }); ok {
		return exec.AsValue(converted.Int())
	}

	raw := strings.TrimSpace(in.String())
	if raw == "" {
		return defaultValue
	}

	switch {
	case base == 16 && strings.HasPrefix(strings.ToLower(raw), "0x"):
		raw = raw[2:]
	case base == 8 && strings.HasPrefix(strings.ToLower(raw), "0o"):
		raw = raw[2:]
	case base == 2 && strings.HasPrefix(strings.ToLower(raw), "0b"):
		raw = raw[2:]
	}

	if parsed, ok := new(big.Int).SetString(raw, base); ok {
		if parsed.IsInt64() {
			parsedInt := parsed.Int64()
			maxInt := int64(^uint(0) >> 1)
			minInt := -maxInt - 1
			if parsedInt >= minInt && parsedInt <= maxInt {
				return exec.AsValue(int(parsedInt))
			}
		}
		return exec.AsValue(parsed)
	}

	if parsed, err := strconv.ParseFloat(raw, 64); err == nil && !math.IsNaN(parsed) && !math.IsInf(parsed, 0) {
		return exec.AsValue(int(parsed))
	}

	return defaultValue
}

func filterJoin(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		d         *exec.Value
		attribute *exec.Value
	)
	if err := params.Take(
		exec.KeywordArgument("d", exec.AsValue(""), takeValueArgument(&d)),
		exec.KeywordArgument("attribute", exec.AsValue(nil), takeValueArgument(&attribute)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if !in.IsIterable() {
		return in
	}
	delimiter := stringifyFilterValue(d)
	if e.Config.AutoEscape && !d.Safe {
		delimiter = escapeFilterValue(d)
	}

	parts := make([]string, 0)
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		item := key
		if !attribute.IsNil() {
			resolved, found := resolveAttributeValue(item, attribute, nil)
			if found {
				item = resolved
			} else {
				item = exec.AsValue(nil)
			}
		}
		if e.Config.AutoEscape {
			if item.Safe {
				parts = append(parts, stringifyFilterValue(item))
			} else {
				parts = append(parts, escapeFilterValue(item))
			}
		} else {
			parts = append(parts, stringifyFilterValue(item))
		}
		return true
	}, func() {})

	joined := strings.Join(parts, delimiter)
	if e.Config.AutoEscape {
		return exec.AsSafeValue(joined)
	}
	return exec.AsValue(joined)
}

func filterLast(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.CanSlice() && in.Len() > 0 {
		return in.Index(in.Len() - 1)
	}
	return exec.AsValue("")
}

func filterLength(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	return exec.AsValue(in.Len())
}

func filterItems(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.IsNil() {
		return exec.AsValue([]tupleValue{})
	}
	if in.IsList() {
		return in
	}
	if !in.IsDict() {
		return exec.AsValue(errors.New("items requires a mapping"))
	}
	items := in.Items()
	out := make([]tupleValue, 0, len(items))
	for _, item := range items {
		out = append(out, tupleValue{item.Key.Interface(), item.Value.Interface()})
	}
	return exec.AsValue(out)
}

func filterList(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.IsString() {
		out := []string{}
		for _, r := range in.String() {
			out = append(out, string(r))
		}
		return exec.AsValue(out)
	}
	out := make([]any, 0)
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		out = append(out, key.Interface())
		return true
	}, func() {})
	return exec.AsValue(out)
}

func filterLower(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	return exec.AsValue(strings.ToLower(in.String()))
}

func filterMap(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if in.IsNil() {
		return exec.AsValue([]any{})
	}

	filterName := ""
	filterArgs := exec.NewVarArgs()
	attribute := exec.AsValue(nil)
	defaultVal := exec.AsValue(nil)

	if len(params.Args) > 0 {
		filterName = params.Args[0].String()
		filterArgs.Args = append(filterArgs.Args, params.Args[1:]...)
		filterArgs.KwArgs = params.KwArgs
	} else {
		if attributeArg, ok := params.KwArgs["attribute"]; ok {
			attribute = attributeArg
		}
		if defaultArg, ok := params.KwArgs["default"]; ok {
			defaultVal = defaultArg
		}
		if len(params.KwArgs) > 2 {
			return exec.AsValue(errors.New("Wrong signature for 'map'"))
		}
	}

	out := make([]any, 0)
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if !attribute.IsNil() {
			attr, found := resolveAttributeValue(val, attribute, defaultVal)
			if found {
				val = attr
			} else {
				return true
			}
		}
		if filterName != "" {
			val = e.ExecuteFilterByName(filterName, val, filterArgs)
		}
		out = append(out, val.Interface())
		return true
	}, func() {})
	return exec.AsValue(out)
}

func filterMax(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		caseSensitive bool
		attribute     *exec.Value
	)
	if err := params.Take(
		exec.KeywordArgument("case_sensitive", exec.AsValue(false), takeBoolArgument(&caseSensitive)),
		exec.KeywordArgument("attribute", exec.AsValue(nil), takeValueArgument(&attribute)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}

	var max *exec.Value
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if attribute != nil && !attribute.IsNil() {
			attr, found := resolveAttributeValue(val, attribute, nil)
			if found {
				val = attr
			} else {
				val = nil
			}
		}
		if max == nil {
			max = val
			return true
		}
		if val == nil || max == nil {
			return true
		}
		switch {
		case max.IsFloat() || max.IsInteger() && val.IsFloat() || val.IsInteger():
			if val.Float() > max.Float() {
				max = val
			}
		case max.IsString() && val.IsString():
			if !caseSensitive && strings.ToLower(val.String()) > strings.ToLower(max.String()) {
				max = val
			} else if caseSensitive && val.String() > max.String() {
				max = val
			}
		default:
			max = exec.AsValue(errors.Errorf(`%s and %s are not comparable`, max.Val.Type(), val.Val.Type()))
		}
		return true
	}, func() {})

	if max == nil {
		return exec.AsValue("")
	}
	return max
}

func filterMin(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		caseSensitive bool
		attribute     *exec.Value
	)
	if err := params.Take(
		exec.KeywordArgument("case_sensitive", exec.AsValue(false), takeBoolArgument(&caseSensitive)),
		exec.KeywordArgument("attribute", exec.AsValue(nil), takeValueArgument(&attribute)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}

	var min *exec.Value
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if attribute != nil && !attribute.IsNil() {
			attr, found := resolveAttributeValue(val, attribute, nil)
			if found {
				val = attr
			} else {
				val = nil
			}
		}
		if min == nil {
			min = val
			return true
		}
		if val == nil || min == nil {
			return true
		}
		switch {
		case min.IsFloat() || min.IsInteger() && val.IsFloat() || val.IsInteger():
			if val.Float() < min.Float() {
				min = val
			}
		case min.IsString() && val.IsString():
			if !caseSensitive && strings.ToLower(val.String()) < strings.ToLower(min.String()) {
				min = val
			} else if caseSensitive && val.String() < min.String() {
				min = val
			}
		default:
			min = exec.AsValue(errors.Errorf(`%s and %s are not comparable`, min.Val.Type(), val.Val.Type()))
		}
		return true
	}, func() {})

	if min == nil {
		return exec.AsValue("")
	}
	return min
}

func filterPPrint(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var verbose bool
	if err := params.Take(
		exec.KeywordArgument("verbose", exec.AsValue(false), takeBoolArgument(&verbose)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	b, err := json.MarshalIndent(in.Interface(), "", "  ")
	if err != nil {
		return exec.AsValue(errors.Wrapf(err, `Unable to pretty print '%s'`, in.String()))
	}
	return exec.AsSafeValue(string(b))
}

func filterRandom(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if !in.CanSlice() || in.Len() <= 0 {
		return in
	}
	i := rand.Intn(in.Len())
	return in.Index(i)
}

func filterReject(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var test func(*exec.Value) bool
	if len(params.Args) == 0 {
		// Reject truthy value
		test = func(in *exec.Value) bool {
			return in.IsTrue()
		}
	} else {
		name := params.First().String()
		testParams := &exec.VarArgs{
			Args:   params.Args[1:],
			KwArgs: params.KwArgs,
		}
		test = func(in *exec.Value) bool {
			out := e.ExecuteTestByName(name, in, testParams)
			return out.IsTrue()
		}
	}

	out := make([]any, 0)

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		if !test(key) {
			out = append(out, key.Interface())
		}
		return true
	}, func() {})

	return exec.AsValue(out)
}

func filterRejectAttr(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if len(params.Args) < 1 {
		return exec.AsValue(errors.New("Wrong signature for 'rejectattr', expect at least an attribute name as argument"))
	}
	attribute := params.First()
	name := ""
	testArgs := []*exec.Value{}
	if len(params.Args) > 2 {
		testArgs = params.Args[2:]
	}
	testParams := &exec.VarArgs{Args: testArgs, KwArgs: params.KwArgs}
	if len(params.Args) > 1 {
		name = params.Args[1].String()
	}

	out := make([]any, 0)

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		attr, _ := resolveAttributeValue(key, attribute, nil)
		keep := false
		if name == "" {
			keep = !attr.IsTrue()
		} else {
			keep = !e.ExecuteTestByName(name, attr, testParams).IsTrue()
		}
		if keep {
			out = append(out, key.Interface())
		}
		return true
	}, func() {})

	return exec.AsValue(out)
}

func filterReplace(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		old   string
		new   string
		count *exec.Value
	)
	if err := params.Take(
		exec.PositionalArgument("old", nil, takeStringArgument(&old)),
		exec.PositionalArgument("new", nil, takeStringArgument(&new)),
		exec.KeywordArgument("count", exec.AsValue(nil), takeValueArgument(&count)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if count.IsNil() {
		return exec.AsValue(strings.ReplaceAll(in.String(), old, new))
	}
	return exec.AsValue(strings.Replace(in.String(), old, new, count.Integer()))
}

func filterReverse(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.IsString() {
		var out strings.Builder
		in.IterateOrder(func(idx, count int, key, value *exec.Value) bool {
			out.WriteString(key.String())
			return true
		}, func() {}, true, false, false)
		return exec.AsValue(out.String())
	}
	out := make([]any, 0)
	in.IterateOrder(func(idx, count int, key, value *exec.Value) bool {
		out = append(out, key.Interface())
		return true
	}, func() {}, true, true, false)
	return exec.AsValue(out)
}

func filterRound(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		precision int
		method    string
	)
	if err := params.Take(
		exec.KeywordArgument("precision", exec.AsValue(0), takeIntArgument(&precision)),
		exec.KeywordArgument("method", exec.AsValue("common"), takeStringArgument(&method)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	var op func(float64) float64
	switch method {
	case "common":
		op = math.Round
	case "floor":
		op = math.Floor
	case "ceil":
		op = math.Ceil
	default:
		return exec.AsValue(errors.Errorf(`Unknown method '%s', mush be one of 'common, 'floor', 'ceil`, method))
	}
	value := in.Float()
	factor := math.Pow(10, float64(precision))
	value = op(value*factor) / factor
	return exec.AsValue(value)
}

func filterSafe(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	in.Safe = true
	return in // nothing to do here, just to keep track of the safe application
}

func filterSelect(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var test func(*exec.Value) bool
	if len(params.Args) == 0 {
		// Reject truthy value
		test = func(in *exec.Value) bool {
			return in.IsTrue()
		}
	} else {
		name := params.First().String()
		testParams := &exec.VarArgs{
			Args:   params.Args[1:],
			KwArgs: params.KwArgs,
		}
		test = func(in *exec.Value) bool {
			out := e.ExecuteTestByName(name, in, testParams)
			return out.IsTrue()
		}
	}

	out := make([]any, 0)

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		if test(key) {
			out = append(out, key.Interface())
		}
		return true
	}, func() {})

	return exec.AsValue(out)
}

func filterSlice(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		slices   int
		fillWith any
	)
	if err := params.Take(
		exec.PositionalArgument("slices", nil, exec.IntArgument(&slices)),
		exec.KeywordArgument("fill_with", exec.AsValue(nil), exec.AnyArgument(&fillWith)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if slices < 1 {
		return exec.AsValue(exec.ErrInvalidCall(fmt.Errorf("slices argument %d must be > 0", slices)))
	}
	if !in.IsList() {
		return exec.AsValue(exec.ErrInvalidCall(fmt.Errorf("%s is not a list", in.String())))
	}
	seq := make([]any, 0, in.Len())
	in.Iterate(func(index, _ int, value, _ *exec.Value) bool {
		seq = append(seq, value.Interface())
		return true
	}, func() {})

	itemsPerSlice := len(seq) / slices
	slicesWithExtra := len(seq) % slices
	offset := 0
	output := make([]any, 0, slices)
	for sliceNumber := 0; sliceNumber < slices; sliceNumber++ {
		start := offset + sliceNumber*itemsPerSlice
		if sliceNumber < slicesWithExtra {
			offset++
		}
		end := offset + (sliceNumber+1)*itemsPerSlice
		column := append([]any{}, seq[start:end]...)
		if fillWith != nil && sliceNumber >= slicesWithExtra {
			column = append(column, fillWith)
		}
		output = append(output, column)
	}
	return exec.AsValue(output)
}

func filterSort(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		reverse       bool
		caseSensitive bool
		attribute     *exec.Value
	)
	if err := params.Take(
		exec.KeywordArgument("reverse", exec.AsValue(false), takeBoolArgument(&reverse)),
		exec.KeywordArgument("case_sensitive", exec.AsValue(false), takeBoolArgument(&caseSensitive)),
		exec.KeywordArgument("attribute", exec.AsValue(nil), takeValueArgument(&attribute)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	items := make([]*exec.Value, 0)
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		items = append(items, key)
		return true
	}, func() {})

	sort.SliceStable(items, func(i, j int) bool {
		if attribute.IsNil() {
			comparison := compareValues(items[i], items[j], caseSensitive)
			if reverse {
				return comparison > 0
			}
			return comparison < 0
		}

		for _, attr := range strings.Split(attribute.String(), ",") {
			left, _ := resolveAttributePath(items[i], attr, nil)
			right, _ := resolveAttributePath(items[j], attr, nil)
			comparison := compareValues(left, right, caseSensitive)
			if comparison == 0 {
				continue
			}
			if reverse {
				return comparison > 0
			}
			return comparison < 0
		}
		return false
	})

	out := make([]any, 0, len(items))
	for _, item := range items {
		out = append(out, item.Interface())
	}
	return exec.AsValue(out)
}

func filterString(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	return exec.AsValue(in.String())
}

var (
	reStripComments = regexp.MustCompile(`(?s)<!--.*?-->`)
	reStriptags     = regexp.MustCompile("<[^>]*?>")
)

func filterStriptags(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	s := in.String()
	s = reStripComments.ReplaceAllString(s, " ")

	// Strip all tags
	s = reStriptags.ReplaceAllString(s, " ")

	return exec.AsValue(strings.Join(strings.Fields(s), " "))
}

func filterSum(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		attribute *exec.Value
		start     float64
	)
	if err := params.Take(
		exec.KeywordArgument("attribute", exec.AsValue(nil), takeValueArgument(&attribute)),
		exec.KeywordArgument("start", exec.AsValue(0), takeFloatArgument(&start)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	sum := start

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if attribute != nil && !attribute.IsNil() {
			resolved, found := resolveAttributeValue(key, attribute, nil)
			if !found {
				return true
			}
			val = resolved
		}
		if val.IsNumber() {
			sum += val.Float()
		}
		return true
	}, func() {})

	if sum == math.Trunc(sum) {
		return exec.AsValue(int64(sum))
	}
	return exec.AsValue(sum)
}

func filterTitle(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.IsNil() {
		return exec.AsValue("")
	}
	return exec.AsValue(cases.Title(language.English).String(strings.ToLower(in.String())))
}

func filterTrim(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	var chars *exec.Value
	if err := params.Take(
		exec.KeywordArgument("chars", exec.AsValue(nil), takeValueArgument(&chars)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if !in.IsString() {
		return exec.AsValue(exec.ErrInvalidCall(fmt.Errorf("%s is not a string", in.String())))
	}
	if chars.IsNil() {
		return exec.AsValue(strings.TrimSpace(in.String()))
	}
	return exec.AsValue(strings.Trim(in.String(), chars.String()))
}

func filterToJSON(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}

	var (
		indent      *exec.Value
		ensureASCII bool
	)
	if err := params.Take(
		exec.KeywordArgument("indent", exec.AsValue(nil), takeValueArgument(&indent)),
		exec.KeywordArgument("ensure_ascii", exec.AsValue(true), takeBoolArgument(&ensureASCII)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	_ = ensureASCII // Accepted for compatibility, ignored (Go handles unicode differently than Python)

	marshalJSON := func(value any) (string, error) {
		if indent.IsNil() {
			b, err := stdjson.Marshal(value)
			if err != nil {
				return "", err
			}
			return strings.ReplaceAll(string(b), "'", `\u0027`), nil
		}
		if !indent.IsInteger() {
			return "", errors.Errorf("Expected an integer for 'indent', got %s", indent.String())
		}
		b, err := stdjson.MarshalIndent(value, "", strings.Repeat(" ", indent.Integer()))
		if err != nil {
			return "", err
		}
		return strings.ReplaceAll(string(b), "'", `\u0027`), nil
	}

	out, err := marshalJSON(normalizeJSONValue(in.Interface()))
	if err != nil {
		casted := in.ToGoSimpleType(false)
		if castErr, ok := casted.(error); ok {
			return exec.AsValue(castErr)
		}
		out, err = marshalJSON(casted)
		if err != nil {
			return exec.AsValue(errors.Wrap(err, "Unable to marhsall to json"))
		}
	}
	return exec.AsSafeValue(out)
}

func filterTruncate(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		length    int
		killWords bool
		end       string
		leeway    int
	)
	if err := params.Take(
		exec.KeywordArgument("length", exec.AsValue(255), takeIntArgument(&length)),
		exec.KeywordArgument("killwords", exec.AsValue(false), takeBoolArgument(&killWords)),
		exec.KeywordArgument("end", exec.AsValue("..."), takeStringArgument(&end)),
		exec.KeywordArgument("leeway", exec.AsValue(5), takeIntArgument(&leeway)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}

	source := in.String()
	rEnd := []rune(end)
	fullLength := length + leeway
	runes := []rune(source)

	if length < len(rEnd) {
		return exec.AsValue(errors.Errorf(`expected length >= %d, got %d`, len(rEnd), length))
	}

	if len(runes) <= fullLength {
		return exec.AsValue(source)
	}

	atLength := string(runes[:length-len(rEnd)])
	if !killWords {
		if split := strings.LastIndexFunc(atLength, unicode.IsSpace); split >= 0 {
			atLength = atLength[:split]
		}
		atLength = strings.TrimRight(atLength, " \n\t")
	}
	return exec.AsValue(fmt.Sprintf("%s%s", atLength, end))
}

func filterUnique(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		caseSensitive bool
		attribute     *exec.Value
	)
	if err := params.Take(
		exec.KeywordArgument("case_sensitive", exec.AsValue(false), takeBoolArgument(&caseSensitive)),
		exec.KeywordArgument("attribute", exec.AsValue(nil), takeValueArgument(&attribute)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}

	out := make([]any, 0)
	tracker := map[any]bool{}

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if !attribute.IsNil() {
			nested, found := resolveAttributeValue(key, attribute, nil)
			if !found {
				return true
			}
			val = nested
		}
		tracked := val.Interface()
		if !caseSensitive && val.IsString() {
			tracked = strings.ToLower(val.String())
		}
		if _, contains := tracker[tracked]; !contains {
			tracker[tracked] = true
			out = append(out, key.Interface())
		}
		return true
	}, func() {})

	return exec.AsValue(out)
}

func filterUpper(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	return exec.AsValue(strings.ToUpper(in.String()))
}

func filterUrlencode(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	switch {
	case in.IsString():
		return exec.AsValue(strings.ReplaceAll(url.PathEscape(in.String()), "%2F", "/"))
	case in.IsDict():
		pairs := make([]string, 0)
		in.IterateOrder(func(idx, count int, key, value *exec.Value) bool {
			pairs = append(pairs, fmt.Sprintf(
				"%s=%s",
				url.QueryEscape(stringifyFilterValue(key)),
				url.QueryEscape(stringifyFilterValue(value)),
			))
			return true
		}, func() {}, false, true, true)
		return exec.AsValue(strings.Join(pairs, "&"))
	case in.IsList():
		pairs := make([]string, 0)
		in.Iterate(func(idx, count int, key, value *exec.Value) bool {
			first, second, ok := urlEncodePair(key)
			if !ok {
				return true
			}
			pairs = append(pairs, fmt.Sprintf(
				"%s=%s",
				url.QueryEscape(stringifyFilterValue(first)),
				url.QueryEscape(stringifyFilterValue(second)),
			))
			return true
		}, func() {})
		return exec.AsValue(strings.Join(pairs, "&"))
	default:
		return exec.AsValue(strings.ReplaceAll(url.PathEscape(stringifyFilterValue(in)), "%2F", "/"))
	}
}

var (
	filterURLizeEmailRegexp  = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	filterURLizeDomainRegexp = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*\.[A-Za-z]{2,}[A-Za-z0-9./_-]*$`)
)

func filterUrlizeHelper(input string, trimURLLimit int, rel string, target string, extraSchemes []string) (string, error) {
	parts := strings.Split(input, " ")
	for idx, part := range parts {
		link, ok := urlizeToken(part, trimURLLimit, rel, target, extraSchemes)
		if ok {
			parts[idx] = link
		}
	}
	return strings.Join(parts, " "), nil
}

func filterUrlize(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var (
		trimURLLimitArgument *exec.Value
		nofollow             bool
		target               *exec.Value
		rel                  *exec.Value
		extraSchemesArgument *exec.Value
	)
	if err := params.Take(
		exec.KeywordArgument("trim_url_limit", exec.AsValue(nil), takeValueArgument(&trimURLLimitArgument)),
		exec.KeywordArgument("nofollow", exec.AsValue(false), takeBoolArgument(&nofollow)),
		exec.KeywordArgument("target", exec.AsValue(nil), takeValueArgument(&target)),
		exec.KeywordArgument("rel", exec.AsValue(nil), takeValueArgument(&rel)),
		exec.KeywordArgument("extra_schemes", exec.AsValue(nil), takeValueArgument(&extraSchemesArgument)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	trimURLLimit := -1
	if trimURLLimitArgument.IsInteger() {
		trimURLLimit = trimURLLimitArgument.Integer()
	}
	extraSchemes := make([]string, 0)
	if !extraSchemesArgument.IsNil() {
		if !extraSchemesArgument.IsList() {
			return exec.AsValue(errors.New("extra_schemes must be a list"))
		}
		for i := 0; i < extraSchemesArgument.Len(); i++ {
			item := extraSchemesArgument.Index(i)
			if !item.IsString() {
				return exec.AsValue(errors.New("extra_schemes must contain strings"))
			}
			extraSchemes = append(extraSchemes, item.String())
		}
	}

	relValue := ""
	if rel.IsNil() {
		relValue = "noopener"
	} else {
		relValue = rel.String()
	}
	if nofollow {
		parts := []string{}
		if relValue != "" {
			parts = append(parts, relValue)
		}
		parts = append(parts, "nofollow")
		relValue = strings.Join(parts, " ")
	}

	s, err := filterUrlizeHelper(in.String(), trimURLLimit, relValue, target.String(), extraSchemes)
	if err != nil {
		return exec.AsValue(err)
	}

	if e.Config.AutoEscape {
		return exec.AsSafeValue(s)
	}
	return exec.AsValue(s)
}

func filterWordcount(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if err := params.Take(); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	return exec.AsValue(len(strings.Fields(in.String())))
}

func filterWordwrap(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var width int
	if err := params.Take(exec.PositionalArgument("width", exec.AsValue(79), takeIntArgument(&width))); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if width <= 0 {
		return in
	}

	lines := make([]string, 0)
	for _, paragraph := range strings.Split(in.String(), "\n") {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}
		current := words[0]
		for _, word := range words[1:] {
			if len([]rune(current))+1+len([]rune(word)) <= width {
				current += " " + word
				continue
			}
			lines = append(lines, current)
			current = word
		}
		lines = append(lines, current)
	}
	return exec.AsValue(strings.Join(lines, "\n"))
}

func filterXMLAttr(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var autospace bool
	if err := params.Take(
		exec.KeywordArgument("autospace", exec.AsValue(true), takeBoolArgument(&autospace)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	kvs := []string{}
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		if !value.IsTrue() {
			return true
		}
		kv := fmt.Sprintf(`%s="%s"`, key.Escaped(), value.Escaped())
		kvs = append(kvs, kv)
		return true
	}, func() {})
	out := strings.Join(kvs, " ")
	if autospace {
		out = " " + out
	}
	return exec.AsValue(out)
}

func filterDefault(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	var (
		defaultValue *exec.Value
		boolean      bool
	)
	if err := params.Take(
		exec.PositionalArgument("default", nil, takeValueArgument(&defaultValue)),
		exec.KeywordArgument("boolean", exec.AsValue(false), exec.BoolArgument(&boolean)),
	); err != nil {
		return exec.AsValue(exec.ErrInvalidCall(err))
	}
	if in.IsError() || in.IsNil() {
		return defaultValue
	}
	if boolean && !in.IsTrue() {
		return defaultValue
	}
	return in
}

func filterSelectAttr(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if len(params.Args) < 1 {
		return exec.AsValue(errors.New("Wrong signature for 'selectattr', expect at least an attribute name as argument"))
	}
	attribute := params.First()
	name := ""
	testArgs := []*exec.Value{}
	if len(params.Args) > 2 {
		testArgs = params.Args[2:]
	}
	testParams := &exec.VarArgs{Args: testArgs, KwArgs: params.KwArgs}
	if len(params.Args) > 1 {
		name = params.Args[1].String()
	}

	out := make([]any, 0)

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		attr, _ := resolveAttributeValue(key, attribute, nil)
		matched := false
		if name == "" {
			matched = attr.IsTrue()
		} else {
			matched = e.ExecuteTestByName(name, attr, testParams).IsTrue()
		}
		if matched {
			out = append(out, key.Interface())
		}
		return true
	}, func() {})

	return exec.AsValue(out)
}
