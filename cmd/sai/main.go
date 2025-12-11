package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sai",
		Short: "A toasty java toolchain",
	}

	rootCmd.AddCommand(newParseCmd())
	rootCmd.AddCommand(newUICmd())
	rootCmd.AddCommand(newScanCmd())
	rootCmd.AddCommand(newLSPCmd())
	rootCmd.AddCommand(newDumpCmd())
	rootCmd.AddCommand(newFmtCmd())
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newClasspathCmd())
	rootCmd.AddCommand(newInitCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
