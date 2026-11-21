package gonja

import "testing"

func TestFloorDivisionAndModuloFollowPython(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"-7 // 2", "-4"},
		{"7 // -2", "-4"},
		{"-7 // -2", "3"},
		{"-7 % 3", "2"},
		{"7 % -3", "-2"},
		{"-7 % -2", "-1"},
		{"10.0 // 3", "3.0"},
		{"-7.0 % 3", "2.0"},
		{"7 % 4.0", "3.0"},
		// Integers beyond float64's exact range must stay exact.
		{"9007199254740993 // 2", "4503599627370496"},
		{"9007199254740993 % 2", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			tmpl, err := FromString("{{ " + tt.expr + " }}")
			if err != nil {
				t.Fatalf("parse %q: %v", tt.expr, err)
			}
			got, err := tmpl.ExecuteToString(DefaultContext)
			if err != nil {
				t.Fatalf("execute %q: %v", tt.expr, err)
			}
			if got != tt.want {
				t.Fatalf("%s = %q, want %q", tt.expr, got, tt.want)
			}
		})
	}
}
