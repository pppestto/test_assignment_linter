package addcheck

import "testing"

func TestLoadConfigDefault(t *testing.T) {
	c, err := loadConfig("")
	if err != nil {
		t.Fatal(err)
	}
	if !c.ruleEnabled("lowercase") {
		t.Fatal("expected all rules on by default")
	}
}

func TestLoadConfigRulesOff(t *testing.T) {
	c := &Config{}
	c.Rules = &struct {
		Lowercase bool `json:"lowercase"`
		English   bool `json:"english"`
		Emoji     bool `json:"emoji"`
		Sensitive bool `json:"sensitive"`
	}{Lowercase: true}
	if !c.ruleEnabled("lowercase") || c.ruleEnabled("english") {
		t.Fatalf("only lowercase should be on")
	}
}
