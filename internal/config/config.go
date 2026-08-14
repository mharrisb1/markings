package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Templates     map[string]string       `yaml:"templates"`
	Data          map[string]any          `yaml:"data"`
	CommentStyles map[string]CommentStyle `yaml:"comments"`
	Rules         []Rule                  `yaml:"rules"`
}

type CommentStyle struct {
	// For line comments (e.g., "// ")
	Prefix string `yaml:"prefix"`

	// For block comments
	Block *BlockStyle `yaml:"block"`
}

type BlockStyle struct {
	Start  string `yaml:"start"`  // e.g., "/*\n"
	Middle string `yaml:"middle"` // e.g., " * "
	End    string `yaml:"end"`    // e.g., "\n */"
}

type Rule struct {
	ID             string `yaml:"id"`
	Match          string `yaml:"match"`
	HeaderTemplate string `yaml:"header_template"`
	FooterTemplate string `yaml:"footer_template"`
	CommentStyle   string `yaml:"comment_style"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
