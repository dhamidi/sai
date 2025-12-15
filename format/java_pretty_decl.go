package format

import (
	"github.com/dhamidi/sai/java/parser"
)

func (p *JavaPrettyPrinter) printPackageDecl(node *parser.Node) {
	// Print any annotations first
	for _, child := range node.Children {
		if child.Kind == parser.KindAnnotation {
			p.writeIndent()
			p.printAnnotation(child)
			p.write("\n")
			p.atLineStart = true
		}
	}
	p.writeIndent()
	p.write("package ")
	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			p.printQualifiedName(child)
		}
	}
	p.write(";\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printImportDecl(node *parser.Node) {
	p.writeIndent()
	p.write("import ")
	isStatic := false
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == "static" {
			isStatic = true
		}
	}
	if isStatic {
		p.write("static ")
	}
	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			p.printQualifiedName(child)
		} else if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == "*" {
			p.write(".*")
		}
	}
	p.write(";\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printQualifiedName(node *parser.Node) {
	p.write(p.qualifiedNameString(node))
}

func (p *JavaPrettyPrinter) qualifiedNameString(node *parser.Node) string {
	var parts []string
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			parts = append(parts, child.Token.Literal)
		}
	}
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += "."
		}
		result += part
	}
	return result
}

func (p *JavaPrettyPrinter) printModuleDecl(node *parser.Node) {
	// Print annotations first
	for _, child := range node.Children {
		if child.Kind == parser.KindAnnotation {
			p.writeIndent()
			p.printAnnotation(child)
			p.write("\n")
			p.atLineStart = true
		}
	}

	p.writeIndent()

	// Check for "open" modifier
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == "open" {
			p.write("open ")
			break
		}
	}

	p.write("module ")

	// Print module name
	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			p.printQualifiedName(child)
			break
		}
	}

	p.write(" {\n")
	p.atLineStart = true
	p.indent++

	// Print directives
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindRequiresDirective,
			parser.KindExportsDirective,
			parser.KindOpensDirective,
			parser.KindUsesDirective,
			parser.KindProvidesDirective:
			p.printNode(child)
		}
	}

	p.indent--
	p.writeIndent()
	p.write("}\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printRequiresDirective(node *parser.Node) {
	p.writeIndent()
	p.write("requires ")

	// Print modifiers (transitive, static)
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			literal := child.Token.Literal
			if literal == "transitive" || literal == "static" {
				p.write(literal + " ")
			}
		}
	}

	// Print module name
	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			p.printQualifiedName(child)
			break
		}
	}

	p.write(";\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printExportsDirective(node *parser.Node) {
	p.writeIndent()
	p.write("exports ")

	// Collect qualified names - first is package, rest are "to" targets
	var names []*parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			names = append(names, child)
		}
	}

	if len(names) > 0 {
		p.printQualifiedName(names[0])
	}

	if len(names) > 1 {
		p.write(" to ")
		for i := 1; i < len(names); i++ {
			if i > 1 {
				p.write(", ")
			}
			p.printQualifiedName(names[i])
		}
	}

	p.write(";\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printOpensDirective(node *parser.Node) {
	p.writeIndent()
	p.write("opens ")

	// Collect qualified names - first is package, rest are "to" targets
	var names []*parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			names = append(names, child)
		}
	}

	if len(names) > 0 {
		p.printQualifiedName(names[0])
	}

	if len(names) > 1 {
		p.write(" to ")
		for i := 1; i < len(names); i++ {
			if i > 1 {
				p.write(", ")
			}
			p.printQualifiedName(names[i])
		}
	}

	p.write(";\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printUsesDirective(node *parser.Node) {
	p.writeIndent()
	p.write("uses ")

	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			p.printQualifiedName(child)
			break
		}
	}

	p.write(";\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printProvidesDirective(node *parser.Node) {
	p.writeIndent()
	p.write("provides ")

	// Collect qualified names - first is service, rest are implementations
	var names []*parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			names = append(names, child)
		}
	}

	if len(names) > 0 {
		p.printQualifiedName(names[0])
	}

	if len(names) > 1 {
		p.write(" with ")
		for i := 1; i < len(names); i++ {
			if i > 1 {
				p.write(", ")
			}
			p.printQualifiedName(names[i])
		}
	}

	p.write(";\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printClassDecl(node *parser.Node) {
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	p.write("class ")

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
			break
		}
	}

	p.printTypeParameters(node)
	p.printExtendsImplements(node)

	p.write(" {\n")
	p.atLineStart = true
	p.indent++

	p.printClassBody(node)

	p.indent--
	p.writeIndent()
	p.write("}\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printInterfaceDecl(node *parser.Node) {
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	p.write("interface ")

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
			break
		}
	}

	p.printTypeParameters(node)
	p.printExtendsImplements(node)

	p.write(" {\n")
	p.atLineStart = true
	p.indent++

	p.printClassBody(node)

	p.indent--
	p.writeIndent()
	p.write("}\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printEnumDecl(node *parser.Node) {
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	p.write("enum ")

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
			break
		}
	}

	p.printExtendsImplements(node)

	p.write(" {\n")
	p.atLineStart = true
	p.indent++

	p.printEnumBody(node)

	p.indent--
	p.writeIndent()
	p.write("}\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printRecordDecl(node *parser.Node) {
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	p.write("record ")

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
			break
		}
	}

	p.printTypeParameters(node)

	params := node.FirstChildOfKind(parser.KindParameters)
	if params != nil {
		p.printParameters(params)
	}

	p.printExtendsImplements(node)

	p.write(" {\n")
	p.atLineStart = true
	p.indent++

	p.printClassBody(node)

	p.indent--
	p.writeIndent()
	p.write("}\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printAnnotationDecl(node *parser.Node) {
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	p.write("@interface ")

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
			break
		}
	}

	p.write(" {\n")
	p.atLineStart = true
	p.indent++

	p.printClassBody(node)

	p.indent--
	p.writeIndent()
	p.write("}\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printModifiers(node *parser.Node) {
	for _, child := range node.Children {
		if child.Kind == parser.KindComment || child.Kind == parser.KindLineComment {
			continue
		}
		if child.Kind == parser.KindAnnotation {
			p.printAnnotation(child)
			// Emit any trailing line comment on the same line as the annotation
			p.emitTrailingLineComment(child.Span.End.Line)
			p.write("\n")
			p.atLineStart = true
			p.writeIndent()
		} else if child.Token != nil {
			p.write(child.Token.Literal)
			p.write(" ")
		}
	}
}

// printModifiersInline prints modifiers inline (without newlines after annotations)
func (p *JavaPrettyPrinter) printModifiersInline(node *parser.Node) {
	for _, child := range node.Children {
		if child.Kind == parser.KindComment || child.Kind == parser.KindLineComment {
			continue
		}
		if child.Kind == parser.KindAnnotation {
			p.printAnnotation(child)
			p.write(" ")
		} else if child.Token != nil {
			p.write(child.Token.Literal)
			p.write(" ")
		}
	}
}

func (p *JavaPrettyPrinter) printAnnotation(node *parser.Node) {
	p.write("@")
	hasValue := false
	for _, child := range node.Children {
		if child.Kind == parser.KindQualifiedName {
			p.printQualifiedName(child)
		} else if child.Kind == parser.KindAnnotationElement {
			if !hasValue {
				p.write("(")
			}
			p.printAnnotationElements(node)
			p.write(")")
			hasValue = true
			break
		} else if child.Kind == parser.KindArrayInit || child.Kind == parser.KindLiteral || child.Kind == parser.KindFieldAccess || child.Kind == parser.KindIdentifier || child.Kind == parser.KindBinaryExpr {
			// Single value annotation like @SuppressWarnings({"unchecked"}), @Value("x"), @Retention(RetentionPolicy.SOURCE), @Retention(SOURCE), or @Name(Type.PREFIX + "Suffix")
			p.write("(")
			p.printAnnotationValue(child)
			p.write(")")
			hasValue = true
		}
	}
}

func (p *JavaPrettyPrinter) printAnnotationElements(node *parser.Node) {
	first := true
	for _, child := range node.Children {
		if child.Kind == parser.KindAnnotationElement {
			if !first {
				p.write(", ")
			}
			p.printAnnotationElement(child)
			first = false
		}
	}
}

func (p *JavaPrettyPrinter) printAnnotationElement(node *parser.Node) {
	var name string
	var valueNode *parser.Node

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			name = child.Token.Literal
		} else {
			valueNode = child
		}
	}

	// Always print the element name if present, for round-trip correctness.
	// Even though "value=" can be omitted in Java, we preserve it to match original source.
	if name != "" {
		p.write(name)
		p.write(" = ")
	}

	if valueNode != nil {
		p.printAnnotationValue(valueNode)
	}
}

func (p *JavaPrettyPrinter) printAnnotationValue(node *parser.Node) {
	switch node.Kind {
	case parser.KindLiteral:
		if node.Token != nil {
			p.write(node.Token.Literal)
		}
	case parser.KindIdentifier:
		if node.Token != nil {
			p.write(node.Token.Literal)
		}
	case parser.KindFieldAccess:
		p.printFieldAccess(node)
	case parser.KindAnnotation:
		p.printAnnotation(node)
	case parser.KindArrayInit:
		p.write("{")
		first := true
		for _, child := range node.Children {
			if !first {
				p.write(", ")
			}
			p.printAnnotationValue(child)
			first = false
		}
		p.write("}")
	default:
		p.printGenericExpr(node)
	}
}

func (p *JavaPrettyPrinter) printTypeParameters(node *parser.Node) {
	tp := node.FirstChildOfKind(parser.KindTypeParameters)
	if tp == nil {
		return
	}

	p.write("<")
	first := true
	for _, child := range tp.Children {
		if child.Kind == parser.KindTypeParameter {
			if !first {
				p.write(", ")
			}
			p.printTypeParameter(child)
			first = false
		}
	}
	p.write(">")
}

func (p *JavaPrettyPrinter) printTypeParameter(node *parser.Node) {
	firstBound := true
	for _, child := range node.Children {
		if child.Kind == parser.KindAnnotation {
			// Type parameter annotations like <@MyAnnotation T>
			p.printAnnotation(child)
			p.write(" ")
		} else if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
		} else if child.Kind == parser.KindType {
			// For intersection types like <T extends A & B>, first bound uses "extends",
			// subsequent bounds use "&".
			if firstBound {
				p.write(" extends ")
				firstBound = false
			} else {
				p.write(" & ")
			}
			p.printType(child)
		}
	}
}

func (p *JavaPrettyPrinter) printExtendsImplements(node *parser.Node) {
	extendsClause := node.FirstChildOfKind(parser.KindExtendsClause)
	implementsClause := node.FirstChildOfKind(parser.KindImplementsClause)

	if extendsClause != nil {
		p.write(" extends ")
		types := extendsClause.ChildrenOfKind(parser.KindType)
		for i, t := range types {
			if i > 0 {
				p.write(", ")
			}
			p.printType(t)
		}
	}

	if implementsClause != nil {
		if node.Kind == parser.KindInterfaceDecl {
			if extendsClause == nil {
				p.write(" extends ")
			} else {
				p.write(", ")
			}
		} else {
			p.write(" implements ")
		}
		types := implementsClause.ChildrenOfKind(parser.KindType)
		for i, t := range types {
			if i > 0 {
				p.write(", ")
			}
			p.printType(t)
		}
	}

	permitsClause := node.FirstChildOfKind(parser.KindPermitsClause)
	if permitsClause != nil {
		types := permitsClause.ChildrenOfKind(parser.KindType)
		if len(types) > 3 {
			// Long permits clause: format on separate lines in groups of 3
			p.write("\n")
			p.write("        permits ")
			for i, t := range types {
				if i > 0 {
					if i%3 == 0 {
						// Start a new line every 3 types
						p.write(",\n                ")
					} else {
						p.write(", ")
					}
				}
				p.printType(t)
			}
		} else {
			// Short permits clause: keep on same line
			p.write(" permits ")
			for i, t := range types {
				if i > 0 {
					p.write(", ")
				}
				p.printType(t)
			}
		}
	}
}

func (p *JavaPrettyPrinter) printType(node *parser.Node) {
	if node.Kind == parser.KindArrayType {
		p.printArrayType(node)
		return
	}

	if node.Token != nil {
		p.write(node.Token.Literal)
		return
	}

	// Track if we've seen type arguments - if so, subsequent type names need a dot separator
	// e.g., Outer<String>.Inner -> Outer, <String>, .Inner
	sawTypeArgs := false
	// Track if we've seen a KindType child (for intersection types in casts)
	sawType := false
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if sawTypeArgs {
				p.write(".")
			}
			if child.Token != nil {
				p.write(child.Token.Literal)
			}
		case parser.KindQualifiedName:
			if sawTypeArgs {
				p.write(".")
			}
			p.printQualifiedName(child)
		case parser.KindParameterizedType:
			// Intersection types in casts: (Type1 & Type2)
			if sawType {
				p.write(" & ")
			}
			p.printParameterizedType(child)
			sawType = true
		case parser.KindArrayType:
			p.printArrayType(child)
		case parser.KindType:
			// Intersection types in casts: (Type1 & Type2)
			if sawType {
				p.write(" & ")
			}
			p.printType(child)
			sawType = true
		case parser.KindTypeArguments:
			p.printTypeArguments(child)
			sawTypeArgs = true
		}
	}
}

func (p *JavaPrettyPrinter) printParameterizedType(node *parser.Node) {
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				p.write(child.Token.Literal)
			}
		case parser.KindQualifiedName:
			p.printQualifiedName(child)
		case parser.KindTypeArguments:
			p.printTypeArguments(child)
		}
	}
}

func (p *JavaPrettyPrinter) printTypeArguments(node *parser.Node) {
	p.write("<")
	first := true
	for _, child := range node.Children {
		if child.Kind == parser.KindTypeArgument || child.Kind == parser.KindType || child.Kind == parser.KindArrayType || child.Kind == parser.KindWildcard {
			if !first {
				p.write(", ")
			}
			p.printTypeArgument(child)
			first = false
		}
	}
	p.write(">")
}

func (p *JavaPrettyPrinter) printTypeArgument(node *parser.Node) {
	if node.Kind == parser.KindType {
		p.printType(node)
		return
	}
	if node.Kind == parser.KindArrayType {
		p.printArrayType(node)
		return
	}
	if node.Kind == parser.KindWildcard {
		p.printWildcard(node)
		return
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			p.printType(child)
		case parser.KindArrayType:
			p.printArrayType(child)
		case parser.KindWildcard:
			p.printWildcard(child)
		}
	}
}

func (p *JavaPrettyPrinter) printWildcard(node *parser.Node) {
	p.write("?")
	for _, child := range node.Children {
		if child.Token != nil {
			if child.Token.Literal == "extends" || child.Token.Literal == "super" {
				p.write(" ")
				p.write(child.Token.Literal)
				p.write(" ")
			}
		}
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			p.printType(child)
		}
	}
}

func (p *JavaPrettyPrinter) printArrayType(node *parser.Node) {
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType, parser.KindArrayType, parser.KindIdentifier, parser.KindQualifiedName:
			p.printType(child)
		case parser.KindFieldAccess:
			p.printFieldAccess(child)
		default:
			continue
		}
		break
	}
	p.write("[]")
}

func (p *JavaPrettyPrinter) printClassBody(node *parser.Node) {
	block := node.FirstChildOfKind(parser.KindBlock)
	if block != nil {
		p.printClassBodyMembers(block)
		return
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindFieldDecl, parser.KindMethodDecl, parser.KindConstructorDecl,
			parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl, parser.KindAnnotationDecl:
			p.printClassBodyMember(child)
		}
	}
}

func (p *JavaPrettyPrinter) printClassBodyMembers(block *parser.Node) {
	for _, child := range block.Children {
		p.printClassBodyMember(child)
	}
}

func (p *JavaPrettyPrinter) printClassBodyMember(child *parser.Node) {
	// Find lines that have annotations - we'll skip line comments on those lines
	// so they can be emitted as trailing comments after the annotation
	annotationLines := make(map[int]bool)
	modifiers := child.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		for _, mod := range modifiers.Children {
			if mod.Kind == parser.KindAnnotation {
				annotationLines[mod.Span.End.Line] = true
			}
		}
	}
	p.emitCommentsBeforeLineSkippingAnnotationLines(child.Span.Start.Line, annotationLines)
	p.printNode(child)
}

func (p *JavaPrettyPrinter) printEnumBody(node *parser.Node) {
	var members []*parser.Node
	var constants []*parser.Node

	// Enum constants and members are direct children of the EnumDecl node.
	// Do NOT look for a Block child - any Block inside an enum is a static initializer,
	// not a body wrapper.
	for _, child := range node.Children {
		if child.Kind == parser.KindFieldDecl && p.isEnumConstant(child) {
			constants = append(constants, child)
		} else if child.Kind == parser.KindFieldDecl || child.Kind == parser.KindMethodDecl ||
			child.Kind == parser.KindConstructorDecl || child.Kind == parser.KindClassDecl ||
			child.Kind == parser.KindInterfaceDecl || child.Kind == parser.KindEnumDecl ||
			child.Kind == parser.KindRecordDecl || child.Kind == parser.KindAnnotationDecl ||
			child.Kind == parser.KindBlock {
			// KindBlock here is a static initializer, treat it as a member
			members = append(members, child)
		}
	}

	for i, c := range constants {
		p.emitCommentsBeforeLine(c.Span.Start.Line)
		p.writeIndent()
		p.printEnumConstant(c)
		if i < len(constants)-1 {
			p.write(",")
			p.emitTrailingLineComment(c.Span.End.Line)
			p.write("\n")
		} else if len(members) > 0 {
			p.write(";")
			p.emitTrailingLineComment(c.Span.End.Line)
			// Only one newline here - printMethodDecl/printFieldDecl adds leading newline
			p.write("\n")
		} else {
			p.emitTrailingLineComment(c.Span.End.Line)
			p.write("\n")
		}
		p.atLineStart = true
	}

	for _, m := range members {
		p.emitCommentsBeforeLine(m.Span.Start.Line)
		p.printNode(m)
	}
}

func (p *JavaPrettyPrinter) isEnumConstant(node *parser.Node) bool {
	if node.Kind != parser.KindFieldDecl {
		return false
	}
	// Enum constants don't have a type - they're just identifiers with optional arguments.
	// Regular fields have either KindType or KindArrayType for array types.
	hasType := node.FirstChildOfKind(parser.KindType) != nil ||
		node.FirstChildOfKind(parser.KindArrayType) != nil
	return !hasType
}

func (p *JavaPrettyPrinter) printEnumConstant(node *parser.Node) {
	// Print annotations first
	for _, child := range node.Children {
		if child.Kind == parser.KindAnnotation {
			p.printAnnotation(child)
			p.write("\n")
			p.atLineStart = true
			p.writeIndent()
		}
	}

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
		} else if child.Kind == parser.KindParameters {
			p.write("(")
			p.printArguments(child)
			p.write(")")
		} else if child.Kind == parser.KindBlock {
			// Enum constant with body (anonymous class-like)
			p.write(" {\n")
			p.indent++
			p.printClassBodyMembers(child)
			p.indent--
			p.writeIndent()
			p.write("}")
		}
	}
}

func (p *JavaPrettyPrinter) printFieldDecl(node *parser.Node) {
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	// Find the type
	var fieldType *parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			fieldType = child
			break
		}
	}

	if fieldType != nil {
		p.printType(fieldType)
		p.write(" ")
	}

	// Print declarators: each is an identifier optionally followed by an initializer.
	// The AST flattens everything, so we need to use source positions to distinguish
	// between variable names and initializers when the initializer is also an Identifier.
	// Between a variable name and its initializer there's '=' in the source.
	// Between two variable names (when there's no initializer) there's ',' in the source.
	first := true
	i := 0
	var prevChild *parser.Node
	prevWasName := false // Track if previous child was a variable name
	for i < len(node.Children) {
		child := node.Children[i]
		if child.Kind == parser.KindModifiers || child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			i++
			continue
		}

		// Check if this child is an initializer for the previous variable by looking
		// for '=' in the source between them. Only check if previous child was an
		// identifier (variable name) - if it was an initializer, current must be new var.
		isInitializer := false
		if prevChild != nil && prevWasName {
			isInitializer = p.hasAssignBetween(prevChild, child)
		}

		if isInitializer {
			p.write(" = ")
			if p.shouldFormatAsChain(child) {
				p.printMethodChain(child, p.indent)
			} else {
				p.printExpr(child)
			}
			prevWasName = false
		} else {
			// This is a new variable name
			if !first {
				p.write(", ")
			}
			first = false
			if child.Kind == parser.KindIdentifier && child.Token != nil {
				p.write(child.Token.Literal)
			} else {
				p.printExpr(child)
			}
			prevWasName = (child.Kind == parser.KindIdentifier)
		}
		prevChild = child
		i++
	}

	p.write(";\n")
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printMethodDecl(node *parser.Node) {
	p.write("\n")
	p.atLineStart = true
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	typeParams := node.FirstChildOfKind(parser.KindTypeParameters)
	if typeParams != nil {
		p.write("<")
		first := true
		for _, child := range typeParams.Children {
			if child.Kind == parser.KindTypeParameter {
				if !first {
					p.write(", ")
				}
				p.printTypeParameter(child)
				first = false
			}
		}
		p.write("> ")
	}

	var returnType *parser.Node
	var name string
	var params *parser.Node
	var throwsList *parser.Node
	var body *parser.Node
	var defaultValue *parser.Node // For annotation methods
	seenParams := false

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType, parser.KindArrayType:
			returnType = child
		case parser.KindIdentifier:
			if child.Token != nil {
				if !seenParams {
					// Before parameters, this is the method name
					name = child.Token.Literal
				} else {
					// After parameters, this is a default value (e.g., constant reference)
					defaultValue = child
				}
			}
		case parser.KindParameters:
			params = child
			seenParams = true
		case parser.KindThrowsList:
			throwsList = child
		case parser.KindBlock:
			body = child
		case parser.KindLiteral, parser.KindArrayInit, parser.KindFieldAccess, parser.KindAnnotation:
			// Annotation method default value
			defaultValue = child
		}
	}

	if returnType != nil {
		p.printType(returnType)
		p.write(" ")
	}

	p.write(name)
	p.printParameters(params)

	if throwsList != nil {
		p.write(" throws ")
		p.printThrowsList(throwsList)
	}

	if body != nil {
		p.write(" ")
		p.printBlock(body)
	} else if defaultValue != nil {
		p.write(" default ")
		p.printExpr(defaultValue)
		p.write(";\n")
	} else {
		p.write(";\n")
	}
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printConstructorDecl(node *parser.Node) {
	p.write("\n")
	p.atLineStart = true
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	typeParams := node.FirstChildOfKind(parser.KindTypeParameters)
	if typeParams != nil {
		p.write("<")
		first := true
		for _, child := range typeParams.Children {
			if child.Kind == parser.KindTypeParameter {
				if !first {
					p.write(", ")
				}
				p.printTypeParameter(child)
				first = false
			}
		}
		p.write("> ")
	}

	var name string
	var params *parser.Node
	var throwsList *parser.Node
	var body *parser.Node

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				name = child.Token.Literal
			}
		case parser.KindParameters:
			params = child
		case parser.KindThrowsList:
			throwsList = child
		case parser.KindBlock:
			body = child
		}
	}

	p.write(name)
	p.printParameters(params)

	if throwsList != nil {
		p.write(" throws ")
		p.printThrowsList(throwsList)
	}

	if body != nil {
		p.write(" ")
		p.printBlock(body)
	}
	p.atLineStart = true
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printParameters(node *parser.Node) {
	if node == nil {
		p.write("()")
		return
	}

	p.write("(")
	first := true
	for _, child := range node.Children {
		if child.Kind == parser.KindParameter {
			if !first {
				p.write(", ")
			}
			p.printParameter(child)
			first = false
		} else if child.Kind == parser.KindIdentifier && child.Token != nil {
			// Lambda parameters can be simple identifiers without types
			if !first {
				p.write(", ")
			}
			p.write(child.Token.Literal)
			first = false
		}
	}
	p.write(")")
}

func (p *JavaPrettyPrinter) printParameter(node *parser.Node) {
	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	var paramType *parser.Node
	var name string
	var isVarargs bool

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType, parser.KindArrayType:
			paramType = child
		case parser.KindIdentifier:
			if child.Token != nil {
				// Varargs ellipsis is stored as an Identifier child with TokenEllipsis
				if child.Token.Kind == parser.TokenEllipsis {
					isVarargs = true
				} else {
					name = child.Token.Literal
				}
			}
		}
	}

	if paramType != nil {
		p.printType(paramType)
	}
	if isVarargs {
		p.write("...")
	}
	if name != "" {
		p.write(" ")
		p.write(name)
	}
}

func (p *JavaPrettyPrinter) printThrowsList(node *parser.Node) {
	first := true
	for _, child := range node.Children {
		if child.Kind == parser.KindType {
			if !first {
				p.write(", ")
			}
			p.printType(child)
			first = false
		}
	}
}
