package pystring

import "testing"

func TestJoin(t *testing.T) {
	tests := []struct {
		s    string
		it   []string
		want string
	}{
		{",", []string{"a", "b", "c"}, "a,b,c"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := JoinString(tt.s, tt.it); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}
