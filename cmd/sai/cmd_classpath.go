package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhamidi/sai/pom"
	"github.com/spf13/cobra"
)

func newClasspathCmd() *cobra.Command {
	var cpLibDir string

	cmd := &cobra.Command{
		Use:   "classpath",
		Short: "Print the classpath from pom.xml or lib/ directory",
		Long: `Print the classpath as a colon-separated list of JAR paths.

If pom.xml exists in the current directory, dependencies are resolved
from it and printed as Maven repository paths (requires downloading).

Otherwise, all .jar files in the lib/ directory (or specified via -l)
are listed.

Examples:
  sai classpath              # Use pom.xml if present, else lib/
  sai classpath -l deps/     # Use deps/ directory`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClasspath(cpLibDir)
		},
	}

	cmd.Flags().StringVarP(&cpLibDir, "lib", "l", "lib", "directory containing JAR files")

	return cmd
}

func runClasspath(libDir string) error {
	if _, err := os.Stat("pom.xml"); err == nil {
		return runClasspathFromPOM()
	}
	return runClasspathFromLib(libDir)
}

func runClasspathFromPOM() error {
	data, err := os.ReadFile("pom.xml")
	if err != nil {
		return fmt.Errorf("read pom.xml: %w", err)
	}

	var project pom.Project
	if err := xml.Unmarshal(data, &project); err != nil {
		return fmt.Errorf("parse pom.xml: %w", err)
	}

	fetcher := pom.NewMavenFetcher()
	resolver := pom.NewResolver(fetcher)
	deps, err := resolver.Resolve(&project)
	if err != nil {
		return fmt.Errorf("resolve dependencies: %w", err)
	}

	var paths []string
	for _, dep := range deps {
		if dep.Type != "" && dep.Type != "jar" {
			continue
		}
		jarPath := fetcher.JarURL(dep.GroupID, dep.ArtifactID, dep.Version, dep.Classifier)
		paths = append(paths, jarPath)
	}

	fmt.Println(strings.Join(paths, ":"))
	return nil
}

func runClasspathFromLib(libDir string) error {
	entries, err := os.ReadDir(libDir)
	if err != nil {
		return fmt.Errorf("read lib directory %s: %w", libDir, err)
	}

	var paths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".jar" {
			paths = append(paths, filepath.Join(libDir, entry.Name()))
		}
	}

	fmt.Println(strings.Join(paths, ":"))
	return nil
}
