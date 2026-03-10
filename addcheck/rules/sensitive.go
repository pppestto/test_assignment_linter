package rules

import (
	"regexp"
	"strings"
)

// CheckSensitiveConcat — литерал содержит ключевое слово (подстрока), рядом через + идёт не-литерал:
// это риск утечки (log.Info("token: " + token)). Статическая фраза "token validated" без + не проверяем здесь.
func ContainsSensitiveSubstring(msg string, keywords []string) bool {
	if msg == "" {
		return false
	}
	lower := strings.ToLower(msg)
	for _, kw := range keywords {
		kw = strings.TrimSpace(strings.ToLower(kw))
		if kw != "" && strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func CheckSensitivePatterns(msg string, patterns []*regexp.Regexp) *Violation {
	if msg == "" || len(patterns) == 0 {
		return nil
	}
	for _, re := range patterns {
		if re != nil && re.MatchString(msg) {
			return NewViolation("sensitive", "message matches sensitive pattern: "+re.String())
		}
	}
	return nil
}
