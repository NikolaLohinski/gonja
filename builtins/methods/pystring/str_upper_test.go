package pystring

import "testing"

func TestUpper(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{"Hello world", "HELLO WORLD"},
		{"they're bill's friends from the UK", "THEY'RE BILL'S FRIENDS FROM THE UK"},
		{"unicode is 😊", "UNICODE IS 😊"},
		{"こんにちは、世界", "こんにちは、世界"},
		{"Привет, мир", "ПРИВЕТ, МИР"},
		{"مرحبا العالم", "مرحبا العالم"},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			result := PyString(test.s).Upper()
			if string(result) != test.expected {
				t.Errorf("Expected %s but got %s", test.expected, result)
			}
		})
	}
}
