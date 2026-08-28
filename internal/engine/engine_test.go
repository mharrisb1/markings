// markings:managed
//
// File: engine_test.go
// Copyright (c) 2026 Michael Harris
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
//
// markings:managed

package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mharrisb1/markings/internal/config"
)

func TestRemoveOldMarkings(t *testing.T) {
	blockStyle := config.CommentStyle{
		Block: &config.BlockStyle{
			Start:  "/*",
			Middle: " *",
			End:    " */",
		},
	}

	lineStyle := config.CommentStyle{
		Prefix: "// ",
	}

	tests := []struct {
		name     string
		style    config.CommentStyle
		content  string
		expected string
	}{
		{
			name:  "Bug: Managed block followed by normal block",
			style: blockStyle,
			content: `/*
 * ` + config.DefaultMarker + `
 * some content
 * ` + config.DefaultMarker + `
 */
/*
 * Normal comment that should NOT be touched
 */`,
			expected: `/*
 * Normal comment that should NOT be touched
 */`,
		},
		{
			name:  "No managed blocks, only normal blocks",
			style: blockStyle,
			content: `/*
 * First comment
 */
/*
 * Second comment
 */`,
			expected: `/*
 * First comment
 */
/*
 * Second comment
 */`,
		},
		{
			name:  "Managed block in middle of file",
			style: blockStyle,
			content: `/*
 * Normal comment
 */
/*
 * ` + config.DefaultMarker + `
 * managed content
 * ` + config.DefaultMarker + `
 */
var x = 1;`,
			expected: `/*
 * Normal comment
 */
var x = 1;`,
		},
		{
			name:  "Line style comments",
			style: lineStyle,
			content: `// ` + config.DefaultMarker + `
// managed content
// ` + config.DefaultMarker + `
// Normal comment
var y = 2;`,
			expected: `// Normal comment
var y = 2;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeOldMarkings(tt.content, tt.style, config.DefaultMarker)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestExecuteTemplate(t *testing.T) {
	cfg := &config.Config{
		Marker: config.DefaultMarker,
		Templates: map[string]string{
			"test": "File: {{ filename }}, Upper: {{ upper .TestVar }}, Replace: {{ replace .TestVar `foo` `bar` }}",
		},
		Data: map[string]interface{}{
			"TestVar": "foobar",
		},
	}

	eng, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	style := config.CommentStyle{
		Prefix: "// ",
	}

	ruleHeader := &config.MarkingConfig{
		Template:      "test",
		NewlineBefore: true,
		NewlineAfter:  true,
	}

	out, err := eng.executeTemplate(ruleHeader, style, "/path/to/my_file.go")
	if err != nil {
		t.Fatalf("executeTemplate failed: %v", err)
	}

	expected := "\n// " + config.DefaultMarker + "\n//\n// File: my_file.go, Upper: FOOBAR, Replace: barbar\n//\n// " + config.DefaultMarker + "\n\n"

	if out != expected {
		t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, out)
	}
}

func TestExecuteTemplate_BlockStyle(t *testing.T) {
	cfg := &config.Config{
		Marker: config.DefaultMarker,
		Templates: map[string]string{
			"test": "File: {{ filename }}",
		},
		Data: map[string]interface{}{},
	}

	eng, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	style := config.CommentStyle{
		Block: &config.BlockStyle{
			Start:  "/*\n",
			Middle: " * ",
			End:    " */",
		},
	}

	ruleHeader := &config.MarkingConfig{
		Template:      "test",
		NewlineBefore: true,
		NewlineAfter:  true,
	}

	out, err := eng.executeTemplate(ruleHeader, style, "/path/to/my_file.js")
	if err != nil {
		t.Fatalf("executeTemplate failed: %v", err)
	}

	expected := "\n/*\n * " + config.DefaultMarker + "\n *\n * File: my_file.js\n *\n * " + config.DefaultMarker + "\n */\n\n"

	if out != expected {
		t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, out)
	}
}

func TestProcessFile_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	emptyFilePath := filepath.Join(tempDir, "empty.go")
	if err := os.WriteFile(emptyFilePath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	tests := []struct {
		name           string
		markEmptyFiles bool
		expectedFound  bool
		expectedValid  bool
	}{
		{
			name:           "Skip empty files",
			markEmptyFiles: false,
			expectedFound:  false,
			expectedValid:  true,
		},
		{
			name:           "Process empty files",
			markEmptyFiles: true,
			expectedFound:  true,
			expectedValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				MarkEmptyFiles: tt.markEmptyFiles,
				Rules: []config.Rule{
					{
						Match:        filepath.ToSlash(filepath.Join(tempDir, "*.go")),
						CommentStyle: "line",
						Header: &config.MarkingConfig{
							Template: "test",
						},
					},
				},
				CommentStyles: map[string]config.CommentStyle{
					"line": {Prefix: "// "},
				},
				Templates: map[string]string{
					"test": "Header",
				},
			}

			eng, err := New(cfg)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			res := eng.ProcessFile(emptyFilePath, false)

			if res.Err != nil {
				t.Fatalf("Unexpected error: %v", res.Err)
			}
			if res.FoundRule != tt.expectedFound {
				t.Errorf("Expected FoundRule=%v, got %v", tt.expectedFound, res.FoundRule)
			}
			if res.IsValid != tt.expectedValid {
				t.Errorf("Expected IsValid=%v, got %v", tt.expectedValid, res.IsValid)
			}
		})
	}
}

// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
