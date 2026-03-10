package rules

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func FixLowercaseFirst(msg string) string {
	if msg == "" {
		return msg
	}
	r, size := utf8.DecodeRuneInString(msg)
	if r == utf8.RuneError {
		return msg
	}
	if r >= 'A' && r <= 'Z' {
		return string(unicode.ToLower(r)) + msg[size:]
	}
	return msg
}

func FixEmojiTZ(msg string) string {
	var b strings.Builder
	for _, r := range msg {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
