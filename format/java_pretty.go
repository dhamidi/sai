package format

import (
	"bytes"
	"io"
	"sort"

	"github.com/dhamidi/sai/java/parser"
)

type JavaPrettyPrinter struct {
	w            io.Writer
	source       []byte
	comments     []parser.Token
	commentIndex int
	indent       int
	indentStr    string
	atLineStart  bool
	lastLine     int
}

func NewJavaPrettyPrinter(w io.Writer) *JavaPrettyPrinter {
	return &JavaPrettyPrinter{
		w:           w,
		indentStr:   "    ",
		atLineStart: true,
		lastLine:    1,
	}
}

func (p *JavaPrettyPrinter) Print(node *parser.Node, source []byte, comments []parser.Token) error {
	p.source = source
	p.comments = comments
	sort.Slice(p.comments, func(i, j int) bool {
		if p.comments[i].Span.Start.Line != p.comments[j].Span.Start.Line {
			return p.comments[i].Span.Start.Line < p.comments[j].Span.Start.Line
		}
		return p.comments[i].Span.Start.Column < p.comments[j].Span.Start.Column
	})
	p.commentIndex = 0

	p.printNode(node)
	p.emitRemainingComments()
	return nil
}

func (p *JavaPrettyPrinter) printNode(node *parser.Node) {
	switch node.Kind {
	case parser.KindCompilationUnit:
		p.printCompilationUnit(node)
	case parser.KindPackageDecl:
		p.printPackageDecl(node)
	case parser.KindImportDecl:
		p.printImportDecl(node)
	case parser.KindClassDecl:
		p.printClassDecl(node)
	case parser.KindInterfaceDecl:
		p.printInterfaceDecl(node)
	case parser.KindEnumDecl:
		p.printEnumDecl(node)
	case parser.KindRecordDecl:
		p.printRecordDecl(node)
	case parser.KindAnnotationDecl:
		p.printAnnotationDecl(node)
	case parser.KindFieldDecl:
		p.printFieldDecl(node)
	case parser.KindMethodDecl:
		p.printMethodDecl(node)
	case parser.KindConstructorDecl:
		p.printConstructorDecl(node)
	case parser.KindBlock:
		if p.isStaticInitializer(node) {
			p.printStaticInitializer(node)
		} else {
			p.printBlock(node)
		}
	default:
		p.printGenericNode(node)
	}
}

// isStaticInitializer checks if a KindBlock represents a static initializer
// (has structure: KindIdentifier("static") + KindBlock)
func (p *JavaPrettyPrinter) isStaticInitializer(node *parser.Node) bool {
	if node.Kind != parser.KindBlock {
		return false
	}
	hasStatic := false
	hasBlock := false
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == "static" {
			hasStatic = true
		} else if child.Kind == parser.KindBlock {
			hasBlock = true
		}
	}
	return hasStatic && hasBlock
}

func (p *JavaPrettyPrinter) printStaticInitializer(node *parser.Node) {
	p.writeIndent()
	p.write("static ")
	for _, child := range node.Children {
		if child.Kind == parser.KindBlock {
			p.printBlock(child)
			return
		}
	}
}

func (p *JavaPrettyPrinter) printCompilationUnit(node *parser.Node) {
	for _, child := range node.Children {
		p.emitCommentsBeforeLine(child.Span.Start.Line)
		p.printNode(child)
	}
}

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
		} else if child.Kind == parser.KindArrayInit || child.Kind == parser.KindLiteral || child.Kind == parser.KindFieldAccess || child.Kind == parser.KindIdentifier {
			// Single value annotation like @SuppressWarnings({"unchecked"}), @Value("x"), @Retention(RetentionPolicy.SOURCE), or @Retention(SOURCE)
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
		p.write(" permits ")
		types := permitsClause.ChildrenOfKind(parser.KindType)
		for i, t := range types {
			if i > 0 {
				p.write(", ")
			}
			p.printType(t)
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
			p.write(",\n")
		} else if len(members) > 0 {
			p.write(";\n\n")
		} else {
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
			p.printExpr(child)
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

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType, parser.KindArrayType:
			returnType = child
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

func (p *JavaPrettyPrinter) printBlock(node *parser.Node) {
	p.write("{\n")
	p.atLineStart = true
	p.indent++

	for _, child := range node.Children {
		// Skip comment nodes - they're handled by emitCommentsBeforeLine
		if child.Kind == parser.KindLineComment || child.Kind == parser.KindComment {
			continue
		}
		p.emitCommentsBeforeLine(child.Span.Start.Line)
		p.printStatement(child)
	}

	// Emit any comments inside the block before the closing brace
	p.emitCommentsBeforeLine(node.Span.End.Line)

	p.indent--
	p.writeIndent()
	p.write("}\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printStatement(node *parser.Node) {
	switch node.Kind {
	case parser.KindBlock:
		p.writeIndent()
		p.printBlock(node)
	case parser.KindLocalVarDecl:
		p.printLocalVarDecl(node)
	case parser.KindExplicitConstructorInvocation:
		p.printExplicitConstructorInvocation(node)
	case parser.KindExprStmt:
		p.printExprStmt(node)
	case parser.KindReturnStmt:
		p.printReturnStmt(node)
	case parser.KindIfStmt:
		p.printIfStmt(node)
	case parser.KindForStmt:
		p.printForStmt(node)
	case parser.KindEnhancedForStmt:
		p.printEnhancedForStmt(node)
	case parser.KindWhileStmt:
		p.printWhileStmt(node)
	case parser.KindDoStmt:
		p.printDoStmt(node)
	case parser.KindSwitchStmt:
		p.printSwitchStmt(node)
	case parser.KindTryStmt:
		p.printTryStmt(node)
	case parser.KindThrowStmt:
		p.printThrowStmt(node)
	case parser.KindBreakStmt:
		p.printBreakStmt(node)
	case parser.KindContinueStmt:
		p.printContinueStmt(node)
	case parser.KindSynchronizedStmt:
		p.printSynchronizedStmt(node)
	case parser.KindAssertStmt:
		p.printAssertStmt(node)
	case parser.KindEmptyStmt:
		p.writeIndent()
		p.write(";\n")
		p.atLineStart = true
	case parser.KindLabeledStmt:
		p.printLabeledStmt(node)
	case parser.KindYieldStmt:
		p.printYieldStmt(node)
	default:
		p.writeIndent()
		p.printGenericNode(node)
		p.write(";\n")
		p.atLineStart = true
	}
	p.lastLine = node.Span.End.Line
}

func (p *JavaPrettyPrinter) printLocalVarDecl(node *parser.Node) {
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	// Find the type
	var varType *parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			varType = child
			break
		}
	}

	if varType != nil {
		p.printType(varType)
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
	prevWasName := false // Track if previous child was a variable name (Identifier)
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
			p.printExpr(child)
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
}

// hasAssignBetween checks if there's an '=' token in the source between two nodes.
// This is used to distinguish variable names from initializers in local var declarations.
func (p *JavaPrettyPrinter) hasAssignBetween(prev, next *parser.Node) bool {
	if p.source == nil {
		return false
	}
	startOffset := prev.Span.End.Offset
	endOffset := next.Span.Start.Offset
	if startOffset < 0 || endOffset < 0 || startOffset >= endOffset {
		return false
	}
	if endOffset > len(p.source) {
		endOffset = len(p.source)
	}
	between := p.source[startOffset:endOffset]
	// Look for '=' that isn't part of '==' or '!=' or '<=' etc.
	for i := 0; i < len(between); i++ {
		if between[i] == '=' {
			// Check it's not == or part of another operator
			if i > 0 && (between[i-1] == '=' || between[i-1] == '!' || between[i-1] == '<' || between[i-1] == '>') {
				continue
			}
			if i+1 < len(between) && between[i+1] == '=' {
				continue
			}
			return true
		}
	}
	return false
}

func (p *JavaPrettyPrinter) printExprStmt(node *parser.Node) {
	p.writeIndent()
	for _, child := range node.Children {
		p.printExpr(child)
	}
	p.write(";\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printReturnStmt(node *parser.Node) {
	p.writeIndent()
	p.write("return")
	for _, child := range node.Children {
		p.write(" ")
		p.printExpr(child)
	}
	p.write(";\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printIfStmt(node *parser.Node) {
	p.printIfStmtInner(node, true)
}

func (p *JavaPrettyPrinter) printIfStmtInner(node *parser.Node, writeIndent bool) {
	if writeIndent {
		p.writeIndent()
	}
	p.write("if (")

	children := node.Children
	if len(children) > 0 {
		p.printExpr(children[0])
	}
	p.write(") ")

	if len(children) > 1 {
		p.printBranchBodyWithElse(children[1], children)
	}
}

func (p *JavaPrettyPrinter) printBranchBodyWithElse(body *parser.Node, ifChildren []*parser.Node) {
	if body.Kind == parser.KindBlock {
		p.write("{\n")
		p.atLineStart = true
		p.indent++

		for _, child := range body.Children {
			p.emitCommentsBeforeLine(child.Span.Start.Line)
			p.printStatement(child)
		}

		p.indent--
		p.writeIndent()
		p.write("}")

		if len(ifChildren) > 2 {
			elseBody := ifChildren[2]
			if elseBody.Kind == parser.KindIfStmt {
				p.write(" else ")
				p.printIfStmtInner(elseBody, false)
			} else {
				p.write(" else ")
				p.printBranchBody(elseBody)
			}
		} else {
			p.write("\n")
			p.atLineStart = true
		}
	} else {
		p.write("\n")
		p.atLineStart = true
		p.indent++
		p.printStatement(body)
		p.indent--

		if len(ifChildren) > 2 {
			p.writeIndent()
			p.write("else ")
			p.printBranchBody(ifChildren[2])
		}
	}
}

func (p *JavaPrettyPrinter) printBranchBody(node *parser.Node) {
	if node.Kind == parser.KindBlock {
		p.printBlock(node)
	} else {
		p.write("\n")
		p.atLineStart = true
		p.indent++
		p.printStatement(node)
		p.indent--
	}
}

func (p *JavaPrettyPrinter) printForStmt(node *parser.Node) {
	p.writeIndent()
	p.write("for (")

	init := node.FirstChildOfKind(parser.KindForInit)
	if init != nil {
		p.printForInit(init)
	}
	p.write("; ")

	// Print condition (any expression that's not ForInit, ForUpdate, or a statement)
	for _, child := range node.Children {
		if child.Kind != parser.KindForInit && child.Kind != parser.KindForUpdate &&
			!p.isStatementKind(child.Kind) {
			p.printExpr(child)
			break
		}
	}
	p.write("; ")

	update := node.FirstChildOfKind(parser.KindForUpdate)
	if update != nil {
		p.printForUpdate(update)
	}
	p.write(") ")

	// Find and print the body statement (last non-ForInit, non-ForUpdate child that's a statement)
	var body *parser.Node
	for _, child := range node.Children {
		if child.Kind != parser.KindForInit && child.Kind != parser.KindForUpdate &&
			p.isStatementKind(child.Kind) {
			body = child
		}
	}

	if body != nil {
		if body.Kind == parser.KindBlock {
			p.printBlock(body)
		} else if body.Kind == parser.KindEmptyStmt {
			p.write(";\n")
			p.atLineStart = true
		} else {
			// Single statement body
			p.write("\n")
			p.atLineStart = true
			p.indent++
			p.printStatement(body)
			p.indent--
		}
	} else {
		p.write(";\n")
		p.atLineStart = true
	}
}

func (p *JavaPrettyPrinter) printForInit(node *parser.Node) {
	first := true
	for _, child := range node.Children {
		if !first {
			p.write(", ")
		}
		if child.Kind == parser.KindLocalVarDecl {
			p.printForLocalVarDecl(child)
		} else {
			p.printExpr(child)
		}
		first = false
	}
}

func (p *JavaPrettyPrinter) printForLocalVarDecl(node *parser.Node) {
	// Print modifiers (including annotations) inline
	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiersInline(modifiers)
	}

	// Find the type
	var varType *parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			varType = child
			break
		}
	}

	if varType != nil {
		p.printType(varType)
		p.write(" ")
	}

	// Print declarators: pattern is varName [= initializer], varName [= initializer], ...
	// The AST structure can be: [Type] [name1] [init1] [name2] [init2] ...
	// OR: [Type] [name1] [name2] [init2] ... (when some vars have no initializer)
	// Key insight: variable names are always KindIdentifier, initializers can be other expressions
	// Two consecutive identifiers means the first is a var without initializer
	first := true
	sawVarName := false
	for i, child := range node.Children {
		if child.Kind == parser.KindModifiers || child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			continue
		}

		if child.Kind == parser.KindIdentifier && child.Token != nil {
			// This is a variable name
			// If we saw a var name before and didn't see an initializer, the previous var had no initializer
			if sawVarName {
				// Previous variable had no initializer, this is a new variable
				p.write(", ")
			} else if !first {
				p.write(", ")
			}
			first = false
			p.write(child.Token.Literal)
			sawVarName = true

			// Check if next child is also an identifier (meaning this var has no initializer)
			// or if this is the last child (also no initializer)
			hasInitializer := false
			if i+1 < len(node.Children) {
				next := node.Children[i+1]
				if next.Kind != parser.KindIdentifier {
					hasInitializer = true
				}
			}
			if !hasInitializer {
				sawVarName = true // still expecting to see an initializer or next var
			}
		} else {
			// This is an initializer expression
			p.write(" = ")
			p.printExpr(child)
			sawVarName = false
		}
	}
}

func (p *JavaPrettyPrinter) printLocalVarDeclInline(node *parser.Node) {
	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		for _, child := range modifiers.Children {
			if child.Kind == parser.KindIdentifier && child.Token != nil {
				p.write(child.Token.Literal)
				p.write(" ")
			}
		}
	}

	// Find the type
	var varType *parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			varType = child
			break
		}
	}

	if varType != nil {
		p.printType(varType)
		p.write(" ")
	}

	// Print declarators using same logic as printLocalVarDecl
	first := true
	i := 0
	var prevChild *parser.Node
	for i < len(node.Children) {
		child := node.Children[i]
		if child.Kind == parser.KindModifiers || child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			i++
			continue
		}

		// Check if this child is an initializer for the previous variable
		isInitializer := false
		if prevChild != nil {
			isInitializer = p.hasAssignBetween(prevChild, child)
		}

		if isInitializer {
			p.write(" = ")
			p.printExpr(child)
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
		}
		prevChild = child
		i++
	}
}

func (p *JavaPrettyPrinter) printForUpdate(node *parser.Node) {
	first := true
	for _, child := range node.Children {
		if !first {
			p.write(", ")
		}
		p.printExpr(child)
		first = false
	}
}

func (p *JavaPrettyPrinter) printEnhancedForStmt(node *parser.Node) {
	p.writeIndent()
	p.write("for (")

	var modifiers *parser.Node
	var paramType *parser.Node
	var name string
	var iterable *parser.Node
	var body *parser.Node
	foundName := false

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindModifiers:
			modifiers = child
		case parser.KindType, parser.KindArrayType:
			paramType = child
		case parser.KindIdentifier:
			if child.Token != nil {
				if !foundName {
					name = child.Token.Literal
					foundName = true
				} else {
					iterable = child
				}
			}
		default:
			// The body can be a Block or any other statement type (e.g., single-statement body)
			if paramType != nil && foundName {
				if iterable == nil {
					iterable = child
				} else {
					body = child
				}
			}
		}
	}

	// Print modifiers (e.g., "final") before type
	if modifiers != nil {
		p.printModifiersInline(modifiers)
	}
	if paramType != nil {
		p.printType(paramType)
		p.write(" ")
	}
	p.write(name)
	p.write(" : ")
	if iterable != nil {
		p.printExpr(iterable)
	}
	p.write(") ")

	if body != nil {
		p.printBranchBody(body)
	}
}

func (p *JavaPrettyPrinter) printWhileStmt(node *parser.Node) {
	p.writeIndent()
	p.write("while (")

	children := node.Children
	if len(children) > 0 {
		p.printExpr(children[0])
	}
	p.write(") ")

	if len(children) > 1 {
		p.printBranchBody(children[1])
	}
}

func (p *JavaPrettyPrinter) printDoStmt(node *parser.Node) {
	p.writeIndent()
	p.write("do ")

	var body *parser.Node
	var condition *parser.Node

	for i, child := range node.Children {
		if child.Kind == parser.KindBlock {
			body = child
		} else if i == len(node.Children)-1 {
			condition = child
		}
	}

	if body != nil {
		p.printBlock(body)
	}

	p.writeIndent()
	p.write("while (")
	if condition != nil {
		p.printExpr(condition)
	}
	p.write(");\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printSwitchStmt(node *parser.Node) {
	p.writeIndent()
	p.write("switch (")

	children := node.Children
	if len(children) > 0 {
		p.printExpr(children[0])
	}
	p.write(") {\n")
	p.atLineStart = true

	for _, child := range children[1:] {
		if child.Kind == parser.KindSwitchCase {
			p.printSwitchCase(child)
		}
	}

	p.writeIndent()
	p.write("}\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printSwitchCase(node *parser.Node) {
	labels := node.ChildrenOfKind(parser.KindSwitchLabel)
	for _, label := range labels {
		p.writeIndent()
		p.printSwitchLabel(label)
	}

	p.indent++
	for _, child := range node.Children {
		if child.Kind != parser.KindSwitchLabel {
			p.printStatement(child)
		}
	}
	p.indent--
}

func (p *JavaPrettyPrinter) printSwitchLabel(node *parser.Node) {
	var caseExprs []*parser.Node
	var guard *parser.Node
	var hasArrow bool
	var isDefault bool

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindGuard:
			guard = child
		case parser.KindIdentifier:
			if child.Token != nil {
				switch child.Token.Kind {
				case parser.TokenArrow:
					hasArrow = true
				case parser.TokenDefault:
					isDefault = true
				default:
					caseExprs = append(caseExprs, child)
				}
			} else {
				caseExprs = append(caseExprs, child)
			}
		default:
			caseExprs = append(caseExprs, child)
		}
	}

	if len(caseExprs) == 0 {
		p.write("default")
	} else if isDefault {
		p.write("case ")
		for i, expr := range caseExprs {
			if i > 0 {
				p.write(", ")
			}
			p.printCaseExpr(expr)
		}
		p.write(", default")
	} else {
		p.write("case ")
		for i, expr := range caseExprs {
			if i > 0 {
				p.write(", ")
			}
			p.printCaseExpr(expr)
		}
	}

	if guard != nil {
		p.write(" when ")
		if len(guard.Children) > 0 {
			p.printExpr(guard.Children[0])
		}
	}

	if hasArrow {
		p.write(" ->")
	} else {
		p.write(":")
	}
	p.write("\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printCaseExpr(node *parser.Node) {
	switch node.Kind {
	case parser.KindTypePattern:
		p.printTypePattern(node)
	case parser.KindRecordPattern:
		p.printRecordPattern(node)
	default:
		p.printExpr(node)
	}
}

func (p *JavaPrettyPrinter) printTypePattern(node *parser.Node) {
	typeNode := node.FirstChildOfKind(parser.KindType)
	if typeNode == nil {
		typeNode = node.FirstChildOfKind(parser.KindArrayType)
	}
	if typeNode != nil {
		p.printType(typeNode)
	}
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(" ")
			p.write(child.Token.Literal)
		}
	}
}

func (p *JavaPrettyPrinter) printRecordPattern(node *parser.Node) {
	typeNode := node.FirstChildOfKind(parser.KindType)
	if typeNode != nil {
		p.printType(typeNode)
	}
	p.write("(")
	first := true
	for _, child := range node.Children {
		if child.Kind == parser.KindTypePattern || child.Kind == parser.KindRecordPattern {
			if !first {
				p.write(", ")
			}
			p.printCaseExpr(child)
			first = false
		}
	}
	p.write(")")
}

func (p *JavaPrettyPrinter) printTryStmt(node *parser.Node) {
	p.writeIndent()
	p.write("try ")

	// Collect try-with-resources: can be LocalVarDecl or just Identifier/FieldAccess
	// Java 9+ allows effectively final variables as resources
	var resources []*parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindLocalVarDecl ||
			child.Kind == parser.KindIdentifier ||
			child.Kind == parser.KindFieldAccess {
			resources = append(resources, child)
		}
	}

	if len(resources) > 0 {
		p.write("(")
		for i, res := range resources {
			if i > 0 {
				p.write("; ")
			}
			if res.Kind == parser.KindLocalVarDecl {
				p.printLocalVarDeclInline(res)
			} else if res.Kind == parser.KindIdentifier && res.Token != nil {
				p.write(res.Token.Literal)
			} else {
				p.printExpr(res)
			}
		}
		p.write(") ")
	}

	block := node.FirstChildOfKind(parser.KindBlock)
	if block != nil {
		p.printBlock(block)
	}

	for _, child := range node.Children {
		if child.Kind == parser.KindCatchClause {
			p.printCatchClause(child)
		} else if child.Kind == parser.KindFinallyClause {
			p.printFinallyClause(child)
		}
	}
}

func (p *JavaPrettyPrinter) printCatchClause(node *parser.Node) {
	p.writeIndent()
	p.write("catch (")

	param := node.FirstChildOfKind(parser.KindParameter)
	if param != nil {
		p.printParameter(param)
	} else {
		// Catch parameter structure from parser:
		// - KindModifiers (optional, e.g., "final")
		// - KindType (wrapper containing one or more types for multi-catch)
		//   - KindType (first exception type)
		//   - KindType (second exception type, if multi-catch)
		//   - ...
		// - KindIdentifier (variable name) or KindUnnamedVariable (_)

		// Print modifiers first (e.g., "final")
		modifiers := node.FirstChildOfKind(parser.KindModifiers)
		if modifiers != nil {
			p.printModifiersInline(modifiers)
		}

		var name string
		var isUnnamed bool
		typeWrapper := node.FirstChildOfKind(parser.KindType)
		for _, child := range node.Children {
			if child.Kind == parser.KindIdentifier && child.Token != nil {
				name = child.Token.Literal
			} else if child.Kind == parser.KindUnnamedVariable {
				isUnnamed = true
			}
		}

		if typeWrapper != nil {
			// Check for nested Type children (multi-catch case)
			nestedTypes := typeWrapper.ChildrenOfKind(parser.KindType)
			if len(nestedTypes) > 0 {
				// Multi-catch: print each type separated by |
				for i, t := range nestedTypes {
					if i > 0 {
						p.write(" | ")
					}
					p.printType(t)
				}
			} else {
				// Single type catch
				p.printType(typeWrapper)
			}
		}

		if name != "" {
			p.write(" ")
			p.write(name)
		} else if isUnnamed {
			p.write(" _")
		}
	}

	p.write(") ")

	block := node.FirstChildOfKind(parser.KindBlock)
	if block != nil {
		p.printBlock(block)
	}
}

func (p *JavaPrettyPrinter) printFinallyClause(node *parser.Node) {
	p.writeIndent()
	p.write("finally ")

	block := node.FirstChildOfKind(parser.KindBlock)
	if block != nil {
		p.printBlock(block)
	}
}

func (p *JavaPrettyPrinter) printThrowStmt(node *parser.Node) {
	p.writeIndent()
	p.write("throw ")
	for _, child := range node.Children {
		p.printExpr(child)
	}
	p.write(";\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printBreakStmt(node *parser.Node) {
	p.writeIndent()
	p.write("break")
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(" ")
			p.write(child.Token.Literal)
		}
	}
	p.write(";\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printContinueStmt(node *parser.Node) {
	p.writeIndent()
	p.write("continue")
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(" ")
			p.write(child.Token.Literal)
		}
	}
	p.write(";\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printSynchronizedStmt(node *parser.Node) {
	p.writeIndent()
	p.write("synchronized (")

	children := node.Children
	if len(children) > 0 {
		p.printExpr(children[0])
	}
	p.write(") ")

	block := node.FirstChildOfKind(parser.KindBlock)
	if block != nil {
		p.printBlock(block)
	}
}

func (p *JavaPrettyPrinter) printAssertStmt(node *parser.Node) {
	p.writeIndent()
	p.write("assert ")
	first := true
	for _, child := range node.Children {
		if !first {
			p.write(" : ")
		}
		p.printExpr(child)
		first = false
	}
	p.write(";\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printLabeledStmt(node *parser.Node) {
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.writeIndent()
			p.write(child.Token.Literal)
			p.write(":\n")
			p.atLineStart = true
		} else {
			p.printStatement(child)
		}
	}
}

func (p *JavaPrettyPrinter) printYieldStmt(node *parser.Node) {
	p.writeIndent()
	p.write("yield ")
	for _, child := range node.Children {
		p.printExpr(child)
	}
	p.write(";\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printExplicitConstructorInvocation(node *parser.Node) {
	p.writeIndent()
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			// Qualifier for super call: outer.super(...)
			if child.Token != nil {
				p.write(child.Token.Literal)
				p.write(".")
			}
		case parser.KindThis:
			p.write("this")
		case parser.KindSuper:
			p.write("super")
		case parser.KindParameters:
			p.write("(")
			p.printArguments(child)
			p.write(")")
		}
	}
	p.write(";\n")
	p.atLineStart = true
}

func (p *JavaPrettyPrinter) printExpr(node *parser.Node) {
	switch node.Kind {
	case parser.KindLiteral:
		if node.Token != nil {
			p.write(node.Token.Literal)
		}
	case parser.KindIdentifier:
		if node.Token != nil {
			p.write(node.Token.Literal)
		}
	case parser.KindThis:
		p.write("this")
	case parser.KindSuper:
		p.write("super")
	case parser.KindQualifiedName:
		p.printQualifiedName(node)
	case parser.KindBinaryExpr:
		p.printBinaryExpr(node)
	case parser.KindUnaryExpr:
		p.printUnaryExpr(node)
	case parser.KindPostfixExpr:
		p.printPostfixExpr(node)
	case parser.KindAssignExpr:
		p.printAssignExpr(node)
	case parser.KindTernaryExpr:
		p.printTernaryExpr(node)
	case parser.KindCallExpr:
		p.printCallExpr(node)
	case parser.KindNewExpr:
		p.printNewExpr(node)
	case parser.KindNewArrayExpr:
		p.printNewArrayExpr(node)
	case parser.KindArrayInit:
		p.printArrayInit(node)
	case parser.KindFieldAccess:
		p.printFieldAccess(node)
	case parser.KindArrayAccess:
		p.printArrayAccess(node)
	case parser.KindCastExpr:
		p.printCastExpr(node)
	case parser.KindInstanceofExpr:
		p.printInstanceofExpr(node)
	case parser.KindParenExpr:
		p.printParenExpr(node)
	case parser.KindLambdaExpr:
		p.printLambdaExpr(node)
	case parser.KindMethodRef:
		p.printMethodRef(node)
	case parser.KindClassLiteral:
		p.printClassLiteral(node)
	case parser.KindSwitchExpr:
		p.printSwitchExpr(node)
	default:
		p.printGenericExpr(node)
	}
}

func (p *JavaPrettyPrinter) printBinaryExpr(node *parser.Node) {
	// BinaryExpr always has 3 children: [left, operator(Identifier), right]
	children := node.Children
	if len(children) < 3 {
		return
	}
	p.printExpr(children[0])
	p.write(" ")
	p.write(children[1].TokenLiteral())
	p.write(" ")
	p.printExpr(children[2])
}

func (p *JavaPrettyPrinter) printUnaryExpr(node *parser.Node) {
	// UnaryExpr has 2 children: [operator(Identifier), operand]
	children := node.Children
	if len(children) < 2 {
		return
	}
	p.write(children[0].TokenLiteral())
	p.printExpr(children[1])
}

func (p *JavaPrettyPrinter) printPostfixExpr(node *parser.Node) {
	// PostfixExpr has 2 children: [operand, operator(Identifier)]
	children := node.Children
	if len(children) < 2 {
		return
	}
	p.printExpr(children[0])
	p.write(children[1].TokenLiteral())
}

func (p *JavaPrettyPrinter) printAssignExpr(node *parser.Node) {
	// AssignExpr always has 3 children: [left, operator(Identifier), right]
	children := node.Children
	if len(children) < 3 {
		return
	}
	p.printExpr(children[0])
	p.write(" ")
	p.write(children[1].TokenLiteral())
	p.write(" ")
	p.printExpr(children[2])
}

func (p *JavaPrettyPrinter) printTernaryExpr(node *parser.Node) {
	children := node.Children
	if len(children) >= 3 {
		p.printExpr(children[0])
		p.write(" ? ")
		p.printExpr(children[1])
		p.write(" : ")
		p.printExpr(children[2])
	}
}

func (p *JavaPrettyPrinter) printCallExpr(node *parser.Node) {
	// CallExpr has 2 children: [target, Parameters]
	// target can be Identifier, FieldAccess, or other expression
	children := node.Children
	if len(children) < 2 {
		return
	}

	target := children[0]
	args := children[1]

	p.printExpr(target)
	p.write("(")
	if args != nil {
		p.printArguments(args)
	}
	p.write(")")
}

func (p *JavaPrettyPrinter) printArguments(node *parser.Node) {
	first := true
	for _, child := range node.Children {
		if !first {
			p.write(", ")
		}
		p.printExpr(child)
		first = false
	}
}

func (p *JavaPrettyPrinter) printNewExpr(node *parser.Node) {
	var outer *parser.Node
	var classType *parser.Node
	var typeArgs *parser.Node
	var args *parser.Node
	var body *parser.Node

	var className string

	// Check for qualified class instance creation (outer.new Inner(...))
	// This is detected when the first child is an expression (Identifier, FieldAccess, etc.)
	// followed by another Identifier for the inner class name
	if len(node.Children) >= 2 {
		first := node.Children[0]
		second := node.Children[1]

		// Qualified instance creation: first child is an expression, second is the class identifier
		isQualified := false
		switch first.Kind {
		case parser.KindIdentifier, parser.KindFieldAccess, parser.KindThis, parser.KindCallExpr, parser.KindParenExpr:
			// If the second child is also an Identifier, this is a qualified creation
			if second.Kind == parser.KindIdentifier {
				isQualified = true
			}
		}

		if isQualified {
			outer = first
		}
	}

	for _, child := range node.Children {
		if child == outer {
			continue // Skip outer, we'll print it specially
		}
		switch child.Kind {
		case parser.KindType, parser.KindParameterizedType:
			classType = child
		case parser.KindQualifiedName:
			className = p.qualifiedNameString(child)
		case parser.KindTypeArguments:
			typeArgs = child
		case parser.KindParameters:
			args = child
		case parser.KindBlock:
			body = child
		case parser.KindIdentifier:
			if child.Token != nil {
				className = child.Token.Literal
			}
		}
	}

	// Print outer reference for qualified class instance creation
	if outer != nil {
		p.printExpr(outer)
		p.write(".")
	}

	p.write("new ")

	if className != "" {
		p.write(className)
	}

	if typeArgs != nil {
		p.printTypeArguments(typeArgs)
	}

	if classType != nil {
		p.printType(classType)
	}

	p.write("(")
	if args != nil {
		p.printArguments(args)
	}
	p.write(")")

	if body != nil {
		p.write(" ")
		p.printAnonymousClassBody(body)
	}
}

func (p *JavaPrettyPrinter) printAnonymousClassBody(node *parser.Node) {
	p.write("{\n")
	p.atLineStart = true
	p.indent++

	for _, child := range node.Children {
		p.emitCommentsBeforeLine(child.Span.Start.Line)
		p.printClassBodyMember(child)
	}

	// Emit any comments inside the block before the closing brace
	p.emitCommentsBeforeLine(node.Span.End.Line)

	p.indent--
	p.writeIndent()
	p.write("}")
}

func (p *JavaPrettyPrinter) printNewArrayExpr(node *parser.Node) {
	p.write("new ")

	var elemType *parser.Node
	var elemName *parser.Node
	var dims []*parser.Node
	var init *parser.Node

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType, parser.KindArrayType:
			elemType = child
		case parser.KindQualifiedName:
			elemName = child
		case parser.KindArrayInit:
			init = child
		case parser.KindAnnotation:
			// Skip annotations for now
		default:
			dims = append(dims, child)
		}
	}

	if elemType != nil {
		p.printType(elemType)
	} else if elemName != nil {
		p.printQualifiedName(elemName)
	}

	// For array initializers without explicit dimensions, we need []
	if init != nil && len(dims) == 0 {
		p.write("[]")
	}

	for _, dim := range dims {
		p.write("[")
		p.printExpr(dim)
		p.write("]")
	}

	if init != nil {
		p.printArrayInit(init)
	}
}

func (p *JavaPrettyPrinter) printArrayInit(node *parser.Node) {
	p.write("{")
	first := true
	for _, child := range node.Children {
		if !first {
			p.write(", ")
		}
		p.printExpr(child)
		first = false
	}
	p.write("}")
}

func (p *JavaPrettyPrinter) printFieldAccess(node *parser.Node) {
	first := true
	prevWasTypeArgs := false
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindTypeArguments:
			// Type arguments appear before method name in generic method calls
			// e.g., Collections.<String>emptyList()
			if !first {
				p.write(".")
			}
			p.printTypeArguments(child)
			first = false
			prevWasTypeArgs = true
		case parser.KindIdentifier:
			// Don't add dot if previous was type arguments (already part of same access)
			if !first && !prevWasTypeArgs {
				p.write(".")
			}
			if child.Token != nil {
				p.write(child.Token.Literal)
			}
			first = false
			prevWasTypeArgs = false
		default:
			if !first {
				p.write(".")
			}
			p.printExpr(child)
			first = false
			prevWasTypeArgs = false
		}
	}
}

func (p *JavaPrettyPrinter) printArrayAccess(node *parser.Node) {
	children := node.Children
	if len(children) >= 2 {
		p.printExpr(children[0])
		p.write("[")
		p.printExpr(children[1])
		p.write("]")
	}
}

func (p *JavaPrettyPrinter) printCastExpr(node *parser.Node) {
	p.write("(")
	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			p.printType(child)
			break
		}
	}
	p.write(") ")
	for _, child := range node.Children {
		if child.Kind != parser.KindType && child.Kind != parser.KindArrayType {
			p.printExpr(child)
			break
		}
	}
}

func (p *JavaPrettyPrinter) printInstanceofExpr(node *parser.Node) {
	children := node.Children
	if len(children) >= 2 {
		p.printExpr(children[0])
		p.write(" instanceof ")

		// Parse remaining children: [final?] type [patternVar?]
		idx := 1

		// Check for optional 'final' modifier
		if idx < len(children) && children[idx].Kind == parser.KindIdentifier &&
			children[idx].Token != nil && children[idx].Token.Literal == "final" {
			p.write("final ")
			idx++
		}

		// Print the type
		if idx < len(children) {
			if children[idx].Kind == parser.KindType || children[idx].Kind == parser.KindArrayType {
				p.printType(children[idx])
			} else {
				p.printExpr(children[idx])
			}
			idx++
		}

		// Check for optional pattern variable (Java 16+ pattern matching)
		if idx < len(children) && children[idx].Kind == parser.KindIdentifier {
			p.write(" ")
			p.write(children[idx].Token.Literal)
		}
	}
}

func (p *JavaPrettyPrinter) printParenExpr(node *parser.Node) {
	p.write("(")
	for _, child := range node.Children {
		p.printExpr(child)
	}
	p.write(")")
}

func (p *JavaPrettyPrinter) printLambdaExpr(node *parser.Node) {
	params := node.FirstChildOfKind(parser.KindParameters)
	var paramIdentifier *parser.Node

	if params != nil {
		p.printParameters(params)
	} else {
		// Single parameter without parentheses: x -> body
		for _, child := range node.Children {
			if child.Kind == parser.KindIdentifier && child.Token != nil {
				paramIdentifier = child
				p.write(child.Token.Literal)
				break
			}
		}
	}

	p.write(" -> ")

	// Find and print the lambda body. Skip Parameters and the single param identifier.
	for _, child := range node.Children {
		if child.Kind == parser.KindParameters {
			continue
		}
		if child == paramIdentifier {
			continue
		}
		// This is the body
		if child.Kind == parser.KindBlock {
			p.printBlock(child)
		} else {
			p.printExpr(child)
		}
		return
	}
}

func (p *JavaPrettyPrinter) printMethodRef(node *parser.Node) {
	first := true
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if !first {
				p.write("::")
			}
			p.write(child.Token.Literal)
			first = false
		} else if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			p.printType(child)
			first = false
		} else {
			p.printExpr(child)
			first = false
		}
	}
}

func (p *JavaPrettyPrinter) printClassLiteral(node *parser.Node) {
	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			p.printType(child)
		} else if child.Kind == parser.KindQualifiedName {
			p.printQualifiedName(child)
		} else if child.Kind == parser.KindFieldAccess {
			// For qualified class names like java.lang.Enum.class
			p.printFieldAccess(child)
		} else if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
		}
	}
	p.write(".class")
}

func (p *JavaPrettyPrinter) printSwitchExpr(node *parser.Node) {
	p.write("switch (")
	children := node.Children
	if len(children) > 0 {
		p.printExpr(children[0])
	}
	p.write(") {\n")
	p.atLineStart = true
	p.indent++

	for _, child := range children[1:] {
		if child.Kind == parser.KindSwitchCase {
			p.printSwitchCase(child)
		}
	}

	p.indent--
	p.writeIndent()
	p.write("}")
}

func (p *JavaPrettyPrinter) printGenericExpr(node *parser.Node) {
	if node.Token != nil {
		p.write(node.Token.Literal)
	}
	for _, child := range node.Children {
		p.printExpr(child)
	}
}

func (p *JavaPrettyPrinter) printGenericNode(node *parser.Node) {
	if node.Token != nil {
		p.write(node.Token.Literal)
	}
	for _, child := range node.Children {
		p.printNode(child)
	}
}

func (p *JavaPrettyPrinter) emitCommentsBeforeLine(line int) {
	p.emitCommentsBeforeLineSkippingAnnotationLines(line, nil)
}

func (p *JavaPrettyPrinter) emitCommentsBeforeLineSkippingAnnotationLines(line int, skipLineCommentLines map[int]bool) {
	for p.commentIndex < len(p.comments) {
		comment := p.comments[p.commentIndex]
		if comment.Span.Start.Line >= line {
			break
		}
		// Skip line comments on annotation lines - they'll be emitted as trailing comments
		if comment.Kind == parser.TokenLineComment && skipLineCommentLines[comment.Span.Start.Line] {
			break
		}
		if comment.Span.Start.Line > p.lastLine+1 {
			p.write("\n")
		}
		p.writeIndent()
		p.write(comment.Literal)
		p.write("\n")
		p.atLineStart = true
		p.lastLine = comment.Span.End.Line
		p.commentIndex++
	}
}

func (p *JavaPrettyPrinter) emitRemainingComments() {
	for p.commentIndex < len(p.comments) {
		comment := p.comments[p.commentIndex]
		if comment.Span.Start.Line > p.lastLine+1 {
			p.write("\n")
		}
		p.writeIndent()
		p.write(comment.Literal)
		p.write("\n")
		p.atLineStart = true
		p.lastLine = comment.Span.End.Line
		p.commentIndex++
	}
}

// emitTrailingLineComment emits a line comment on the given line if one exists.
// This is used to emit trailing comments that appear at the end of a line after code.
func (p *JavaPrettyPrinter) emitTrailingLineComment(line int) {
	if p.commentIndex >= len(p.comments) {
		return
	}
	comment := p.comments[p.commentIndex]
	// Only emit if it's a line comment on the exact same line
	if comment.Kind == parser.TokenLineComment && comment.Span.Start.Line == line {
		p.write(" ")
		p.write(comment.Literal)
		p.lastLine = comment.Span.End.Line
		p.commentIndex++
	}
}

func (p *JavaPrettyPrinter) writeIndent() {
	if !p.atLineStart {
		return
	}
	for i := 0; i < p.indent; i++ {
		p.write(p.indentStr)
	}
	p.atLineStart = false
}

func (p *JavaPrettyPrinter) write(s string) {
	p.w.Write([]byte(s))
}

// isStatementKind returns true if the node kind represents a statement
func (p *JavaPrettyPrinter) isStatementKind(kind parser.NodeKind) bool {
	switch kind {
	case parser.KindBlock, parser.KindEmptyStmt, parser.KindExprStmt, parser.KindIfStmt,
		parser.KindForStmt, parser.KindEnhancedForStmt, parser.KindWhileStmt, parser.KindDoStmt,
		parser.KindSwitchStmt, parser.KindSwitchExpr, parser.KindReturnStmt, parser.KindBreakStmt,
		parser.KindContinueStmt, parser.KindThrowStmt, parser.KindTryStmt, parser.KindSynchronizedStmt,
		parser.KindAssertStmt, parser.KindYieldStmt, parser.KindLocalVarDecl, parser.KindLocalClassDecl,
		parser.KindLabeledStmt:
		return true
	}
	return false
}

func PrettyPrintJava(source []byte) ([]byte, error) {
	opts := []parser.Option{parser.WithComments()}
	pr := parser.ParseCompilationUnit(bytes.NewReader(source), opts...)
	node := pr.Finish()
	if node == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	pp := NewJavaPrettyPrinter(&buf)
	if err := pp.Print(node, source, pr.Comments()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
