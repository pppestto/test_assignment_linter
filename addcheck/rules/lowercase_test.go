package rules

import "testing"

func TestCheckLowerCase(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"empty", "", false},
		{"lowercase ok", "hello world", false},
		{"starts with upper", "Hello", true},
		{"single upper", "A", true},
		{"number first", "123", false},
		{"space first", " hello", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := CheckLowerCase(tt.msg)
			if tt.want && v == nil {
				t.Fatalf("expected violation for %q", tt.msg)
			}
			if !tt.want && v != nil {
				t.Fatalf("unexpected violation: %v", v)
			}
			if v != nil && v.Rule != "LowerCase" {
				t.Fatalf("Rule = %q, want LowerCase", v.Rule)
			}
		})
	}
}
