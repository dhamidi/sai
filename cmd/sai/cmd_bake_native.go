package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dhamidi/sai/project"
	"github.com/spf13/cobra"
)

func newBakeNativeCmd() *cobra.Command {
	var (
		verbose     bool
		output      string
		installMise bool
		extraArgs   []string
	)

	// Discover entrypoints for help text and validation
	entrypoints := discoverEntrypoints()
	entrypointMap := make(map[string]project.Entrypoint)
	var slugs []string
	for _, ep := range entrypoints {
		entrypointMap[ep.Slug] = ep
		slugs = append(slugs, ep.Slug)
	}

	var entrypointsHelp string
	if len(slugs) > 0 {
		entrypointsHelp = fmt.Sprintf("\n\nAvailable entrypoints: %s\nDefault: cli", strings.Join(slugs, ", "))
	}

	cmd := &cobra.Command{
		Use:   "native [entrypoint]",
		Short: "Create a native executable using GraalVM native-image",
		Long: `Create a native executable using GraalVM native-image.

This command:
  1. Ensures GraalVM with native-image is available
  2. Compiles the project
  3. Creates modular JARs in mlib/
  4. Runs native-image to create the executable

If GraalVM is not found, use --install to install it via mise.

Examples:
  sai bake native                        # Default: dist/cli
  sai bake native drop-zone-app          # Output: dist/drop-zone-app
  sai bake native --install              # Install GraalVM via mise if needed
  sai bake native cli -o dist/myapp` + entrypointsHelp,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return slugs, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			entrypoint := "cli"
			if len(args) > 0 {
				entrypoint = args[0]
			}
			return runBakeNative(entrypoint, output, verbose, installMise, extraArgs, entrypointMap, slugs)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "output file (default: dist/<entrypoint>)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print exact commands being executed")
	cmd.Flags().BoolVar(&installMise, "install", false, "install GraalVM via mise if not available")
	cmd.Flags().StringArrayVar(&extraArgs, "native-arg", nil, "additional arguments to pass to native-image")

	return cmd
}

func runBakeNative(entrypoint, output string, verbose, installMise bool, extraArgs []string, entrypointMap map[string]project.Entrypoint, validSlugs []string) error {
	// Step 1: Check for native-image
	nativeImagePath, err := findNativeImage()
	if err != nil {
		if !installMise {
			return fmt.Errorf("native-image not found: %w\n\nRun with --install to install GraalVM via mise", err)
		}
		fmt.Println("native-image not found, installing GraalVM via mise...")
		if err := installGraalVMViaMise(verbose); err != nil {
			return fmt.Errorf("install GraalVM: %w", err)
		}
		nativeImagePath, err = findNativeImage()
		if err != nil {
			return fmt.Errorf("native-image still not found after installation: %w", err)
		}
	}

	if verbose {
		fmt.Printf("Using native-image: %s\n", nativeImagePath)
	}

	// Step 2: Compile the project
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

	// Resolve entrypoint to main class
	ep, ok := entrypointMap[entrypoint]
	if !ok {
		if len(validSlugs) == 0 {
			return fmt.Errorf("entrypoint %q not found (no entrypoints discovered in project)", entrypoint)
		}
		return fmt.Errorf("entrypoint %q not found\n\nAvailable entrypoints: %s", entrypoint, strings.Join(validSlugs, ", "))
	}
	mainClass := ep.FullName

	// Set defaults - use entrypoint slug for output name
	if output == "" {
		output = filepath.Join("dist", entrypoint)
	}

	// Ensure dist directory exists
	distDir := filepath.Dir(output)
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("create dist: %w", err)
	}

	// Step 3: Create mlib/ directory and modular JARs
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

	// Step 4: Run native-image
	// Remove output path if it exists (may be a directory from bake jar)
	if err := os.RemoveAll(output); err != nil {
		return fmt.Errorf("clean output path: %w", err)
	}

	fmt.Println("Building native image (this may take several minutes)...")

	// Build list of all modules to include (except test)
	var moduleNames []string
	for _, mod := range proj.ModulesInOrder() {
		if mod.Name == "test" {
			continue
		}
		moduleNames = append(moduleNames, mod.FullName())
	}

	args := []string{
		"--enable-preview",
		"--module-path", mlibDir,
		"--add-modules", strings.Join(moduleNames, ","),
		"--module", mainMod.FullName() + "/" + mainClass,
		"-o", output,
		"--no-fallback",
	}

	// Add extra arguments
	args = append(args, extraArgs...)

	if verbose {
		fmt.Printf("+ %s %s\n", nativeImagePath, formatArgs(args))
	}

	nativeCmd := exec.Command(nativeImagePath, args...)
	nativeCmd.Stdout = os.Stdout
	nativeCmd.Stderr = os.Stderr
	if err := nativeCmd.Run(); err != nil {
		return fmt.Errorf("native-image: %w", err)
	}

	fmt.Printf("\nNative executable created: %s\n", output)
	return nil
}

func findNativeImage() (string, error) {
	// First try PATH
	if path, err := exec.LookPath("native-image"); err == nil {
		return path, nil
	}

	// Check GRAALVM_HOME
	if graalHome := os.Getenv("GRAALVM_HOME"); graalHome != "" {
		path := filepath.Join(graalHome, "bin", "native-image")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Check JAVA_HOME (might be GraalVM)
	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		path := filepath.Join(javaHome, "bin", "native-image")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("native-image not found in PATH, GRAALVM_HOME, or JAVA_HOME")
}

func installGraalVMViaMise(verbose bool) error {
	// Check if mise is available
	if _, err := exec.LookPath("mise"); err != nil {
		return fmt.Errorf("mise not found in PATH; install mise first: https://mise.jdx.dev")
	}

	// Install GraalVM globally (latest version)
	args := []string{"use", "-g", "graalvm@latest"}

	if verbose {
		fmt.Printf("+ mise %s\n", formatArgs(args))
	}

	fmt.Println("Installing GraalVM via mise (this may take a few minutes)...")
	cmd := exec.Command("mise", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
