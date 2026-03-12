package rules

import "testing"

func TestCheckEnglish(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"empty", "", false},
		{"ascii only", "hello world 123 !", false},
		{"cyrillic", "привет", true},
		{"emoji only", "😀", true}, // > 127
		{"mixed ascii cyrillic", "hello мир", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := CheckEnglish(tt.msg)
			if tt.want && v == nil {
				t.Fatalf("expected violation for %q", tt.msg)
			}
			if !tt.want && v != nil {
				t.Fatalf("unexpected violation: %v", v)
			}
			if v != nil && v.Rule != "English" {
				t.Fatalf("Rule = %q", v.Rule)
			}
		})
	}
}
