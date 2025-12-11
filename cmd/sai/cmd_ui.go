package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dhamidi/sai/ui"
	"github.com/spf13/cobra"
)

func newUICmd() *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Start the web UI server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := ui.NewServer()
			if err != nil {
				return fmt.Errorf("create server: %w", err)
			}
			displayAddr := addr
			if strings.HasPrefix(addr, ":") {
				displayAddr = "localhost" + addr
			}
			fmt.Printf("Starting server at http://%s\n", displayAddr)
			return http.ListenAndServe(addr, server)
		},
	}

	cmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "address to listen on")

	return cmd
}
