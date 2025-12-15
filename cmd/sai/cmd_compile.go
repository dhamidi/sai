package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dhamidi/sai/project"
	"github.com/spf13/cobra"
)

func newCompileCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "compile",
		Short: "Compile the Java project",
		Long: `Compile the Java project using javac.

This command:
  - Creates output directories for each module
  - Compiles modules in dependency order (core, then main)

The project structure is detected from src/<project>/<module>/module-info.java files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompile(verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print exact commands being executed")

	return cmd
}

func runCompile(verbose bool) error {
	proj, err := project.Load()
	if err != nil {
		return err
	}

	// Compile modules in dependency order
	for _, mod := range proj.ModulesInOrder() {
		// Skip test module during normal compilation
		if mod.Name == "test" {
			continue
		}

		if err := compileModule(proj, mod, verbose); err != nil {
			return err
		}
	}

	fmt.Println("Compilation successful!")
	return nil
}

func compileModule(proj *project.Project, mod *project.Module, verbose bool) error {
	if err := mod.EnsureOutDir(); err != nil {
		return err
	}

	javaFiles, err := mod.JavaFiles(true)
	if err != nil {
		return err
	}

	args := []string{
		"--enable-preview",
		"--source", "25",
		"-p", proj.ModulePath(true),
		"-d", mod.OutDir,
	}
	args = append(args, javaFiles...)

	fmt.Printf("Compiling %s...\n", mod.FullName())
	if verbose {
		fmt.Printf("+ javac %s\n", formatArgs(args))
	}

	javac := exec.Command("javac", args...)
	javac.Stdout = os.Stdout
	javac.Stderr = os.Stderr
	if err := javac.Run(); err != nil {
		return fmt.Errorf("compile %s: %w", mod.Name, err)
	}

	return nil
}

func formatArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		if needsQuoting(arg) {
			result += fmt.Sprintf("%q", arg)
		} else {
			result += arg
		}
	}
	return result
}

func needsQuoting(s string) bool {
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '"' || r == '\'' || r == '\\' {
			return true
		}
	}
	return false
}
