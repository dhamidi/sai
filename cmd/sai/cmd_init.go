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

var coreModuleInfoTemplate = `module {{PROJECT_ID}}.core {
    exports {{PROJECT_ID}}.core;
}
`

var mainModuleInfoTemplate = `module {{PROJECT_ID}}.main {
    requires {{PROJECT_ID}}.core;
}
`

var helloJavaTemplate = `package {{PROJECT_ID}}.core;

public class Hello {
    public static String greet(String name) {
        return "Hello, " + name + "!";
    }
}
`

var cliJavaTemplate = `package {{PROJECT_ID}}.main;

import {{PROJECT_ID}}.core.Hello;

public class Cli {
    public static void main(String[] args) {
        System.out.println(Hello.greet("world"));
    }
}
`

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

	fmt.Println("\nProject initialized! Next steps:")
	fmt.Println("  - Compile: sai compile")
	fmt.Println("  - Run: sai run")
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
