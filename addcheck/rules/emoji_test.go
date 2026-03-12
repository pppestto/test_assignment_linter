package rules

import "testing"

func TestCheckEmoji(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"empty", "", false},
		{"letters and spaces", "server started", false},
		{"digits", "code 200 ok", false},
		{"exclamation", "started!", true},
		{"ellipsis", "went wrong...", true},
		{"triple bang", "failed!!!", true},
		{"emoji", "hi\U0001F680", true},
		{"colon", "a: b", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := CheckEmoji(tt.msg)
			if tt.want && v == nil {
				t.Fatalf("expected violation for %q", tt.msg)
			}
			if !tt.want && v != nil {
				t.Fatalf("unexpected violation: %v", v)
			}
		})
	}
}

func TestCheckEmojiSymbolOther(t *testing.T) {
	if v := CheckEmoji("snow"); v != nil {
		t.Fatalf("unexpected: %v", v)
	}
}
