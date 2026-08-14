package comments

import (
	"markings/internal/config"
	"strings"
)

const Marker = "markings:managed"

// Format applies the given comment style to the raw text and adds demarcation markers.
func Format(text string, style config.CommentStyle) string {
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")

	// Drop the last empty string if text ends with a newline
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	var b strings.Builder

	if style.Block != nil {
		b.WriteString(style.Block.Start)

		// Start marker
		b.WriteString(style.Block.Middle)
		b.WriteString(Marker)
		b.WriteString("\n")

		// Content
		for _, line := range lines {
			b.WriteString(style.Block.Middle)
			b.WriteString(line)
			b.WriteString("\n")
		}

		// End marker
		b.WriteString(style.Block.Middle)
		b.WriteString(Marker)
		b.WriteString("\n")

		b.WriteString(style.Block.End)
		return b.String()
	}

	// Fallback to line comments

	// Start marker
	b.WriteString(style.Prefix)
	b.WriteString(Marker)
	b.WriteString("\n")

	// Content
	for _, line := range lines {
		b.WriteString(style.Prefix)
		b.WriteString(line)
		b.WriteString("\n")
	}

	// End marker
	b.WriteString(style.Prefix)
	b.WriteString(Marker)
	b.WriteString("\n")

	return b.String()
}
