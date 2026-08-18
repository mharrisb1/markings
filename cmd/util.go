// markings:managed
//
// File: util.go
// Copyright (c) 2026 Michael Harris
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
//
// markings:managed

package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

type FileCallback func(path string)

func isExcluded(path string, globs []string) bool {
	for _, glob := range globs {
		if matched, _ := doublestar.Match(glob, path); matched {
			return true
		}
	}
	return false
}

func isIncluded(path string, globs []string) bool {
	if len(globs) == 0 {
		return true
	}
	for _, glob := range globs {
		if matched, _ := doublestar.Match(glob, path); matched {
			return true
		}
	}
	return false
}

func findFiles(root string, excludeGlobs []string, includeGlobs []string, callback FileCallback) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if isExcluded(path, excludeGlobs) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			if isIncluded(path, includeGlobs) {
				callback(path)
			}
		}

		return nil
	})
}

func expandArgs(args []string, excludeGlobs []string, includeGlobs []string) ([]string, error) {
	collected := []string{}

	appendCallback := func(path string) {
		collected = append(collected, path)
	}

	if len(args) == 0 {
		args = append(args, ".")
	}

	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return []string{}, err
		}

		if info.IsDir() {
			err = findFiles(arg, excludeGlobs, includeGlobs, appendCallback)
			if err != nil {
				return []string{}, err
			}
		} else {
			if !isExcluded(arg, excludeGlobs) && isIncluded(arg, includeGlobs) {
				appendCallback(arg)
			}
		}
	}

	return collected, nil
}

// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
