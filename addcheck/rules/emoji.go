package rules

import "unicode"

func CheckEmoji(msg string) *Violation {
	if msg == "" {
		return nil
	}
	for _, r := range msg {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}
		return NewViolation("Emoji", "Message must contain only letters, digits and spaces")
	}
	return nil
}
