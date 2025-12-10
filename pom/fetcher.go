package pom

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultMavenRepoURL = "https://repo1.maven.org/maven2"
	EnvMavenRepoURL     = "MAVEN_REPO_URL"
)

type MavenFetcher struct {
	RepoURL    string
	httpClient *http.Client
}

func NewMavenFetcher() *MavenFetcher {
	repoURL := os.Getenv(EnvMavenRepoURL)
	if repoURL == "" {
		repoURL = DefaultMavenRepoURL
	}
	repoURL = strings.TrimSuffix(repoURL, "/")

	return &MavenFetcher{
		RepoURL:    repoURL,
		httpClient: &http.Client{},
	}
}

func (f *MavenFetcher) FetchPOM(groupID, artifactID, version string) (*Project, error) {
	url := f.pomURL(groupID, artifactID, version)
	resp, err := f.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch POM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch POM: HTTP %d for %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read POM: %w", err)
	}

	project, err := f.parsePOM(data)
	if err != nil {
		return nil, fmt.Errorf("parse POM: %w", err)
	}

	if err := f.resolveParent(project); err != nil {
		return nil, err
	}

	f.interpolateProperties(project)

	return project, nil
}

func (f *MavenFetcher) parsePOM(data []byte) (*Project, error) {
	var project Project
	if err := xml.Unmarshal(data, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (f *MavenFetcher) resolveParent(project *Project) error {
	if project.Parent == nil {
		return nil
	}

	parent, err := f.FetchPOM(project.Parent.GroupID, project.Parent.ArtifactID, project.Parent.Version)
	if err != nil {
		return fmt.Errorf("fetch parent POM: %w", err)
	}

	if project.GroupID == "" {
		project.GroupID = parent.GroupID
	}
	if project.Version == "" {
		project.Version = parent.Version
	}

	if project.Properties == nil {
		project.Properties = &Properties{Entries: make(map[string]string)}
	}
	if parent.Properties != nil {
		for k, v := range parent.Properties.Entries {
			if _, exists := project.Properties.Entries[k]; !exists {
				project.Properties.Entries[k] = v
			}
		}
	}

	if project.DependencyManagement == nil && parent.DependencyManagement != nil {
		project.DependencyManagement = parent.DependencyManagement
	} else if project.DependencyManagement != nil && parent.DependencyManagement != nil {
		existingDeps := make(map[string]bool)
		for _, d := range project.DependencyManagement.Dependencies {
			key := d.GroupID + ":" + d.ArtifactID
			existingDeps[key] = true
		}
		for _, d := range parent.DependencyManagement.Dependencies {
			key := d.GroupID + ":" + d.ArtifactID
			if !existingDeps[key] {
				project.DependencyManagement.Dependencies = append(project.DependencyManagement.Dependencies, d)
			}
		}
	}

	return nil
}

func (f *MavenFetcher) interpolateProperties(project *Project) {
	props := make(map[string]string)
	props["project.groupId"] = project.GroupID
	props["project.artifactId"] = project.ArtifactID
	props["project.version"] = project.Version
	props["pom.groupId"] = project.GroupID
	props["pom.artifactId"] = project.ArtifactID
	props["pom.version"] = project.Version

	if project.Properties != nil {
		for k, v := range project.Properties.Entries {
			props[k] = v
		}
	}

	interpolate := func(s string) string {
		for k, v := range props {
			s = strings.ReplaceAll(s, "${"+k+"}", v)
		}
		return s
	}

	for i := range project.Dependencies {
		project.Dependencies[i].GroupID = interpolate(project.Dependencies[i].GroupID)
		project.Dependencies[i].ArtifactID = interpolate(project.Dependencies[i].ArtifactID)
		project.Dependencies[i].Version = interpolate(project.Dependencies[i].Version)
	}

	if project.DependencyManagement != nil {
		for i := range project.DependencyManagement.Dependencies {
			project.DependencyManagement.Dependencies[i].GroupID = interpolate(project.DependencyManagement.Dependencies[i].GroupID)
			project.DependencyManagement.Dependencies[i].ArtifactID = interpolate(project.DependencyManagement.Dependencies[i].ArtifactID)
			project.DependencyManagement.Dependencies[i].Version = interpolate(project.DependencyManagement.Dependencies[i].Version)
		}
	}
}

func (f *MavenFetcher) pomURL(groupID, artifactID, version string) string {
	groupPath := strings.ReplaceAll(groupID, ".", "/")
	return fmt.Sprintf("%s/%s/%s/%s/%s-%s.pom", f.RepoURL, groupPath, artifactID, version, artifactID, version)
}

func (f *MavenFetcher) JarURL(groupID, artifactID, version, classifier string) string {
	groupPath := strings.ReplaceAll(groupID, ".", "/")
	if classifier != "" {
		return fmt.Sprintf("%s/%s/%s/%s/%s-%s-%s.jar", f.RepoURL, groupPath, artifactID, version, artifactID, version, classifier)
	}
	return fmt.Sprintf("%s/%s/%s/%s/%s-%s.jar", f.RepoURL, groupPath, artifactID, version, artifactID, version)
}

func (f *MavenFetcher) DownloadJar(groupID, artifactID, version, classifier, destDir string) (string, error) {
	url := f.JarURL(groupID, artifactID, version, classifier)
	resp, err := f.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("download JAR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download JAR: HTTP %d for %s", resp.StatusCode, url)
	}

	var filename string
	if classifier != "" {
		filename = fmt.Sprintf("%s-%s-%s.jar", artifactID, version, classifier)
	} else {
		filename = fmt.Sprintf("%s-%s.jar", artifactID, version)
	}
	destPath := filepath.Join(destDir, filename)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("create directory: %w", err)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		os.Remove(destPath)
		return "", fmt.Errorf("write file: %w", err)
	}

	return destPath, nil
}

func ParseCoordinate(coord string) (groupID, artifactID, version, classifier string, err error) {
	parts := strings.Split(coord, ":")
	switch len(parts) {
	case 3:
		return parts[0], parts[1], parts[2], "", nil
	case 4:
		return parts[0], parts[1], parts[3], parts[2], nil
	default:
		return "", "", "", "", fmt.Errorf("invalid Maven coordinate: %s (expected groupId:artifactId:version or groupId:artifactId:classifier:version)", coord)
	}
}
