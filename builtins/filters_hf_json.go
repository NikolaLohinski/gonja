package builtins

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/nikolalohinski/gonja/v2/exec"
)

// filterHFToJSON implements Hugging Face's no-argument tojson filter.
func filterHFToJSON(_ *exec.Evaluator, in *exec.Value, args *exec.VarArgs) *exec.Value {
	if err := args.Take(); err != nil {
		return exec.AsValue(err)
	}
	out, err := encodeHF(in)
	if err != nil {
		return exec.AsValue(err)
	}
	return exec.AsSafeValue(out)
}

// quoteHF preserves Unicode (including U+2028/U+2029) and HTML characters.
func quoteHF(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"', '\\':
			b.WriteByte('\\')
			b.WriteRune(r)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 32 {
				fmt.Fprintf(&b, `\u%04x`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
	return b.String()
}

func encodeHF(v *exec.Value) (string, error) {
	switch {
	case v.IsNil():
		return "null", nil
	case v.IsString():
		return quoteHF(v.String()), nil
	case v.IsBool():
		return strconv.FormatBool(v.Bool()), nil
	case v.IsInteger():
		return strconv.Itoa(v.Integer()), nil
	case v.IsFloat():
		f := v.Float()
		switch {
		case math.IsNaN(f):
			return "NaN", nil
		case math.IsInf(f, 1):
			return "Infinity", nil
		case math.IsInf(f, -1):
			return "-Infinity", nil
		}
		if f != 0 && (math.Abs(f) < 1e-4 || math.Abs(f) >= 1e16) {
			return strconv.FormatFloat(f, 'e', -1, 64), nil
		}
		s := strconv.FormatFloat(f, 'f', -1, 64)
		if !strings.Contains(s, ".") {
			s += ".0"
		}
		return s, nil
	}
	switch v.Interface().(type) {
	case exec.Dict, *exec.Dict:
		pairs := v.Items()
		parts := make([]string, 0, len(pairs))
		for _, p := range pairs {
			if !p.Key.IsString() {
				return "", fmt.Errorf("JSON object keys must be strings")
			}
			s, err := encodeHF(p.Value)
			if err != nil {
				return "", err
			}
			parts = append(parts, quoteHF(p.Key.String())+": "+s)
		}
		return "{" + strings.Join(parts, ", ") + "}", nil
	}
	if v.IsList() {
		parts := make([]string, v.Len())
		for i := range parts {
			s, err := encodeHF(v.Index(i))
			if err != nil {
				return "", err
			}
			parts[i] = s
		}
		return "[" + strings.Join(parts, ", ") + "]", nil
	}
	return "", fmt.Errorf("unsupported JSON type %T; objects must retain their order", v.Interface())
}
