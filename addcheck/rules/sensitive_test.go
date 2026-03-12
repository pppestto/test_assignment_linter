package rules

import (
	"regexp"
	"testing"
)

func TestContainsSensitiveSubstring(t *testing.T) {
	kw := []string{"password", "token", "secret"}
	if !ContainsSensitiveSubstring("token: ", kw) {
		t.Fatal("token prefix should match")
	}
	if !ContainsSensitiveSubstring("user password here", kw) {
		t.Fatal("password substring")
	}
	if ContainsSensitiveSubstring("all good", kw) {
		t.Fatal("no keyword")
	}
}

func TestCheckSensitivePatterns(t *testing.T) {
	re := regexp.MustCompile(`api[_-]?key`)
	if v := CheckSensitivePatterns("my api_key here", []*regexp.Regexp{re}); v == nil {
		t.Fatal("expected match")
	}
	if v := CheckSensitivePatterns("no match", []*regexp.Regexp{re}); v != nil {
		t.Fatal("unexpected")
	}
}
