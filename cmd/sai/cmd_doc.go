package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/dhamidi/sai/java"
	"github.com/dhamidi/sai/java/parser"
	"github.com/spf13/cobra"
)

func newDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc <name>",
		Short: "Show documentation for a Java class, method, or field",
		Long: `Show documentation for a Java class, method, or field.

The name can be:
  - A package name (e.g., java.util lists all classes in the package)
  - A fully qualified class name (e.g., java.util.List)
  - A class name with method/field (e.g., java.util.List.add)
  - A simple class name (e.g., String for java.lang.String)

With no arguments, lists all available packages from src/ and lib/.

For JDK classes, documentation is extracted from src.zip in the JDK installation.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return listAvailablePackages()
			}
			return runDoc(args[0])
		},
	}

	return cmd
}

func runDoc(name string) error {
	// Check if this looks like a package name (all parts lowercase)
	if isPackageName(name) {
		return listPackageContentsAll(name)
	}

	className, memberName := parseDocName(name)

	// For bare identifiers or partial paths, check if it's a package first
	if isLowercase(name) {
		if isLocalPackage(name) || isJDKPackage(name) || isLibPackage(name) {
			return listPackageContentsAll(name)
		}
	}

	// For bare identifiers, try to find in src/ first before defaulting to java.lang
	if !strings.Contains(name, ".") {
		if source, foundClass, err := findLocalClassBySimpleName(name); err == nil {
			return showDocFromSource(source, foundClass, memberName)
		}
	}

	// Try local sources first (src/ directory)
	source, err := findLocalClassSource(className)
	if err == nil {
		return showDocFromSource(source, className, memberName)
	}

	// Try JDK src.zip
	javaHome, err := findJavaHome()
	if err != nil {
		return fmt.Errorf("find java home: %w", err)
	}

	srcZipPath := filepath.Join(javaHome, "lib", "src.zip")
	if _, err := os.Stat(srcZipPath); os.IsNotExist(err) {
		srcZipPath = filepath.Join(javaHome, "src.zip")
		if _, err := os.Stat(srcZipPath); os.IsNotExist(err) {
			return fmt.Errorf("src.zip not found in JDK at %s", javaHome)
		}
	}

	source, err = findClassSourceInZip(srcZipPath, className)
	if err != nil {
		return fmt.Errorf("class %s not found in src/ or JDK", className)
	}

	return showDocFromSource(source, className, memberName)
}

func showDocFromSource(source []byte, className, memberName string) error {
	models, err := java.ClassModelsFromSource(source, parser.WithFile(className+".java"))
	if err != nil {
		return fmt.Errorf("parse source: %w", err)
	}

	if len(models) == 0 {
		return fmt.Errorf("no class found for %s", className)
	}

	model := models[0]

	if memberName == "" {
		printClassDoc(model)
	} else {
		if err := printMemberDoc(model, memberName); err != nil {
			return err
		}
	}

	return nil
}

func findLocalClassSource(className string) ([]byte, error) {
	// Convert class name to path: com.example.Foo -> com/example/Foo.java
	classPath := strings.ReplaceAll(className, ".", "/") + ".java"

	// Check in src/ directory
	candidates := []string{
		filepath.Join("src", classPath),
	}

	// Also check for modular structure: src/module.name/package/Class.java
	srcDir, err := os.Open("src")
	if err == nil {
		entries, err := srcDir.ReadDir(-1)
		srcDir.Close()
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					candidates = append(candidates, filepath.Join("src", entry.Name(), classPath))
				}
			}
		}
	}

	for _, candidate := range candidates {
		if data, err := os.ReadFile(candidate); err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("not found in src/")
}

// findLocalClassBySimpleName searches src/ for a class by its simple name (e.g., "Visor")
// Returns the source, fully qualified class name, and error
func findLocalClassBySimpleName(simpleName string) ([]byte, string, error) {
	srcDir, err := os.Open("src")
	if err != nil {
		return nil, "", err
	}
	defer srcDir.Close()

	entries, err := srcDir.ReadDir(-1)
	if err != nil {
		return nil, "", err
	}

	fileName := simpleName + ".java"

	// Search in module directories: src/module.name/**/SimpleName.java
	for _, entry := range entries {
		if entry.IsDir() {
			if source, className, err := findClassInDir(filepath.Join("src", entry.Name()), fileName, ""); err == nil {
				return source, className, nil
			}
		}
	}

	// Also check direct: src/**/SimpleName.java
	if source, className, err := findClassInDir("src", fileName, ""); err == nil {
		return source, className, nil
	}

	return nil, "", fmt.Errorf("not found")
}

// findClassInDir recursively searches for a .java file and returns its content and package-qualified name
func findClassInDir(dir, fileName, pkgPrefix string) ([]byte, string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "", err
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			// Build package prefix for subdirectories
			newPrefix := name
			if pkgPrefix != "" {
				newPrefix = pkgPrefix + "." + name
			}
			if source, className, err := findClassInDir(filepath.Join(dir, name), fileName, newPrefix); err == nil {
				return source, className, nil
			}
		} else if name == fileName {
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return nil, "", err
			}
			className := strings.TrimSuffix(fileName, ".java")
			if pkgPrefix != "" {
				className = pkgPrefix + "." + className
			}
			return data, className, nil
		}
	}

	return nil, "", fmt.Errorf("not found")
}

func listPackageContentsAll(packageName string) error {
	packagePath := strings.ReplaceAll(packageName, ".", "/")
	var classes []string
	seenClasses := make(map[string]bool)
	subpackages := make(map[string]bool)

	// Check src/ directory first
	findPackageContentsInDir("src", packagePath, seenClasses, &classes, subpackages)

	// Check JDK src.zip
	javaHome, _ := findJavaHome()
	if javaHome != "" {
		srcZipPath := filepath.Join(javaHome, "lib", "src.zip")
		if _, err := os.Stat(srcZipPath); os.IsNotExist(err) {
			srcZipPath = filepath.Join(javaHome, "src.zip")
		}
		if _, err := os.Stat(srcZipPath); err == nil {
			findPackageContentsInZip(srcZipPath, packagePath, seenClasses, &classes, subpackages)
		}
	}

	if len(classes) == 0 && len(subpackages) == 0 {
		return fmt.Errorf("package %s not found", packageName)
	}

	fmt.Printf("package %s\n", packageName)

	// List subpackages first
	if len(subpackages) > 0 {
		var subpkgList []string
		for sp := range subpackages {
			subpkgList = append(subpkgList, sp)
		}
		sort.Strings(subpkgList)
		fmt.Println("\nSubpackages:")
		for _, sp := range subpkgList {
			fmt.Printf("    %s.%s\n", packageName, sp)
		}
	}

	// List classes
	if len(classes) > 0 {
		sort.Strings(classes)
		fmt.Println("\nClasses:")
		for _, c := range classes {
			fmt.Printf("    %s\n", c)
		}
	}

	return nil
}

func findPackageContentsInDir(baseDir, packagePath string, seen map[string]bool, classes *[]string, subpackages map[string]bool) {
	// Check direct path: src/com/example/
	checkDir := filepath.Join(baseDir, packagePath)
	scanPackageDir(checkDir, seen, classes, subpackages)

	// Check modular structure: src/module.name/com/example/
	srcDir, err := os.Open(baseDir)
	if err != nil {
		return
	}
	entries, err := srcDir.ReadDir(-1)
	srcDir.Close()
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			modPath := filepath.Join(baseDir, entry.Name(), packagePath)
			scanPackageDir(modPath, seen, classes, subpackages)
		}
	}
}

func scanPackageDir(dir string, seen map[string]bool, classes *[]string, subpackages map[string]bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			// Subpackage if lowercase
			if len(name) > 0 && name[0] >= 'a' && name[0] <= 'z' {
				subpackages[name] = true
			}
		} else if strings.HasSuffix(name, ".java") {
			className := strings.TrimSuffix(name, ".java")
			if !seen[className] && className != "package-info" && className != "module-info" {
				*classes = append(*classes, className)
				seen[className] = true
			}
		}
	}
}

func listAvailablePackages() error {
	var srcPackages []string
	var libPackages []string

	// Find top-level packages in src/
	srcDir, err := os.Open("src")
	if err == nil {
		entries, err := srcDir.ReadDir(-1)
		srcDir.Close()
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					name := entry.Name()
					// Skip hidden directories
					if strings.HasPrefix(name, ".") {
						continue
					}
					// If it's a lowercase name, it's a package
					if isLowercase(name) {
						srcPackages = append(srcPackages, name)
					}
				}
			}
		}
	}

	// Find packages from lib/ JARs
	libDir, err := os.Open("lib")
	if err == nil {
		entries, err := libDir.ReadDir(-1)
		libDir.Close()
		if err == nil {
			seen := make(map[string]bool)
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jar") {
					pkgs := findPackagesInJar(filepath.Join("lib", entry.Name()))
					for _, pkg := range pkgs {
						if !seen[pkg] {
							libPackages = append(libPackages, pkg)
							seen[pkg] = true
						}
					}
				}
			}
		}
	}

	if len(srcPackages) == 0 && len(libPackages) == 0 {
		return fmt.Errorf("no packages found in src/ or lib/")
	}

	if len(srcPackages) > 0 {
		sort.Strings(srcPackages)
		fmt.Println("Source packages (src/):")
		for _, pkg := range srcPackages {
			fmt.Printf("    %s\n", pkg)
		}
	}

	if len(libPackages) > 0 {
		sort.Strings(libPackages)
		if len(srcPackages) > 0 {
			fmt.Println()
		}
		fmt.Println("Library packages (lib/):")
		for _, pkg := range libPackages {
			fmt.Printf("    %s\n", pkg)
		}
	}

	return nil
}

func isModuleDir(dir string) bool {
	// Check if directory contains module-info.java
	if _, err := os.Stat(filepath.Join(dir, "module-info.java")); err == nil {
		return true
	}
	// Check if it has subdirectories that look like packages
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() && isLowercase(entry.Name()) {
			return true
		}
	}
	return false
}

func findTopLevelPackages(moduleDir string) []string {
	var packages []string
	moduleName := filepath.Base(moduleDir)
	entries, err := os.ReadDir(moduleDir)
	if err != nil {
		return packages
	}
	for _, entry := range entries {
		if entry.IsDir() && isLowercase(entry.Name()) {
			// Use module name as prefix if it matches first package component
			pkgName := entry.Name()
			// Check if this package is the module's main package (e.g., visor.core module has visor package)
			if strings.HasPrefix(moduleName, pkgName+".") || moduleName == pkgName {
				packages = append(packages, pkgName)
			} else {
				// For modular projects like visor.core containing visor/ subdir
				packages = append(packages, pkgName)
			}
		}
	}
	return packages
}

func findPackagesInJar(jarPath string) []string {
	r, err := zip.OpenReader(jarPath)
	if err != nil {
		return nil
	}
	defer r.Close()

	packages := make(map[string]bool)
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".class") {
			// Get directory part
			dir := filepath.Dir(f.Name)
			if dir != "." && dir != "" && dir != "META-INF" && !strings.HasPrefix(dir, "META-INF/") {
				// Convert path to package name
				pkg := strings.ReplaceAll(dir, "/", ".")
				// Only add top-level package
				if idx := strings.Index(pkg, "."); idx > 0 {
					pkg = pkg[:idx]
				}
				packages[pkg] = true
			}
		}
	}

	var result []string
	for pkg := range packages {
		result = append(result, pkg)
	}
	return result
}

func isPackageName(name string) bool {
	if name == "" || !strings.Contains(name, ".") {
		return false
	}
	parts := strings.Split(name, ".")
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
		// Package names are all lowercase
		if part[0] >= 'A' && part[0] <= 'Z' {
			return false
		}
	}
	return true
}

func isLowercase(name string) bool {
	for _, c := range name {
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func isLocalPackage(name string) bool {
	packagePath := strings.ReplaceAll(name, ".", "/")

	// Check src/packagePath
	if info, err := os.Stat(filepath.Join("src", packagePath)); err == nil && info.IsDir() {
		return true
	}

	// Check src/module.name/packagePath
	srcDir, err := os.Open("src")
	if err != nil {
		return false
	}
	defer srcDir.Close()

	entries, err := srcDir.ReadDir(-1)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			modPath := filepath.Join("src", entry.Name(), packagePath)
			if info, err := os.Stat(modPath); err == nil && info.IsDir() {
				return true
			}
		}
	}

	return false
}

func isJDKPackage(name string) bool {
	javaHome, err := findJavaHome()
	if err != nil {
		return false
	}

	srcZipPath := filepath.Join(javaHome, "lib", "src.zip")
	if _, err := os.Stat(srcZipPath); os.IsNotExist(err) {
		srcZipPath = filepath.Join(javaHome, "src.zip")
		if _, err := os.Stat(srcZipPath); os.IsNotExist(err) {
			return false
		}
	}

	packagePath := strings.ReplaceAll(name, ".", "/")

	r, err := zip.OpenReader(srcZipPath)
	if err != nil {
		return false
	}
	defer r.Close()

	// Check if any file matches the package path pattern
	for _, f := range r.File {
		// Match patterns like "java.base/com/example/" or "com/example/"
		if strings.Contains(f.Name, packagePath+"/") {
			return true
		}
	}

	return false
}

func isLibPackage(name string) bool {
	packagePath := strings.ReplaceAll(name, ".", "/")

	libDir, err := os.Open("lib")
	if err != nil {
		return false
	}
	defer libDir.Close()

	entries, err := libDir.ReadDir(-1)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jar") {
			if hasPackageInJar(filepath.Join("lib", entry.Name()), packagePath) {
				return true
			}
		}
	}

	return false
}

func hasPackageInJar(jarPath, packagePath string) bool {
	r, err := zip.OpenReader(jarPath)
	if err != nil {
		return false
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, packagePath+"/") {
			return true
		}
	}
	return false
}

func findPackageContentsInZip(srcZipPath, packagePath string, seen map[string]bool, classes *[]string, subpackages map[string]bool) {
	r, err := zip.OpenReader(srcZipPath)
	if err != nil {
		return
	}
	defer r.Close()

	for _, f := range r.File {
		// Match files in the package directory
		// Handle both "java.base/java/util/List.java" and "java/util/List.java" patterns
		name := f.Name
		var remaining string
		matched := false

		if idx := strings.Index(name, "/"); idx >= 0 {
			// Check if it matches module/package/... pattern
			afterModule := name[idx+1:]
			if strings.HasPrefix(afterModule, packagePath+"/") {
				remaining = afterModule[len(packagePath)+1:]
				matched = true
			}
		}
		// Also check direct package/... pattern
		if !matched && strings.HasPrefix(name, packagePath+"/") {
			remaining = name[len(packagePath)+1:]
			matched = true
		}

		if matched && remaining != "" {
			if slashIdx := strings.Index(remaining, "/"); slashIdx >= 0 {
				// This is a subpackage
				subpkg := remaining[:slashIdx]
				// Only add if it looks like a package (lowercase)
				if len(subpkg) > 0 && subpkg[0] >= 'a' && subpkg[0] <= 'z' {
					subpackages[subpkg] = true
				}
			} else if strings.HasSuffix(remaining, ".java") {
				// Direct class in this package
				className := strings.TrimSuffix(remaining, ".java")
				if !seen[className] && className != "package-info" && className != "module-info" {
					*classes = append(*classes, className)
					seen[className] = true
				}
			}
		}
	}
}

func findJavaHome() (string, error) {
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		return jh, nil
	}

	cmd := exec.Command("java", "-XshowSettings:properties", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("run java: %w", err)
	}

	re := regexp.MustCompile(`java\.home\s*=\s*(.+)`)
	matches := re.FindSubmatch(output)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find java.home in output")
	}

	return strings.TrimSpace(string(matches[1])), nil
}

func parseDocName(name string) (className, memberName string) {
	if name == "" {
		return "", ""
	}

	if !strings.Contains(name, ".") {
		return "java.lang." + name, ""
	}

	parts := strings.Split(name, ".")
	for i := len(parts) - 1; i > 0; i-- {
		part := parts[i]
		if len(part) > 0 && part[0] >= 'A' && part[0] <= 'Z' {
			// Check if this looks like a class name (PascalCase) vs a constant (ALL_CAPS)
			// Constants are all uppercase with possible underscores
			if isConstantName(part) && i > 0 {
				// This is likely a constant, check previous part
				continue
			}
			if i == len(parts)-1 {
				return name, ""
			}
			return strings.Join(parts[:i+1], "."), strings.Join(parts[i+1:], ".")
		}
	}

	return name, ""
}

// isConstantName returns true if the name looks like a Java constant (ALL_CAPS_WITH_UNDERSCORES)
func isConstantName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for _, c := range name {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

func findClassSourceInZip(srcZipPath, className string) ([]byte, error) {
	r, err := zip.OpenReader(srcZipPath)
	if err != nil {
		return nil, fmt.Errorf("open src.zip: %w", err)
	}
	defer r.Close()

	classPath := strings.ReplaceAll(className, ".", "/") + ".java"

	candidates := []string{
		classPath,
		"java.base/" + classPath,
	}

	for _, f := range r.File {
		for _, candidate := range candidates {
			if f.Name == candidate || strings.HasSuffix(f.Name, "/"+candidate) {
				rc, err := f.Open()
				if err != nil {
					return nil, fmt.Errorf("open %s: %w", f.Name, err)
				}
				defer rc.Close()
				return io.ReadAll(rc)
			}
		}
	}

	return nil, fmt.Errorf("source for %s not found in src.zip", className)
}

func printClassDoc(model *java.ClassModel) {
	if model.Javadoc != "" {
		fmt.Println(formatJavadoc(model.Javadoc))
		fmt.Println()
	}

	printClassSignature(model)
	fmt.Println()

	publicMethods := filterPublicMethods(model.Methods)
	if len(publicMethods) > 0 {
		fmt.Println("\nMethods:")
		for _, m := range publicMethods {
			fmt.Printf("    %s\n", formatMethodSignature(m))
		}
	}

	publicFields := filterPublicFields(model.Fields)
	if len(publicFields) > 0 {
		fmt.Println("\nFields:")
		for _, f := range publicFields {
			fmt.Printf("    %s\n", formatFieldSignature(f))
		}
	}
}

func printMemberDoc(model *java.ClassModel, memberName string) error {
	for _, m := range model.Methods {
		if m.Name == memberName || m.Name == "<init>" && memberName == model.SimpleName {
			if m.Javadoc != "" {
				fmt.Println(formatJavadoc(m.Javadoc))
				fmt.Println()
			}
			fmt.Printf("func %s\n", formatMethodSignature(m))
			return nil
		}
	}

	for _, f := range model.Fields {
		if f.Name == memberName {
			if f.Javadoc != "" {
				fmt.Println(formatJavadoc(f.Javadoc))
				fmt.Println()
			}
			fmt.Printf("field %s\n", formatFieldSignature(f))
			return nil
		}
	}

	return fmt.Errorf("member %s not found in %s", memberName, model.Name)
}

func printClassSignature(model *java.ClassModel) {
	var sb strings.Builder

	if model.Visibility == java.VisibilityPublic {
		sb.WriteString("public ")
	}
	if model.IsAbstract && model.Kind != java.ClassKindInterface {
		sb.WriteString("abstract ")
	}
	if model.IsFinal && model.Kind != java.ClassKindRecord {
		sb.WriteString("final ")
	}

	switch model.Kind {
	case java.ClassKindInterface:
		sb.WriteString("interface ")
	case java.ClassKindEnum:
		sb.WriteString("enum ")
	case java.ClassKindRecord:
		sb.WriteString("record ")
	case java.ClassKindAnnotation:
		sb.WriteString("@interface ")
	default:
		sb.WriteString("class ")
	}

	sb.WriteString(model.Name)

	if len(model.TypeParameters) > 0 {
		sb.WriteString("<")
		for i, tp := range model.TypeParameters {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(tp.Name)
			if len(tp.Bounds) > 0 {
				sb.WriteString(" extends ")
				for j, b := range tp.Bounds {
					if j > 0 {
						sb.WriteString(" & ")
					}
					sb.WriteString(formatTypeModel(b))
				}
			}
		}
		sb.WriteString(">")
	}

	if model.SuperClass != "" && model.SuperClass != "java.lang.Object" {
		sb.WriteString(" extends ")
		sb.WriteString(model.SuperClass)
	}

	if len(model.Interfaces) > 0 {
		if model.Kind == java.ClassKindInterface {
			sb.WriteString(" extends ")
		} else {
			sb.WriteString(" implements ")
		}
		sb.WriteString(strings.Join(model.Interfaces, ", "))
	}

	fmt.Println(sb.String())
}

func formatMethodSignature(m java.MethodModel) string {
	var sb strings.Builder

	if len(m.TypeParameters) > 0 {
		sb.WriteString("<")
		for i, tp := range m.TypeParameters {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(tp.Name)
		}
		sb.WriteString("> ")
	}

	sb.WriteString(formatTypeModel(m.ReturnType))
	sb.WriteString(" ")
	if m.Name == "<init>" {
		sb.WriteString("(constructor)")
	} else {
		sb.WriteString(m.Name)
	}
	sb.WriteString("(")

	for i, p := range m.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(formatTypeModel(p.Type))
		if p.Name != "" {
			sb.WriteString(" ")
			sb.WriteString(p.Name)
		}
	}
	sb.WriteString(")")

	if len(m.Exceptions) > 0 {
		sb.WriteString(" throws ")
		sb.WriteString(strings.Join(m.Exceptions, ", "))
	}

	return sb.String()
}

func formatFieldSignature(f java.FieldModel) string {
	var sb strings.Builder

	if f.IsStatic {
		sb.WriteString("static ")
	}
	if f.IsFinal {
		sb.WriteString("final ")
	}

	sb.WriteString(formatTypeModel(f.Type))
	sb.WriteString(" ")
	sb.WriteString(f.Name)

	return sb.String()
}

func formatTypeModel(t java.TypeModel) string {
	s := t.Name
	if len(t.TypeArguments) > 0 {
		s += "<"
		for i, ta := range t.TypeArguments {
			if i > 0 {
				s += ", "
			}
			if ta.IsWildcard {
				s += "?"
				if ta.BoundKind != "" && ta.Bound != nil {
					s += " " + ta.BoundKind + " " + formatTypeModel(*ta.Bound)
				}
			} else if ta.Type != nil {
				s += formatTypeModel(*ta.Type)
			}
		}
		s += ">"
	}
	for i := 0; i < t.ArrayDepth; i++ {
		s += "[]"
	}
	return s
}

func formatJavadoc(javadoc string) string {
	lines := strings.Split(javadoc, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "/**" || line == "*/" {
			continue
		}
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimPrefix(line, "*")
		result = append(result, line)
	}

	text := strings.Join(result, "\n")
	text = strings.TrimSpace(text)

	text = removeHTMLTags(text)

	return text
}

func removeHTMLTags(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")

	codeRe := regexp.MustCompile(`{@code\s+([^}]*)}`)
	s = codeRe.ReplaceAllString(s, "`$1`")

	linkRe := regexp.MustCompile(`{@link\s+([^}]*)}`)
	s = linkRe.ReplaceAllString(s, "$1")

	var buf bytes.Buffer
	for i := 0; i < len(s); i++ {
		if s[i] == '&' {
			end := strings.Index(s[i:], ";")
			if end > 0 && end < 10 {
				entity := s[i : i+end+1]
				switch entity {
				case "&lt;":
					buf.WriteString("<")
				case "&gt;":
					buf.WriteString(">")
				case "&amp;":
					buf.WriteString("&")
				case "&quot;":
					buf.WriteString("\"")
				case "&nbsp;":
					buf.WriteString(" ")
				default:
					buf.WriteString(entity)
				}
				i += end
				continue
			}
		}
		buf.WriteByte(s[i])
	}

	return buf.String()
}

func filterPublicMethods(methods []java.MethodModel) []java.MethodModel {
	var result []java.MethodModel
	for _, m := range methods {
		if m.Visibility == java.VisibilityPublic && !m.IsSynthetic && !m.IsBridge {
			result = append(result, m)
		}
	}
	return result
}

func filterPublicFields(fields []java.FieldModel) []java.FieldModel {
	var result []java.FieldModel
	for _, f := range fields {
		if f.Visibility == java.VisibilityPublic && !f.IsSynthetic {
			result = append(result, f)
		}
	}
	return result
}
