package pystring

import "testing"

func TestCenter(t *testing.T) {
	tests := []struct {
		in    string
		width int
		out   PyString
	}{
		{in: "", width: 15, out: "               "},
		{in: "grüßen", width: 15, out: "     grüßen    "},
		{in: "ÄÖÜ", width: 15, out: "      ÄÖÜ      "},
		{in: "ǅABCDǄ", width: 15, out: "     ǅABCDǄ    "},
		{in: "ǅABCDǄ", width: 14, out: "    ǅABCDǄ    "},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			pys := PyString(tt.in)
			if got := pys.Center(tt.width, 0); got != tt.out {
				t.Errorf("Center() = '%v', want '%v'", got, tt.out)
			}
		})
	}
}
