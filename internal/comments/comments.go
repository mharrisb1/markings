// markings:managed
//
// Copyright (c) 2026 Michael Harris
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
//
// markings:managed

package comments

import (
	"markings/internal/config"
	"strings"
)

// Format applies the given comment style to the raw text and adds demarcation markers.
func Format(text string, style config.CommentStyle, marker string) string {
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
		b.WriteString(marker)
		b.WriteString("\n")

		b.WriteString(strings.TrimRight(style.Block.Middle, " \t"))
		b.WriteString("\n")

		// Content
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				b.WriteString(strings.TrimRight(style.Block.Middle, " \t"))
			} else {
				b.WriteString(style.Block.Middle)
				b.WriteString(line)
			}
			b.WriteString("\n")
		}

		b.WriteString(strings.TrimRight(style.Block.Middle, " \t"))
		b.WriteString("\n")

		// End marker
		b.WriteString(style.Block.Middle)
		b.WriteString(marker)
		b.WriteString("\n")

		b.WriteString(style.Block.End)
		return b.String()
	}

	// Fallback to line comments

	// Start marker
	b.WriteString(style.Prefix)
	b.WriteString(marker)
	b.WriteString("\n")

	b.WriteString(strings.TrimRight(style.Prefix, " \t"))
	b.WriteString("\n")

	// Content
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			b.WriteString(strings.TrimRight(style.Prefix, " \t"))
		} else {
			b.WriteString(style.Prefix)
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	b.WriteString(strings.TrimRight(style.Prefix, " \t"))
	b.WriteString("\n")

	// End marker
	b.WriteString(style.Prefix)
	b.WriteString(marker)
	b.WriteString("\n")

	return b.String()
}

// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
