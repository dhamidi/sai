package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dhamidi/javalyzer/format"
	"github.com/dhamidi/javalyzer/java"
	"github.com/dhamidi/javalyzer/java/parser"
	"github.com/dhamidi/javalyzer/ui"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "javalyzer",
		Short: "Java class file analyzer",
	}

	var outputFormat string
	parseCmd := &cobra.Command{
		Use:   "parse <file>",
		Short: "Parse a .class or .java file and dump the result",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			ext := filepath.Ext(filename)

			var result any
			switch ext {
			case ".class":
				class, err := java.ParseClassFile(filename)
				if err != nil {
					return fmt.Errorf("parse class file: %w", err)
				}
				result = class
			case ".java":
				data, err := os.ReadFile(filename)
				if err != nil {
					return fmt.Errorf("read java file: %w", err)
				}
				p := parser.ParseCompilationUnit(parser.WithFile(filename))
				p.Push(data)
				node := p.Finish()
				if node == nil {
					return fmt.Errorf("parse java file: incomplete or invalid syntax")
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				if err := enc.Encode(node); err != nil {
					return fmt.Errorf("encode: %w", err)
				}
				return nil
			default:
				return fmt.Errorf("unsupported file extension: %s (expected .class or .java)", ext)
			}

			var encoder format.Encoder
			switch outputFormat {
			case "json":
				encoder = format.NewJSONEncoder(os.Stdout)
			case "java":
				encoder = format.NewJavaEncoder(os.Stdout)
			default:
				return fmt.Errorf("unknown format: %s", outputFormat)
			}

			if err := encoder.Encode(result.(*java.Class)); err != nil {
				return fmt.Errorf("encode: %w", err)
			}

			return nil
		},
	}
	parseCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "output format (json, java)")

	var addr string
	uiCmd := &cobra.Command{
		Use:   "ui",
		Short: "Start the web UI server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := ui.NewServer()
			if err != nil {
				return fmt.Errorf("create server: %w", err)
			}
			fmt.Printf("Starting server at http://%s\n", addr)
			return http.ListenAndServe(addr, server)
		},
	}
	uiCmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "address to listen on")

	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(uiCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
