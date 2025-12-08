package java

import (
	"github.com/dhamidi/javalyzer/java/parser"
)

// sourceClass holds class information extracted from source code
type sourceClass struct {
	name        string
	simpleName  string
	pkg         string
	superClass  string
	interfaces  []string
	visibility  string
	kind        string
	isAbstract  bool
	isFinal     bool
	fields      []sourceField
	methods     []sourceMethod
	sourceFile  string
	annotations []Annotation
}

type sourceField struct {
	name       string
	typeName   string
	arrayDepth int
	visibility string
	isStatic   bool
	isFinal    bool
}

type sourceMethod struct {
	name       string
	returnType Type
	parameters []Parameter
	visibility string
	isStatic   bool
	isFinal    bool
	isAbstract bool
}

// ClassFromNode creates a Class from a parser.Node (compilation unit).
func ClassFromNode(node *parser.Node) *Class {
	if node == nil || node.Kind != parser.KindCompilationUnit {
		return nil
	}

	sc := &sourceClass{}

	// Extract package
	if pkgDecl := node.FirstChildOfKind(parser.KindPackageDecl); pkgDecl != nil {
		sc.pkg = extractQualifiedName(pkgDecl)
	}

	// Find the first type declaration
	var typeDecl *parser.Node
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl,
			parser.KindRecordDecl, parser.KindAnnotationDecl:
			typeDecl = child
			break
		}
	}

	if typeDecl == nil {
		return nil
	}

	sc.extractTypeDecl(typeDecl)
	sc.sourceFile = node.Span.Start.File

	return &Class{source: sc}
}

func (sc *sourceClass) extractTypeDecl(node *parser.Node) {
	// Extract name from identifier child
	if ident := node.FirstChildOfKind(parser.KindIdentifier); ident != nil {
		sc.simpleName = ident.TokenLiteral()
	}

	// Build full name
	if sc.pkg != "" {
		sc.name = sc.pkg + "." + sc.simpleName
	} else {
		sc.name = sc.simpleName
	}

	// Determine kind
	switch node.Kind {
	case parser.KindClassDecl:
		sc.kind = "class"
	case parser.KindInterfaceDecl:
		sc.kind = "interface"
	case parser.KindEnumDecl:
		sc.kind = "enum"
	case parser.KindRecordDecl:
		sc.kind = "record"
	case parser.KindAnnotationDecl:
		sc.kind = "annotation"
	}

	// Extract modifiers
	if mods := node.FirstChildOfKind(parser.KindModifiers); mods != nil {
		sc.visibility, sc.isAbstract, sc.isFinal = extractModifiers(mods)
	}
	if sc.visibility == "" {
		sc.visibility = "package"
	}

	// Extract extends (superclass)
	sc.superClass = extractExtends(node)
	if sc.superClass == "" && sc.kind == "class" {
		sc.superClass = "java.lang.Object"
	}

	// Extract implements (interfaces)
	sc.interfaces = extractImplements(node)

	// Extract fields
	for _, child := range node.ChildrenOfKind(parser.KindFieldDecl) {
		sc.fields = append(sc.fields, extractField(child))
	}

	// Extract methods
	for _, child := range node.ChildrenOfKind(parser.KindMethodDecl) {
		sc.methods = append(sc.methods, extractMethod(child))
	}
	for _, child := range node.ChildrenOfKind(parser.KindConstructorDecl) {
		sc.methods = append(sc.methods, extractConstructor(child, sc.simpleName))
	}

	// Extract annotations
	for _, child := range node.ChildrenOfKind(parser.KindAnnotation) {
		sc.annotations = append(sc.annotations, extractAnnotation(child))
	}
}

func extractQualifiedName(node *parser.Node) string {
	// Look for identifier or qualified name
	if qn := node.FirstChildOfKind(parser.KindQualifiedName); qn != nil {
		return qn.TokenLiteral()
	}
	if ident := node.FirstChildOfKind(parser.KindIdentifier); ident != nil {
		return ident.TokenLiteral()
	}
	return ""
}

func extractModifiers(mods *parser.Node) (visibility string, isAbstract, isFinal bool) {
	for _, child := range mods.Children {
		if child.Token == nil {
			continue
		}
		switch child.Token.Literal {
		case "public":
			visibility = "public"
		case "protected":
			visibility = "protected"
		case "private":
			visibility = "private"
		case "abstract":
			isAbstract = true
		case "final":
			isFinal = true
		}
	}
	return
}

func extractExtends(node *parser.Node) string {
	for _, child := range node.Children {
		if child.Kind == parser.KindType {
			// Check if this is an extends clause by looking at previous siblings
			// For now, we look for the first Type after the class name
			return extractTypeName(child)
		}
	}
	return ""
}

func extractImplements(node *parser.Node) []string {
	var interfaces []string
	// Interfaces appear as Type nodes after extends
	types := node.ChildrenOfKind(parser.KindType)
	if len(types) > 1 {
		for _, t := range types[1:] {
			interfaces = append(interfaces, extractTypeName(t))
		}
	}
	return interfaces
}

func extractTypeName(node *parser.Node) string {
	if node == nil {
		return ""
	}
	if node.Token != nil {
		return node.Token.Literal
	}
	// Check for qualified name or identifier
	if qn := node.FirstChildOfKind(parser.KindQualifiedName); qn != nil {
		return qn.TokenLiteral()
	}
	if ident := node.FirstChildOfKind(parser.KindIdentifier); ident != nil {
		return ident.TokenLiteral()
	}
	// Check for parameterized type
	if pt := node.FirstChildOfKind(parser.KindParameterizedType); pt != nil {
		return extractTypeName(pt)
	}
	return ""
}

func extractField(node *parser.Node) sourceField {
	f := sourceField{visibility: "package"}

	if mods := node.FirstChildOfKind(parser.KindModifiers); mods != nil {
		f.visibility, _, f.isFinal = extractModifiers(mods)
		for _, child := range mods.Children {
			if child.Token != nil && child.Token.Literal == "static" {
				f.isStatic = true
			}
		}
	}
	if f.visibility == "" {
		f.visibility = "package"
	}

	if typeNode := node.FirstChildOfKind(parser.KindType); typeNode != nil {
		f.typeName = extractTypeName(typeNode)
		// Check for array type
		if at := typeNode.FirstChildOfKind(parser.KindArrayType); at != nil {
			f.arrayDepth = countArrayDimensions(at)
		}
	}

	if ident := node.FirstChildOfKind(parser.KindIdentifier); ident != nil {
		f.name = ident.TokenLiteral()
	}

	return f
}

func countArrayDimensions(node *parser.Node) int {
	// Count nested array types
	count := 1
	for _, child := range node.Children {
		if child.Kind == parser.KindArrayType {
			count += countArrayDimensions(child)
		}
	}
	return count
}

func extractMethod(node *parser.Node) sourceMethod {
	m := sourceMethod{visibility: "package"}

	if mods := node.FirstChildOfKind(parser.KindModifiers); mods != nil {
		m.visibility, m.isAbstract, m.isFinal = extractModifiers(mods)
		for _, child := range mods.Children {
			if child.Token != nil && child.Token.Literal == "static" {
				m.isStatic = true
			}
		}
	}
	if m.visibility == "" {
		m.visibility = "package"
	}

	// Return type is the first Type child
	if typeNode := node.FirstChildOfKind(parser.KindType); typeNode != nil {
		m.returnType = Type{Name: extractTypeName(typeNode)}
		if at := typeNode.FirstChildOfKind(parser.KindArrayType); at != nil {
			m.returnType.ArrayDepth = countArrayDimensions(at)
		}
	}

	if ident := node.FirstChildOfKind(parser.KindIdentifier); ident != nil {
		m.name = ident.TokenLiteral()
	}

	if params := node.FirstChildOfKind(parser.KindParameters); params != nil {
		m.parameters = extractParameters(params)
	}

	return m
}

func extractConstructor(node *parser.Node, className string) sourceMethod {
	m := sourceMethod{
		name:       "<init>",
		returnType: Type{Name: "void"},
		visibility: "package",
	}

	if mods := node.FirstChildOfKind(parser.KindModifiers); mods != nil {
		m.visibility, _, _ = extractModifiers(mods)
	}
	if m.visibility == "" {
		m.visibility = "package"
	}

	if params := node.FirstChildOfKind(parser.KindParameters); params != nil {
		m.parameters = extractParameters(params)
	}

	return m
}

func extractParameters(params *parser.Node) []Parameter {
	var result []Parameter
	for _, p := range params.ChildrenOfKind(parser.KindParameter) {
		param := Parameter{}

		if typeNode := p.FirstChildOfKind(parser.KindType); typeNode != nil {
			param.Type = Type{Name: extractTypeName(typeNode)}
			if at := typeNode.FirstChildOfKind(parser.KindArrayType); at != nil {
				param.Type.ArrayDepth = countArrayDimensions(at)
			}
		}

		if ident := p.FirstChildOfKind(parser.KindIdentifier); ident != nil {
			param.Name = ident.TokenLiteral()
		}

		result = append(result, param)
	}
	return result
}

func extractAnnotation(node *parser.Node) Annotation {
	ann := Annotation{}
	if ident := node.FirstChildOfKind(parser.KindIdentifier); ident != nil {
		ann.Type = ident.TokenLiteral()
	} else if qn := node.FirstChildOfKind(parser.KindQualifiedName); qn != nil {
		ann.Type = qn.TokenLiteral()
	}
	return ann
}
