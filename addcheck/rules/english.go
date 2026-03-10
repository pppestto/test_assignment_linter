package rules

func CheckEnglish(msg string) *Violation {
	if msg == "" {
		return nil
	}

	for _, r := range msg {
		if !isASCII(r) {
			return NewViolation("English", "Message contains non-English characters")
		}
	}

	return nil
}

func isASCII(r rune) bool {
	return r <= 0x7F
}
