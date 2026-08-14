package engine

import (
	"bytes"
	"fmt"
	"text/template"

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
