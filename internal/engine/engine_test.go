package engine

import (
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
 * markings:managed
 * some content
 * markings:managed
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
 * markings:managed
 * managed content
 * markings:managed
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
			content: `// markings:managed
// managed content
// markings:managed
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
