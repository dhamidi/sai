package main

import (
	"fmt"
	"strings"

	"github.com/dhamidi/sai/project"
	"github.com/spf13/cobra"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Show project structure",
		Long:  `Display the detected project structure including all modules and their dependencies.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProject()
		},
	}

	return cmd
}

func runProject() error {
	proj, err := project.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Project: %s\n", proj.ID)
	fmt.Printf("Root:    %s\n", proj.RootDir)
	fmt.Printf("Source:  %s\n", proj.SrcDir)
	fmt.Printf("Output:  %s\n", proj.OutDir)
	fmt.Printf("Libs:    %s\n", proj.LibDir)
	fmt.Printf("\nModules:\n")

	for _, mod := range proj.Modules {
		fmt.Printf("  %s\n", mod.FullName())
		fmt.Printf("    src:  %s\n", mod.SrcDir)
		fmt.Printf("    out:  %s\n", mod.OutDir)

		if len(mod.Dependencies) > 0 {
			deps := make([]string, len(mod.Dependencies))
			for i, d := range mod.Dependencies {
				deps[i] = proj.ID + "." + d
			}
			fmt.Printf("    deps: %s\n", strings.Join(deps, ", "))
		}

		files, err := mod.JavaFiles(false)
		if err != nil {
			fmt.Printf("    files: error: %v\n", err)
		} else {
			fmt.Printf("    files: %d java files\n", len(files)+1) // +1 for module-info.java
		}
	}

	fmt.Printf("\nCompilation order:\n")
	for i, mod := range proj.ModulesInOrder() {
		fmt.Printf("  %d. %s\n", i+1, mod.FullName())
	}

	return nil
}
