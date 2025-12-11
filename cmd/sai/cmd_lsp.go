package main

import (
	"github.com/dhamidi/sai/java/codebase"
	"github.com/spf13/cobra"
)

func newLSPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lsp",
		Short: "Start the Language Server Protocol server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server := codebase.NewLSPServer("0.1.0")
			return server.RunStdio()
		},
	}
}
