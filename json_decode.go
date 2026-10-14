package gonja

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/nikolalohinski/gonja/v2/exec"
)

// DecodeOrdered decodes JSON into ordered Gonja dictionaries.
func DecodeOrdered(raw []byte) (any, error) {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	v, err := decodeValue(d)
	if err != nil {
		return nil, err
	}
	if _, err = d.Token(); err != io.EOF {
		return nil, fmt.Errorf("trailing JSON: %v", err)
	}
	return v, nil
}

func decodeValue(d *json.Decoder) (any, error) {
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	switch t {
	case json.Delim('{'):
		obj := exec.NewDict()
		indexes := make(map[string]int)
		for d.More() {
			k, err := d.Token()
			if err != nil {
				return nil, err
			}
			value, err := decodeValue(d)
			if err != nil {
				return nil, err
			}
			key := k.(string) // The decoder validates object keys are strings.
			if i, found := indexes[key]; found {
				obj.Pairs[i].Value = exec.AsValue(value)
			} else {
				indexes[key] = len(obj.Pairs)
				obj.Pairs = append(obj.Pairs, &exec.Pair{Key: exec.AsValue(key), Value: exec.AsValue(value)})
			}
		}
		_, err = d.Token()
		return obj, err
	case json.Delim('['):
		a := []any{}
		for d.More() {
			x, err := decodeValue(d)
			if err != nil {
				return nil, err
			}
			a = append(a, x)
		}
		_, err = d.Token()
		return a, err
	default:
		if n, ok := t.(json.Number); ok {
			if !strings.ContainsAny(string(n), ".eE") {
				return n.Int64()
			}
			return n.Float64()
		}
		return t, nil
	}
}
