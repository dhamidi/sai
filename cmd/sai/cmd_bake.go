package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dhamidi/sai/project"
	"github.com/spf13/cobra"
)

func newBakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bake",
		Short: "Package the project for distribution",
		Long: `Package the Java project for distribution.

Subcommands:
  jar     - Create a custom runtime image using jlink
  native  - Create a native executable using GraalVM native-image`,
	}

	cmd.AddCommand(newBakeJarCmd())
	cmd.AddCommand(newBakeNativeCmd())

	return cmd
}

func newBakeJarCmd() *cobra.Command {
	var (
		mainClass string
		verbose   bool
		output    string
	)

	cmd := &cobra.Command{
		Use:   "jar",
		Short: "Create a custom runtime image using jlink",
		Long: `Create a custom Java runtime image using jlink.

This command:
  1. Compiles the project
  2. Creates modular JARs in mlib/
  3. Uses jlink to create a minimal runtime image

The output is a self-contained directory with a launcher script.

Examples:
  sai bake jar                           # Default output to dist/<project-id>/
  sai bake jar --main-class myapp.main.App
  sai bake jar -o dist/release`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBakeJar(mainClass, output, verbose)
		},
	}

	cmd.Flags().StringVarP(&mainClass, "main-class", "m", "", "main class (default: <project>.main.Cli)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output directory (default: dist/<project-id>)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print exact commands being executed")

	return cmd
}

func runBakeJar(mainClass, output string, verbose bool) error {
	// Step 1: Compile the project
	if err := runCompile(verbose); err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	proj, err := project.Load()
	if err != nil {
		return err
	}

	mainMod := proj.Module("main")
	if mainMod == nil {
		return fmt.Errorf("no main module found in project %s", proj.ID)
	}

	// Set defaults
	if mainClass == "" {
		mainClass = mainMod.FullName() + ".Cli"
	}
	if output == "" {
		output = filepath.Join("dist", proj.ID)
	}

	// Step 2: Create mlib/ directory and modular JARs
	mlibDir := filepath.Join(proj.RootDir, "mlib")
	if err := os.MkdirAll(mlibDir, 0755); err != nil {
		return fmt.Errorf("create mlib: %w", err)
	}

	fmt.Println("Creating modular JARs...")
	for _, mod := range proj.ModulesInOrder() {
		if mod.Name == "test" {
			continue
		}
		if err := createModuleJar(proj, mod, mlibDir, mainClass, verbose); err != nil {
			return err
		}
	}

	// Copy library JARs to mlib/
	if err := copyLibraryJars(proj.LibDir, mlibDir, verbose); err != nil {
		return err
	}

	// Step 3: Run jlink
	fmt.Println("Creating runtime image with jlink...")
	if err := runJlink(proj, mainMod, mlibDir, output, mainClass, verbose); err != nil {
		return err
	}

	fmt.Printf("\nRuntime image created at: %s\n", output)
	fmt.Printf("Run with: %s/bin/%s\n", output, proj.ID)
	return nil
}

func createModuleJar(proj *project.Project, mod *project.Module, mlibDir, mainClass string, verbose bool) error {
	jarName := mod.FullName() + ".jar"
	jarPath := filepath.Join(mlibDir, jarName)

	args := []string{
		"--create",
		"--file=" + jarPath,
	}

	// Set main class for the main module
	if mod.Name == "main" {
		args = append(args, "--main-class="+mainClass)
	}

	args = append(args, "-C", mod.OutDir, ".")

	if verbose {
		fmt.Printf("+ jar %s\n", formatArgs(args))
	}

	jarCmd := exec.Command("jar", args...)
	jarCmd.Stdout = os.Stdout
	jarCmd.Stderr = os.Stderr
	if err := jarCmd.Run(); err != nil {
		return fmt.Errorf("create JAR for %s: %w", mod.Name, err)
	}

	fmt.Printf("  Created %s\n", jarName)
	return nil
}

func copyLibraryJars(libDir, mlibDir string, verbose bool) error {
	entries, err := os.ReadDir(libDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No lib directory is OK
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jar") {
			continue
		}

		src := filepath.Join(libDir, entry.Name())
		dst := filepath.Join(mlibDir, entry.Name())

		if verbose {
			fmt.Printf("+ cp %s %s\n", src, dst)
		}

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("copy %s: %w", entry.Name(), err)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func runJlink(proj *project.Project, mainMod *project.Module, mlibDir, output, mainClass string, verbose bool) error {
	// Clean output directory
	if err := os.RemoveAll(output); err != nil {
		return fmt.Errorf("clean output: %w", err)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Build list of modules to include
	var addModules []string
	for _, mod := range proj.Modules {
		if mod.Name != "test" {
			addModules = append(addModules, mod.FullName())
		}
	}

	args := []string{
		"--module-path", mlibDir,
		"--add-modules", strings.Join(addModules, ","),
		"--launcher", proj.ID + "=" + mainMod.FullName() + "/" + mainClass,
		"--output", output,
		"--no-header-files",
		"--no-man-pages",
		"--strip-debug",
		"--compress=zip-6",
	}

	if verbose {
		fmt.Printf("+ jlink %s\n", formatArgs(args))
	}

	jlink := exec.Command("jlink", args...)
	jlink.Stdout = os.Stdout
	jlink.Stderr = os.Stderr
	return jlink.Run()
}
