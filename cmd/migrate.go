// markings:managed
//
// File: migrate.go
// Copyright (c) 2026 Michael Harris
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
//
// markings:managed

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mharrisb1/markings/internal/config"

	"github.com/spf13/cobra"
)

var migrateMarkerCmd = &cobra.Command{
	Use:   "migrate-marker <old-marker> [paths...]",
	Short: "Migrate existing markings to a new marker defined in the config",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		oldMarker := args[0]
		files, err := expandArgs(args[1:], cfg.Exclude, cfg.Include)
		if err != nil {
			return fmt.Errorf("failed to expand files: %w", err)
		}

		newMarker := cfg.GetMarker()

		if oldMarker == newMarker {
			fmt.Println("Old marker and new marker are the same. Nothing to do.")
			return nil
		}

		for _, file := range files {
			if err := migrateFile(file, oldMarker, newMarker); err != nil {
				fmt.Printf("Error migrating %s: %v\n", file, err)
			}
		}

		return nil
	},
}

func migrateFile(path, oldMarker, newMarker string) error {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	if !strings.Contains(content, oldMarker) {
		return nil
	}

	newContent := strings.ReplaceAll(content, oldMarker, newMarker)
	err = os.WriteFile(path, []byte(newContent), 0644)
	if err == nil {
		fmt.Printf("Migrated %s\n", path)
	}
	return err
}

func init() {
	rootCmd.AddCommand(migrateMarkerCmd)
}

// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
