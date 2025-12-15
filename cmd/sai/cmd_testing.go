package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dhamidi/sai/project"
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

	proj, err := project.Load()
	if err != nil {
		return err
	}

	testMod := proj.Module("test")
	if testMod == nil {
		return fmt.Errorf("no test module found in project %s", proj.ID)
	}

	// Compile test module
	if err := compileModule(proj, testMod, verbose); err != nil {
		return err
	}

	// Run tests using the programmatic TestRunner
	// TODO: make test runner class configurable
	testRunnerClass := testMod.FullName() + ".TestRunner"

	javaArgs := []string{
		"-p", proj.ModulePath(true),
		"--add-modules", "org.junit.jupiter.engine",
		"-m", testMod.FullName() + "/" + testRunnerClass,
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
