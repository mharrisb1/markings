// markings:managed
//
// File: engine.go
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
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"markings/internal/comments"
	"markings/internal/config"

	"github.com/bmatcuk/doublestar/v4"
)

type Engine struct {
	config    *config.Config
	templates map[string]*template.Template
}

// New creates a new Engine and pre-compiles all templates into their final formatted strings.
func New(cfg *config.Config) (*Engine, error) {
	eng := &Engine{
		config:    cfg,
		templates: make(map[string]*template.Template),
	}

	funcMap := template.FuncMap{
		"upper":     strings.ToUpper,
		"lower":     strings.ToLower,
		"title":     strings.ToTitle,
		"trim":      strings.Trim,
		"trimLeft":   strings.TrimLeft,
		"trimRight":  strings.TrimRight,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"split":      strings.Split,
		"join":       strings.Join,
		"replace":    strings.ReplaceAll,
		"time":      func(layout string) string { return time.Now().Format(layout) },
		"year":      func() string { return time.Now().Format("2006") },
		"month":     func() string { return time.Now().Format("01") },
		"day":       func() string { return time.Now().Format("02") },
		"date":      func() string { return time.Now().Format("2006-01-02") },
		"filename":  func() string { return "" }, // placeholder
	}

	for name, tmplStr := range cfg.Templates {
		tmpl, err := template.New(name).Funcs(funcMap).Parse(tmplStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %q: %w", name, err)
		}
		eng.templates[name] = tmpl
	}

	return eng, nil
}

func (e *Engine) matchRule(path string) (*config.Rule, error) {
	for i := range e.config.Rules {
		matched, err := doublestar.Match(e.config.Rules[i].Match, path)
		if err != nil {
			return nil, err
		}
		if matched {
			return &e.config.Rules[i], nil
		}
	}
	return nil, nil // No match
}

func (e *Engine) ProcessFile(path string, fix bool) (bool, bool, error) {
	rule, err := e.matchRule(path)
	if err != nil {
		return false, false, err
	}
	if rule == nil {
		return false, true, nil // No rule applies, so it's technically valid
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return true, false, err
	}
	content := string(contentBytes)

	style, ok := e.config.CommentStyles[rule.CommentStyle]
	if !ok {
		return true, false, fmt.Errorf("comment style %q not found", rule.CommentStyle)
	}

	expectedHeader, err := e.executeTemplate(rule.Header, style, path)
	if err != nil {
		return true, false, err
	}

	expectedFooter, err := e.executeTemplate(rule.Footer, style, path)
	if err != nil {
		return true, false, err
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
		return true, isValid, nil
	}

	// Fix logic: remove old markers and inject new ones
	newContent := removeOldMarkings(content, style, e.config.GetMarker())

	if expectedHeader != "" {
		newContent = expectedHeader + strings.TrimLeft(newContent, "\n")
	}
	if expectedFooter != "" {
		newContent = strings.TrimRight(newContent, "\n")
		newContent += "\n"
		newContent += expectedFooter
	}

	err = os.WriteFile(path, []byte(newContent), 0644)
	return true, err == nil, err
}

func (e *Engine) executeTemplate(c *config.MarkingConfig, style config.CommentStyle, path string) (string, error) {
	if c == nil || c.Template == "" {
		return "", nil
	}

	tmpl, ok := e.templates[c.Template]
	if !ok {
		return "", nil
	}

	tmplClone, err := tmpl.Clone()
	if err != nil {
		return "", err
	}

	tmplClone.Funcs(template.FuncMap{
		"filename": func() string { return filepath.Base(path) },
	})

	var buf bytes.Buffer
	if err := tmplClone.Execute(&buf, e.config.Data); err != nil {
		return "", err
	}

	out := comments.Format(buf.String(), style, e.config.GetMarker())
	
	out = strings.TrimRight(out, "\n")
	if c.NewlineBefore {
		out = "\n" + out
	}
	if c.NewlineAfter {
		out += "\n\n"
	} else {
		out += "\n"
	}

	return out, nil
}

type parseState int

const (
	stateNormal parseState = iota
	stateInsideMarking
	stateExpectingBlockEnd
)

// State machine for stripping out existing marking blocks
func removeOldMarkings(content string, style config.CommentStyle, marker string) string {
	lines := strings.Split(content, "\n")
	var result []string
	state := stateNormal

	for _, line := range lines {
		switch state {
		case stateNormal:
			if strings.Contains(line, marker) {
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
			if strings.Contains(line, marker) {
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
// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
