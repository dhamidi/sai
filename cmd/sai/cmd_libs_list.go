package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newLibsListCmd() *cobra.Command {
	var libDir string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List dependencies in lib/ directory",
		Long: `List all JAR dependencies with their Maven coordinates.

Reads the MANIFEST.MF from each JAR file to extract coordinates.
Falls back to parsing the filename if manifest information is unavailable.

Examples:
  sai libs list
  sai libs list --lib vendor/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLibsList(libDir)
		},
	}

	cmd.Flags().StringVarP(&libDir, "lib", "l", "lib", "directory containing JAR files")

	return cmd
}

type JarInfo struct {
	Path       string
	GroupID    string
	ArtifactID string
	Version    string
}

func runLibsList(libDir string) error {
	entries, err := os.ReadDir(libDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No lib directory found at %s\n", libDir)
			return nil
		}
		return err
	}

	var jars []JarInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jar") {
			continue
		}

		jarPath := filepath.Join(libDir, entry.Name())
		info, err := extractJarInfo(jarPath)
		if err != nil {
			info = parseJarFilename(jarPath)
		}
		jars = append(jars, info)
	}

	if len(jars) == 0 {
		fmt.Printf("No JAR files found in %s\n", libDir)
		return nil
	}

	fmt.Printf("Found %d JARs in %s:\n\n", len(jars), libDir)
	for _, jar := range jars {
		if jar.GroupID != "" && jar.ArtifactID != "" && jar.Version != "" {
			fmt.Printf("  %s:%s:%s\n", jar.GroupID, jar.ArtifactID, jar.Version)
		} else if jar.ArtifactID != "" && jar.Version != "" {
			fmt.Printf("  %s:%s (group unknown)\n", jar.ArtifactID, jar.Version)
		} else {
			fmt.Printf("  %s (coordinates unknown)\n", filepath.Base(jar.Path))
		}
	}

	return nil
}

func extractJarInfo(jarPath string) (JarInfo, error) {
	info := JarInfo{Path: jarPath}

	r, err := zip.OpenReader(jarPath)
	if err != nil {
		return info, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "META-INF/MANIFEST.MF" {
			rc, err := f.Open()
			if err != nil {
				return info, err
			}
			defer rc.Close()

			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "Bundle-SymbolicName:") {
					info.GroupID = strings.TrimSpace(strings.TrimPrefix(line, "Bundle-SymbolicName:"))
					if idx := strings.Index(info.GroupID, ";"); idx != -1 {
						info.GroupID = info.GroupID[:idx]
					}
				} else if strings.HasPrefix(line, "Implementation-Vendor-Id:") {
					info.GroupID = strings.TrimSpace(strings.TrimPrefix(line, "Implementation-Vendor-Id:"))
				} else if strings.HasPrefix(line, "Implementation-Title:") {
					info.ArtifactID = strings.TrimSpace(strings.TrimPrefix(line, "Implementation-Title:"))
				} else if strings.HasPrefix(line, "Bundle-Name:") && info.ArtifactID == "" {
					info.ArtifactID = strings.TrimSpace(strings.TrimPrefix(line, "Bundle-Name:"))
				} else if strings.HasPrefix(line, "Implementation-Version:") {
					info.Version = strings.TrimSpace(strings.TrimPrefix(line, "Implementation-Version:"))
				} else if strings.HasPrefix(line, "Bundle-Version:") && info.Version == "" {
					info.Version = strings.TrimSpace(strings.TrimPrefix(line, "Bundle-Version:"))
				}
			}

			if info.GroupID != "" || info.ArtifactID != "" || info.Version != "" {
				return info, nil
			}
			break
		}

		if strings.HasPrefix(f.Name, "META-INF/maven/") && strings.HasSuffix(f.Name, "/pom.properties") {
			rc, err := f.Open()
			if err != nil {
				continue
			}

			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "groupId=") {
					info.GroupID = strings.TrimPrefix(line, "groupId=")
				} else if strings.HasPrefix(line, "artifactId=") {
					info.ArtifactID = strings.TrimPrefix(line, "artifactId=")
				} else if strings.HasPrefix(line, "version=") {
					info.Version = strings.TrimPrefix(line, "version=")
				}
			}
			rc.Close()

			if info.GroupID != "" && info.ArtifactID != "" && info.Version != "" {
				return info, nil
			}
		}
	}

	if info.GroupID == "" && info.ArtifactID == "" && info.Version == "" {
		return info, fmt.Errorf("no coordinate information found")
	}

	return info, nil
}

func parseJarFilename(jarPath string) JarInfo {
	info := JarInfo{Path: jarPath}
	name := filepath.Base(jarPath)
	name = strings.TrimSuffix(name, ".jar")

	lastDash := strings.LastIndex(name, "-")
	if lastDash == -1 {
		info.ArtifactID = name
		return info
	}

	possibleVersion := name[lastDash+1:]
	if len(possibleVersion) > 0 && (possibleVersion[0] >= '0' && possibleVersion[0] <= '9') {
		info.ArtifactID = name[:lastDash]
		info.Version = possibleVersion
	} else {
		info.ArtifactID = name
	}

	return info
}
