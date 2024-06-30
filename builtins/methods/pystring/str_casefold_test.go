package pystring

import "testing"

func TestCaseFold(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		in  string
		out PyString
	}{
		{in: "", out: ""},
		{in: "grüßen", out: "grüssen"},
		{in: "ÄÖÜ", out: "äöü"},
		{in: "ǅABCDǄ", out: "ǆabcdǆ"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			pys := PyString(tt.in)
			if got := pys.Casefold(); got != tt.out {
				t.Errorf("Casefold() = %v, want %v", got, tt.out)
			}
		})
	}
}
