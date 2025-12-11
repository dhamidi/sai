package main

import (
	"fmt"
	"path/filepath"

	"github.com/dhamidi/sai/pom"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var libDir string

	cmd := &cobra.Command{
		Use:   "add <groupId:artifactId:version>",
		Short: "Download a Maven dependency and its transitive dependencies to lib/",
		Long: `Download a Maven dependency and its transitive dependencies.

The coordinate format is: groupId:artifactId:version
Or with classifier: groupId:artifactId:classifier:version

Examples:
  sai add com.google.guava:guava:32.1.3-jre
  sai add org.slf4j:slf4j-api:2.0.9

Environment variables:
  MAVEN_REPO_URL - Override the Maven repository URL (default: https://repo1.maven.org/maven2)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(args[0], libDir)
		},
	}

	cmd.Flags().StringVarP(&libDir, "lib", "l", "lib", "directory to download JARs to")

	return cmd
}

func runAdd(coord string, libDir string) error {
	groupID, artifactID, version, classifier, err := pom.ParseCoordinate(coord)
	if err != nil {
		return err
	}

	fetcher := pom.NewMavenFetcher()
	fmt.Printf("Using repository: %s\n", fetcher.RepoURL)
	fmt.Printf("Resolving %s:%s:%s\n", groupID, artifactID, version)

	project, err := fetcher.FetchPOM(groupID, artifactID, version)
	var deps []pom.ResolvedDependency

	if err != nil {
		fmt.Printf("Warning: could not fetch POM (%v), downloading JAR only\n", err)
		deps = []pom.ResolvedDependency{{
			GroupID:    groupID,
			ArtifactID: artifactID,
			Version:    version,
			Classifier: classifier,
		}}
	} else {
		project.GroupID = groupID
		project.ArtifactID = artifactID
		project.Version = version
		project.Dependencies = append([]pom.Dependency{{
			GroupID:    groupID,
			ArtifactID: artifactID,
			Version:    version,
			Classifier: classifier,
			Scope:      "compile",
		}}, project.Dependencies...)

		resolver := pom.NewResolver(fetcher)
		deps, err = resolver.Resolve(project)
		if err != nil {
			return fmt.Errorf("resolve dependencies: %w", err)
		}
	}

	fmt.Printf("Found %d dependencies\n", len(deps))

	downloaded := make(map[string]bool)
	var errors []string
	for _, dep := range deps {
		if dep.Type != "" && dep.Type != "jar" {
			continue
		}

		key := fmt.Sprintf("%s:%s:%s:%s", dep.GroupID, dep.ArtifactID, dep.Version, dep.Classifier)
		if downloaded[key] {
			continue
		}
		downloaded[key] = true

		fmt.Printf("  Downloading %s:%s:%s", dep.GroupID, dep.ArtifactID, dep.Version)
		if dep.Classifier != "" {
			fmt.Printf(":%s", dep.Classifier)
		}
		fmt.Print("...")

		path, err := fetcher.DownloadJar(dep.GroupID, dep.ArtifactID, dep.Version, dep.Classifier, libDir)
		if err != nil {
			fmt.Printf(" FAILED: %v\n", err)
			errors = append(errors, fmt.Sprintf("%s:%s:%s: %v", dep.GroupID, dep.ArtifactID, dep.Version, err))
			continue
		}
		fmt.Printf(" OK (%s)\n", filepath.Base(path))
	}

	fmt.Printf("\nDownloaded %d JARs to %s/\n", len(downloaded)-len(errors), libDir)
	if len(errors) > 0 {
		fmt.Printf("Failed to download %d JARs:\n", len(errors))
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
	}

	return nil
}
