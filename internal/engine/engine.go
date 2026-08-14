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
	if rule.HeaderTemplate != "" {
		if raw, ok := e.headers[rule.HeaderTemplate]; ok {
			expectedHeader = comments.Format(raw, style)
			if style.NewlineAfter {
				expectedHeader += "\n"
			}
		}
	}

	expectedFooter := ""
	if rule.FooterTemplate != "" {
		if raw, ok := e.footers[rule.FooterTemplate]; ok {
			expectedFooter = comments.Format(raw, style)
			if style.NewlineBefore {
				expectedFooter = "\n" + expectedFooter
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

// removeOldMarkings strips out existing marking blocks using the markers
func removeOldMarkings(content string, style config.CommentStyle) string {
	lines := strings.Split(content, "\n")
	var result []string
	inMarking := false

	for _, line := range lines {
		if strings.Contains(line, comments.Marker) {
			if !inMarking {
				inMarking = true
				if style.Block != nil && len(result) > 0 {
					// Remove the preceding block start (e.g. `/*`)
					if strings.HasPrefix(strings.TrimSpace(result[len(result)-1]), strings.TrimSpace(style.Block.Start)) {
						result = result[:len(result)-1]
					}
				}
				continue
			} else {
				inMarking = false
				continue
			}
		}

		if inMarking {
			continue
		}

		// Remove the trailing block end (e.g. `*/`) immediately following a marking block
		if !inMarking && style.Block != nil && strings.HasPrefix(strings.TrimSpace(line), strings.TrimSpace(style.Block.End)) {
			// Only skip if the previous line was the end of a marking block (we just flipped inMarking to false on the prev loop)
			// A simpler heuristic: if we aren't in a marking, and this is an end block, and we haven't added anything since the marking ended...
			// For brevity, we just drop it if it's right here.
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
