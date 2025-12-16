package project

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhamidi/sai/java/parser"
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
	Name         string
	SrcDir       string
	OutDir       string
	ModuleInfo   string
	Project      *Project
	Dependencies []string // module names this module requires
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

		// Parse dependencies from module-info.java files
		for _, m := range proj.Modules {
			deps, err := parseModuleDependencies(m.ModuleInfo, proj.ID)
			if err != nil {
				// Non-fatal: continue without dependencies
				continue
			}
			m.Dependencies = deps
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

// parseModuleDependencies parses a module-info.java file and extracts
// the names of required modules that belong to this project.
func parseModuleDependencies(moduleInfoPath string, projectID string) ([]string, error) {
	f, err := os.Open(moduleInfoPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := parser.ParseCompilationUnit(f)
	root := p.Finish()
	if root == nil {
		return nil, fmt.Errorf("failed to parse %s", moduleInfoPath)
	}

	var deps []string
	prefix := projectID + "."

	// Find module declaration
	moduleDecl := root.FirstChildOfKind(parser.KindModuleDecl)
	if moduleDecl == nil {
		return nil, fmt.Errorf("no module declaration in %s", moduleInfoPath)
	}

	// Find all requires directives
	for _, child := range moduleDecl.Children {
		if child.Kind != parser.KindRequiresDirective {
			continue
		}

		// Get the qualified name (last child that is a QualifiedName)
		var qualName *parser.Node
		for _, c := range child.Children {
			if c.Kind == parser.KindQualifiedName {
				qualName = c
			}
		}
		if qualName == nil {
			continue
		}

		// Build the module name from the qualified name parts
		moduleName := qualifiedNameToString(qualName)

		// Only include dependencies on modules within this project
		if strings.HasPrefix(moduleName, prefix) {
			// Extract the short module name (e.g., "myproject.core" -> "core")
			shortName := strings.TrimPrefix(moduleName, prefix)
			deps = append(deps, shortName)
		}
	}

	return deps, nil
}

// qualifiedNameToString converts a QualifiedName node to a string.
func qualifiedNameToString(node *parser.Node) string {
	if node.Kind == parser.KindIdentifier {
		return node.TokenLiteral()
	}

	var parts []string
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier {
			parts = append(parts, child.TokenLiteral())
		} else if child.Kind == parser.KindQualifiedName {
			parts = append(parts, qualifiedNameToString(child))
		}
	}
	return strings.Join(parts, ".")
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

// ModulesInOrder returns modules sorted in dependency order (dependencies first).
// Modules with no dependencies come first, then modules that depend only on
// already-listed modules.
func (p *Project) ModulesInOrder() []*Module {
	// Build a set of module names in this project
	moduleSet := make(map[string]bool)
	for _, m := range p.Modules {
		moduleSet[m.Name] = true
	}

	// Topological sort using Kahn's algorithm
	inDegree := make(map[string]int)
	for _, m := range p.Modules {
		inDegree[m.Name] = 0
	}

	// Count in-degrees (only for project-internal dependencies)
	for _, m := range p.Modules {
		for _, dep := range m.Dependencies {
			if moduleSet[dep] {
				inDegree[m.Name]++
			}
		}
	}

	// Start with modules that have no dependencies
	var queue []string
	for _, m := range p.Modules {
		if inDegree[m.Name] == 0 {
			queue = append(queue, m.Name)
		}
	}

	var result []*Module
	for len(queue) > 0 {
		// Pop first element
		name := queue[0]
		queue = queue[1:]

		mod := p.Module(name)
		if mod != nil {
			result = append(result, mod)
		}

		// Reduce in-degree of modules that depend on this one
		for _, m := range p.Modules {
			for _, dep := range m.Dependencies {
				if dep == name {
					inDegree[m.Name]--
					if inDegree[m.Name] == 0 {
						queue = append(queue, m.Name)
					}
				}
			}
		}
	}

	// If we didn't get all modules, there's a cycle - return original order
	if len(result) != len(p.Modules) {
		return p.Modules
	}

	return result
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

// Entrypoint represents a class with a main method that can be run.
type Entrypoint struct {
	ClassName string // Simple class name (e.g., "Cli", "DropZoneApp")
	FullName  string // Fully qualified name (e.g., "sai.main.Cli")
	Slug      string // Kebab-case name for CLI (e.g., "cli", "drop-zone-app")
}

// FindEntrypoints returns all classes in this module that have a main method.
func (m *Module) FindEntrypoints() ([]Entrypoint, error) {
	javaFiles, err := m.JavaFiles(false)
	if err != nil {
		return nil, err
	}

	var entrypoints []Entrypoint
	for _, file := range javaFiles {
		ep, ok, err := findEntrypointInFile(file, m)
		if err != nil {
			continue // skip files that fail to parse
		}
		if ok {
			entrypoints = append(entrypoints, ep)
		}
	}

	return entrypoints, nil
}

// findEntrypointInFile checks if a Java file contains a class with a main method.
func findEntrypointInFile(path string, m *Module) (Entrypoint, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return Entrypoint{}, false, err
	}
	defer f.Close()

	p := parser.ParseCompilationUnit(f)
	root := p.Finish()
	if root == nil {
		return Entrypoint{}, false, fmt.Errorf("failed to parse %s", path)
	}

	// Find the class declaration
	classDecl := root.FirstChildOfKind(parser.KindClassDecl)
	if classDecl == nil {
		return Entrypoint{}, false, nil
	}

	// Get class name from identifier
	classIdent := classDecl.FirstChildOfKind(parser.KindIdentifier)
	if classIdent == nil {
		return Entrypoint{}, false, nil
	}
	className := classIdent.TokenLiteral()

	// Check if there's a main method
	if !hasMainMethod(classDecl) {
		return Entrypoint{}, false, nil
	}

	// Build the package name from the file path relative to module source dir
	relPath, err := filepath.Rel(m.SrcDir, path)
	if err != nil {
		return Entrypoint{}, false, err
	}
	dir := filepath.Dir(relPath)
	var packagePrefix string
	if dir != "." {
		packagePrefix = strings.ReplaceAll(dir, string(filepath.Separator), ".") + "."
	}

	fullName := m.FullName() + "." + packagePrefix + className

	return Entrypoint{
		ClassName: className,
		FullName:  fullName,
		Slug:      classNameToSlug(className),
	}, true, nil
}

// hasMainMethod checks if a class declaration contains a public static void main method.
func hasMainMethod(classDecl *parser.Node) bool {
	// Methods are inside a Block child of the class declaration
	classBody := classDecl.FirstChildOfKind(parser.KindBlock)
	if classBody == nil {
		return false
	}

	for _, child := range classBody.Children {
		if child.Kind != parser.KindMethodDecl {
			continue
		}

		// Check method name
		methodIdent := child.FirstChildOfKind(parser.KindIdentifier)
		if methodIdent == nil || methodIdent.TokenLiteral() != "main" {
			continue
		}

		// Check modifiers for public and static
		modifiers := child.FirstChildOfKind(parser.KindModifiers)
		if modifiers == nil {
			continue
		}

		hasPublic := false
		hasStatic := false
		for _, mod := range modifiers.Children {
			if mod.TokenLiteral() == "public" {
				hasPublic = true
			}
			if mod.TokenLiteral() == "static" {
				hasStatic = true
			}
		}
		if !hasPublic || !hasStatic {
			continue
		}

		// Check return type is void
		typeNode := child.FirstChildOfKind(parser.KindType)
		if typeNode == nil {
			continue
		}
		typeIdent := typeNode.FirstChildOfKind(parser.KindIdentifier)
		if typeIdent == nil || typeIdent.TokenLiteral() != "void" {
			continue
		}

		// Found a valid main method
		return true
	}
	return false
}

// classNameToSlug converts a PascalCase class name to kebab-case.
// e.g., "DropZoneApp" -> "drop-zone-app", "Cli" -> "cli"
func classNameToSlug(name string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
