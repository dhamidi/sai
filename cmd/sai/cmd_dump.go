package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhamidi/sai/format"
	"github.com/dhamidi/sai/java"
	"github.com/dhamidi/sai/java/parser"
	"github.com/spf13/cobra"
)

func newDumpCmd() *cobra.Command {
	var dumpFormat string

	cmd := &cobra.Command{
		Use:   "dump <file>",
		Short: "Dump the class model from a .class or .java file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			ext := filepath.Ext(filename)

			var models []*java.ClassModel
			var err error

			switch ext {
			case ".class":
				model, e := java.ClassModelFromFile(filename)
				if e != nil {
					return fmt.Errorf("parse class file: %w", e)
				}
				models = []*java.ClassModel{model}
			case ".java":
				data, e := os.ReadFile(filename)
				if e != nil {
					return fmt.Errorf("read java file: %w", e)
				}
				models, err = java.ClassModelsFromSource(data, parser.WithFile(filename), parser.WithSourcePath(filename))
				if err != nil {
					return fmt.Errorf("parse java file: %w", err)
				}
			default:
				return fmt.Errorf("unsupported file extension: %s (expected .class or .java)", ext)
			}

			for _, model := range models {
				switch dumpFormat {
				case "json":
					enc := format.NewJSONModelEncoder(os.Stdout)
					if err := enc.Encode(model); err != nil {
						return fmt.Errorf("encode json: %w", err)
					}
					fmt.Println()
				case "java":
					enc := format.NewJavaModelEncoder(os.Stdout)
					if err := enc.Encode(model); err != nil {
						return fmt.Errorf("encode java: %w", err)
					}
				case "line":
					enc := format.NewLineModelEncoder(os.Stdout)
					if err := enc.Encode(model); err != nil {
						return fmt.Errorf("encode line: %w", err)
					}
				default:
					return fmt.Errorf("unknown format: %s (expected json, java, or line)", dumpFormat)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dumpFormat, "format", "f", "line", "output format (json, java, line)")

	return cmd
}
