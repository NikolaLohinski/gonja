package exec

import (
	"strings"
	"testing"
)

func TestStringRepeatBoundsCheck(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		count   int
		wantErr bool
	}{
		{
			name:    "reasonable repeat",
			str:     "A",
			count:   1000,
			wantErr: false,
		},
		{
			name:    "at limit",
			str:     strings.Repeat("A", 1024),
			count:   100 * 1024,
			wantErr: false,
		},
		{
			name:    "exceeds limit",
			str:     "A",
			count:   101 * 1024 * 1024,
			wantErr: true,
		},
		{
			name:    "massive repeat causes OOM without check",
			str:     "A",
			count:   8000000000000,
			wantErr: true,
		},
		{
			name:    "negative count treated as zero",
			str:     "A",
			count:   -1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := AsValue(tt.str)
			right := AsValue(tt.count)

			count := right.Integer()
			if count < 0 {
				count = 0
			}
			resultLen := int64(len(left.String())) * int64(count)
			exceeds := resultLen > maxStringRepeatBytes

			if exceeds != tt.wantErr {
				t.Errorf("expected exceeds=%v, got %v (resultLen=%d)", tt.wantErr, exceeds, resultLen)
			}
		})
	}
}
