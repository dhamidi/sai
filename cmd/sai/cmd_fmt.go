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
		Use:   "fmt [file or directory]",
		Short: "Pretty-print a .java file, preserving comments",
		Long: `Pretty-print a .java file to stdout.

If a file is provided, it must have a .java extension.
If a directory is provided, formats all .java files in it (implies -w).
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
				info, err := os.Stat(filename)
				if err != nil {
					return fmt.Errorf("stat %s: %w", filename, err)
				}

				if info.IsDir() {
					return formatDirectory(filename)
				}

				ext := filepath.Ext(filename)
				if ext != ".java" {
					return fmt.Errorf("expected .java file, got %s", ext)
				}
				source, err = os.ReadFile(filename)
				if err != nil {
					return fmt.Errorf("read file: %w", err)
				}
			}

			output, err := format.PrettyPrintJavaFile(source, filename)
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

func formatDirectory(dir string) error {
	var formatted int

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(d.Name()) != ".java" {
			return nil
		}

		source, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		output, err := format.PrettyPrintJavaFile(source, path)
		if err != nil {
			return fmt.Errorf("format %s: %w", path, err)
		}

		if err := os.WriteFile(path, output, 0644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		formatted++
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("Formatted %d file(s)\n", formatted)
	return nil
}
