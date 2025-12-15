package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed init/gitignore
var gitignoreContent string

//go:embed init/AGENTS.md
var agentsMDContent string

//go:embed init/core-module-info.java.tmpl
var coreModuleInfoTemplate string

//go:embed init/main-module-info.java.tmpl
var mainModuleInfoTemplate string

//go:embed init/Hello.java.tmpl
var helloJavaTemplate string

//go:embed init/Cli.java.tmpl
var cliJavaTemplate string

//go:embed init/test-module-info.java.tmpl
var testModuleInfoTemplate string

//go:embed init/HelloTest.java.tmpl
var helloTestJavaTemplate string

//go:embed init/TestRunner.java.tmpl
var testRunnerJavaTemplate string

//go:embed init/LinePerTestReporter.java.tmpl
var linePerTestReporterJavaTemplate string

func newInitCmd() *cobra.Command {
	var projectID string

	cmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize a new sai Java project",
		Long: `Initialize a new sai Java project.

If a directory is provided, creates it and initializes the project there.
Otherwise, initializes in the current directory.

The project identifier defaults to the directory basename if it's a valid
Java identifier. Use -p to override.

This command:
  - Ensures a git repository exists
  - Creates the modular source structure: src/<project-id>/{core,main,test}/
  - Creates .gitignore with out/, mlib/, and dist/
  - Creates AGENTS.md with sai workflow instructions
  - Creates CLAUDE.md as a symlink to AGENTS.md
  - Installs JUnit 5 dependencies for testing

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

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "project identifier (defaults to directory name)")

	return cmd
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

	testPath := filepath.Join(dir, "src", projectID, "test")
	if err := os.MkdirAll(testPath, 0755); err != nil {
		return fmt.Errorf("create src/%s/test: %w", projectID, err)
	}
	fmt.Printf("Created src/%s/test/\n", projectID)

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

	coreModuleInfo := filepath.Join(corePath, "module-info.java")
	if _, err := os.Stat(coreModuleInfo); os.IsNotExist(err) {
		content := strings.ReplaceAll(coreModuleInfoTemplate, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(coreModuleInfo, []byte(content), 0644); err != nil {
			return fmt.Errorf("create core module-info.java: %w", err)
		}
		fmt.Printf("Created src/%s/core/module-info.java\n", projectID)
	}

	helloJava := filepath.Join(corePath, "Hello.java")
	if _, err := os.Stat(helloJava); os.IsNotExist(err) {
		content := strings.ReplaceAll(helloJavaTemplate, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(helloJava, []byte(content), 0644); err != nil {
			return fmt.Errorf("create Hello.java: %w", err)
		}
		fmt.Printf("Created src/%s/core/Hello.java\n", projectID)
	}

	mainModuleInfo := filepath.Join(mainPath, "module-info.java")
	if _, err := os.Stat(mainModuleInfo); os.IsNotExist(err) {
		content := strings.ReplaceAll(mainModuleInfoTemplate, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(mainModuleInfo, []byte(content), 0644); err != nil {
			return fmt.Errorf("create main module-info.java: %w", err)
		}
		fmt.Printf("Created src/%s/main/module-info.java\n", projectID)
	}

	cliJava := filepath.Join(mainPath, "Cli.java")
	if _, err := os.Stat(cliJava); os.IsNotExist(err) {
		content := strings.ReplaceAll(cliJavaTemplate, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(cliJava, []byte(content), 0644); err != nil {
			return fmt.Errorf("create Cli.java: %w", err)
		}
		fmt.Printf("Created src/%s/main/Cli.java\n", projectID)
	}

	testModuleInfo := filepath.Join(testPath, "module-info.java")
	if _, err := os.Stat(testModuleInfo); os.IsNotExist(err) {
		content := strings.ReplaceAll(testModuleInfoTemplate, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(testModuleInfo, []byte(content), 0644); err != nil {
			return fmt.Errorf("create test module-info.java: %w", err)
		}
		fmt.Printf("Created src/%s/test/module-info.java\n", projectID)
	}

	helloTestJava := filepath.Join(testPath, "HelloTest.java")
	if _, err := os.Stat(helloTestJava); os.IsNotExist(err) {
		content := strings.ReplaceAll(helloTestJavaTemplate, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(helloTestJava, []byte(content), 0644); err != nil {
			return fmt.Errorf("create HelloTest.java: %w", err)
		}
		fmt.Printf("Created src/%s/test/HelloTest.java\n", projectID)
	}

	testRunnerJava := filepath.Join(testPath, "TestRunner.java")
	if _, err := os.Stat(testRunnerJava); os.IsNotExist(err) {
		content := strings.ReplaceAll(testRunnerJavaTemplate, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(testRunnerJava, []byte(content), 0644); err != nil {
			return fmt.Errorf("create TestRunner.java: %w", err)
		}
		fmt.Printf("Created src/%s/test/TestRunner.java\n", projectID)
	}

	linePerTestReporterJava := filepath.Join(testPath, "LinePerTestReporter.java")
	if _, err := os.Stat(linePerTestReporterJava); os.IsNotExist(err) {
		content := strings.ReplaceAll(linePerTestReporterJavaTemplate, "{{PROJECT_ID}}", projectID)
		if err := os.WriteFile(linePerTestReporterJava, []byte(content), 0644); err != nil {
			return fmt.Errorf("create LinePerTestReporter.java: %w", err)
		}
		fmt.Printf("Created src/%s/test/LinePerTestReporter.java\n", projectID)
	}

	// Install JUnit 5 dependencies
	fmt.Println("\nInstalling JUnit 5 dependencies...")
	junitDeps := []string{
		"org.junit.jupiter:junit-jupiter:5.13.0",
		"org.junit.platform:junit-platform-launcher:1.13.0",
	}
	for _, dep := range junitDeps {
		cmd := exec.Command("sai", "add", dep)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("install %s: %w", dep, err)
		}
	}

	fmt.Println("\nProject initialized! Next steps:")
	fmt.Println("  - Compile: sai compile")
	fmt.Println("  - Run: sai run")
	fmt.Println("  - Test: sai test")
	fmt.Println("  - Add dependencies: sai add <groupId:artifactId:version>")
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
