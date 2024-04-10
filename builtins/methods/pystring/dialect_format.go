package pystring

import (
	"fmt"
	"strconv"
	"strings"
)

// Any value which implements formatter will itself decide how a formatting
// string should be interpreted. If not the default formatter is used.
type Formatter interface {
	Format(format string) (string, error)
}

func (d Dialect) Format(s string, vargs []any, kwarg map[string]any) (string, error) {

	var res strings.Builder

	scan := NewScanner(s, d)
	for {
		t, s, err := scan.Next()
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrInternal, err.Error())
		}

		switch {
		case t == Characters:
			res.WriteString(s)

		case t == ReplacementBlock:
			err := d.parseReplacementField(s, &res, vargs, kwarg)
			if err != nil {
				return "", err
			}

		case t == EOF:
			return res.String(), nil

		default:
			return "", fmt.Errorf("%w: Unknown token type %s", ErrInternal, t.String())

		}
	}
}

// parseReplacementField numbers automatic replacement blocks "{}" and normalizes
// all accessor patters (e.g. m['sub'] -> m.sub). Then splits up value part from
// formatting specifications; extracts the value from vargs and kwargs and applies
// the specified formatting.
func (d Dialect) parseReplacementField(s string, res *strings.Builder, vargs []any, kwarg map[string]any) error {
	// Strip initial and final braces
	if s[0] != '{' && s[len(s)-1] != '}' {
		return fmt.Errorf("%w: format didn't have format directives in: %s", ErrInternal, s)
	}
	s = s[1 : len(s)-1]

	// Extract string value of and format directives.
	formatDelimiterIdx := strings.Index(string(s), ":")
	var value, formatSpecifier string
	if formatDelimiterIdx == -1 {
		value = s
	} else {
		value = s[:formatDelimiterIdx]
		formatSpecifier = s[formatDelimiterIdx+1:]
	}

	// formatSpecifier may itself be based on replacement fields
	format, err := d.Format(formatSpecifier, vargs, kwarg)
	if err != nil {
		return fmt.Errorf("%w: failed subformat on value '%s'", err, format)
	}
	_ = format

	// Extract vargs
	if v, err := strconv.Atoi(string(value)); err == nil && v >= 0 {
		if v >= len(vargs) || v < 0 {
			return fmt.Errorf("%w: Replacement index %d out of range for positional args tuple", ErrIndex, v)
		}
		return d.formatReplacementFieldValue(res, vargs[v], string(format))
	}

	// python separates attributes and getAttr but for us they will be one and the same.
	// translate [] accessors to dot notation; e.g. m['sub'] -> m.sub;
	// TODO:needs rewrite.
	//mangledKeys := strings.Replace(string(value), "[", ".", -1)
	//mangledKeys = strings.Replace(mangledKeys, "]", ".", -1)
	//mangledKeys = strings.Replace(mangledKeys, "..", ".", -1)
	//mangledKeys = strings.Trim(mangledKeys, ". ")

	anyVal, err := getNestedKwArgs(simpleJSONPathSplit(string(value)), KwArgs(kwarg))
	if err != nil {
		return err
	}

	return d.formatReplacementFieldValue(res, anyVal, string(format))
}

func (d Dialect) formatReplacementFieldValue(res *strings.Builder, value any, formatStr string) error {
	if dataTypeFormatter, ok := value.(Formatter); ok {
		s, err := dataTypeFormatter.Format(formatStr)
		if err != nil {
			return err
		}
		_, err = res.WriteString(s)
		return err
	}

	dialectFormatter, err := d.NewFormatterSpecFromStr(formatStr)
	if err != nil {
		return err
	}

	formattedString, err := dialectFormatter.Format(value)
	if err != nil {
		return err
	}

	_, err = res.WriteString(formattedString)
	return err
}

// TODO: make more robust
func simpleJSONPathSplit(path string) []string {
	var parts []string
	start := 0

	for idx, char := range path {
		if path[start:idx] == "" {
			continue
		}

		withinBrackets := false

		switch char {
		case '"':
		case '.':
			part := strings.Trim(path[start:idx], "'\"[].")
			if withinBrackets && part != "" {
				parts = append(parts, part)
				start = idx
			}
		case '[', ']':
			withinBrackets = char == '['
			part := strings.Trim(path[start:idx], "'\"[].")
			if part != "" {
				parts = append(parts, part)
				start = idx
			}
		}
	}

	part := strings.Trim(path[start:], "'\"[].")
	if part != "" {
		parts = append(parts, part)
	}
	return parts
}
