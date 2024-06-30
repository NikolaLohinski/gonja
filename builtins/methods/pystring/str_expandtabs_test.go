package pystring

import "testing"

func TestExpandTabs(t *testing.T) {
	tests := []struct {
		s       string
		tabSize *int
		want    string
	}{
		{
			s:       "01\t012\t0123\t01234",
			tabSize: nil,
			want:    "01      012     0123    01234",
		},
		{
			s:       "01\t012\t0123\t01234",
			tabSize: intP(4),
			want:    "01  012 0123    01234",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := ExpandTabs(tt.s, tt.tabSize); got != tt.want {
				t.Errorf("ExpandTabs() = %v, want %v", got, tt.want)
			}
		})
	}
}
