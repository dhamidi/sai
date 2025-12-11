package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newCompileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "Compile the Java project",
		Long: `Compile the Java project using javac.

This command:
  - Creates output directories (out/<project>.core, out/<project>.main)
  - Compiles the core module first
  - Compiles the main module (which depends on core)

The project identifier is detected from the src/ directory structure.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompile()
		},
	}
	return cmd
}

func runCompile() error {
	projectID, err := detectProjectID()
	if err != nil {
		return err
	}

	coreOutDir := filepath.Join("out", projectID+".core")
	mainOutDir := filepath.Join("out", projectID+".main")
	if err := os.MkdirAll(coreOutDir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", coreOutDir, err)
	}
	if err := os.MkdirAll(mainOutDir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", mainOutDir, err)
	}

	coreSrcDir := filepath.Join("src", projectID, "core")
	coreModuleInfo := filepath.Join(coreSrcDir, "module-info.java")

	coreJavaFiles, err := filepath.Glob(filepath.Join(coreSrcDir, "*.java"))
	if err != nil {
		return fmt.Errorf("glob core java files: %w", err)
	}

	coreArgs := []string{"-p", "lib", "-d", coreOutDir, coreModuleInfo}
	for _, f := range coreJavaFiles {
		if f != coreModuleInfo {
			coreArgs = append(coreArgs, f)
		}
	}

	fmt.Printf("Compiling %s.core...\n", projectID)
	javacCore := exec.Command("javac", coreArgs...)
	javacCore.Stdout = os.Stdout
	javacCore.Stderr = os.Stderr
	if err := javacCore.Run(); err != nil {
		return fmt.Errorf("compile core: %w", err)
	}

	mainSrcDir := filepath.Join("src", projectID, "main")
	mainModuleInfo := filepath.Join(mainSrcDir, "module-info.java")

	mainJavaFiles, err := filepath.Glob(filepath.Join(mainSrcDir, "*.java"))
	if err != nil {
		return fmt.Errorf("glob main java files: %w", err)
	}

	mainArgs := []string{"-p", "lib:out", "-d", mainOutDir, mainModuleInfo}
	for _, f := range mainJavaFiles {
		if f != mainModuleInfo {
			mainArgs = append(mainArgs, f)
		}
	}

	fmt.Printf("Compiling %s.main...\n", projectID)
	javacMain := exec.Command("javac", mainArgs...)
	javacMain.Stdout = os.Stdout
	javacMain.Stderr = os.Stderr
	if err := javacMain.Run(); err != nil {
		return fmt.Errorf("compile main: %w", err)
	}

	fmt.Println("Compilation successful!")
	return nil
}

func detectProjectID() (string, error) {
	entries, err := os.ReadDir("src")
	if err != nil {
		return "", fmt.Errorf("read src directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			corePath := filepath.Join("src", entry.Name(), "core")
			mainPath := filepath.Join("src", entry.Name(), "main")

			coreExists := false
			mainExists := false

			if info, err := os.Stat(corePath); err == nil && info.IsDir() {
				coreExists = true
			}
			if info, err := os.Stat(mainPath); err == nil && info.IsDir() {
				mainExists = true
			}

			if coreExists && mainExists {
				return entry.Name(), nil
			}
		}
	}

	return "", fmt.Errorf("could not detect project ID: no src/<project>/{core,main} structure found")
}
