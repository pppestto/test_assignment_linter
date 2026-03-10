package rules

type Violation struct {
	Rule    string
	Message string
}

func NewViolation(rule, message string) *Violation {
	return &Violation{
		Rule:    rule,
		Message: message,
	}
}
