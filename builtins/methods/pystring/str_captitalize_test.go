package pystring

import "testing"

func TestCapitalize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		in  string
		out PyString
	}{
		{in: "", out: ""},
		{in: "hello", out: "Hello"},
		{in: "HELLO", out: "Hello"},
		{in: "hELLO", out: "Hello"},
		{in: "hello world", out: "Hello world"},
		{in: "hello world", out: "Hello world"},
		{in: "ätö", out: "Ätö"},
		{in: "işğüı", out: "Işğüı"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			pys := PyString(tt.in)
			if got := pys.Capitalize(); got != tt.out {
				t.Errorf("Capitalize() = %v, want %v", got, tt.out)
			}
		})
	}
}
