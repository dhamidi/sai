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
			p.write("\n")
			p.atLineStart = true
			p.writeIndent()
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
		} else if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
		} else if child.Kind == parser.KindAnnotationElement {
			if !hasValue {
				p.write("(")
			}
			p.printAnnotationElements(node)
			p.write(")")
			hasValue = true
			break
		} else if child.Kind == parser.KindArrayInit || child.Kind == parser.KindLiteral {
			// Single value annotation like @SuppressWarnings({"unchecked"}) or @Value("x")
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

	if name != "" && name != "value" {
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
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
		} else if child.Kind == parser.KindType {
			p.write(" extends ")
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

	var hasTypeArgs bool
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				p.write(child.Token.Literal)
			}
		case parser.KindQualifiedName:
			p.printQualifiedName(child)
		case parser.KindParameterizedType:
			p.printParameterizedType(child)
		case parser.KindArrayType:
			p.printArrayType(child)
		case parser.KindType:
			p.printType(child)
		case parser.KindTypeArguments:
			p.printTypeArguments(child)
			hasTypeArgs = true
		}
	}
	_ = hasTypeArgs
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
		if child.Kind == parser.KindTypeArgument || child.Kind == parser.KindType || child.Kind == parser.KindWildcard {
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
	if node.Kind == parser.KindWildcard {
		p.printWildcard(node)
		return
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			p.printType(child)
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
		if child.Kind == parser.KindType {
			p.printType(child)
		}
	}
}

func (p *JavaPrettyPrinter) printArrayType(node *parser.Node) {
	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType || child.Kind == parser.KindIdentifier || child.Kind == parser.KindQualifiedName {
			p.printType(child)
			break
		}
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
	p.emitCommentsBeforeLine(child.Span.Start.Line)
	p.printNode(child)
}

func (p *JavaPrettyPrinter) printEnumBody(node *parser.Node) {
	block := node.FirstChildOfKind(parser.KindBlock)
	var members []*parser.Node
	var constants []*parser.Node

	var children []*parser.Node
	if block != nil {
		children = block.Children
	} else {
		children = node.Children
	}

	for _, child := range children {
		if child.Kind == parser.KindFieldDecl && p.isEnumConstant(child) {
			constants = append(constants, child)
		} else if child.Kind == parser.KindFieldDecl || child.Kind == parser.KindMethodDecl ||
			child.Kind == parser.KindConstructorDecl || child.Kind == parser.KindClassDecl ||
			child.Kind == parser.KindInterfaceDecl || child.Kind == parser.KindEnumDecl ||
			child.Kind == parser.KindRecordDecl || child.Kind == parser.KindAnnotationDecl {
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
	return node.FirstChildOfKind(parser.KindType) == nil
}

func (p *JavaPrettyPrinter) printEnumConstant(node *parser.Node) {
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			p.write(child.Token.Literal)
		} else if child.Kind == parser.KindParameters {
			p.write("(")
			p.printArguments(child)
			p.write(")")
		}
	}
}

func (p *JavaPrettyPrinter) printFieldDecl(node *parser.Node) {
	p.writeIndent()

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		p.printModifiers(modifiers)
	}

	var fieldType *parser.Node
	var names []string
	var initializer *parser.Node

	for _, child := range node.Children {
		if child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			fieldType = child
		} else if child.Kind == parser.KindIdentifier && child.Token != nil {
			names = append(names, child.Token.Literal)
		} else if child.Kind == parser.KindModifiers {
			// Already handled above
		} else {
			// Any other child is likely an initializer expression
			initializer = child
		}
	}

	if fieldType != nil {
		p.printType(fieldType)
		p.write(" ")
	}

	for i, name := range names {
		if i > 0 {
			p.write(", ")
		}
		p.write(name)
	}

	if initializer != nil {
		p.write(" = ")
		p.printExpr(initializer)
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
				name = child.Token.Literal
			}
		}
	}

	if node.Token != nil && node.Token.Literal == "..." {
		isVarargs = true
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

	// Print declarators: each is an identifier optionally followed by an initializer
	// Structure after type: identifier [initializer], identifier [initializer], ...
	first := true
	i := 0
	for i < len(node.Children) {
		child := node.Children[i]
		if child.Kind == parser.KindModifiers || child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			i++
			continue
		}

		// This should be an identifier (variable name)
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if !first {
				p.write(", ")
			}
			first = false
			p.write(child.Token.Literal)
			i++

			// Check if next child is an initializer (not an identifier and not modifiers/type)
			if i < len(node.Children) {
				next := node.Children[i]
				if next.Kind != parser.KindIdentifier && next.Kind != parser.KindModifiers &&
					next.Kind != parser.KindType && next.Kind != parser.KindArrayType {
					p.write(" = ")
					p.printExpr(next)
					i++
				}
			}
		} else {
			// Unexpected child, skip
			i++
		}
	}

	p.write(";\n")
	p.atLineStart = true
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

	// Print declarators: each is an identifier optionally followed by an initializer
	first := true
	i := 0
	for i < len(node.Children) {
		child := node.Children[i]
		if child.Kind == parser.KindModifiers || child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			i++
			continue
		}

		// This should be an identifier (variable name)
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if !first {
				p.write(", ")
			}
			first = false
			p.write(child.Token.Literal)
			i++

			// Check if next child is an initializer (not an identifier and not modifiers/type)
			if i < len(node.Children) {
				next := node.Children[i]
				if next.Kind != parser.KindIdentifier && next.Kind != parser.KindModifiers &&
					next.Kind != parser.KindType && next.Kind != parser.KindArrayType {
					p.write(" = ")
					p.printExpr(next)
					i++
				}
			}
		} else {
			// Unexpected child, skip
			i++
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

	// Print declarators: each is an identifier optionally followed by an initializer
	first := true
	i := 0
	for i < len(node.Children) {
		child := node.Children[i]
		if child.Kind == parser.KindModifiers || child.Kind == parser.KindType || child.Kind == parser.KindArrayType {
			i++
			continue
		}

		// This should be an identifier (variable name)
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if !first {
				p.write(", ")
			}
			first = false
			p.write(child.Token.Literal)
			i++

			// Check if next child is an initializer (not an identifier and not modifiers/type)
			if i < len(node.Children) {
				next := node.Children[i]
				if next.Kind != parser.KindIdentifier && next.Kind != parser.KindModifiers &&
					next.Kind != parser.KindType && next.Kind != parser.KindArrayType {
					p.write(" = ")
					p.printExpr(next)
					i++
				}
			}
		} else {
			// Unexpected child, skip
			i++
		}
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

	var paramType *parser.Node
	var name string
	var iterable *parser.Node
	var body *parser.Node
	foundName := false

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindModifiers:
			continue
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
		case parser.KindBlock:
			body = child
		default:
			if paramType != nil && foundName && iterable == nil {
				iterable = child
			}
		}
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
		p.printBlock(body)
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

	// Collect try-with-resources declarations
	var resources []*parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindLocalVarDecl {
			resources = append(resources, child)
		}
	}

	if len(resources) > 0 {
		p.write("(")
		for i, res := range resources {
			if i > 0 {
				p.write("; ")
			}
			p.printLocalVarDeclInline(res)
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
		// Catch parameter is inlined as Type and Identifier children
		var types []*parser.Node
		var name string
		for _, child := range node.Children {
			if child.Kind == parser.KindType {
				types = append(types, child)
			} else if child.Kind == parser.KindIdentifier && child.Token != nil {
				name = child.Token.Literal
			}
		}
		for i, t := range types {
			if i > 0 {
				p.write(" | ")
			}
			p.printType(t)
		}
		if name != "" {
			p.write(" ")
			p.write(name)
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
	p.write("new ")

	var classType *parser.Node
	var typeArgs *parser.Node
	var args *parser.Node
	var body *parser.Node

	var className string

	for _, child := range node.Children {
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
		p.printBlock(body)
	}
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
	if params != nil {
		p.printParameters(params)
	} else {
		for _, child := range node.Children {
			if child.Kind == parser.KindIdentifier && child.Token != nil {
				p.write(child.Token.Literal)
				break
			}
		}
	}

	p.write(" -> ")

	for _, child := range node.Children {
		if child.Kind == parser.KindBlock {
			p.printBlock(child)
			return
		}
		if child.Kind != parser.KindParameters && child.Kind != parser.KindIdentifier {
			p.printExpr(child)
			return
		}
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
		} else if child.Kind == parser.KindType {
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
	for p.commentIndex < len(p.comments) {
		comment := p.comments[p.commentIndex]
		if comment.Span.Start.Line >= line {
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
