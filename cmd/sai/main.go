package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhamidi/sai/format"
	"github.com/dhamidi/sai/java"
	"github.com/dhamidi/sai/java/codebase"
	"github.com/dhamidi/sai/java/parser"
	"github.com/dhamidi/sai/pom"
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

	var libDir string
	addCmd := &cobra.Command{
		Use:   "add <groupId:artifactId:version>",
		Short: "Download a Maven dependency and its transitive dependencies to lib/",
		Long: `Download a Maven dependency and its transitive dependencies.

The coordinate format is: groupId:artifactId:version
Or with classifier: groupId:artifactId:classifier:version

Examples:
  sai add com.google.guava:guava:32.1.3-jre
  sai add org.slf4j:slf4j-api:2.0.9

Environment variables:
  MAVEN_REPO_URL - Override the Maven repository URL (default: https://repo1.maven.org/maven2)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(args[0], libDir)
		},
	}
	addCmd.Flags().StringVarP(&libDir, "lib", "l", "lib", "directory to download JARs to")

	var cpLibDir string
	classpathCmd := &cobra.Command{
		Use:   "classpath",
		Short: "Print the classpath from pom.xml or lib/ directory",
		Long: `Print the classpath as a colon-separated list of JAR paths.

If pom.xml exists in the current directory, dependencies are resolved
from it and printed as Maven repository paths (requires downloading).

Otherwise, all .jar files in the lib/ directory (or specified via -l)
are listed.

Examples:
  sai classpath              # Use pom.xml if present, else lib/
  sai classpath -l deps/     # Use deps/ directory`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClasspath(cpLibDir)
		},
	}
	classpathCmd.Flags().StringVarP(&cpLibDir, "lib", "l", "lib", "directory containing JAR files")

	var projectID string
	initCmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize a new sai Java project",
		Long: `Initialize a new sai Java project.

If a directory is provided, creates it and initializes the project there.
Otherwise, initializes in the current directory.

The project identifier defaults to the directory basename if it's a valid
Java identifier. Use -p to override.

This command:
  - Ensures a git repository exists
  - Creates the modular source structure: src/<project-id>/{core,main}/
  - Creates .gitignore with out/ and mlib/
  - Creates AGENTS.md with sai workflow instructions
  - Creates CLAUDE.md as a symlink to AGENTS.md

Examples:
  sai init myapp                 # Create myapp/ with project id "myapp"
  sai init -p jq myproject       # Create myproject/ with project id "jq"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			pid := projectID
			if pid == "" {
				absDir, err := filepath.Abs(dir)
				if err != nil {
					return fmt.Errorf("resolve directory: %w", err)
				}
				basename := filepath.Base(absDir)
				if isValidJavaIdentifier(basename) {
					pid = basename
				} else {
					return fmt.Errorf("directory name %q is not a valid Java identifier; use -p to specify project id", basename)
				}
			}
			return runInit(dir, pid)
		},
	}
	initCmd.Flags().StringVarP(&projectID, "project", "p", "", "project identifier (defaults to directory name)")

	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(uiCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(lspCmd)
	rootCmd.AddCommand(dumpCmd)
	rootCmd.AddCommand(fmtCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(classpathCmd)
	rootCmd.AddCommand(initCmd)

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

func runInit(dir string, projectID string) error {
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}
		fmt.Printf("Created %s/\n", dir)
	}

	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "init")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git init: %w", err)
		}
	} else {
		fmt.Println("Git repository already exists")
	}

	corePath := filepath.Join(dir, "src", projectID, "core")
	if err := os.MkdirAll(corePath, 0755); err != nil {
		return fmt.Errorf("create src/%s/core: %w", projectID, err)
	}
	fmt.Printf("Created src/%s/core/\n", projectID)

	mainPath := filepath.Join(dir, "src", projectID, "main")
	if err := os.MkdirAll(mainPath, 0755); err != nil {
		return fmt.Errorf("create src/%s/main: %w", projectID, err)
	}
	fmt.Printf("Created src/%s/main/\n", projectID)

	libPath := filepath.Join(dir, "lib")
	if err := os.MkdirAll(libPath, 0755); err != nil {
		return fmt.Errorf("create lib: %w", err)
	}
	fmt.Println("Created lib/")

	gitignorePath := filepath.Join(dir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
			return fmt.Errorf("create .gitignore: %w", err)
		}
		fmt.Println("Created .gitignore")
	} else {
		fmt.Println(".gitignore already exists")
	}

	agentsPath := filepath.Join(dir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		content := strings.ReplaceAll(agentsMDContent, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(agentsPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("create AGENTS.md: %w", err)
		}
		fmt.Println("Created AGENTS.md")
	} else {
		fmt.Println("AGENTS.md already exists")
	}

	claudePath := filepath.Join(dir, "CLAUDE.md")
	if _, err := os.Lstat(claudePath); os.IsNotExist(err) {
		if err := os.Symlink("AGENTS.md", claudePath); err != nil {
			return fmt.Errorf("create CLAUDE.md symlink: %w", err)
		}
		fmt.Println("Created CLAUDE.md -> AGENTS.md")
	} else {
		fmt.Println("CLAUDE.md already exists")
	}

	fmt.Println("\nProject initialized! Next steps:")
	fmt.Println("  - Add dependencies: sai add <groupId:artifactId:version>")
	fmt.Printf("  - Add source files to src/%s/core/ and src/%s/main/\n", projectID, projectID)
	return nil
}

func isValidJavaIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !isJavaIdentifierStart(r) {
				return false
			}
		} else {
			if !isJavaIdentifierPart(r) {
				return false
			}
		}
	}
	return !isJavaKeyword(s)
}

func isJavaIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '$'
}

func isJavaIdentifierPart(r rune) bool {
	return isJavaIdentifierStart(r) || (r >= '0' && r <= '9')
}

func isJavaKeyword(s string) bool {
	keywords := map[string]bool{
		"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true,
		"case": true, "catch": true, "char": true, "class": true, "const": true,
		"continue": true, "default": true, "do": true, "double": true, "else": true,
		"enum": true, "extends": true, "final": true, "finally": true, "float": true,
		"for": true, "goto": true, "if": true, "implements": true, "import": true,
		"instanceof": true, "int": true, "interface": true, "long": true, "native": true,
		"new": true, "package": true, "private": true, "protected": true, "public": true,
		"return": true, "short": true, "static": true, "strictfp": true, "super": true,
		"switch": true, "synchronized": true, "this": true, "throw": true, "throws": true,
		"transient": true, "try": true, "void": true, "volatile": true, "while": true,
		"true": true, "false": true, "null": true,
	}
	return keywords[s]
}

const gitignoreContent = `out/
mlib/
`

const agentsMDContent = `# Sai Java Project

This project uses **sai** instead of Maven or Gradle for Java development.

## Commands

### Adding Dependencies

` + "```" + `bash
sai add <groupId:artifactId:version>
# Example: sai add com.google.guava:guava:32.1.3-jre
` + "```" + `

Dependencies are downloaded to ` + "`lib/`" + `.

### Getting the Classpath

` + "```" + `bash
sai classpath
` + "```" + `

Returns a colon-separated list of JAR paths for use with ` + "`javac`" + ` or ` + "`java`" + `.

### Formatting

` + "```" + `bash
sai fmt src/Main.java        # Print formatted output
sai fmt -w src/Main.java     # Overwrite file in place
` + "```" + `

## Project Structure

` + "```" + `
src/
  {{PROJECT_ID}}/
    core/                    # Logic module
      module-info.java
    main/                    # Entrypoints module
      module-info.java
lib/                         # Dependencies (JARs)
out/                         # Compiled classes
mlib/                        # Modular JARs for jlink/jpackage
` + "```" + `

## Build Commands

### Compile

` + "```" + `sh
mkdir -p out/{{PROJECT_ID}}.core out/{{PROJECT_ID}}.main mlib

# Compile core
javac -p lib -d out/{{PROJECT_ID}}.core \
  src/{{PROJECT_ID}}/core/module-info.java \
  src/{{PROJECT_ID}}/core/*.java

# Compile main
javac -p lib:out -d out/{{PROJECT_ID}}.main \
  src/{{PROJECT_ID}}/main/module-info.java \
  src/{{PROJECT_ID}}/main/*.java
` + "```" + `

### Create Modular JARs

` + "```" + `sh
# Core module
jar --create --file=mlib/{{PROJECT_ID}}.core.jar -C out/{{PROJECT_ID}}.core .

# Main module with entrypoint
jar --create --file=mlib/{{PROJECT_ID}}.main.jar --main-class={{PROJECT_ID}}.main.Main -C out/{{PROJECT_ID}}.main .

# Include dependencies
cp lib/*.jar mlib/
` + "```" + `

### Run During Development

` + "```" + `sh
# Using module path (after creating modular JARs)
java -p mlib -m {{PROJECT_ID}}.main

# Using exploded classes
java -p lib:out -m {{PROJECT_ID}}.main/{{PROJECT_ID}}.main.Main
` + "```" + `

## Notes for AI Agents

- Do NOT use Maven (` + "`mvn`" + `) or Gradle (` + "`gradle`" + `, ` + "`./gradlew`" + `)
- Use ` + "`sai add`" + ` to add dependencies
- Use ` + "`sai classpath`" + ` to get the classpath for compilation/execution
- Use ` + "`sai fmt -w`" + ` to format Java files
- The ` + "`lib/`" + ` directory contains all JAR dependencies
`

func runClasspath(libDir string) error {
	if _, err := os.Stat("pom.xml"); err == nil {
		return runClasspathFromPOM()
	}
	return runClasspathFromLib(libDir)
}

func runClasspathFromPOM() error {
	data, err := os.ReadFile("pom.xml")
	if err != nil {
		return fmt.Errorf("read pom.xml: %w", err)
	}

	var project pom.Project
	if err := xml.Unmarshal(data, &project); err != nil {
		return fmt.Errorf("parse pom.xml: %w", err)
	}

	fetcher := pom.NewMavenFetcher()
	resolver := pom.NewResolver(fetcher)
	deps, err := resolver.Resolve(&project)
	if err != nil {
		return fmt.Errorf("resolve dependencies: %w", err)
	}

	var paths []string
	for _, dep := range deps {
		if dep.Type != "" && dep.Type != "jar" {
			continue
		}
		jarPath := fetcher.JarURL(dep.GroupID, dep.ArtifactID, dep.Version, dep.Classifier)
		paths = append(paths, jarPath)
	}

	fmt.Println(strings.Join(paths, ":"))
	return nil
}

func runClasspathFromLib(libDir string) error {
	entries, err := os.ReadDir(libDir)
	if err != nil {
		return fmt.Errorf("read lib directory %s: %w", libDir, err)
	}

	var paths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".jar" {
			paths = append(paths, filepath.Join(libDir, entry.Name()))
		}
	}

	fmt.Println(strings.Join(paths, ":"))
	return nil
}

func runAdd(coord string, libDir string) error {
	groupID, artifactID, version, classifier, err := pom.ParseCoordinate(coord)
	if err != nil {
		return err
	}

	fetcher := pom.NewMavenFetcher()
	fmt.Printf("Using repository: %s\n", fetcher.RepoURL)
	fmt.Printf("Resolving %s:%s:%s\n", groupID, artifactID, version)

	project, err := fetcher.FetchPOM(groupID, artifactID, version)
	var deps []pom.ResolvedDependency

	if err != nil {
		fmt.Printf("Warning: could not fetch POM (%v), downloading JAR only\n", err)
		deps = []pom.ResolvedDependency{{
			GroupID:    groupID,
			ArtifactID: artifactID,
			Version:    version,
			Classifier: classifier,
		}}
	} else {
		project.GroupID = groupID
		project.ArtifactID = artifactID
		project.Version = version
		project.Dependencies = append([]pom.Dependency{{
			GroupID:    groupID,
			ArtifactID: artifactID,
			Version:    version,
			Classifier: classifier,
			Scope:      "compile",
		}}, project.Dependencies...)

		resolver := pom.NewResolver(fetcher)
		deps, err = resolver.Resolve(project)
		if err != nil {
			return fmt.Errorf("resolve dependencies: %w", err)
		}
	}

	fmt.Printf("Found %d dependencies\n", len(deps))

	downloaded := make(map[string]bool)
	var errors []string
	for _, dep := range deps {
		if dep.Type != "" && dep.Type != "jar" {
			continue
		}

		key := fmt.Sprintf("%s:%s:%s:%s", dep.GroupID, dep.ArtifactID, dep.Version, dep.Classifier)
		if downloaded[key] {
			continue
		}
		downloaded[key] = true

		fmt.Printf("  Downloading %s:%s:%s", dep.GroupID, dep.ArtifactID, dep.Version)
		if dep.Classifier != "" {
			fmt.Printf(":%s", dep.Classifier)
		}
		fmt.Print("...")

		path, err := fetcher.DownloadJar(dep.GroupID, dep.ArtifactID, dep.Version, dep.Classifier, libDir)
		if err != nil {
			fmt.Printf(" FAILED: %v\n", err)
			errors = append(errors, fmt.Sprintf("%s:%s:%s: %v", dep.GroupID, dep.ArtifactID, dep.Version, err))
			continue
		}
		fmt.Printf(" OK (%s)\n", filepath.Base(path))
	}

	fmt.Printf("\nDownloaded %d JARs to %s/\n", len(downloaded)-len(errors), libDir)
	if len(errors) > 0 {
		fmt.Printf("Failed to download %d JARs:\n", len(errors))
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
	}

	return nil
}
