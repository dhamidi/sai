package project

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Project represents a Java project with multiple modules.
type Project struct {
	ID      string
	RootDir string
	SrcDir  string
	OutDir  string
	LibDir  string
	Modules []*Module
}

// Module represents a single Java module within a project.
type Module struct {
	Name       string
	SrcDir     string
	OutDir     string
	ModuleInfo string
	Project    *Project
}

// Load scans the current directory for a Java project structure.
// It looks for src/<project>/<module>/module-info.java patterns.
func Load() (*Project, error) {
	return LoadFrom(".")
}

// LoadFrom scans the given directory for a Java project structure.
func LoadFrom(rootDir string) (*Project, error) {
	srcDir := filepath.Join(rootDir, "src")
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("read src directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectDir := filepath.Join(srcDir, entry.Name())
		modules, err := scanModules(projectDir)
		if err != nil {
			continue
		}

		if len(modules) == 0 {
			continue
		}

		proj := &Project{
			ID:      entry.Name(),
			RootDir: rootDir,
			SrcDir:  srcDir,
			OutDir:  filepath.Join(rootDir, "out"),
			LibDir:  filepath.Join(rootDir, "lib"),
			Modules: modules,
		}

		// Link modules back to project and set output directories
		for _, m := range proj.Modules {
			m.Project = proj
			m.OutDir = filepath.Join(proj.OutDir, proj.ID+"."+m.Name)
		}

		return proj, nil
	}

	return nil, fmt.Errorf("could not detect project: no src/<project>/<module>/module-info.java structure found")
}

func scanModules(projectDir string) ([]*Module, error) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, err
	}

	var modules []*Module
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		moduleDir := filepath.Join(projectDir, entry.Name())
		moduleInfo := filepath.Join(moduleDir, "module-info.java")

		if _, err := os.Stat(moduleInfo); err != nil {
			continue
		}

		modules = append(modules, &Module{
			Name:       entry.Name(),
			SrcDir:     moduleDir,
			ModuleInfo: moduleInfo,
		})
	}

	return modules, nil
}

// Module returns the module with the given name, or nil if not found.
func (p *Project) Module(name string) *Module {
	for _, m := range p.Modules {
		if m.Name == name {
			return m
		}
	}
	return nil
}

// ModulePath returns the module path for java/javac commands.
// It includes the lib directory and optionally the out directory.
func (p *Project) ModulePath(includeOut bool) string {
	if includeOut {
		return p.LibDir + ":" + p.OutDir
	}
	return p.LibDir
}

// FullName returns the fully qualified module name (e.g., "myproject.core").
func (m *Module) FullName() string {
	return m.Project.ID + "." + m.Name
}

// JavaFiles returns all .java files in this module, recursively.
// The module-info.java is always first if includeModuleInfo is true.
func (m *Module) JavaFiles(includeModuleInfo bool) ([]string, error) {
	var files []string

	err := filepath.WalkDir(m.SrcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".java") {
			return nil
		}
		if path == m.ModuleInfo {
			return nil // handle separately
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan java files in %s: %w", m.SrcDir, err)
	}

	if includeModuleInfo {
		// module-info.java must come first
		files = append([]string{m.ModuleInfo}, files...)
	}

	return files, nil
}

// EnsureOutDir creates the output directory for this module if it doesn't exist.
func (m *Module) EnsureOutDir() error {
	if err := os.MkdirAll(m.OutDir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", m.OutDir, err)
	}
	return nil
}
