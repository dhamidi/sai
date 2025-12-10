package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhamidi/sai/format"
	"github.com/dhamidi/sai/java"
	"github.com/dhamidi/sai/java/codebase"
	"github.com/dhamidi/sai/java/parser"
	"github.com/dhamidi/sai/ui"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sai",
		Short: "A toasty java toolchain",
	}

	var outputFormat string
	var includeComments bool
	var includePositions bool
	parseCmd := &cobra.Command{
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
	parseCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "output format (json, java)")
	parseCmd.Flags().BoolVar(&includeComments, "comments", true, "include comments in output for .java files")
	parseCmd.Flags().BoolVar(&includePositions, "positions", true, "include token positions in output for .java files")

	var addr string
	uiCmd := &cobra.Command{
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
	uiCmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "address to listen on")

	var timeout time.Duration
	scanCmd := &cobra.Command{
		Use:   "scan <path>",
		Short: "Scan a directory, jar, or zip file for Java classes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return runScan(path, timeout)
		},
	}
	scanCmd.Flags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "timeout per file")

	lspCmd := &cobra.Command{
		Use:   "lsp",
		Short: "Start the Language Server Protocol server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server := codebase.NewLSPServer("0.1.0")
			return server.RunStdio()
		},
	}

	var dumpFormat string
	dumpCmd := &cobra.Command{
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
	dumpCmd.Flags().StringVarP(&dumpFormat, "format", "f", "line", "output format (json, java, line)")

	var fmtOverwrite bool
	fmtCmd := &cobra.Command{
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
	fmtCmd.Flags().BoolVarP(&fmtOverwrite, "write", "w", false, "overwrite the file in place")

	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(lspCmd)
	rootCmd.AddCommand(dumpCmd)
	rootCmd.AddCommand(fmtCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runScan(path string, timeout time.Duration) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	var classes []*java.ClassModel
	var errors []string

	if info.IsDir() {
		classes, errors = scanDirectory(path, timeout)
	} else {
		ext := filepath.Ext(path)
		if ext == ".jar" || ext == ".zip" {
			classes, errors = scanZipFile(path, timeout)
		} else if ext == ".class" || ext == ".java" {
			classes, errors = scanSingleFile(path, timeout)
		} else {
			return fmt.Errorf("unsupported file type: %s", ext)
		}
	}

	fmt.Printf("\n=== SCAN COMPLETE ===\n")
	fmt.Printf("Classes found: %d\n", len(classes))
	fmt.Printf("Errors: %d\n", len(errors))
	for _, e := range errors {
		fmt.Printf("  - %s\n", e)
	}
	return nil
}

func scanSingleFile(path string, timeout time.Duration) ([]*java.ClassModel, []string) {
	ext := filepath.Ext(path)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	var classes []*java.ClassModel
	var parseErr error

	go func() {
		defer close(done)
		switch ext {
		case ".class":
			class, err := java.ClassModelFromFile(path)
			if err != nil {
				parseErr = err
			} else if class != nil {
				classes = []*java.ClassModel{class}
			}
		case ".java":
			data, err := os.ReadFile(path)
			if err != nil {
				parseErr = err
			} else {
				classes, parseErr = java.ClassModelsFromSource(data, parser.WithFile(path), parser.WithSourcePath(path))
			}
		}
	}()

	select {
	case <-done:
		if parseErr != nil {
			return nil, []string{fmt.Sprintf("parse %s: %v", path, parseErr)}
		}
		fmt.Printf("[OK] %s (%d classes)\n", path, len(classes))
		return classes, nil
	case <-ctx.Done():
		return nil, []string{fmt.Sprintf("timeout parsing %s", path)}
	}
}

func scanDirectory(path string, timeout time.Duration) ([]*java.ClassModel, []string) {
	var files []string
	var errors []string

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, fmt.Sprintf("walk %s: %v", p, err))
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(p)
			if ext == ".class" || ext == ".java" {
				files = append(files, p)
			}
		}
		return nil
	})
	if err != nil {
		errors = append(errors, fmt.Sprintf("walk %s: %v", path, err))
	}

	fmt.Printf("Found %d files to scan\n", len(files))

	var classes []*java.ClassModel
	for i, file := range files {
		fmt.Printf("[%d/%d] ", i+1, len(files))
		fileClasses, fileErrors := scanSingleFile(file, timeout)
		classes = append(classes, fileClasses...)
		errors = append(errors, fileErrors...)
	}

	return classes, errors
}

func scanZipFile(path string, timeout time.Duration) ([]*java.ClassModel, []string) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, []string{fmt.Sprintf("open zip: %v", err)}
	}
	defer r.Close()

	var sourceFiles []*zip.File
	var jarFiles []*zip.File
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name)
		switch ext {
		case ".class", ".java":
			sourceFiles = append(sourceFiles, f)
		case ".jar":
			jarFiles = append(jarFiles, f)
		}
	}

	total := len(sourceFiles)
	for _, jarFile := range jarFiles {
		total += countFilesInJar(jarFile)
	}

	fmt.Printf("Found %d files to scan (%d source files, %d jars)\n", total, len(sourceFiles), len(jarFiles))

	var classes []*java.ClassModel
	var errors []string
	progress := 0

	for _, f := range sourceFiles {
		progress++
		fmt.Printf("[%d/%d] ", progress, total)
		fileClasses, fileErrors := scanZipEntry(f, path, timeout)
		classes = append(classes, fileClasses...)
		errors = append(errors, fileErrors...)
	}

	for _, jarFile := range jarFiles {
		jarClasses, jarErrors := scanJarInZip(jarFile, timeout, &progress, total)
		classes = append(classes, jarClasses...)
		errors = append(errors, jarErrors...)
	}

	return classes, errors
}

func scanZipEntry(f *zip.File, zipPath string, timeout time.Duration) ([]*java.ClassModel, []string) {
	rc, err := f.Open()
	if err != nil {
		return nil, []string{fmt.Sprintf("open %s: %v", f.Name, err)}
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, []string{fmt.Sprintf("read %s: %v", f.Name, err)}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	var classes []*java.ClassModel
	var parseErr error

	go func() {
		defer close(done)
		ext := filepath.Ext(f.Name)
		switch ext {
		case ".class":
			class, err := java.ClassModelFromReader(bytes.NewReader(data))
			if err != nil {
				parseErr = err
			} else if class != nil {
				classes = []*java.ClassModel{class}
			}
		case ".java":
			classes, parseErr = java.ClassModelsFromSource(data, parser.WithFile(f.Name), parser.WithSourcePath(f.Name))
		}
	}()

	select {
	case <-done:
		if parseErr != nil {
			fmt.Printf("[ERROR] %s: %v\n", f.Name, parseErr)
			return nil, []string{fmt.Sprintf("parse %s: %v", f.Name, parseErr)}
		}
		fmt.Printf("[OK] %s (%d classes)\n", f.Name, len(classes))
		return classes, nil
	case <-ctx.Done():
		fmt.Printf("[TIMEOUT] %s\n", f.Name)
		return nil, []string{fmt.Sprintf("timeout parsing %s", f.Name)}
	}
}

func countFilesInJar(jarFile *zip.File) int {
	rc, err := jarFile.Open()
	if err != nil {
		return 0
	}
	defer rc.Close()

	jarData, err := io.ReadAll(rc)
	if err != nil {
		return 0
	}

	jarReader, err := zip.NewReader(bytes.NewReader(jarData), int64(len(jarData)))
	if err != nil {
		return 0
	}

	count := 0
	for _, f := range jarReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name)
		if ext == ".class" || ext == ".java" {
			count++
		}
	}
	return count
}

func scanJarInZip(jarFile *zip.File, timeout time.Duration, progress *int, total int) ([]*java.ClassModel, []string) {
	rc, err := jarFile.Open()
	if err != nil {
		return nil, []string{fmt.Sprintf("open jar %s: %v", jarFile.Name, err)}
	}
	defer rc.Close()

	jarData, err := io.ReadAll(rc)
	if err != nil {
		return nil, []string{fmt.Sprintf("read jar %s: %v", jarFile.Name, err)}
	}

	jarReader, err := zip.NewReader(bytes.NewReader(jarData), int64(len(jarData)))
	if err != nil {
		return nil, []string{fmt.Sprintf("open jar %s as zip: %v", jarFile.Name, err)}
	}

	var classes []*java.ClassModel
	var errors []string
	for _, f := range jarReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name)
		if ext != ".class" && ext != ".java" {
			continue
		}

		*progress++
		fmt.Printf("[%d/%d] %s: ", *progress, total, jarFile.Name)

		fileRC, err := f.Open()
		if err != nil {
			fmt.Printf("[ERROR] open %s: %v\n", f.Name, err)
			errors = append(errors, fmt.Sprintf("open %s in %s: %v", f.Name, jarFile.Name, err))
			continue
		}

		data, err := io.ReadAll(fileRC)
		fileRC.Close()
		if err != nil {
			fmt.Printf("[ERROR] read %s: %v\n", f.Name, err)
			errors = append(errors, fmt.Sprintf("read %s in %s: %v", f.Name, jarFile.Name, err))
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		done := make(chan struct{})
		var fileClasses []*java.ClassModel
		var parseErr error

		go func() {
			defer close(done)
			switch ext {
			case ".class":
				class, err := java.ClassModelFromReader(bytes.NewReader(data))
				if err != nil {
					parseErr = err
				} else if class != nil {
					fileClasses = []*java.ClassModel{class}
				}
			case ".java":
				fileClasses, parseErr = java.ClassModelsFromSource(data, parser.WithFile(f.Name), parser.WithSourcePath(f.Name))
			}
		}()

		select {
		case <-done:
			cancel()
			if parseErr != nil {
				fmt.Printf("[ERROR] %s: %v\n", f.Name, parseErr)
				errors = append(errors, fmt.Sprintf("parse %s in %s: %v", f.Name, jarFile.Name, parseErr))
			} else {
				fmt.Printf("[OK] %s (%d classes)\n", f.Name, len(fileClasses))
				classes = append(classes, fileClasses...)
			}
		case <-ctx.Done():
			cancel()
			fmt.Printf("[TIMEOUT] %s\n", f.Name)
			errors = append(errors, fmt.Sprintf("timeout parsing %s in %s", f.Name, jarFile.Name))
		}
	}

	return classes, errors
}
