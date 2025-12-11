package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dhamidi/sai/format"
	"github.com/spf13/cobra"
)

func newFmtCmd() *cobra.Command {
	var fmtOverwrite bool

	cmd := &cobra.Command{
		Use:   "fmt [file]",
		Short: "Pretty-print a .java file, preserving comments",
		Long: `Pretty-print a .java file to stdout.

If a file is provided, it must have a .java extension.
If no file is provided, reads Java source from stdin.

Use -w to overwrite the file in place (requires a file argument).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var source []byte
			var err error
			var filename string

			if len(args) == 0 {
				if fmtOverwrite {
					return fmt.Errorf("-w requires a file argument")
				}
				source, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("read stdin: %w", err)
				}
			} else {
				filename = args[0]
				ext := filepath.Ext(filename)
				if ext != ".java" {
					return fmt.Errorf("expected .java file, got %s", ext)
				}
				source, err = os.ReadFile(filename)
				if err != nil {
					return fmt.Errorf("read file: %w", err)
				}
			}

			output, err := format.PrettyPrintJava(source)
			if err != nil {
				return fmt.Errorf("format: %w", err)
			}

			if fmtOverwrite {
				return os.WriteFile(filename, output, 0644)
			}
			_, err = os.Stdout.Write(output)
			return err
		},
	}

	cmd.Flags().BoolVarP(&fmtOverwrite, "write", "w", false, "overwrite the file in place")

	return cmd
}
