package main

import (
	"fmt"

	"github.com/dhamidi/sai/project"
	"github.com/spf13/cobra"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Show project structure",
		Long:  `Display the detected project structure including all modules.`,
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
		fmt.Printf("    src: %s\n", mod.SrcDir)
		fmt.Printf("    out: %s\n", mod.OutDir)

		files, err := mod.JavaFiles(false)
		if err != nil {
			fmt.Printf("    files: error: %v\n", err)
		} else {
			fmt.Printf("    files: %d java files\n", len(files)+1) // +1 for module-info.java
		}
	}

	return nil
}
