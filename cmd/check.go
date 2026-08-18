// markings:managed
//
// File: check.go
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

	"markings/internal/config"
	"markings/internal/engine"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [paths...]",
	Short: "Checks if the specified files have the correct markings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		eng, err := engine.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize engine: %w", err)
		}

		hasErrors := false
		files, err := expandArgs(args, cfg.Exclude, cfg.Include)
		if err != nil {
			return fmt.Errorf("failed to expand files: %w", err)
		}

		for _, file := range files {
			_, valid, err := eng.ProcessFile(file, false)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", file, err)
				hasErrors = true
				continue
			}
			if !valid {
				fmt.Printf("File %s is missing required markings or has incorrect markings\n", file)
				hasErrors = true
			}
		}

		if hasErrors {
			os.Exit(1)
		}

		fmt.Println("All files have the correct markings.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

// markings:managed
//
// Copyright (c) 2026 Michael Harris
//
// markings:managed
