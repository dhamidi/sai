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
	"time"

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
	var jsonOutput bool
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

				if jsonOutput {
					models, err := java.ClassModelsFromSource(data, parser.WithFile(filename))
					if err != nil {
						return fmt.Errorf("parse java file: %w", err)
					}
					if len(models) == 0 {
						return fmt.Errorf("parse java file: no classes found")
					}
					encoder := format.NewJSONModelEncoder(os.Stdout)
					for _, model := range models {
						if err := encoder.Encode(model); err != nil {
							return fmt.Errorf("encode: %w", err)
						}
						fmt.Println()
					}
				} else {
					p := parser.ParseCompilationUnit(bytes.NewReader(data), parser.WithFile(filename))
					node := p.Finish()
					if node == nil {
						return fmt.Errorf("parse java file: incomplete or invalid syntax")
					}
					fmt.Println(node.String())
				}
			default:
				return fmt.Errorf("unsupported file extension: %s (expected .class or .java)", ext)
			}

			return nil
		},
	}
	parseCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "output format (json, java)")
	parseCmd.Flags().BoolVar(&jsonOutput, "json", false, "output JSON for .java files")

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

	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.AddCommand(scanCmd)

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
				classes, parseErr = java.ClassModelsFromSource(data, parser.WithFile(path))
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
			classes, parseErr = java.ClassModelsFromSource(data, parser.WithFile(f.Name))
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
				fileClasses, parseErr = java.ClassModelsFromSource(data, parser.WithFile(f.Name))
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
