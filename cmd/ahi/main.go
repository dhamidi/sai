package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ahi",
		Short: "Development tools for sai",
	}

	rootCmd.AddCommand(newEbnfCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
