package helpers

import "testing"

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "valid lowercase", in: "abcdefg", want: true},
		{name: "valid mixed", in: "aB3xY9z", want: true},
		{name: "too short", in: "abc", want: false},
		{name: "too long", in: "abcdefgh", want: false},
		{name: "empty", in: "", want: false},
		{name: "special chars", in: "abc-efg", want: false},
		{name: "space", in: "abc efg", want: false},
		{name: "unicode", in: "abcdefé", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid(tt.in); got != tt.want {
				t.Errorf("Valid(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNewCode(t *testing.T) {
	seen := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		code, err := NewCode()
		if err != nil {
			t.Fatalf("NewCode() error: %v", err)
		}
		if !Valid(code) {
			t.Fatalf("NewCode() produced invalid code %q", code)
		}
		if seen[code] {
			// 62^7 keyspace: a duplicate in 1000 draws means the generator is broken.
			t.Fatalf("NewCode() produced duplicate %q within 1000 draws", code)
		}
		seen[code] = true
	}
}
