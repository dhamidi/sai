package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dhamidi/sai/project"
	"github.com/spf13/cobra"
)

func newProjectCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "project",
		Short: "Show project structure",
		Long:  `Display the detected project structure including all modules and their dependencies.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProject(jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output as JSON")

	return cmd
}

type projectJSON struct {
	ID               string       `json:"id"`
	RootDir          string       `json:"rootDir"`
	SrcDir           string       `json:"srcDir"`
	OutDir           string       `json:"outDir"`
	LibDir           string       `json:"libDir"`
	Modules          []moduleJSON `json:"modules"`
	CompilationOrder []string     `json:"compilationOrder"`
}

type moduleJSON struct {
	Name         string   `json:"name"`
	FullName     string   `json:"fullName"`
	SrcDir       string   `json:"srcDir"`
	OutDir       string   `json:"outDir"`
	ModuleInfo   string   `json:"moduleInfo"`
	Dependencies []string `json:"dependencies"`
	FileCount    int      `json:"fileCount"`
}

func runProject(jsonOutput bool) error {
	proj, err := project.Load()
	if err != nil {
		return err
	}

	if jsonOutput {
		return runProjectJSON(proj)
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

func runProjectJSON(proj *project.Project) error {
	modules := make([]moduleJSON, len(proj.Modules))
	for i, mod := range proj.Modules {
		fileCount := 1 // module-info.java
		files, err := mod.JavaFiles(false)
		if err == nil {
			fileCount += len(files)
		}

		// Convert short dependency names to full names
		deps := make([]string, len(mod.Dependencies))
		for j, d := range mod.Dependencies {
			deps[j] = proj.ID + "." + d
		}

		modules[i] = moduleJSON{
			Name:         mod.Name,
			FullName:     mod.FullName(),
			SrcDir:       mod.SrcDir,
			OutDir:       mod.OutDir,
			ModuleInfo:   mod.ModuleInfo,
			Dependencies: deps,
			FileCount:    fileCount,
		}
	}

	compOrder := make([]string, 0, len(proj.Modules))
	for _, mod := range proj.ModulesInOrder() {
		compOrder = append(compOrder, mod.FullName())
	}

	out := projectJSON{
		ID:               proj.ID,
		RootDir:          proj.RootDir,
		SrcDir:           proj.SrcDir,
		OutDir:           proj.OutDir,
		LibDir:           proj.LibDir,
		Modules:          modules,
		CompilationOrder: compOrder,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
