package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhamidi/sai/format"
	"github.com/dhamidi/sai/java"
	"github.com/dhamidi/sai/java/parser"
	"github.com/spf13/cobra"
)

func newParseCmd() *cobra.Command {
	var outputFormat string
	var includeComments bool
	var includePositions bool

	cmd := &cobra.Command{
		Use:   "parse <file>",
		Short: "Parse a .class or .java file and dump the result",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			ext := filepath.Ext(filename)

			switch ext {
			case ".class":
				class, err := java.ParseClassFile(filename)
				if err != nil {
					return fmt.Errorf("parse class file: %w", err)
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

				if err := encoder.Encode(class); err != nil {
					return fmt.Errorf("encode: %w", err)
				}
			case ".java":
				data, err := os.ReadFile(filename)
				if err != nil {
					return fmt.Errorf("read java file: %w", err)
				}

				opts := []parser.Option{parser.WithFile(filename)}
				if includeComments {
					opts = append(opts, parser.WithComments())
				}
				if includePositions {
					opts = append(opts, parser.WithPositions())
				}
				p := parser.ParseCompilationUnit(bytes.NewReader(data), opts...)
				node := p.Finish()
				if node == nil {
					return fmt.Errorf("parse java file: incomplete or invalid syntax")
				}

				switch outputFormat {
				case "json":
					enc := format.NewASTJSONEncoder(os.Stdout)
					if err := enc.Encode(node); err != nil {
						return fmt.Errorf("encode json: %w", err)
					}
					fmt.Println()
				case "java":
					if p.IncludesPositions() {
						fmt.Println(node.StringWithPositions())
					} else {
						fmt.Println(node.String())
					}
				default:
					return fmt.Errorf("unknown format: %s", outputFormat)
				}
			default:
				return fmt.Errorf("unsupported file extension: %s (expected .class or .java)", ext)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "output format (json, java)")
	cmd.Flags().BoolVar(&includeComments, "comments", true, "include comments in output for .java files")
	cmd.Flags().BoolVar(&includePositions, "positions", true, "include token positions in output for .java files")

	return cmd
}
