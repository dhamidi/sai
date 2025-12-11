package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newTestCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "test [-- args...]",
		Short: "Run tests using JUnit",
		Long: `Run tests using JUnit with a custom test runner.

This command:
  - Compiles the project (core and main modules)
  - Compiles the test module
  - Runs the TestRunner which uses JUnit with a custom reporter

Examples:
  sai test                           # Run all tests
  sai test -v                        # Show exact commands`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTest(args, verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print exact commands being executed")

	return cmd
}

func runTest(junitArgs []string, verbose bool) error {
	// First compile the main project
	if err := runCompile(verbose); err != nil {
		return err
	}

	projectID, err := detectProjectID()
	if err != nil {
		return err
	}

	// Check if test module exists
	testSrcDir := filepath.Join("src", projectID, "test")
	if _, err := os.Stat(testSrcDir); os.IsNotExist(err) {
		return fmt.Errorf("no test module found at %s", testSrcDir)
	}

	// Create test output directory
	testOutDir := filepath.Join("out", projectID+".test")
	if err := os.MkdirAll(testOutDir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", testOutDir, err)
	}

	// Compile test module
	testModuleInfo := filepath.Join(testSrcDir, "module-info.java")
	testJavaFiles, err := filepath.Glob(filepath.Join(testSrcDir, "*.java"))
	if err != nil {
		return fmt.Errorf("glob test java files: %w", err)
	}

	testArgs := []string{"-p", "lib:out", "-d", testOutDir, testModuleInfo}
	for _, f := range testJavaFiles {
		if f != testModuleInfo {
			testArgs = append(testArgs, f)
		}
	}

	fmt.Printf("Compiling %s.test...\n", projectID)
	if verbose {
		fmt.Printf("+ javac %s\n", formatArgs(testArgs))
	}
	javacTest := exec.Command("javac", testArgs...)
	javacTest.Stdout = os.Stdout
	javacTest.Stderr = os.Stderr
	if err := javacTest.Run(); err != nil {
		return fmt.Errorf("compile test: %w", err)
	}

	// Run tests using the programmatic TestRunner
	javaArgs := []string{
		"-p", "out:lib",
		"--add-modules", "org.junit.jupiter.engine",
		"-m", fmt.Sprintf("%s.test/%s.test.TestRunner", projectID, projectID),
	}

	fmt.Println("\nRunning tests...")
	if verbose {
		fmt.Printf("+ java %s\n", formatArgs(javaArgs))
	}

	javaCmd := exec.Command("java", javaArgs...)
	javaCmd.Stdout = os.Stdout
	javaCmd.Stderr = os.Stderr

	if err := javaCmd.Run(); err != nil {
		os.Exit(1)
	}
	return nil
}
