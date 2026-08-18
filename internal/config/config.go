// markings:managed
//
// File: config.go
// Copyright (c) 2026 Michael Harris
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
//
// markings:managed

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Templates     map[string]string       `yaml:"templates"`
	Data          map[string]any          `yaml:"data"`
	CommentStyles map[string]CommentStyle `yaml:"comment_styles"`
	Rules         []Rule                  `yaml:"rules"`
	Marker        string                  `yaml:"marker"`
	Exclude       []string                `yaml:"exclude"`
	Include       []string                `yaml:"include"`
}

type CommentStyle struct {
	Prefix string      `yaml:"prefix"`
	Block  *BlockStyle `yaml:"block"`
}

type BlockStyle struct {
	Start  string `yaml:"start"`
	Middle string `yaml:"middle"`
	End    string `yaml:"end"`
}

type MarkingConfig struct {
	Template      string `yaml:"template"`
	NewlineBefore bool   `yaml:"newline_before"`
	NewlineAfter  bool   `yaml:"newline_after"`
}

type Rule struct {
	ID           string         `yaml:"id"`
	Match        string         `yaml:"match"`
	Header       *MarkingConfig `yaml:"header"`
	Footer       *MarkingConfig `yaml:"footer"`
	CommentStyle string         `yaml:"comment_style"`
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

const DefaultMarker = "markings:" + "managed"

func (c *Config) GetMarker() string {
	if c.Marker == "" {
		return DefaultMarker
	}
	return c.Marker
}

// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
