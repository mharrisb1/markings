package cmd

import (
	"fmt"
	
	"markings/internal/config"
	"markings/internal/engine"

	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix [files...]",
	Short: "Adds or updates the markings on the specified files",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		eng, err := engine.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize engine: %w", err)
		}

		for _, file := range args {
			_, err := eng.ProcessFile(file, true)
			if err != nil {
				fmt.Printf("Error fixing %s: %v\n", file, err)
			} else {
				fmt.Printf("Fixed %s\n", file)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
