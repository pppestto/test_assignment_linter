package addcheck

import (
	"encoding/json"
	"os"
)

type Config struct {
	Rules *struct {
		Lowercase bool `json:"lowercase"`
		English   bool `json:"english"`
		Emoji     bool `json:"emoji"`
		Sensitive bool `json:"sensitive"`
	} `json:"rules"`

	SensitiveKeywords []string `json:"sensitive_keywords"`

	SensitivePatterns []string `json:"sensitive_patterns"`
}

func defaultConfig() *Config {
	return &Config{
		SensitiveKeywords: []string{"password", "token", "secret", "apikey", "bearer"},
	}
}

func loadConfig(path string) (*Config, error) {
	if path == "" {
		return defaultConfig(), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if len(c.SensitiveKeywords) == 0 {
		c.SensitiveKeywords = defaultConfig().SensitiveKeywords
	}
	return &c, nil
}

func (c *Config) ruleEnabled(name string) bool {
	if c.Rules == nil {
		return true
	}
	switch name {
	case "lowercase":
		return c.Rules.Lowercase
	case "english":
		return c.Rules.English
	case "emoji":
		return c.Rules.Emoji
	case "sensitive":
		return c.Rules.Sensitive
	default:
		return true
	}
}
