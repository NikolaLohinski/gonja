//go:build integration
// +build integration

package pystring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestValidExpressionsMatchesPython_311(t *testing.T) {
	pythonCmd := "python3.11"
	dialect := DialectPython3_11

	expressions := getValidExpressions(dialect)
	t.Logf("found %d valid expressions; evaluating against dialect %#v", len(expressions), dialect)
	for _, expr := range expressions {
		expr.dialect = dialect
		if expr.ExpectFloatType() {
			v := 1.123456789
			s, goErr := expr.Format(v)
			ps, _, pyErr := getPythonRes(pythonCmd, fmt.Sprintf("{0:%s}", expr.String()), []any{v})
			if goErr != nil && pyErr == nil || goErr == nil && pyErr != nil {
				t.Errorf("Error miss match between golang and python of format '%s': goErr: %v, pyErr: %v", expr.String(), goErr, pyErr)
			}
			if goErr != nil || pyErr != nil {
				continue
			}
			if s != ps {
				t.Errorf("Different response from Go/Python. Go:'%s' != Py:'%s' for expression: %s (%#v)", s, ps, expr.String(), v)
			}

			//t.Logf("Go: '%s' == Py:'%s' for expression: %s ", s, ps, expr.String())

		} else if expr.ExpectIntType() {
			v := 16789
			s, goErr := expr.Format(v)
			ps, _, pyErr := getPythonRes(pythonCmd, fmt.Sprintf("{0:%s}", expr.String()), []any{v})
			if goErr != nil && pyErr == nil || goErr == nil && pyErr != nil {
				t.Errorf("Error miss match between golang and python of format '%s': goErr: %v, pyErr: %v", expr.String(), goErr, pyErr)
			}
			if goErr != nil || pyErr != nil {
				continue
			}
			if s != ps {
				t.Errorf("Different response from Go/Python. Go:'%s' != Py:'%s' for expression: %s (%#v)", s, ps, expr.String(), v)
			}

			//t.Logf("Go: '%s' == Py:'%s' for expression: %s ", s, ps, expr.String())

		} else if expr.ExpectNumericType() {
			v := 16789
			s, goErr := expr.Format(v)
			ps, _, pyErr := getPythonRes(pythonCmd, fmt.Sprintf("{0:%s}", expr.String()), []any{v})
			if goErr != nil && pyErr == nil || goErr == nil && pyErr != nil {
				t.Errorf("Error miss match between golang and python of format '%s': goErr: %v, pyErr: %v", expr.String(), goErr, pyErr)
			}
			if goErr != nil || pyErr != nil {
				continue
			}
			if s != ps {
				t.Errorf("Different response from Go/Python. Go:'%s' != Py:'%s' for expression: %s (%#v)", s, ps, expr.String(), v)
			}

			//t.Logf("Go: '%s' == Py:'%s' for expression: %s ", s, ps, expr.String())

		} else {
			v := "foobar"
			s, goErr := expr.Format(v)
			ps, _, pyErr := getPythonRes(pythonCmd, fmt.Sprintf("{0:%s}", expr.String()), []any{v})
			if goErr != nil && pyErr == nil || goErr == nil && pyErr != nil {
				t.Fatalf("Error miss match between golang and python of format '%s': goErr: %v, pyErr: %v", expr.String(), goErr, pyErr)
			}
			if goErr != nil || pyErr != nil {
				continue
			}
			if s != ps {
				t.Fatalf("Different response from Go/Python. Go:'%s' != Py:'%s' for expression: %s (%#v)", s, ps, expr.String(), v)
			}

			//t.Logf("Go: '%s' == Py:'%s' for expression: %s ", s, ps, expr.String())

		}
	}
}

func FuzzIntTestingWithPython_311(t *testing.F) {
	t.Add('<', rune(0), rune(0), false, false, uint(0), uint(0), rune(0), 1.0)
	t.Add('>', rune(0), rune(0), false, false, uint(0), uint(0), rune(0), 1.0)
	t.Add('^', rune(0), rune(0), false, false, uint(0), uint(0), rune(0), 1.0)
	t.Add('=', rune(0), rune(0), false, false, uint(0), uint(0), rune(0), 1.0)

	t.Add(rune(0), ' ', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), '>', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), '.', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), 'g', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), '0', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), 'O', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), '#', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), '<', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), '^', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)
	t.Add(rune(0), '=', rune(0), false, false, uint(0), uint(0), rune(0), 10.0)

	t.Add(rune(0), rune(0), '+', false, false, uint(0), uint(0), rune(0), 2.0)
	t.Add(rune(0), rune(0), '-', false, false, uint(0), uint(0), rune(0), 2.0)
	t.Add(rune(0), rune(0), ' ', false, false, uint(0), uint(0), rune(0), 2.0)
	t.Add(rune(0), rune(0), rune(0), true, false, uint(0), uint(0), rune(0), 3.0)
	t.Add(rune(0), rune(0), rune(0), false, true, uint(0), uint(0), rune(0), 3.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(1), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(2), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(3), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(4), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(5), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(6), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(7), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(8), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(9), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(10), uint(0), rune(0), 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(1), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(2), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(3), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(4), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(5), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(6), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(7), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(8), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(9), rune(0), 20.01123)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(10), rune(0), -123.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'b', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'c', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'd', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'o', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'x', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'X', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'e', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'E', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'f', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'F', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'g', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'G', 1.0)
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), '%', 1.0)

	t.Fuzz(func(
		t *testing.T,
		Fill rune,
		Align rune,
		Sign rune,
		Alternate bool,
		ZeroPadding bool,
		MinWidth uint,
		Precision uint,
		Type rune,
		rnd float64,
	) {

		// Let's skip utf fill characters for now. Not sure how the process execution
		// to python cmd works with utf characters.
		if Fill < 33 || Fill > 126 {
			return
		}
		if Fill != 0 && Align != '<' && Align != '>' && Align != '^' && Align != '=' && Align != 0 {
			return
		}
		if Sign != 0 && Sign != '+' && Sign != '-' && Sign != ' ' {
			return
		}
		if Type != 0 && Type != 'b' && Type != 'c' && Type != 'd' && Type != 'o' && Type != 'x' && Type != 'X' && Type != 'e' && Type != 'E' && Type != 'f' && Type != 'F' && Type != 'g' && Type != 'G' && Type != '%' {
			return
		}
		if Align == 0 && (Fill == '<' || Fill == '>' || Fill == '^' || Fill == '=') {
			return
		}
		if Type == 'c' && (rnd < 33 || rnd > 126) {
			// don't test unprintable characters.
			return
		}

		// Golang behaves correctly here but python is a bit random in how many
		// digits it returns; sometimes trailing zeros are wipes other times they
		// are not.
		// e.g.
		//>>> '{0:.17g}'.format(1.0001)
		//'1.0001'
		//>>> '{0:.18g}'.format(1.0001)
		//'1.00009999999999999'

		if Precision > 5 {
			return
		}

		spec := FormatSpec{
			Fill:        Fill,
			Align:       Align,
			Sign:        Sign,
			Alternate:   Alternate,
			ZeroPadding: ZeroPadding,
			MinWidth:    MinWidth,
			Precision:   Precision,
			Type:        Type,
		}

		if spec.Validate() != nil {
			return
		}

		template := "{0:" + spec.String() + "}"

		var val []any

		switch {
		case spec.ExpectFloatType() || spec.Precision > 0:
			val = []any{rnd + 0.0001}
		case spec.ExpectIntType() || spec.ExpectNumericType():
			val = []any{int(rnd)}
		default:
			val = []any{fmt.Sprintf("%v", rnd)}
		}

		var goRes string
		spec, goErr := NewFormatterSpecFromStr(spec.String())
		if goErr == nil {
			goRes, goErr = spec.Format(val[0])
		}
		pythonRes, pythonCmd, pyErr := getPythonRes("python3.11", template, val)

		if goErr != nil && pyErr == nil {
			t.Errorf("PythonErr: nil != GoErr: '%s'; Spec: %s", goErr.Error(), template)
		}
		if goErr == nil && pyErr != nil {
			t.Errorf("PythonErr: '%s' != GoErr: nil; Spec: %s", pyErr.Error(), template)
		}

		if goRes != pythonRes {
			// python is random with its precision, so we can't compare the strings directly.
			if strings.HasPrefix(goRes, pythonRes) {
				return
			}
			// This also the case with padding
			// '{0:,=39.3G}'.format(10.0001)
			// python ',,,,,10'
			// golang ',,,10.0'

			t.Errorf("Python: '%s' != Go: '%s' \nTemplate:%#v => %s.format(%v); used Python CMD: %s", pythonRes, goRes, spec, template, val, pythonCmd)
		}
	})
}

// Test utils

func getPythonRes(pythonCmd string, template string, args []any) (string, string, error) {
	b, err := json.Marshal(args)
	if err != nil {
		panic(err)
	}

	template = strings.ReplaceAll(template, "'", "\\'")

	cmdString := fmt.Sprintf("print( '%s'.format(%s) ); exit(0)", template, string(b[1:len(b)-1]))
	cmd := exec.Command(pythonCmd, "-c", cmdString)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		return "", cmdString, fmt.Errorf("%w: after cmd:python3 -c \"%s\"   stdErr: %s", err, cmdString, stderr.String())
	}
	cmd.Process.Kill()

	return strings.Trim(string(output), "\n\r"), cmdString, nil
}
