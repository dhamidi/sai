package java

import (
	"strings"
)

// ResolveInnerClassReferences fixes type references to inner classes (Option 2: Post-processing fixup).
//
// Problem: When scanning source code files, inner classes like HeaderInfo in Authentication.HeaderInfo
// may be incorrectly resolved as org.eclipse.jetty.client.HeaderInfo instead of
// org.eclipse.jetty.client.Authentication.HeaderInfo if they are referenced from files that
// don't have access to the full context of the outer class.
//
// Solution: After all files in a scan are parsed, this function iterates through all ClassModel types
// and fixes references like "pkg.SimpleName" to "pkg.OuterClass.InnerClass" when the simple name
// matches a known inner class in that package.
//
// The strategy:
//  1. Build a map of package -> (simpleName -> fully qualified inner class name) for all known inner classes
//  2. For each ClassModel, check all type references (methods, fields, parameters, superclass, interfaces)
//  3. If a type matches the pattern "pkg.SimpleName" and SimpleName is a known inner class in that package,
//     replace it with the correct fully qualified name.
func ResolveInnerClassReferences(classes []*ClassModel) {
	// Build a map of package -> (simpleName -> fully qualified inner class name)
	innerClassMap := buildInnerClassMap(classes)

	// For each class, fix all type references
	for _, model := range classes {
		fixClassModelTypes(model, innerClassMap)
	}
}

// buildInnerClassMap creates a map of package -> (simpleName -> fullName) for all inner classes
// Maps from package to a map of simple name -> fully qualified name of inner class
func buildInnerClassMap(classes []*ClassModel) map[string]map[string]string {
	result := make(map[string]map[string]string)

	// Collect all inner classes
	allInnerClasses := make(map[string]string) // fullName -> exists

	for _, model := range classes {
		// Add inner classes from the InnerClasses list
		for _, inner := range model.InnerClasses {
			allInnerClasses[inner.InnerClass] = inner.InnerClass
		}
	}

	// Also add inner classes that appear as separate ClassModels
	for _, model := range classes {
		if isInnerClass(model) {
			allInnerClasses[model.Name] = model.Name
		}
	}

	// Now build the result map: for each inner class, we need to register it
	// under its package (NOT the package of the outer class)
	// This allows us to match pkg.SimpleName -> pkg.Outer.Inner
	for fullName := range allInnerClasses {
		// Extract the simple name (last component after last dot)
		simpleName := extractSimpleName(fullName)

		// Extract the package (everything up to the first component after package)
		// For "org.eclipse.jetty.client.Authentication.HeaderInfo":
		// We want to extract "org.eclipse.jetty.client" as the package
		// where this inner class might be mistakenly resolved
		pkg := extractPackageFromInnerClassName(fullName)

		if pkg == "" {
			continue
		}

		// Register under the package
		if result[pkg] == nil {
			result[pkg] = make(map[string]string)
		}
		result[pkg][simpleName] = fullName
	}

	return result
}

// extractPackageFromInnerClassName extracts the package from an inner class name.
// For "org.eclipse.jetty.client.Authentication.HeaderInfo", it returns "org.eclipse.jetty.client"
func extractPackageFromInnerClassName(innerClassName string) string {
	// Find the first occurrence of the class name part (has upper case)
	// This is a simple heuristic: find the last dot before finding an uppercase letter
	// Actually, we need to find where the package ends and the classes begin
	// The package consists of all lowercase components separated by dots
	// Then comes OuterClass.InnerClass...

	// A simpler approach: just find all dots and assume the first one that's followed
	// by an uppercase letter is the start of the class name
	parts := strings.Split(innerClassName, ".")

	// Find where the class names start (first part with uppercase)
	for i, part := range parts {
		if len(part) > 0 && part[0] >= 'A' && part[0] <= 'Z' {
			// Found start of class name, package is everything before
			if i == 0 {
				return ""
			}
			return strings.Join(parts[:i], ".")
		}
	}

	// If no uppercase found, return all but last part
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ".")
	}
	return ""
}

// isInnerClass checks if a ClassModel is an inner class
// An inner class has its outer class as part of its fully qualified name (e.g., Outer.Inner)
func isInnerClass(model *ClassModel) bool {
	// An inner class will have a name like "package.OuterClass.InnerClass"
	// When we remove the package, we get "OuterClass.InnerClass" which contains a dot
	if model.Package == "" {
		return false
	}
	// Remove package prefix
	nameWithoutPkg := strings.TrimPrefix(model.Name, model.Package+".")
	// If there's still a dot, it's an inner class
	return strings.Contains(nameWithoutPkg, ".")
}

// fixClassModelTypes fixes all type references in a ClassModel
func fixClassModelTypes(model *ClassModel, innerClassMap map[string]map[string]string) {
	// Fix superclass
	if model.SuperClass != "" {
		model.SuperClass = fixTypeName(model.SuperClass, innerClassMap)
	}

	// Fix interfaces
	for i := range model.Interfaces {
		model.Interfaces[i] = fixTypeName(model.Interfaces[i], innerClassMap)
	}

	// Fix fields
	for i := range model.Fields {
		model.Fields[i].Type.Name = fixTypeName(model.Fields[i].Type.Name, innerClassMap)
	}

	// Fix methods
	for i := range model.Methods {
		model.Methods[i].ReturnType.Name = fixTypeName(model.Methods[i].ReturnType.Name, innerClassMap)

		// Fix method parameters
		for j := range model.Methods[i].Parameters {
			model.Methods[i].Parameters[j].Type.Name = fixTypeName(model.Methods[i].Parameters[j].Type.Name, innerClassMap)
		}

		// Fix method exceptions
		for j := range model.Methods[i].Exceptions {
			model.Methods[i].Exceptions[j] = fixTypeName(model.Methods[i].Exceptions[j], innerClassMap)
		}

		// Fix method type parameters
		for j := range model.Methods[i].TypeParameters {
			for k := range model.Methods[i].TypeParameters[j].Bounds {
				model.Methods[i].TypeParameters[j].Bounds[k].Name = fixTypeName(
					model.Methods[i].TypeParameters[j].Bounds[k].Name,
					innerClassMap,
				)
			}
		}
	}

	// Fix type parameters
	for i := range model.TypeParameters {
		for j := range model.TypeParameters[i].Bounds {
			model.TypeParameters[i].Bounds[j].Name = fixTypeName(
				model.TypeParameters[i].Bounds[j].Name,
				innerClassMap,
			)
		}
	}

	// Fix record components
	for i := range model.RecordComponents {
		model.RecordComponents[i].Type.Name = fixTypeName(model.RecordComponents[i].Type.Name, innerClassMap)
	}
}

// fixTypeName checks if a type name matches pkg.SimpleName pattern where SimpleName
// is an inner class in that package, and fixes it to the fully qualified inner class name.
func fixTypeName(typeName string, innerClassMap map[string]map[string]string) string {
	if typeName == "" || !strings.Contains(typeName, ".") {
		return typeName
	}

	// Extract package and simple name
	lastDot := strings.LastIndex(typeName, ".")
	if lastDot == -1 {
		return typeName
	}

	pkg := typeName[:lastDot]
	simpleName := typeName[lastDot+1:]

	// Check if this simple name is a known inner class in this package
	if innerClasses, ok := innerClassMap[pkg]; ok {
		if fullName, ok := innerClasses[simpleName]; ok {
			return fullName
		}
	}

	return typeName
}
