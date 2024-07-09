package pystring

import (
	"testing"
	"unicode"
)

func FuzzFormatSpec(t *testing.F) {
	t.Add('<', rune(0), rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add('>', rune(0), rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add('^', rune(0), rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add('=', rune(0), rune(0), false, false, uint(0), uint(0), rune(0))

	t.Add(rune(0), ' ', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '>', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '.', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), 'g', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '0', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), 'O', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '#', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '<', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '^', rune(0), false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), '=', rune(0), false, false, uint(0), uint(0), rune(0))

	t.Add(rune(0), rune(0), '+', false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), '-', false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), ' ', false, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), true, false, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, true, uint(0), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(1), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(2), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(3), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(4), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(5), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(6), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(7), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(8), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(9), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(10), uint(0), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(1), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(2), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(3), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(4), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(5), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(6), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(7), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(8), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(9), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(10), rune(0))
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'b')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'c')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'd')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'o')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'x')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'X')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'e')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'E')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'f')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'F')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'g')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), 'G')
	t.Add(rune(0), rune(0), rune(0), false, false, uint(0), uint(0), '%')

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
	) {

		// Some fuzz responses aren't really valid test cases.
		if !unicode.IsPrint(Fill) {
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
		if Fill != 0 && (Fill != '<' && Fill != '>' && Fill != '^' && Fill != '=') {
			return
		}
		if Align == 0 && (Fill == '<' || Fill == '>' || Fill == '^' || Fill == '=') {
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
		orig := spec.String()

		specRecovered, err := NewFormatterSpecFromStr(orig)
		if err != nil {
			t.Errorf("Expected to parse '%s' but got error: %v on data %#v", orig, err, spec)
		}
		if orig != specRecovered.String() {
			t.Errorf("Expected '%s' = '%s' on data %#v", specRecovered.String(), orig, spec)
		}
	})
}
