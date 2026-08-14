package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix [files...]",
	Short: "Adds or updates banner markings for target files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fix command called with args:", args)
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
