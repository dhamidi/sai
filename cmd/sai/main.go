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

	rootCmd.AddCommand(newDumpCmd())
	rootCmd.AddCommand(newFmtCmd())
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newClasspathCmd())
	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newLibsCmd())
	rootCmd.AddCommand(newCompileCmd())
	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newTestCmd())
	rootCmd.AddCommand(newProjectCmd())
	rootCmd.AddCommand(newBakeCmd())
	rootCmd.AddCommand(newDocCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
