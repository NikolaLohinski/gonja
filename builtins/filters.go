package builtins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	json "github.com/json-iterator/go"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/yargevad/filepathx"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/utils"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Filters export all builtin filters
var Filters = exec.FilterSet{
	"abs":            filterAbs,
	"add":            filterAdd,
	"append":         filterAppend,
	"attr":           filterAttr,
	"basename":       filterBasename,
	"batch":          filterBatch,
	"bool":           filterBool,
	"capitalize":     filterCapitalize,
	"center":         filterCenter,
	"concat":         filterConcat,
	"default":        filterDefault,
	"d":              filterDefault,
	"dictsort":       filterDictSort,
	"dir":            filterDir,
	"e":              filterEscape,
	"escape":         filterEscape,
	"fail":           filterFail,
	"file":           filterFile,
	"fileset":        filterFileset,
	"filesizeformat": filterFileSize,
	"first":          filterFirst,
	"flatten":        filterFlatten,
	"float":          filterFloat,
	"forceescape":    filterForceEscape,
	"format":         filterFormat,
	"fromjson":       filterFromJSON,
	"fromyaml":       filterFromYAML,
	"fromtoml":       filterFromTOML,
	"get":            filterGet,
	"groupby":        filterGroupBy,
	"ifelse":         filterIfElse,
	"indent":         filterIndent,
	"insert":         filterInsert,
	"int":            filterInteger,
	"join":           filterJoin,
	"keys":           filterKeys,
	"last":           filterLast,
	"length":         filterLength,
	"list":           filterList,
	"lower":          filterLower,
	"map":            filterMap,
	"match":          filterMatch,
	"max":            filterMax,
	"min":            filterMin,
	"panic":          filterPanic,
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
	"split":          filterSplit,
	"string":         filterString,
	"striptags":      filterStriptags,
	"sum":            filterSum,
	"title":          filterTitle,
	"tojson":         filterToJSON,
	"totoml":         filterToToml,
	"toyaml":         filterToYAML,
	"trim":           filterTrim,
	"truncate":       filterTruncate,
	"try":            filterTry,
	"unique":         filterUnique,
	"unset":          filterUnset,
	"upper":          filterUpper,
	"urlencode":      filterUrlencode,
	"urlize":         filterUrlize,
	"values":         filterValues,
	"wordcount":      filterWordcount,
	"wordwrap":       filterWordwrap,
	"xmlattr":        filterXMLAttr,
}

func filterAbs(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'abs'"))
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
	p := params.ExpectArgs(1)
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'attr'"))
	}
	attr := p.First().String()
	value, _ := in.Getattr(attr)
	return value
}

func filterBatch(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(1, []*exec.KwArg{{"fill_with", nil}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'batch'"))
	}
	size := p.First().Integer()
	out := make([]interface{}, 0)
	var row []interface{}
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		if math.Mod(float64(idx), float64(size)) == 0 {
			if row != nil {
				out = append(out, exec.AsValue(row).Interface())
			}
			row = make([]interface{}, 0)
		}
		row = append(row, key.Interface())
		return true
	}, func() {})
	if len(row) > 0 {
		fillWith := p.KwArgs["fill_with"]
		if !fillWith.IsNil() {
			for len(row) < size {
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
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'capitalize'"))
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
	p := params.ExpectArgs(1)
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'center'"))
	}
	width := p.First().Integer()
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

func sortByKey(in *exec.Value, caseSensitive bool, reverse bool) [][2]interface{} {
	out := make([][2]interface{}, 0)
	in.IterateOrder(func(idx, count int, key, value *exec.Value) bool {
		out = append(out, [2]interface{}{key.Interface(), value.Interface()})
		return true
	}, func() {}, reverse, true, caseSensitive)
	return out
}

func sortByValue(in *exec.Value, caseSensitive, reverse bool) [][2]interface{} {
	out := make([][2]interface{}, 0)
	items := in.Items()
	var sorter func(i, j int) bool
	switch {
	case caseSensitive && reverse:
		sorter = func(i, j int) bool {
			return items[i].Value.String() > items[j].Value.String()
		}
	case caseSensitive && !reverse:
		sorter = func(i, j int) bool {
			return items[i].Value.String() < items[j].Value.String()
		}
	case !caseSensitive && reverse:
		sorter = func(i, j int) bool {
			return strings.ToLower(items[i].Value.String()) > strings.ToLower(items[j].Value.String())
		}
	case !caseSensitive && !reverse:
		sorter = func(i, j int) bool {
			return strings.ToLower(items[i].Value.String()) < strings.ToLower(items[j].Value.String())
		}
	}
	sort.Slice(items, sorter)
	for _, item := range items {
		out = append(out, [2]interface{}{item.Key.Interface(), item.Value.Interface()})
	}
	return out
}

func filterDictSort(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{
		{"case_sensitive", false},
		{"by", "key"},
		{"reverse", false},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'dictsort'"))
	}

	caseSensitive := p.KwArgs["case_sensitive"].Bool()
	by := p.KwArgs["by"].String()
	reverse := p.KwArgs["reverse"].Bool()

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
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'escape'"))
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
	p := params.Expect(0, []*exec.KwArg{{"binary", false}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'filesizeformat'"))
	}
	bytes := in.Float()
	binary := p.KwArgs["binary"].Bool()
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
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'first'"))
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
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'float'"))
	}
	return exec.AsValue(in.Float())
}

func filterForceEscape(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'forceescape'"))
	}
	return exec.AsSafeValue(in.Escaped())
}

func filterFormat(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	args := []interface{}{}
	for _, arg := range params.Args {
		args = append(args, arg.Interface())
	}
	return exec.AsValue(fmt.Sprintf(in.String(), args...))
}

func filterGroupBy(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.ExpectArgs(1)
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'groupby"))
	}
	field := p.First().String()
	groups := make(map[interface{}][]interface{})
	groupers := []interface{}{}

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		attr, found := key.Get(field)
		if !found {
			return true
		}
		lst, exists := groups[attr.Interface()]
		if !exists {
			lst = make([]interface{}, 0)
			groupers = append(groupers, attr.Interface())
		}
		lst = append(lst, key.Interface())
		groups[attr.Interface()] = lst
		return true
	}, func() {})

	out := make([]map[string]interface{}, 0)
	for _, grouper := range groupers {
		out = append(out, map[string]interface{}{
			"grouper": exec.AsValue(grouper).Interface(),
			"list":    exec.AsValue(groups[grouper]).Interface(),
		})
	}
	return exec.AsValue(out)
}

func filterIndent(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{
		{"width", 4},
		{"first", false},
		{"blank", false},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'indent'"))
	}
	width := p.KwArgs["width"].Integer()
	first := p.KwArgs["first"].Bool()
	blank := p.KwArgs["blank"].Bool()
	indent := strings.Repeat(" ", width)
	lines := strings.Split(in.String(), "\n")
	// start := 1
	// if first {start = 0}
	var out strings.Builder
	for idx, line := range lines {
		if line == "" && !blank {
			out.WriteByte('\n')
			continue
		}
		if idx > 0 || first {
			out.WriteString(indent)
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}
	return exec.AsValue(out.String())
}

func filterInteger(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'int'"))
	}
	return exec.AsValue(in.Integer())
}

func filterBool(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'bool'"))
	}
	switch {
	case in.IsBool():
		return exec.AsValue(in.Bool())
	case in.IsString():
		trues := []string{"true", "yes", "on", "1"}
		falses := []string{"false", "no", "off", "0", ""}
		loweredString := strings.ToLower(in.String())
		if slices.Contains(trues, loweredString) {
			return exec.AsValue(true)
		} else if slices.Contains(falses, loweredString) {
			return exec.AsValue(false)
		} else {
			return exec.AsValue(fmt.Errorf("\"%s\" can not be cast to boolean as it's not in [\"%s\"] nor [\"%s\"]", in.String(), strings.Join(trues, "\",\""), strings.Join(falses, "\",\"")))
		}
	case in.IsInteger():
		if in.Integer() == 1 {
			return exec.AsValue(true)
		} else if in.Integer() == 0 {
			return exec.AsValue(false)
		} else {
			return exec.AsValue(fmt.Errorf("%d can not be cast to boolean as it's not in [0,1]", in.Integer()))
		}
	case in.IsFloat():
		if in.Float() == 1.0 {
			return exec.AsValue(true)
		} else if in.Float() == 0.0 {
			return exec.AsValue(false)
		} else {
			return exec.AsValue(fmt.Errorf("%f can not be cast to boolean as it's not in [0.0,1.0]", in.Float()))
		}
	case in.IsNil():
		return exec.AsValue(false)
	default:
		return exec.AsValue(fmt.Errorf("filter 'bool' failed to cast: %s", in.String()))
	}
}

func filterJoin(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{
		{"d", ""},
		{"attribute", nil},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'join'"))
	}
	if !in.CanSlice() {
		return in
	}
	sep := p.KwArgs["d"].String()
	sl := make([]string, 0, in.Len())
	for i := 0; i < in.Len(); i++ {
		sl = append(sl, in.Index(i).String())
	}
	return exec.AsValue(strings.Join(sl, sep))
}

func filterLast(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'last'"))
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
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'length'"))
	}
	return exec.AsValue(in.Len())
}

func filterList(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'list'"))
	}
	if in.IsString() {
		out := []string{}
		for _, r := range in.String() {
			out = append(out, string(r))
		}
		return exec.AsValue(out)
	}
	out := make([]interface{}, 0)
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
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'lower'"))
	}
	return exec.AsValue(strings.ToLower(in.String()))
}

func filterMap(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{
		{"filter", ""},
		{"attribute", nil},
		{"default", nil},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'map'"))
	}
	filter := p.KwArgs["filter"].String()
	attribute := p.KwArgs["attribute"].String()
	defaultVal := p.KwArgs["default"]
	out := make([]interface{}, 0)
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if len(attribute) > 0 {
			attr, found := val.Get(attribute)
			if found {
				val = attr
			} else if defaultVal != nil {
				val = defaultVal
			} else {
				return true
			}
		}
		if len(filter) > 0 {
			val = e.ExecuteFilterByName(filter, val, exec.NewVarArgs())
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
	p := params.Expect(0, []*exec.KwArg{
		{"case_sensitive", false},
		{"attribute", nil},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'max'"))
	}
	caseSensitive := p.KwArgs["case_sensitive"].Bool()
	attribute := p.KwArgs["attribute"].String()

	var max *exec.Value
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if len(attribute) > 0 {
			attr, found := val.Get(attribute)
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
	p := params.Expect(0, []*exec.KwArg{
		{"case_sensitive", false},
		{"attribute", nil},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'min'"))
	}
	caseSensitive := p.KwArgs["case_sensitive"].Bool()
	attribute := p.KwArgs["attribute"].String()

	var min *exec.Value
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if len(attribute) > 0 {
			attr, found := val.Get(attribute)
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
	p := params.Expect(0, []*exec.KwArg{{"verbose", false}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'pprint'"))
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
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'random'"))
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

	out := make([]interface{}, 0)

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
	var test func(*exec.Value) *exec.Value
	if len(params.Args) < 1 {
		return exec.AsValue(errors.New("Wrong signature for 'rejectattr', expect at least an attribute name as argument"))
	}
	attribute := params.First().String()
	if len(params.Args) == 1 {
		// Reject truthy value
		test = func(in *exec.Value) *exec.Value {
			attr, found := in.Get(attribute)
			if !found {
				return exec.AsValue(errors.Errorf(`%s has no attribute '%s'`, in.String(), attribute))
			}
			return attr
		}
	} else {
		name := params.Args[1].String()
		testParams := &exec.VarArgs{
			Args:   params.Args[2:],
			KwArgs: params.KwArgs,
		}
		test = func(in *exec.Value) *exec.Value {
			attr, found := in.Get(attribute)
			if !found {
				return exec.AsValue(errors.Errorf(`%s has no attribute '%s'`, in.String(), attribute))
			}
			out := e.ExecuteTestByName(name, attr, testParams)
			return out
		}
	}

	out := make([]interface{}, 0)
	var err *exec.Value

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		result := test(key)
		if result.IsError() {
			err = result
			return false
		}
		if !result.IsTrue() {
			out = append(out, key.Interface())
		}
		return true
	}, func() {})

	if err != nil {
		return err
	}
	return exec.AsValue(out)
}

func filterReplace(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(2, []*exec.KwArg{{"count", nil}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'replace'"))
	}
	old := p.Args[0].String()
	new := p.Args[1].String()
	count := p.KwArgs["count"]
	if count.IsNil() {
		return exec.AsValue(strings.ReplaceAll(in.String(), old, new))
	}
	return exec.AsValue(strings.Replace(in.String(), old, new, count.Integer()))
}

func filterReverse(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'safe'"))
	}
	if in.IsString() {
		var out strings.Builder
		in.IterateOrder(func(idx, count int, key, value *exec.Value) bool {
			out.WriteString(key.String())
			return true
		}, func() {}, true, false, false)
		return exec.AsValue(out.String())
	}
	out := make([]interface{}, 0)
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
	p := params.Expect(0, []*exec.KwArg{{"precision", 0}, {"method", "common"}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'round'"))
	}
	method := p.KwArgs["method"].String()
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
	factor := float64(10 * p.KwArgs["precision"].Integer())
	if factor > 0 {
		value = value * factor
	}
	value = op(value)
	if factor > 0 {
		value = value / factor
	}
	return exec.AsValue(value)
}

func filterSafe(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'safe'"))
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

	out := make([]interface{}, 0)

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
	comp := strings.Split(params.Args[0].String(), ":")
	if len(comp) != 2 {
		return exec.AsValue(errors.New("Slice string must have the format 'from:to' [from/to can be omitted, but the ':' is required]"))
	}

	if !in.CanSlice() {
		return in
	}

	from := exec.AsValue(comp[0]).Integer()
	to := in.Len()

	if from > to {
		from = to
	}

	vto := exec.AsValue(comp[1]).Integer()
	if vto >= from && vto <= in.Len() {
		to = vto
	}

	return in.Slice(from, to)
}

func filterSort(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{{"reverse", false}, {"case_sensitive", false}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'sort'"))
	}
	reverse := p.KwArgs["reverse"].Bool()
	caseSensitive := p.KwArgs["case_sensitive"].Bool()
	out := make([]interface{}, 0)
	in.IterateOrder(func(idx, count int, key, value *exec.Value) bool {
		out = append(out, key.Interface())
		return true
	}, func() {}, reverse, true, caseSensitive)
	return exec.AsValue(out)
}

func filterString(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'string'"))
	}
	return exec.AsValue(in.String())
}

var reStriptags = regexp.MustCompile("<[^>]*?>")

func filterStriptags(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'striptags'"))
	}
	s := in.String()

	// Strip all tags
	s = reStriptags.ReplaceAllString(s, "")

	return exec.AsValue(strings.TrimSpace(s))
}

func filterSum(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{{"attribute", nil}, {"start", 0}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'sum'"))
	}

	attribute := p.KwArgs["attribute"]
	sum := p.KwArgs["start"].Float()
	var err error

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		if attribute.IsString() {
			val := key
			found := true
			for _, attr := range strings.Split(attribute.String(), ".") {
				val, found = val.Get(attr)
				if !found {
					err = errors.Errorf("'%s' has no attribute '%s'", key.String(), attribute.String())
					return false
				}
			}
			if found && val.IsNumber() {
				sum += val.Float()
			}
		} else if attribute.IsInteger() {
			value, found := key.Getitem(attribute.Integer())
			if found {
				sum += value.Float()
			}
		} else {
			sum += key.Float()
		}
		return true
	}, func() {})

	if err != nil {
		return exec.AsValue(err)
	} else if sum == math.Trunc(sum) {
		return exec.AsValue(int64(sum))
	}
	return exec.AsValue(sum)
}

func filterTitle(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'title'"))
	}
	if !in.IsString() {
		return exec.AsValue("")
	}
	return exec.AsValue(strings.Title(strings.ToLower(in.String())))
}

func filterTrim(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	charsParam := exec.KwArg{
		Name:    "chars",
		Default: " ",
	}
	p := params.ExpectKwArgs([]*exec.KwArg{&charsParam})
	if p.IsError() || !in.IsString() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'trim'"))
	}
	chars := p.GetKwarg(charsParam.Name, charsParam.Default).String()
	return exec.AsValue(strings.Trim(in.String(), chars))
}

func filterToJSON(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	// Done not mess around with trying to marshall error pipelines
	if in.IsError() {
		return in
	}

	// Monkey patching because arrays handling is broken
	if in.IsList() {
		inCast := make([]interface{}, in.Len())
		for index := range inCast {
			item := exec.ToValue(in.Index(index).Val)
			inCast[index] = item.Val.Interface()
		}
		in = exec.AsValue(inCast)
	}

	p := params.Expect(0, []*exec.KwArg{{"indent", nil}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'tojson'"))
	}

	casted := in.ToGoSimpleType(true)
	if err, ok := casted.(error); ok {
		return exec.AsValue(err)
	}

	indent := p.KwArgs["indent"]
	var out string
	if indent.IsNil() {
		b, err := json.ConfigCompatibleWithStandardLibrary.Marshal(casted)
		if err != nil {
			return exec.AsValue(errors.Wrap(err, "Unable to marhsall to json"))
		}
		out = string(b)
	} else if indent.IsInteger() {
		b, err := json.ConfigCompatibleWithStandardLibrary.MarshalIndent(casted, "", strings.Repeat(" ", indent.Integer()))
		if err != nil {
			return exec.AsValue(errors.Wrap(err, "Unable to marhsall to json"))
		}
		out = string(b)
	} else {
		return exec.AsValue(errors.Errorf("Expected an integer for 'indent', got %s", indent.String()))
	}
	return exec.AsSafeValue(out)
}

func filterTruncate(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{
		{"length", 255},
		{"killwords", false},
		{"end", "..."},
		{"leeway", 0},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'truncate'"))
	}

	source := in.String()
	length := p.KwArgs["length"].Integer()
	leeway := p.KwArgs["leeway"].Integer()
	killwords := p.KwArgs["killwords"].Bool()
	end := p.KwArgs["end"].String()
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
	if !killwords {
		atLength = strings.TrimRightFunc(atLength, func(r rune) bool {
			return !unicode.IsSpace(r)
		})
		atLength = strings.TrimRight(atLength, " \n\t")
	}
	return exec.AsValue(fmt.Sprintf("%s%s", atLength, end))
}

func filterUnique(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{{"case_sensitive", false}, {"attribute", nil}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'unique'"))
	}

	caseSensitive := p.KwArgs["case_sensitive"].Bool()
	attribute := p.KwArgs["attribute"]

	out := make([]interface{}, 0)
	tracker := map[interface{}]bool{}
	var err error

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		val := key
		if attribute.IsString() {
			attr := attribute.String()
			nested, found := key.Get(attr)
			if !found {
				err = errors.Errorf(`%s has no attribute %s`, key.String(), attr)
				return false
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

	if err != nil {
		return exec.AsValue(err)
	}
	return exec.AsValue(out)
}

func filterUpper(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'upper'"))
	}
	return exec.AsValue(strings.ToUpper(in.String()))
}

func filterUrlencode(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'urlencode'"))
	}
	return exec.AsValue(url.QueryEscape(in.String()))
}

// TODO: This regexp could do some work
var filterUrlizeURLRegexp = regexp.MustCompile(`((((http|https)://)|www\.|((^|[ ])[0-9A-Za-z_\-]+(\.com|\.net|\.org|\.info|\.biz|\.de))))(?U:.*)([ ]+|$)`)
var filterUrlizeEmailRegexp = regexp.MustCompile(`(\w+@\w+\.\w{2,4})`)

func filterUrlizeHelper(input string, trunc int, rel string, target string) (string, error) {
	var soutErr error
	sout := filterUrlizeURLRegexp.ReplaceAllStringFunc(input, func(raw_url string) string {
		var prefix string
		var suffix string
		if strings.HasPrefix(raw_url, " ") {
			prefix = " "
		}
		if strings.HasSuffix(raw_url, " ") {
			suffix = " "
		}

		raw_url = strings.TrimSpace(raw_url)

		url := utils.IRIEncode(raw_url)

		if !strings.HasPrefix(url, "http") {
			url = fmt.Sprintf("http://%s", url)
		}

		title := raw_url

		if trunc > 3 && len(title) > trunc {
			title = fmt.Sprintf("%s...", title[:trunc-3])
		}

		title = utils.Escape(title)

		attrs := ""
		if len(target) > 0 {
			attrs = fmt.Sprintf(` target="%s"`, target)
		}

		rels := []string{}
		cleanedRel := strings.Trim(strings.Replace(rel, "noopener", "", -1), " ")
		if len(cleanedRel) > 0 {
			rels = append(rels, cleanedRel)
		}
		rels = append(rels, "noopener")
		rel = strings.Join(rels, " ")

		return fmt.Sprintf(`%s<a href="%s" rel="%s"%s>%s</a>%s`, prefix, url, rel, attrs, title, suffix)
	})
	if soutErr != nil {
		return "", soutErr
	}

	sout = filterUrlizeEmailRegexp.ReplaceAllStringFunc(sout, func(mail string) string {
		title := mail

		if trunc > 3 && len(title) > trunc {
			title = fmt.Sprintf("%s...", title[:trunc-3])
		}

		return fmt.Sprintf(`<a href="mailto:%s">%s</a>`, mail, title)
	})
	return sout, nil
}

func filterUrlize(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(0, []*exec.KwArg{
		{"trim_url_limit", nil},
		{"nofollow", false},
		{"target", nil},
		{"rel", nil},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'urlize'"))
	}
	truncate := -1
	if param := p.KwArgs["trim_url_limit"]; param.IsInteger() {
		truncate = param.Integer()
	}
	rel := p.KwArgs["rel"]
	target := p.KwArgs["target"]

	s, err := filterUrlizeHelper(in.String(), truncate, rel.String(), target.String())
	if err != nil {
		return exec.AsValue(err)
	}

	return exec.AsValue(s)
}

func filterWordcount(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'wordcount'"))
	}
	return exec.AsValue(len(strings.Fields(in.String())))
}

func filterWordwrap(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	words := strings.Fields(in.String())
	wordsLen := len(words)
	wrapAt := params.Args[0].Integer()
	if wrapAt <= 0 {
		return in
	}

	linecount := wordsLen/wrapAt + wordsLen%wrapAt
	lines := make([]string, 0, linecount)
	for i := 0; i < linecount; i++ {
		lines = append(lines, strings.Join(words[wrapAt*i:utils.Min(wrapAt*(i+1), wordsLen)], " "))
	}
	return exec.AsValue(strings.Join(lines, "\n"))
}

func filterXMLAttr(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.ExpectKwArgs([]*exec.KwArg{{"autospace", true}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'xmlattr'"))
	}
	autospace := p.KwArgs["autospace"].Bool()
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

func filterIfElse(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.ExpectArgs(2)
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'ifelse'"))
	}
	if in.IsTrue() {
		return p.Args[0]
	} else {
		return p.Args[1]
	}
}

func filterGet(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.Expect(1, []*exec.KwArg{
		{Name: "strict", Default: false},
		{Name: "default", Default: nil},
	})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'get'"))
	}
	if !in.IsDict() {
		return exec.AsValue(errors.New("Filter 'get' was passed a non-dict type"))
	}
	item := p.First().String()
	value, ok := in.Getitem(item)
	if !ok {
		if fallback := p.GetKwarg("default", nil); !fallback.IsNil() {
			return fallback
		}
		if p.GetKwarg("strict", false).Bool() {
			return exec.AsValue(fmt.Errorf("item '%s' not found in: %s", item, in.String()))
		}
	}
	return value
}

func filterValues(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'values'"))
	}

	if !in.IsDict() {
		return exec.AsValue(errors.New("Filter 'values' was passed a non-dict type"))
	}

	out := make([]interface{}, 0)
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		out = append(out, value.Interface())
		return true
	}, func() {})

	return exec.AsValue(out)
}

func filterKeys(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'keys'"))
	}
	if !in.IsDict() {
		return exec.AsValue(errors.New("Filter 'keys' was passed a non-dict type"))
	}
	out := make([]interface{}, 0)
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		out = append(out, key.Interface())
		return true
	}, func() {})
	return exec.AsValue(out)
}

func filterTry(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'try'"))
	}
	if in == nil || in.IsError() || !in.IsTrue() {
		return exec.AsValue(nil)
	}
	return in
}

func filterFromJSON(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'fromjson'"))
	}

	if !in.IsString() || in.String() == "" {
		return exec.AsValue(errors.New("Filter 'fromjson' was passed an empty or non-string type"))
	}
	object := new(interface{})
	// first check if it's a JSON indeed
	if err := json.Unmarshal([]byte(in.String()), object); err != nil {
		return exec.AsValue(fmt.Errorf("failed to unmarshal %s: %s", in.String(), err))
	}
	// then use YAML because native JSON lib does not handle integers properly
	if err := yaml.Unmarshal([]byte(in.String()), object); err != nil {
		return exec.AsValue(fmt.Errorf("failed to unmarshal %s: %s", in.String(), err))
	}
	return exec.AsValue(*object)
}

func filterConcat(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsList() {
		return exec.AsValue(errors.New("Filter 'concat' was passed a non-list type"))
	}
	out := make([]interface{}, 0)
	in.Iterate(func(idx, count int, item, _ *exec.Value) bool {
		out = append(out, item.Interface())
		return true
	}, func() {})
	for index, argument := range params.Args {
		if !argument.IsList() {
			return exec.AsValue(fmt.Errorf("%s argument passed to filter 'concat' is not a list: %s", humanize.Ordinal(index+1), argument))
		}
		argument.Iterate(func(idx, count int, item, _ *exec.Value) bool {
			out = append(out, item.Interface())
			return true
		}, func() {})
	}
	return exec.AsValue(out)
}

func filterSplit(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsString() {
		return exec.AsValue(errors.New("Filter 'split' was passed a non-string type"))
	}
	p := params.ExpectArgs(1)
	if p.IsError() || !p.First().IsString() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'split'"))
	}
	delimiter := p.First().String()

	list := strings.Split(in.String(), delimiter)

	out := make([]interface{}, len(list))
	for index, item := range list {
		out[index] = item
	}

	return exec.AsValue(out)
}

func filterAdd(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}

	if in.IsList() {
		return filterAppend(e, in, params)
	}

	if in.IsDict() {
		return filterInsert(e, in, params)
	}

	return exec.AsValue(errors.New("Filter 'add' was passed a non-dict nor list type"))
}

func filterFail(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return exec.AsValue(fmt.Errorf("%s: %s", in.String(), in.Error()))
	}
	if p := params.ExpectNothing(); p.IsError() || !in.IsString() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'fail'"))
	}

	return exec.AsValue(errors.New(in.String()))
}

func filterInsert(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsDict() {
		return exec.AsValue(errors.New("Filter 'insert' was passed a non-dict type"))
	}
	p := params.ExpectArgs(2)
	if p.IsError() || len(p.Args) != 2 {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'insert'"))
	}
	newKey := p.Args[0]
	newValue := p.Args[1]

	out := make(map[string]interface{})
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		out[key.String()] = value.Interface()
		return true
	}, func() {})
	out[newKey.String()] = newValue.Interface()
	return exec.AsValue(out)
}

func filterUnset(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsDict() {
		return exec.AsValue(errors.New("Filter 'unset' was passed a non-dict type"))
	}
	p := params.ExpectArgs(1)
	if p.IsError() || len(p.Args) != 1 {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'unset'"))
	}
	toRemove := p.Args[0]

	out := make(map[string]interface{})
	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		if key.String() == toRemove.String() {
			return true
		}
		out[key.String()] = value.Interface()
		return true
	}, func() {})

	return exec.AsValue(out)
}

func filterAppend(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsList() {
		return exec.AsValue(errors.New("Filter 'append' was passed a non-list type"))
	}

	p := params.ExpectArgs(1)
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'append'"))
	}
	newItem := p.First()

	out := make([]interface{}, 0)
	in.Iterate(func(idx, count int, item, _ *exec.Value) bool {
		out = append(out, item.Interface())
		return true
	}, func() {})
	out = append(out, newItem)

	return exec.AsValue(out)
}

func filterFlatten(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsList() {
		return exec.AsValue(errors.New("Filter 'flatten' was passed a non-list type"))
	}

	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'flatten'"))
	}

	out := make([]interface{}, 0)
	in.Iterate(func(_, _ int, item, _ *exec.Value) bool {
		if !item.IsList() {
			out = append(out, item.Interface())
		} else {
			item.Iterate(func(_, _ int, subItem, _ *exec.Value) bool {
				out = append(out, subItem.Interface())
				return true
			}, func() {})
		}
		return true
	}, func() {})

	return exec.AsValue(out)
}

func filterFileset(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsString() {
		return exec.AsValue(errors.New("Filter 'fileset' was passed a non-string type"))
	}

	p := params.ExpectNothing()
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'fileset'"))
	}

	base, err := e.Loader.Path(".")
	if err != nil {
		return exec.AsValue(fmt.Errorf("failed to resolve path %s with loader: %s", in.String(), err))
	}
	out, err := filepathx.Glob(path.Join(base, in.String()))
	if err != nil {
		return exec.AsValue(fmt.Errorf("failed to traverse %s: %s", in.String(), err))
	}
	return exec.AsValue(out)
}

func filterFile(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsString() {
		return exec.AsValue(errors.New("Filter 'file' was passed a non-string type"))
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'file'"))
	}

	path := in.String()
	if !filepath.IsAbs(path) {
		base, err := e.Loader.Path(".")
		if err != nil {
			return exec.AsValue(fmt.Errorf("failed to get current path with loader: %s", err))
		}
		path, err = filepath.Abs(filepath.Join(base, path))
		if err != nil {
			return exec.AsValue(fmt.Errorf("failed to resolve path %s with loader: %s", path, err))
		}
	}

	out, err := ioutil.ReadFile(path)
	if err != nil {
		return exec.AsValue(fmt.Errorf("failed to read file at path %s: %s", path, err))
	}

	return exec.AsValue(string(out))
}

func filterBasename(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsString() {
		return exec.AsValue(errors.New("Filter 'basename' was passed a non-string type"))
	}

	p := params.ExpectNothing()
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'basename'"))
	}

	return exec.AsValue(filepath.Base(in.String()))
}

func filterDir(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if !in.IsString() {
		return exec.AsValue(errors.New("Filter 'dir' was passed a non-string type"))
	}

	p := params.ExpectNothing()
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'dir'"))
	}

	return exec.AsValue(filepath.Dir(in.String()))
}

func filterPanic(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	panic("panic filter was called")
}

func filterDefault(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	p := params.Expect(1, []*exec.KwArg{{
		Name:    "boolean",
		Default: false,
	}})
	if p.IsError() || !p.GetKwarg("boolean", false).IsBool() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'default'"))
	}
	if in.IsError() || in.IsNil() {
		return p.First()
	}
	if p.GetKwarg("boolean", false).Bool() && !in.Bool() {
		return p.First()
	}
	return in
}

func filterFromYAML(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'fromyaml'"))
	}
	if !in.IsString() || in.String() == "" {
		return exec.AsValue(errors.New("Filter 'fromyaml' was passed an empty or non-string type"))
	}
	object := new(interface{})
	if err := yaml.Unmarshal([]byte(in.String()), object); err != nil {
		return exec.AsValue(fmt.Errorf("failed to unmarshal %s: %s", in.String(), err))
	}
	return exec.AsValue(*object)
}

func filterToYAML(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	const defaultIndent = 2

	p := params.Expect(0, []*exec.KwArg{{Name: "indent", Default: defaultIndent}})
	if p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'toyaml'"))
	}

	indent, ok := p.KwArgs["indent"]
	if !ok || indent.IsNil() {
		indent = exec.AsValue(defaultIndent)
	}

	if !indent.IsInteger() {
		return exec.AsValue(errors.Errorf("Expected an integer for 'indent', got %s", indent.String()))
	}
	if in.IsNil() {
		return exec.AsValue(errors.New("Filter 'toyaml' was called with a nil object"))
	}
	output := bytes.NewBuffer(nil)
	encoder := yaml.NewEncoder(output)
	encoder.SetIndent(indent.Integer())

	// Monkey patching because the pipeline input parser is broken when the input is a list
	if in.IsList() {
		inCast := make([]interface{}, in.Len())
		for index := range inCast {
			item := exec.ToValue(in.Index(index).Val)
			inCast[index] = item.Val.Interface()
		}
		in = exec.AsValue(inCast)
	}

	castedType := in.ToGoSimpleType(true)
	if err, ok := castedType.(error); ok {
		return exec.AsValue(err)
	}

	if err := encoder.Encode(castedType); err != nil {
		return exec.AsValue(fmt.Errorf("unable to marshal to yaml: %s: %s", in.String(), err))
	}

	return exec.AsValue(output.String())
}

func filterSelectAttr(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	var test func(*exec.Value) *exec.Value
	if len(params.Args) < 1 {
		return exec.AsValue(errors.New("Wrong signature for 'selectattr', expect at least an attribute name as argument"))
	}
	attribute := params.First().String()
	if len(params.Args) == 1 {
		// Reject truthy value
		test = func(in *exec.Value) *exec.Value {
			attr, found := in.Get(attribute)
			if !found {
				return exec.AsValue(errors.Errorf(`%s has no attribute '%s'`, in.String(), attribute))
			}
			return attr
		}
	} else {
		name := params.Args[1].String()
		testParams := &exec.VarArgs{
			Args:   params.Args[2:],
			KwArgs: params.KwArgs,
		}
		test = func(in *exec.Value) *exec.Value {
			attr, found := in.Get(attribute)
			if !found {
				return exec.AsValue(errors.Errorf(`%s has no attribute '%s'`, in.String(), attribute))
			}
			out := e.ExecuteTestByName(name, attr, testParams)
			return out
		}
	}

	out := make([]interface{}, 0)
	var err *exec.Value

	in.Iterate(func(idx, count int, key, value *exec.Value) bool {
		result := test(key)
		if result.IsError() {
			err = result
			return false
		}
		if result.IsTrue() {
			out = append(out, key.Interface())
		}
		return true
	}, func() {})

	if err != nil {
		return err
	}
	return exec.AsValue(out)
}

func filterToToml(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	// Done not mess around with trying to marshall error pipelines
	if in.IsError() {
		return in
	}

	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'totoml'"))
	}

	casted := in.ToGoSimpleType(false)
	if err, ok := casted.(error); ok {
		return exec.AsValue(err)
	}

	out, err := toml.Marshal(casted)
	if err != nil {
		return exec.AsValue(errors.Wrap(err, "Unable to marhsall to toml"))
	}

	return exec.AsSafeValue(string(out))
}

func filterFromTOML(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	if p := params.ExpectNothing(); p.IsError() {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'fromtoml'"))
	}
	if !in.IsString() || in.String() == "" {
		return exec.AsValue(errors.New("Filter 'fromtoml' was passed an empty or non-string type"))
	}
	object := new(interface{})
	if err := toml.Unmarshal([]byte(in.String()), object); err != nil {
		return exec.AsValue(fmt.Errorf("failed to unmarshal from toml %s: %s", in.String(), err))
	}
	return exec.AsValue(*object)
}

func filterMatch(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
	if in.IsError() {
		return in
	}
	p := params.ExpectArgs(1)
	if p.IsError() || len(p.Args) != 1 {
		return exec.AsValue(errors.Wrap(p, "Wrong signature for 'match'"))
	}
	if !in.IsString() {
		return exec.AsValue(errors.New("Filter 'match' was passed a non-string type input"))
	}
	if !p.Args[0].IsString() {
		return exec.AsValue(errors.New("Filter 'match' was passed a non-string type argument"))
	}
	expression := p.Args[0].String()
	matcher, err := regexp.Compile(expression)
	if err != nil {
		return exec.AsValue(fmt.Errorf("failed to compile: %s: %s", expression, err))
	}

	return exec.AsValue(matcher.MatchString(in.String()))
}
