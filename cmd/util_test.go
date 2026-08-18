// markings:managed
//
// File: util_test.go
// Copyright (c) 2026 Michael Harris
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
//
// markings:managed

package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestExpandArgs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test directory structure
	// tmpDir/
	// ├── file1.txt
	// ├── file3.txt
	// └── subdir/
	//     ├── file2.txt
	//     └── .git/
	//         └── hidden.txt

	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file3.txt"), []byte("test"), 0644)

	subdir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subdir, 0755)
	os.WriteFile(filepath.Join(subdir, "file2.txt"), []byte("test"), 0644)

	// Mock a .git directory to test exclusion behavior
	gitDir := filepath.Join(subdir, ".git")
	os.MkdirAll(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "hidden.txt"), []byte("test"), 0644)

	tests := []struct {
		name     string
		args     []string
		exclude  []string
		include  []string
		expected []string
	}{
		{
			name:     "Explicit file without config",
			args:     []string{filepath.Join(tmpDir, "file1.txt")},
			exclude:  nil,
			include:  nil,
			expected: []string{filepath.Join(tmpDir, "file1.txt")},
		},
		{
			name:     "Directory expansion with .git exclusion",
			args:     []string{tmpDir},
			exclude:  []string{"**/.git"},
			include:  nil,
			expected: []string{
				filepath.Join(tmpDir, "file1.txt"),
				filepath.Join(tmpDir, "file3.txt"),
				filepath.Join(tmpDir, "subdir", "file2.txt"),
			},
		},
		{
			name:     "Include only specific files",
			args:     []string{tmpDir},
			exclude:  []string{"**/.git"},
			include:  []string{"**/file1.txt", "**/file2.txt"},
			expected: []string{
				filepath.Join(tmpDir, "file1.txt"),
				filepath.Join(tmpDir, "subdir", "file2.txt"),
			},
		},
		{
			name:     "Exclude specific file even if included",
			args:     []string{tmpDir},
			exclude:  []string{"**/file1.txt"},
			include:  []string{"**/*.txt"},
			expected: []string{
				filepath.Join(tmpDir, "file3.txt"),
				filepath.Join(tmpDir, "subdir", "file2.txt"),
				filepath.Join(tmpDir, "subdir", ".git", "hidden.txt"),
			},
		},
		{
			name:     "Multiple arguments with exclusions",
			args:     []string{
				filepath.Join(tmpDir, "file1.txt"),
				subdir,
			},
			exclude:  []string{"**/.git"},
			include:  nil,
			expected: []string{
				filepath.Join(tmpDir, "file1.txt"),
				filepath.Join(tmpDir, "subdir", "file2.txt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandArgs(tt.args, tt.exclude, tt.include)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Sort slices to prevent flaky tests from arbitrary map/walk ordering
			sort.Strings(got)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("\nExpected:\n%v\nGot:\n%v", tt.expected, got)
			}
		})
	}
}

// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
