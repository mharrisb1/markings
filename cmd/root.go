package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "markings",
	Short: "Enfore source code banner markings",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", ".markings.yaml", "config file")
}
