package java

import (
	"bytes"
	"sort"
	"strings"

	"github.com/dhamidi/sai/java/parser"
)

// javadocFinder helps find Javadoc comments for declarations based on position.
type javadocFinder struct {
	comments []parser.Token // only block comments starting with /**, sorted by start line ascending
	used     map[int]bool   // tracks which comments (by index) have been used
}

func newJavadocFinder(comments []parser.Token) *javadocFinder {
	var javadocs []parser.Token
	for _, c := range comments {
		if c.Kind == parser.TokenComment && strings.HasPrefix(c.Literal, "/**") {
			javadocs = append(javadocs, c)
		}
	}
	// Sort by start line ascending
	sort.Slice(javadocs, func(i, j int) bool {
		return javadocs[i].Span.Start.Line < javadocs[j].Span.Start.Line
	})
	return &javadocFinder{comments: javadocs, used: make(map[int]bool)}
}

// FindForNode returns the Javadoc comment that immediately precedes the given node.
// A Javadoc is considered "preceding" if it ends just before the node starts
// (allowing for annotations and modifiers). Each Javadoc can only be matched once.
func (jf *javadocFinder) FindForNode(node *parser.Node) string {
	if jf == nil || len(jf.comments) == 0 {
		return ""
	}
	startLine := node.Span.Start.Line

	// Find the closest unused Javadoc that ends before this node starts
	bestIdx := -1
	bestDistance := 100 // max distance we'll accept

	for i, c := range jf.comments {
		if jf.used[i] {
			continue
		}
		endLine := c.Span.End.Line
		endCol := c.Span.End.Column
		// The Javadoc must end before the declaration starts (either on an earlier line,
		// or on the same line but before the node's start column for single-line javadocs)
		if endLine > startLine {
			continue
		}
		if endLine == startLine && endCol >= node.Span.Start.Column {
			continue
		}
		distance := startLine - endLine
		// Only match if within a reasonable distance (allowing for annotations/modifiers)
		// and this is closer than any previous match
		if distance < bestDistance {
			bestIdx = i
			bestDistance = distance
		}
	}

	if bestIdx >= 0 {
		jf.used[bestIdx] = true
		return jf.comments[bestIdx].Literal
	}
	return ""
}

func ClassModelsFromSource(source []byte, opts ...parser.Option) ([]*ClassModel, error) {
	opts = append(opts, parser.WithComments())
	p := parser.ParseCompilationUnit(bytes.NewReader(source), opts...)
	node := p.Finish()
	if node == nil {
		return nil, nil
	}
	comments := p.Comments()
	models := classModelsFromCompilationUnit(node, comments)

	if sourcePath := p.SourcePath(); sourcePath != "" {
		sourceURL := FileURL(sourcePath)
		for _, m := range models {
			m.SourceURL = sourceURL
		}
	}

	return models, nil
}

func classModelsFromCompilationUnit(cu *parser.Node, comments []parser.Token) []*ClassModel {
	var models []*ClassModel
	pkg := packageFromCompilationUnit(cu)
	resolver := newTypeResolver(pkg, importsFromCompilationUnit(cu), nil)
	jf := newJavadocFinder(comments)

	for _, child := range cu.Children {
		switch child.Kind {
		case parser.KindClassDecl:
			models = append(models, classModelFromClassDecl(child, pkg, resolver, jf))
			models = append(models, innerClassesFromClassDecl(child, pkg, resolver, jf)...)
		case parser.KindInterfaceDecl:
			models = append(models, classModelFromInterfaceDecl(child, pkg, resolver, jf))
			models = append(models, innerClassesFromInterfaceDecl(child, pkg, resolver, jf)...)
		case parser.KindEnumDecl:
			models = append(models, classModelFromEnumDecl(child, pkg, resolver, jf))
			models = append(models, innerClassesFromEnumDecl(child, pkg, resolver, jf)...)
		case parser.KindRecordDecl:
			models = append(models, classModelFromRecordDecl(child, pkg, resolver, jf))
			models = append(models, innerClassesFromRecordDecl(child, pkg, resolver, jf)...)
		case parser.KindAnnotationDecl:
			models = append(models, classModelFromAnnotationDecl(child, pkg, resolver, jf))
		}
	}
	return models
}

func packageFromCompilationUnit(cu *parser.Node) string {
	pkgDecl := cu.FirstChildOfKind(parser.KindPackageDecl)
	if pkgDecl == nil {
		return ""
	}
	qn := pkgDecl.FirstChildOfKind(parser.KindQualifiedName)
	if qn == nil {
		return ""
	}
	return qualifiedNameToString(qn)
}

func qualifiedNameToString(qn *parser.Node) string {
	var parts []string
	for _, child := range qn.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			parts = append(parts, child.Token.Literal)
		}
	}
	return strings.Join(parts, ".")
}

type importInfo struct {
	qualifiedName string
	isStatic      bool
	isWildcard    bool
}

func importsFromCompilationUnit(cu *parser.Node) []importInfo {
	var imports []importInfo
	for _, child := range cu.Children {
		if child.Kind == parser.KindImportDecl {
			imp := importInfo{}
			for _, ic := range child.Children {
				if ic.Kind == parser.KindIdentifier && ic.Token != nil {
					if ic.Token.Literal == "static" {
						imp.isStatic = true
					} else if ic.Token.Literal == "*" {
						imp.isWildcard = true
					}
				} else if ic.Kind == parser.KindQualifiedName {
					imp.qualifiedName = qualifiedNameToString(ic)
				}
			}
			imports = append(imports, imp)
		}
	}
	return imports
}

type typeResolver struct {
	pkg          string
	imports      []importInfo
	innerClasses map[string]string // map of simpleName -> fully qualified name
	classes      []*ClassModel     // available classes for resolving star imports
}

func newTypeResolver(pkg string, imports []importInfo, classes []*ClassModel) *typeResolver {
	return &typeResolver{
		pkg:          pkg,
		imports:      imports,
		innerClasses: make(map[string]string),
		classes:      classes,
	}
}

// registerInnerClass registers an inner class with its simple name
func (r *typeResolver) registerInnerClass(simpleName, fullName string) {
	r.innerClasses[simpleName] = fullName
}

var javaLangTypes = map[string]bool{
	"Object": true, "String": true, "Class": true, "System": true,
	"Throwable": true, "Exception": true, "RuntimeException": true, "Error": true,
	"Integer": true, "Long": true, "Short": true, "Byte": true,
	"Float": true, "Double": true, "Character": true, "Boolean": true,
	"Number": true, "Comparable": true, "CharSequence": true,
	"Iterable": true, "Cloneable": true, "Runnable": true,
	"Thread": true, "StringBuilder": true, "StringBuffer": true,
	"Math": true, "Enum": true, "Record": true,
	"Override": true, "Deprecated": true, "SuppressWarnings": true, "FunctionalInterface": true,
}

func (r *typeResolver) resolve(simpleName string) string {
	if simpleName == "" {
		return ""
	}

	if strings.Contains(simpleName, ".") {
		return simpleName
	}

	switch simpleName {
	case "boolean", "byte", "char", "short", "int", "long", "float", "double", "void":
		return simpleName
	}

	// Check if this is a known inner class first
	if fullName, ok := r.innerClasses[simpleName]; ok {
		return fullName
	}

	for _, imp := range r.imports {
		if imp.isWildcard || imp.isStatic {
			continue
		}
		parts := strings.Split(imp.qualifiedName, ".")
		if len(parts) > 0 && parts[len(parts)-1] == simpleName {
			return imp.qualifiedName
		}
	}

	// Check star imports against available classes
	for _, imp := range r.imports {
		if !imp.isWildcard || imp.isStatic {
			continue
		}
		// imp.qualifiedName is e.g. "com.example" for "import com.example.*"
		candidate := imp.qualifiedName + "." + simpleName
		for _, cls := range r.classes {
			if cls.Name == candidate {
				return candidate
			}
		}
	}

	if javaLangTypes[simpleName] {
		return "java.lang." + simpleName
	}

	if r.pkg != "" {
		return r.pkg + "." + simpleName
	}

	return simpleName
}

func classModelFromClassDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindClass,
		Package:    pkg,
		Visibility: VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	// First pass: extract class name and collect inner class declarations
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			model.SimpleName = child.Token.Literal
			if pkg != "" {
				model.Name = pkg + "." + model.SimpleName
			} else {
				model.Name = model.SimpleName
			}
			break
		}
	}

	// Register all inner classes with the resolver before processing members
	collectAndRegisterInnerClasses(node, model.Name, resolver)

	// Second pass: process body members
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindType:
			if model.SuperClass == "" {
				model.SuperClass = typeModelFromTypeNode(child, resolver).Name
			} else {
				model.Interfaces = append(model.Interfaces, typeModelFromTypeNode(child, resolver).Name)
			}
		case parser.KindBlock:
			extractClassBodyMembers(child, model, resolver, jf)
		case parser.KindFieldDecl:
			model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver, jf)...)
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver, jf))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver, jf))
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			inner := classModelFromClassDeclNested(child, model.Name, resolver, jf)
			model.InnerClasses = append(model.InnerClasses, InnerClassModel{
				InnerClass: inner.Name,
				OuterClass: model.Name,
				InnerName:  inner.SimpleName,
				Visibility: inner.Visibility,
				IsStatic:   inner.IsStatic,
				IsFinal:    inner.IsFinal,
				IsAbstract: inner.IsAbstract,
			})
		}
	}

	return model
}

func extractClassBodyMembers(block *parser.Node, model *ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range block.Children {
		switch child.Kind {
		case parser.KindFieldDecl:
			model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver, jf)...)
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver, jf))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver, jf))
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			inner := classModelFromClassDeclNested(child, model.Name, resolver, jf)
			model.InnerClasses = append(model.InnerClasses, InnerClassModel{
				InnerClass: inner.Name,
				OuterClass: model.Name,
				InnerName:  inner.SimpleName,
				Visibility: inner.Visibility,
				IsStatic:   inner.IsStatic,
				IsFinal:    inner.IsFinal,
				IsAbstract: inner.IsAbstract,
			})
		}
	}
}

// collectAndRegisterInnerClasses recursively finds and registers all inner classes
func collectAndRegisterInnerClasses(node *parser.Node, outerName string, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			// Get the simple name of the inner class
			var simpleName string
			for _, subchild := range child.Children {
				if subchild.Kind == parser.KindIdentifier && subchild.Token != nil {
					simpleName = subchild.Token.Literal
					break
				}
			}
			if simpleName != "" {
				fullName := outerName + "." + simpleName
				resolver.registerInnerClass(simpleName, fullName)
				// Recursively register nested inner classes
				collectAndRegisterInnerClasses(child, fullName, resolver)
			}
		case parser.KindBlock:
			// Also check for inner classes inside the block
			collectAndRegisterInnerClassesFromBlock(child, outerName, resolver)
		}
	}
}

// collectAndRegisterInnerClassesFromBlock finds inner classes in a block (class body)
func collectAndRegisterInnerClassesFromBlock(block *parser.Node, outerName string, resolver *typeResolver) {
	for _, child := range block.Children {
		switch child.Kind {
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			// Get the simple name of the inner class
			var simpleName string
			for _, subchild := range child.Children {
				if subchild.Kind == parser.KindIdentifier && subchild.Token != nil {
					simpleName = subchild.Token.Literal
					break
				}
			}
			if simpleName != "" {
				fullName := outerName + "." + simpleName
				resolver.registerInnerClass(simpleName, fullName)
				// Recursively register nested inner classes
				collectAndRegisterInnerClasses(child, fullName, resolver)
			}
		}
	}
}

func classModelFromClassDeclNested(node *parser.Node, outerName string, resolver *typeResolver, jf *javadocFinder) *ClassModel {
	model := &ClassModel{
		Visibility: VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	switch node.Kind {
	case parser.KindClassDecl:
		model.Kind = ClassKindClass
	case parser.KindInterfaceDecl:
		model.Kind = ClassKindInterface
	case parser.KindEnumDecl:
		model.Kind = ClassKindEnum
	case parser.KindRecordDecl:
		model.Kind = ClassKindRecord
	case parser.KindAnnotationDecl:
		model.Kind = ClassKindAnnotation
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	// Extract the class name
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			model.SimpleName = child.Token.Literal
			model.Name = outerName + "." + model.SimpleName
			break
		}
	}

	// Extract package from outer name (everything before the last dot of the outer class)
	if idx := strings.Index(outerName, "."); idx != -1 {
		// Find the package portion (before the first uppercase class name)
		parts := strings.Split(outerName, ".")
		for i, part := range parts {
			if len(part) > 0 && part[0] >= 'A' && part[0] <= 'Z' {
				model.Package = strings.Join(parts[:i], ".")
				break
			}
		}
	}

	// Parse full class body - fields, methods, constructors, etc.
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindParameters:
			// Record components
			model.RecordComponents = recordComponentsFromParameters(child, resolver)
		case parser.KindType:
			if model.SuperClass == "" && model.Kind == ClassKindClass {
				model.SuperClass = typeModelFromTypeNode(child, resolver).Name
			} else {
				model.Interfaces = append(model.Interfaces, typeModelFromTypeNode(child, resolver).Name)
			}
		case parser.KindBlock:
			extractClassBodyMembers(child, model, resolver, jf)
		case parser.KindFieldDecl:
			model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver, jf)...)
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver, jf))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver, jf))
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			inner := classModelFromClassDeclNested(child, model.Name, resolver, jf)
			model.InnerClasses = append(model.InnerClasses, InnerClassModel{
				InnerClass: inner.Name,
				OuterClass: model.Name,
				InnerName:  inner.SimpleName,
				Visibility: inner.Visibility,
				IsStatic:   inner.IsStatic,
				IsFinal:    inner.IsFinal,
				IsAbstract: inner.IsAbstract,
			})
		}
	}

	return model
}

func innerClassesFromClassDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) []*ClassModel {
	outerClassName := ""
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if pkg != "" {
				outerClassName = pkg + "." + child.Token.Literal
			} else {
				outerClassName = child.Token.Literal
			}
			break
		}
	}
	if outerClassName == "" {
		return nil
	}
	return collectInnerClasses(node, outerClassName, resolver, jf)
}

func innerClassesFromInterfaceDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) []*ClassModel {
	outerClassName := ""
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if pkg != "" {
				outerClassName = pkg + "." + child.Token.Literal
			} else {
				outerClassName = child.Token.Literal
			}
			break
		}
	}
	if outerClassName == "" {
		return nil
	}
	return collectInnerClasses(node, outerClassName, resolver, jf)
}

func innerClassesFromEnumDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) []*ClassModel {
	outerClassName := ""
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if pkg != "" {
				outerClassName = pkg + "." + child.Token.Literal
			} else {
				outerClassName = child.Token.Literal
			}
			break
		}
	}
	if outerClassName == "" {
		return nil
	}
	return collectInnerClasses(node, outerClassName, resolver, jf)
}

func innerClassesFromRecordDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) []*ClassModel {
	outerClassName := ""
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if pkg != "" {
				outerClassName = pkg + "." + child.Token.Literal
			} else {
				outerClassName = child.Token.Literal
			}
			break
		}
	}
	if outerClassName == "" {
		return nil
	}
	return collectInnerClasses(node, outerClassName, resolver, jf)
}

func collectInnerClasses(node *parser.Node, outerClassName string, resolver *typeResolver, jf *javadocFinder) []*ClassModel {
	var models []*ClassModel

	// Look in direct children
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			innerModel := classModelFromClassDeclNested(child, outerClassName, resolver, jf)
			models = append(models, innerModel)
			// Recursively collect nested inner classes
			models = append(models, collectInnerClasses(child, innerModel.Name, resolver, jf)...)
		case parser.KindBlock:
			// Also look in blocks (class body)
			models = append(models, collectInnerClassesFromBlock(child, outerClassName, resolver, jf)...)
		}
	}

	return models
}

func collectInnerClassesFromBlock(block *parser.Node, outerClassName string, resolver *typeResolver, jf *javadocFinder) []*ClassModel {
	var models []*ClassModel

	for _, child := range block.Children {
		switch child.Kind {
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			innerModel := classModelFromClassDeclNested(child, outerClassName, resolver, jf)
			models = append(models, innerModel)
			// Recursively collect nested inner classes
			models = append(models, collectInnerClasses(child, innerModel.Name, resolver, jf)...)
		}
	}

	return models
}

func classModelFromInterfaceDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindInterface,
		Package:    pkg,
		Visibility: VisibilityPackage,
		IsAbstract: true,
		Javadoc:    jf.FindForNode(node),
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	// First pass: extract interface name and collect inner class declarations
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			model.SimpleName = child.Token.Literal
			if pkg != "" {
				model.Name = pkg + "." + model.SimpleName
			} else {
				model.Name = model.SimpleName
			}
			break
		}
	}

	// Register all inner classes with the resolver before processing members
	collectAndRegisterInnerClasses(node, model.Name, resolver)

	// Second pass: process body members
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindType:
			model.Interfaces = append(model.Interfaces, resolver.resolve(typeNameFromTypeNode(child, resolver)))
		case parser.KindBlock:
			extractInterfaceBodyMembers(child, model, resolver, jf)
		case parser.KindFieldDecl:
			fields := fieldModelsFromFieldDecl(child, resolver, jf)
			for i := range fields {
				fields[i].IsStatic = true
				fields[i].IsFinal = true
				fields[i].Visibility = VisibilityPublic
			}
			model.Fields = append(model.Fields, fields...)
		case parser.KindMethodDecl:
			method := methodModelFromMethodDecl(child, resolver, jf)
			if !method.IsStatic && !method.IsDefault {
				method.IsAbstract = true
			}
			method.Visibility = VisibilityPublic
			model.Methods = append(model.Methods, method)
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}

	return model
}

func extractInterfaceBodyMembers(block *parser.Node, model *ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range block.Children {
		switch child.Kind {
		case parser.KindFieldDecl:
			fields := fieldModelsFromFieldDecl(child, resolver, jf)
			for i := range fields {
				fields[i].IsStatic = true
				fields[i].IsFinal = true
				fields[i].Visibility = VisibilityPublic
			}
			model.Fields = append(model.Fields, fields...)
		case parser.KindMethodDecl:
			method := methodModelFromMethodDecl(child, resolver, jf)
			if !method.IsStatic && !method.IsDefault {
				method.IsAbstract = true
			}
			method.Visibility = VisibilityPublic
			model.Methods = append(model.Methods, method)
		}
	}
}

func classModelFromEnumDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindEnum,
		Package:    pkg,
		Visibility: VisibilityPackage,
		SuperClass: "java.lang.Enum",
		Javadoc:    jf.FindForNode(node),
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	// First pass: extract enum name and collect inner class declarations
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			model.SimpleName = child.Token.Literal
			if pkg != "" {
				model.Name = pkg + "." + model.SimpleName
			} else {
				model.Name = model.SimpleName
			}
			break
		}
	}

	// Register all inner classes with the resolver before processing members
	collectAndRegisterInnerClasses(node, model.Name, resolver)

	// Second pass: process body members
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			model.Interfaces = append(model.Interfaces, resolver.resolve(typeNameFromTypeNode(child, resolver)))
		case parser.KindBlock:
			extractClassBodyMembers(child, model, resolver, jf)
		case parser.KindFieldDecl:
			// Check if this is an enum constant or a regular field
			if isEnumConstant(child) {
				model.EnumConstants = append(model.EnumConstants, enumConstantFromFieldDecl(child))
			} else {
				model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver, jf)...)
			}
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver, jf))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver, jf))
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}

	return model
}

func classModelFromRecordDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindRecord,
		Package:    pkg,
		Visibility: VisibilityPackage,
		SuperClass: "java.lang.Record",
		IsFinal:    true,
		Javadoc:    jf.FindForNode(node),
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	// First pass: extract record name and collect inner class declarations
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			model.SimpleName = child.Token.Literal
			if pkg != "" {
				model.Name = pkg + "." + model.SimpleName
			} else {
				model.Name = model.SimpleName
			}
			break
		}
	}

	// Register all inner classes with the resolver before processing members
	collectAndRegisterInnerClasses(node, model.Name, resolver)

	// Second pass: process body members
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindParameters:
			model.RecordComponents = recordComponentsFromParameters(child, resolver)
		case parser.KindType:
			model.Interfaces = append(model.Interfaces, resolver.resolve(typeNameFromTypeNode(child, resolver)))
		case parser.KindBlock:
			extractClassBodyMembers(child, model, resolver, jf)
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver, jf))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver, jf))
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}

	return model
}

func classModelFromAnnotationDecl(node *parser.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindAnnotation,
		Package:    pkg,
		Visibility: VisibilityPackage,
		IsAbstract: true,
		Interfaces: []string{"java.lang.annotation.Annotation"},
		Javadoc:    jf.FindForNode(node),
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				model.SimpleName = child.Token.Literal
				if pkg != "" {
					model.Name = pkg + "." + model.SimpleName
				} else {
					model.Name = model.SimpleName
				}
			}
		case parser.KindBlock:
			extractAnnotationBodyMembers(child, model, resolver, jf)
		case parser.KindMethodDecl:
			method := methodModelFromMethodDecl(child, resolver, jf)
			method.IsAbstract = true
			method.Visibility = VisibilityPublic
			model.Methods = append(model.Methods, method)
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}

	return model
}

func extractAnnotationBodyMembers(block *parser.Node, model *ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range block.Children {
		if child.Kind == parser.KindMethodDecl {
			method := methodModelFromMethodDecl(child, resolver, jf)
			method.IsAbstract = true
			method.Visibility = VisibilityPublic
			model.Methods = append(model.Methods, method)
		}
	}
}

func applyModifiersToClass(modifiers *parser.Node, model *ClassModel, resolver *typeResolver) {
	for _, child := range modifiers.Children {
		if child.Token == nil {
			if child.Kind == parser.KindAnnotation {
				model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
			}
			continue
		}
		switch child.Token.Literal {
		case "public":
			model.Visibility = VisibilityPublic
		case "protected":
			model.Visibility = VisibilityProtected
		case "private":
			model.Visibility = VisibilityPrivate
		case "abstract":
			model.IsAbstract = true
		case "static":
			model.IsStatic = true
		case "final":
			model.IsFinal = true
		case "sealed":
			model.IsSealed = true
		}
	}
}

func typeParametersFromNode(node *parser.Node, resolver *typeResolver) []TypeParameterModel {
	var params []TypeParameterModel
	for _, child := range node.Children {
		if child.Kind == parser.KindTypeParameter {
			params = append(params, typeParameterFromNode(child, resolver))
		}
	}
	return params
}

func typeParameterFromNode(node *parser.Node, resolver *typeResolver) TypeParameterModel {
	param := TypeParameterModel{}
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				param.Name = child.Token.Literal
			}
		case parser.KindType:
			param.Bounds = append(param.Bounds, typeModelFromTypeNode(child, resolver))
		}
	}
	return param
}

func typeNameFromTypeNode(node *parser.Node, resolver *typeResolver) string {
	tm := typeModelFromTypeNode(node, resolver)
	return tm.Name
}

func typeModelFromTypeNode(node *parser.Node, resolver *typeResolver) TypeModel {
	model := TypeModel{}

	if node.Token != nil {
		model.Name = node.Token.Literal
		if resolver != nil {
			model.Name = resolver.resolve(model.Name)
		}
		return model
	}

	if node.Kind == parser.KindArrayType {
		for _, ac := range node.Children {
			if ac.Kind == parser.KindType || ac.Kind == parser.KindQualifiedName || ac.Kind == parser.KindIdentifier || ac.Kind == parser.KindArrayType {
				inner := typeModelFromTypeNode(ac, resolver)
				model.Name = inner.Name
				model.ArrayDepth = inner.ArrayDepth + 1
				model.TypeArguments = inner.TypeArguments
				return model
			}
		}
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				if model.Name != "" {
					model.Name += "."
				}
				model.Name += child.Token.Literal
			}
		case parser.KindQualifiedName:
			model.Name = qualifiedNameToString(child)
		case parser.KindArrayType:
			for _, ac := range child.Children {
				if ac.Kind == parser.KindType || ac.Kind == parser.KindQualifiedName || ac.Kind == parser.KindIdentifier {
					inner := typeModelFromTypeNode(ac, resolver)
					model.Name = inner.Name
					model.ArrayDepth = inner.ArrayDepth + 1
					model.TypeArguments = inner.TypeArguments
					break
				}
			}
		case parser.KindParameterizedType:
			for _, pc := range child.Children {
				switch pc.Kind {
				case parser.KindQualifiedName:
					model.Name = qualifiedNameToString(pc)
				case parser.KindIdentifier:
					if pc.Token != nil {
						if model.Name != "" {
							model.Name += "."
						}
						model.Name += pc.Token.Literal
					}
				case parser.KindTypeArguments:
					model.TypeArguments = typeArgumentsFromNode(pc, resolver)
				}
			}
		case parser.KindType:
			inner := typeModelFromTypeNode(child, resolver)
			model.Name = inner.Name
			model.ArrayDepth = inner.ArrayDepth
			model.TypeArguments = inner.TypeArguments
		}
	}

	if resolver != nil {
		model.Name = resolver.resolve(model.Name)
	}
	return model
}

func typeArgumentsFromNode(node *parser.Node, resolver *typeResolver) []TypeArgumentModel {
	var args []TypeArgumentModel
	for _, child := range node.Children {
		if child.Kind == parser.KindTypeArgument || child.Kind == parser.KindType {
			args = append(args, typeArgumentFromNode(child, resolver))
		}
	}
	return args
}

func typeArgumentFromNode(node *parser.Node, resolver *typeResolver) TypeArgumentModel {
	arg := TypeArgumentModel{}

	if node.Kind == parser.KindType {
		tm := typeModelFromTypeNode(node, resolver)
		arg.Type = &tm
		return arg
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType, parser.KindArrayType:
			tm := typeModelFromTypeNode(child, resolver)
			if arg.BoundKind != "" {
				arg.Bound = &tm
			} else {
				arg.Type = &tm
			}
		case parser.KindWildcard:
			arg.IsWildcard = true
			for _, wc := range child.Children {
				if wc.Token != nil {
					switch wc.Token.Literal {
					case "extends":
						arg.BoundKind = "extends"
					case "super":
						arg.BoundKind = "super"
					}
				}
				if wc.Kind == parser.KindType {
					tm := typeModelFromTypeNode(wc, resolver)
					arg.Bound = &tm
				}
			}
		}
	}
	return arg
}

func fieldModelsFromFieldDecl(node *parser.Node, resolver *typeResolver, jf *javadocFinder) []FieldModel {
	var fields []FieldModel
	baseField := FieldModel{
		Visibility: VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToField(modifiers, &baseField, resolver)
	}

	var fieldType TypeModel
	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			fieldType = typeModelFromTypeNode(child, resolver)
			break
		}
	}

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			field := baseField
			field.Name = child.Token.Literal
			field.Type = fieldType
			fields = append(fields, field)
		}
	}

	return fields
}

func applyModifiersToField(modifiers *parser.Node, field *FieldModel, resolver *typeResolver) {
	for _, child := range modifiers.Children {
		if child.Token == nil {
			if child.Kind == parser.KindAnnotation {
				field.Annotations = append(field.Annotations, annotationModelFromNode(child, resolver))
			}
			continue
		}
		switch child.Token.Literal {
		case "public":
			field.Visibility = VisibilityPublic
		case "protected":
			field.Visibility = VisibilityProtected
		case "private":
			field.Visibility = VisibilityPrivate
		case "static":
			field.IsStatic = true
		case "final":
			field.IsFinal = true
		case "volatile":
			field.IsVolatile = true
		case "transient":
			field.IsTransient = true
		}
	}
}

func methodModelFromMethodDecl(node *parser.Node, resolver *typeResolver, jf *javadocFinder) MethodModel {
	model := MethodModel{
		Visibility: VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToMethod(modifiers, &model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType, parser.KindArrayType:
			model.ReturnType = typeModelFromTypeNode(child, resolver)
		case parser.KindIdentifier:
			if child.Token != nil {
				model.Name = child.Token.Literal
			}
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindParameters:
			model.Parameters = parametersFromNode(child, resolver)
		case parser.KindThrowsList:
			model.Exceptions = exceptionsFromThrowsList(child, resolver)
		}
	}

	return model
}

func methodModelFromConstructorDecl(node *parser.Node, className string, resolver *typeResolver, jf *javadocFinder) MethodModel {
	model := MethodModel{
		Name:       "<init>",
		Visibility: VisibilityPackage,
		ReturnType: TypeModel{Name: "void"},
		Javadoc:    jf.FindForNode(node),
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToMethod(modifiers, &model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindParameters:
			model.Parameters = parametersFromNode(child, resolver)
		case parser.KindThrowsList:
			model.Exceptions = exceptionsFromThrowsList(child, resolver)
		}
	}

	return model
}

func applyModifiersToMethod(modifiers *parser.Node, method *MethodModel, resolver *typeResolver) {
	for _, child := range modifiers.Children {
		if child.Token == nil {
			if child.Kind == parser.KindAnnotation {
				method.Annotations = append(method.Annotations, annotationModelFromNode(child, resolver))
			}
			continue
		}
		switch child.Token.Literal {
		case "public":
			method.Visibility = VisibilityPublic
		case "protected":
			method.Visibility = VisibilityProtected
		case "private":
			method.Visibility = VisibilityPrivate
		case "static":
			method.IsStatic = true
		case "final":
			method.IsFinal = true
		case "abstract":
			method.IsAbstract = true
		case "synchronized":
			method.IsSynchronized = true
		case "native":
			method.IsNative = true
		case "default":
			method.IsDefault = true
		}
	}
}

func parametersFromNode(node *parser.Node, resolver *typeResolver) []ParameterModel {
	var params []ParameterModel
	for _, child := range node.Children {
		if child.Kind == parser.KindParameter {
			params = append(params, parameterFromNode(child, resolver))
		}
	}
	return params
}

func parameterFromNode(node *parser.Node, resolver *typeResolver) ParameterModel {
	param := ParameterModel{}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		for _, child := range modifiers.Children {
			if child.Token != nil && child.Token.Literal == "final" {
				param.IsFinal = true
			}
			if child.Kind == parser.KindAnnotation {
				param.Annotations = append(param.Annotations, annotationModelFromNode(child, resolver))
			}
		}
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			param.Type = typeModelFromTypeNode(child, resolver)
		case parser.KindArrayType:
			param.Type = typeModelFromTypeNode(child, resolver)
		case parser.KindIdentifier:
			if child.Token != nil {
				param.Name = child.Token.Literal
			}
		}
	}

	return param
}

func exceptionsFromThrowsList(node *parser.Node, resolver *typeResolver) []string {
	var exceptions []string
	for _, child := range node.Children {
		if child.Kind == parser.KindType {
			exceptions = append(exceptions, resolver.resolve(typeNameFromTypeNode(child, resolver)))
		}
	}
	return exceptions
}

func recordComponentsFromParameters(node *parser.Node, resolver *typeResolver) []RecordComponentModel {
	var components []RecordComponentModel
	for _, child := range node.Children {
		if child.Kind == parser.KindParameter {
			comp := recordComponentFromParameter(child, resolver)
			components = append(components, comp)
		}
	}
	return components
}

func recordComponentFromParameter(node *parser.Node, resolver *typeResolver) RecordComponentModel {
	comp := RecordComponentModel{}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		for _, child := range modifiers.Children {
			if child.Kind == parser.KindAnnotation {
				comp.Annotations = append(comp.Annotations, annotationModelFromNode(child, resolver))
			}
		}
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			comp.Type = typeModelFromTypeNode(child, resolver)
		case parser.KindIdentifier:
			if child.Token != nil {
				comp.Name = child.Token.Literal
			}
		}
	}

	return comp
}

func annotationModelFromNode(node *parser.Node, resolver *typeResolver) AnnotationModel {
	ann := AnnotationModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindQualifiedName:
			ann.Type = qualifiedNameToString(child)
		case parser.KindIdentifier:
			if child.Token != nil {
				ann.Type = child.Token.Literal
			}
		case parser.KindAnnotationElement:
			if ann.Values == nil {
				ann.Values = make(map[string]interface{})
			}
			name, value := annotationElementFromNode(child)
			ann.Values[name] = value
		}
	}

	if resolver != nil {
		ann.Type = resolver.resolve(ann.Type)
	}

	return ann
}

func annotationElementFromNode(node *parser.Node) (string, interface{}) {
	name := "value"
	var value interface{}

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			name = child.Token.Literal
		} else {
			value = annotationValueFromNode(child)
		}
	}

	return name, value
}

func annotationValueFromNode(node *parser.Node) interface{} {
	if node.Kind == parser.KindLiteral && node.Token != nil {
		return node.Token.Literal
	}
	if node.Kind == parser.KindIdentifier && node.Token != nil {
		return node.Token.Literal
	}
	if node.Kind == parser.KindAnnotation {
		return annotationModelFromNode(node, nil)
	}
	if node.Kind == parser.KindArrayInit {
		var values []interface{}
		for _, child := range node.Children {
			values = append(values, annotationValueFromNode(child))
		}
		return values
	}
	if node.Kind == parser.KindFieldAccess {
		var parts []string
		for _, child := range node.Children {
			if child.Token != nil {
				parts = append(parts, child.Token.Literal)
			}
		}
		return strings.Join(parts, ".")
	}
	return nil
}

// isEnumConstant checks if a KindFieldDecl node is an enum constant
// Enum constants have no Type child, only Identifier, optional Arguments, and optional ClassBody
func isEnumConstant(node *parser.Node) bool {
	if node.Kind != parser.KindFieldDecl {
		return false
	}
	// If it has a Type child, it's a regular field
	return node.FirstChildOfKind(parser.KindType) == nil
}

// enumConstantFromFieldDecl extracts enum constant data from a KindFieldDecl node
func enumConstantFromFieldDecl(node *parser.Node) EnumConstantModel {
	ec := EnumConstantModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				ec.Name = child.Token.Literal
			}
		case parser.KindParameters:
			// Enum constants can have arguments in a Parameters node
			ec.Arguments = argumentsFromParametersNode(child)
		}
	}

	return ec
}

// argumentsFromParametersNode extracts string representations of arguments from a Parameters node
func argumentsFromParametersNode(node *parser.Node) []string {
	var args []string
	for _, child := range node.Children {
		// In Parameters nodes, children can be literals or other expression nodes
		arg := argumentToString(child)
		if arg != "" {
			args = append(args, arg)
		}
	}
	return args
}

// argumentToString converts an argument expression node to a string representation
func argumentToString(node *parser.Node) string {
	if node.Token != nil {
		return node.Token.Literal
	}

	// For literals and simple identifiers
	if node.Kind == parser.KindLiteral || node.Kind == parser.KindIdentifier {
		if node.Token != nil {
			return node.Token.Literal
		}
	}

	// For qualified names and complex expressions
	var parts []string
	for _, child := range node.Children {
		s := argumentToString(child)
		if s != "" {
			parts = append(parts, s)
		}
	}

	if len(parts) > 0 {
		return strings.Join(parts, ".")
	}
	return ""
}
