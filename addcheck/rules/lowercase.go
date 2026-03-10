package rules

func CheckLowerCase(msg string) *Violation {
	if msg == "" {
		return nil
	}

	r := []rune(msg)[0]
	if r >= 'A' && r <= 'Z' {
		return NewViolation("LowerCase", "Message contains uppercase letters")
	}

	return nil
}
