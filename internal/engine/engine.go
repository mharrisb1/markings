// markings:managed
//
// Copyright (c) 2026 Michael Harris
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
//
// markings:managed

package engine

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"markings/internal/comments"
	"markings/internal/config"

	"github.com/bmatcuk/doublestar/v4"
)

type Engine struct {
	config  *config.Config
	headers map[string]string
	footers map[string]string
}

// New creates a new Engine and pre-compiles all templates into their final formatted strings.
func New(cfg *config.Config) (*Engine, error) {
	e := &Engine{
		config:  cfg,
		headers: make(map[string]string),
		footers: make(map[string]string),
	}

	for name, tplStr := range cfg.Templates {
		t, err := template.New(name).Parse(tplStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %q: %w", name, err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, cfg.Data); err != nil {
			return nil, fmt.Errorf("error executing template %q: %w", name, err)
		}

		// The raw rendered template string
		rawText := buf.String()

		// Note: We'll format the raw text with the comment style when we process a specific rule,
		// because different rules might use the same template but different comment styles (e.g., // vs #).
		// For now, we just store the raw rendered text.
		e.headers[name] = rawText
		e.footers[name] = rawText
	}

	return e, nil
}

// matchRule finds the first rule that matches the given file path.
func (e *Engine) matchRule(path string) (*config.Rule, error) {
	for _, rule := range e.config.Rules {
		matched, err := doublestar.Match(rule.Match, path)
		if err != nil {
			return nil, err
		}
		if matched {
			return &rule, nil
		}
	}
	return nil, nil // No match
}

// ProcessFile checks a file against the rules and fixes it if requested.
// Returns true if the file has the correct markings (or was successfully fixed).
func (e *Engine) ProcessFile(path string, fix bool) (bool, error) {
	rule, err := e.matchRule(path)
	if err != nil || rule == nil {
		return true, err // No rule applies, so it's technically valid
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	content := string(contentBytes)

	style, ok := e.config.CommentStyles[rule.CommentStyle]
	if !ok {
		return false, fmt.Errorf("comment style %q not found", rule.CommentStyle)
	}

	expectedHeader := ""
	if rule.Header != nil && rule.Header.Template != "" {
		if raw, ok := e.headers[rule.Header.Template]; ok {
			expectedHeader = comments.Format(raw, style)
			if rule.Header.NewlineBefore {
				expectedHeader = "\n" + expectedHeader
			}
			if rule.Header.NewlineAfter {
				expectedHeader += "\n"
			}
		}
	}

	expectedFooter := ""
	if rule.Footer != nil && rule.Footer.Template != "" {
		if raw, ok := e.footers[rule.Footer.Template]; ok {
			expectedFooter = comments.Format(raw, style)
			if rule.Footer.NewlineBefore {
				expectedFooter = "\n" + expectedFooter
			}
			if rule.Footer.NewlineAfter {
				expectedFooter += "\n"
			}
		}
	}

	// Simple check: does the file start and end with the exact expected blocks?
	isValid := true
	if expectedHeader != "" && !strings.HasPrefix(content, expectedHeader) {
		isValid = false
	}
	if expectedFooter != "" && !strings.HasSuffix(strings.TrimSpace(content), strings.TrimSpace(expectedFooter)) {
		isValid = false
	}

	if isValid || !fix {
		return isValid, nil
	}

	// Fix logic: remove old markers and inject new ones
	newContent := removeOldMarkings(content, style)

	if expectedHeader != "" {
		newContent = expectedHeader + newContent
	}
	if expectedFooter != "" {
		// Ensure there's a newline before appending the footer
		if !strings.HasSuffix(newContent, "\n") {
			newContent += "\n"
		}
		newContent += expectedFooter
	}

	err = os.WriteFile(path, []byte(newContent), 0644)
	return err == nil, err
}

type parseState int

const (
	stateNormal parseState = iota
	stateInsideMarking
	stateExpectingBlockEnd
)

// State machine for stripping out existing marking blocks
func removeOldMarkings(content string, style config.CommentStyle) string {
	lines := strings.Split(content, "\n")
	var result []string
	state := stateNormal

	for _, line := range lines {
		switch state {
		case stateNormal:
			if strings.Contains(line, comments.Marker) {
				state = stateInsideMarking
				if style.Block != nil && len(result) > 0 {
					if strings.HasPrefix(strings.TrimSpace(result[len(result)-1]), strings.TrimSpace(style.Block.Start)) {
						result = result[:len(result)-1]
					}
				}
			} else {
				result = append(result, line)
			}

		case stateInsideMarking:
			if strings.Contains(line, comments.Marker) {
				if style.Block != nil {
					state = stateExpectingBlockEnd
				} else {
					state = stateNormal
				}
			}

		case stateExpectingBlockEnd:
			state = stateNormal

			if strings.HasPrefix(strings.TrimSpace(line), strings.TrimSpace(style.Block.End)) {
				// Drop line
			} else {
				result = append(result, line)
			}
		}
	}

	return strings.Join(result, "\n")
}
