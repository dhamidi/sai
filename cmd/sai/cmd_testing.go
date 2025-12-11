package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func newTestCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "test [-- args...]",
		Short: "Run tests using JUnit",
		Long: `Run tests using JUnit Platform Console Standalone.

This command:
  - Compiles the project (core and main modules)
  - Compiles the test module
  - Runs JUnit with the Jupiter engine

Any arguments after -- are forwarded to the JUnit console runner.

Examples:
  sai test                           # Run all tests
  sai test -- --select-class=MyTest  # Run specific test class
  sai test -- --fail-if-no-tests     # Fail if no tests found`,
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

	// Find junit-platform-console-standalone jar
	junitJar, err := findJUnitConsoleLauncher()
	if err != nil {
		return err
	}

	// Build classpath for test execution
	classpath := fmt.Sprintf("out/%s.core:out/%s.test:lib/*", projectID, projectID)

	// Run JUnit
	javaArgs := []string{
		"-jar", junitJar,
		"execute",
		"--disable-banner",
		"-cp", classpath,
		"--scan-classpath",
		"-e", "junit-jupiter",
	}
	javaArgs = append(javaArgs, junitArgs...)

	fmt.Println("\nRunning tests...")
	if verbose {
		fmt.Printf("+ java %s\n", formatArgs(javaArgs))
	}

	javaCmd := exec.Command("java", javaArgs...)
	javaCmd.Stderr = os.Stderr

	stdout, err := javaCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("create stdout pipe: %w", err)
	}

	if err := javaCmd.Start(); err != nil {
		return fmt.Errorf("start junit: %w", err)
	}

	summary := filterJUnitOutput(stdout, os.Stdout)

	exitErr := javaCmd.Wait()

	// Print summary line
	if summary.failed > 0 {
		fmt.Printf("\n%d/%d tests failed\n", summary.failed, summary.total)
	} else if summary.total > 0 {
		fmt.Printf("\n%d tests passed\n", summary.total)
	}

	if exitErr != nil {
		os.Exit(1)
	}
	return nil
}

type testSummary struct {
	total  int
	failed int
}

func filterJUnitOutput(r io.Reader, w io.Writer) testSummary {
	scanner := bufio.NewScanner(r)
	summary := testSummary{}
	inSummary := false

	testsFoundRe := regexp.MustCompile(`\[\s*(\d+) tests found\s*\]`)
	testsFailedRe := regexp.MustCompile(`\[\s*(\d+) tests failed\s*\]`)

	for scanner.Scan() {
		line := scanner.Text()

		// Detect start of summary block
		if strings.HasPrefix(line, "Test run finished") {
			inSummary = true
			continue
		}

		if inSummary {
			// Parse summary lines
			if m := testsFoundRe.FindStringSubmatch(line); m != nil {
				summary.total, _ = strconv.Atoi(m[1])
			}
			if m := testsFailedRe.FindStringSubmatch(line); m != nil {
				summary.failed, _ = strconv.Atoi(m[1])
			}
			continue
		}

		// Skip sponsorship message
		if strings.Contains(line, "Thanks for using JUnit") {
			continue
		}

		fmt.Fprintln(w, line)
	}

	return summary
}

func findJUnitConsoleLauncher() (string, error) {
	matches, err := filepath.Glob("lib/junit-platform-console-standalone-*.jar")
	if err != nil {
		return "", fmt.Errorf("search for junit console launcher: %w", err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("junit-platform-console-standalone not found in lib/; run: sai add org.junit.platform:junit-platform-console-standalone:1.13.0")
	}
	return matches[0], nil
}
