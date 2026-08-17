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
	"markings/internal/comments"
	"markings/internal/config"
	"testing"
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
 * ` + comments.Marker + `
 * some content
 * ` + comments.Marker + `
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
 * ` + comments.Marker + `
 * managed content
 * ` + comments.Marker + `
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
			content: `// ` + comments.Marker + `
// managed content
// ` + comments.Marker + `
// Normal comment
var y = 2;`,
			expected: `// Normal comment
var y = 2;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeOldMarkings(tt.content, tt.style)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
