package pystring

import "testing"

func intP(i int) *int {
	return &i
}

func TestIndex(t *testing.T) {
	tests := []struct {
		s     string
		start *int
		end   *int

		want string
	}{
		{
			s:     "hello",
			start: nil,
			end:   nil,
			want:  "hello",
		},
		{
			s:     "",
			start: nil,
			end:   nil,
			want:  "",
		},
		{
			s:     "hello",
			start: intP(0),
			end:   intP(-100),
			want:  "",
		},
		{
			s:     "hello",
			start: intP(-100),
			end:   nil,
			want:  "hello",
		},
		{
			s:     "hello",
			start: intP(2),
			end:   nil,
			want:  "llo",
		},
		{
			s:     "hello",
			start: intP(-1),
			end:   nil,
			want:  "o",
		},
		{
			s:     "hello",
			start: intP(-5),
			end:   nil,
			want:  "hello",
		},
		{
			s:     "hello",
			start: intP(-2),
			end:   nil,
			want:  "lo",
		},
		{
			s:     "hello",
			start: intP(-2),
			end:   intP(100),
			want:  "lo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, _ := PyString(tt.s).Idx(tt.start, tt.end); string(got) != tt.want {
				switch {
				case tt.start == nil && tt.end == nil:
					t.Errorf("%q.Index(nil, nil) = %v, want %v", tt.s, got, tt.want)
				case tt.start == nil:
					t.Errorf("%q.Index(nil, %d) = %v, want %v", tt.s, *tt.end, got, tt.want)
				case tt.end == nil:
					t.Errorf("%q.Index(%d, nil) = %v, want %v", tt.s, *tt.start, got, tt.want)
				}
			}
		})
	}
}
