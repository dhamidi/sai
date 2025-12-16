package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/dhamidi/sai/project"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var verbose bool

	// Load project and discover entrypoints at command creation time
	// for validation and help text
	entrypoints := discoverEntrypoints()
	entrypointMap := make(map[string]project.Entrypoint)
	var slugs []string
	for _, ep := range entrypoints {
		entrypointMap[ep.Slug] = ep
		slugs = append(slugs, ep.Slug)
	}
	sort.Strings(slugs)

	// Build usage string showing available entrypoints
	usageStr := "run [entrypoint] [args...]"

	// Build the valid entrypoints help text
	var entrypointsHelp string
	if len(slugs) > 0 {
		entrypointsHelp = fmt.Sprintf("\n\nAvailable entrypoints: %s\nDefault: cli (if available)", strings.Join(slugs, ", "))
	}

	cmd := &cobra.Command{
		Use:   usageStr,
		Short: "Run the Java project",
		Long: `Run the Java project using java.

This command runs a class with a main method from the main module.
If no entrypoint is specified, it defaults to 'cli'.

Any additional arguments are passed to the Java program.

The project must be compiled first with 'sai compile'.` + entrypointsHelp,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				// Complete entrypoint names
				return slugs, cobra.ShellCompDirectiveNoFileComp
			}
			// After entrypoint, no completion
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRun(args, verbose, entrypointMap, slugs)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print exact command being executed")
	cmd.Flags().SetInterspersed(false) // Stop parsing flags after entrypoint

	return cmd
}

// discoverEntrypoints loads the project and finds all classes with main methods.
func discoverEntrypoints() []project.Entrypoint {
	proj, err := project.Load()
	if err != nil {
		return nil
	}

	mainMod := proj.Module("main")
	if mainMod == nil {
		return nil
	}

	entrypoints, err := mainMod.FindEntrypoints()
	if err != nil {
		return nil
	}

	return entrypoints
}

func runRun(args []string, verbose bool, entrypointMap map[string]project.Entrypoint, validSlugs []string) error {
	proj, err := project.Load()
	if err != nil {
		return err
	}

	mainMod := proj.Module("main")
	if mainMod == nil {
		return fmt.Errorf("no main module found in project %s", proj.ID)
	}

	// Determine which entrypoint to run and what args to pass
	var entrypointSlug string
	var programArgs []string

	if len(args) == 0 {
		// No args: default to "cli"
		entrypointSlug = "cli"
	} else {
		// Check if first arg is a known entrypoint
		if _, ok := entrypointMap[args[0]]; ok {
			entrypointSlug = args[0]
			programArgs = args[1:]
		} else {
			// First arg is not a valid entrypoint
			if len(validSlugs) == 0 {
				return fmt.Errorf("unknown entrypoint %q (no entrypoints discovered in project)", args[0])
			}
			return fmt.Errorf("unknown entrypoint %q\n\nAvailable entrypoints: %s", args[0], strings.Join(validSlugs, ", "))
		}
	}

	// Look up the entrypoint
	ep, ok := entrypointMap[entrypointSlug]
	if !ok {
		if len(validSlugs) == 0 {
			return fmt.Errorf("no entrypoints discovered in project (ensure main module has classes with main methods)")
		}
		return fmt.Errorf("entrypoint %q not found\n\nAvailable entrypoints: %s", entrypointSlug, strings.Join(validSlugs, ", "))
	}

	javaArgs := []string{
		"--enable-preview",
		"-p", proj.ModulePath(true),
		"-m", mainMod.FullName() + "/" + ep.FullName,
	}
	javaArgs = append(javaArgs, programArgs...)

	if verbose {
		fmt.Printf("+ java %s\n", formatArgs(javaArgs))
	}

	javaCmd := exec.Command("java", javaArgs...)
	javaCmd.Stdout = os.Stdout
	javaCmd.Stderr = os.Stderr
	javaCmd.Stdin = os.Stdin

	return javaCmd.Run()
}
