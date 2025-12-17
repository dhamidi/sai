package format

import (
	"github.com/dhamidi/sai/java/parser"
)

func (p *JavaPrettyPrinter) printBlock(node *parser.Node) {
	p.write("{\n")
	p.atLineStart = true
	p.indent++
	p.lastLine = node.Span.Start.Line

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

// printBlockInline prints a block without a trailing newline, for use in try/catch/finally chains
func (p *JavaPrettyPrinter) printBlockInline(node *parser.Node) {
	p.write("{\n")
	p.atLineStart = true
	p.indent++
	p.lastLine = node.Span.Start.Line

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
	p.write("}")
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
	prevWasName := false // Track if previous child was a variable name (Identifier or UnnamedVariable)
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
			} else if child.Kind == parser.KindUnnamedVariable {
				p.write("_")
			} else {
				p.printExpr(child)
			}
			prevWasName = (child.Kind == parser.KindIdentifier || child.Kind == parser.KindUnnamedVariable)
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
		if p.shouldFormatAsChain(child) {
			p.printMethodChain(child, p.indent)
		} else {
			p.printExpr(child)
		}
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
			if p.shouldFormatAsChain(child) {
				p.printMethodChain(child, p.indent)
			} else {
				p.printExpr(child)
			}
		} else {
			// This is a new variable name
			if !first {
				p.write(", ")
			}
			first = false
			if child.Kind == parser.KindIdentifier && child.Token != nil {
				p.write(child.Token.Literal)
			} else if child.Kind == parser.KindUnnamedVariable {
				p.write("_")
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

	// The last child is always the condition, everything before is the body
	if len(node.Children) >= 2 {
		body = node.Children[0]
		condition = node.Children[len(node.Children)-1]
	} else if len(node.Children) == 1 {
		condition = node.Children[0]
	}

	if body != nil {
		if body.Kind == parser.KindBlock {
			p.printBlock(body)
		} else {
			// Single statement body (no braces)
			p.printStatement(body)
		}
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

	// Collect non-label children (the case body)
	var bodyChildren []*parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindSwitchLabel {
			continue
		}
		if child.Kind == parser.KindLineComment || child.Kind == parser.KindComment {
			continue
		}
		bodyChildren = append(bodyChildren, child)
	}

	// Check if this is an arrow case with a single expression statement
	hasArrow := false
	for _, label := range labels {
		for _, child := range label.Children {
			if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Kind == parser.TokenArrow {
				hasArrow = true
				break
			}
		}
	}

	// For arrow cases with a single expression, print on same line
	if hasArrow && len(bodyChildren) == 1 && p.isSingleLineArrowBody(bodyChildren[0]) {
		for _, label := range labels {
			p.writeIndent()
			p.printSwitchLabelInline(label)
		}
		p.write(" ")
		p.printArrowBodyInline(bodyChildren[0])
		p.write("\n")
		p.atLineStart = true
		return
	}

	// Default behavior: print label then body on separate lines
	for _, label := range labels {
		p.writeIndent()
		p.printSwitchLabel(label)
	}

	p.indent++
	for _, child := range bodyChildren {
		p.emitCommentsBeforeLine(child.Span.Start.Line)
		p.printStatement(child)
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

// printSwitchLabelInline prints the switch label without trailing newline
func (p *JavaPrettyPrinter) printSwitchLabelInline(node *parser.Node) {
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
}

// isSingleLineArrowBody returns true if the body should be printed on the same line as the arrow
func (p *JavaPrettyPrinter) isSingleLineArrowBody(node *parser.Node) bool {
	switch node.Kind {
	case parser.KindExprStmt:
		// Expression statements are single-line friendly
		return true
	case parser.KindThrowStmt:
		// throw statements are single-line friendly
		return true
	case parser.KindBlock:
		// Empty blocks or blocks with content should go on new lines
		return false
	case parser.KindYieldStmt:
		// yield statements are single-line friendly
		return true
	default:
		return false
	}
}

// printArrowBodyInline prints the arrow case body without indent/newlines
func (p *JavaPrettyPrinter) printArrowBodyInline(node *parser.Node) {
	switch node.Kind {
	case parser.KindExprStmt:
		if len(node.Children) > 0 {
			p.printExpr(node.Children[0])
		}
		p.write(";")
	case parser.KindThrowStmt:
		p.write("throw ")
		if len(node.Children) > 0 {
			p.printExpr(node.Children[0])
		}
		p.write(";")
	case parser.KindYieldStmt:
		p.write("yield ")
		if len(node.Children) > 0 {
			p.printExpr(node.Children[0])
		}
		p.write(";")
	default:
		// Fallback: print as statement (shouldn't reach here)
		p.printStatement(node)
	}
}

func (p *JavaPrettyPrinter) printCaseExpr(node *parser.Node) {
	switch node.Kind {
	case parser.KindTypePattern:
		p.printTypePattern(node)
	case parser.KindRecordPattern:
		p.printRecordPattern(node)
	case parser.KindMatchAllPattern:
		p.write("_")
	default:
		p.printExpr(node)
	}
}

func (p *JavaPrettyPrinter) printTypePattern(node *parser.Node) {
	// Check for optional 'final' modifier (first identifier child with literal "final")
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == "final" {
			p.write("final ")
			break
		}
	}

	typeNode := node.FirstChildOfKind(parser.KindType)
	if typeNode == nil {
		typeNode = node.FirstChildOfKind(parser.KindArrayType)
	}
	if typeNode != nil {
		p.printType(typeNode)
	}

	// Print the pattern variable (identifier that is not 'final')
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal != "final" {
			p.write(" ")
			p.write(child.Token.Literal)
		}
	}
}

func (p *JavaPrettyPrinter) printRecordPattern(node *parser.Node) {
	// Check for optional 'final' modifier
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == "final" {
			p.write("final ")
			break
		}
	}

	typeNode := node.FirstChildOfKind(parser.KindType)
	if typeNode != nil {
		p.printType(typeNode)
	}
	p.write("(")
	first := true
	for _, child := range node.Children {
		if child.Kind == parser.KindTypePattern || child.Kind == parser.KindRecordPattern || child.Kind == parser.KindMatchAllPattern {
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

	// Collect catch and finally clauses
	var catchClauses []*parser.Node
	var finallyClause *parser.Node
	for _, child := range node.Children {
		if child.Kind == parser.KindCatchClause {
			catchClauses = append(catchClauses, child)
		} else if child.Kind == parser.KindFinallyClause {
			finallyClause = child
		}
	}

	block := node.FirstChildOfKind(parser.KindBlock)
	if block != nil {
		// Print try block inline (without trailing newline) so we can append catch/finally
		p.printBlockInline(block)
	}

	// Print catch clauses on same line as closing brace
	for _, catchClause := range catchClauses {
		p.emitCommentsBeforeLine(catchClause.Span.Start.Line)
		p.write(" ")
		p.printCatchClauseInline(catchClause)
	}

	// Print finally clause on same line as closing brace
	if finallyClause != nil {
		p.emitCommentsBeforeLine(finallyClause.Span.Start.Line)
		p.write(" ")
		p.printFinallyClauseInline(finallyClause)
	}

	p.write("\n")
	p.atLineStart = true
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

// printCatchClauseInline prints a catch clause without leading indent and with inline block
func (p *JavaPrettyPrinter) printCatchClauseInline(node *parser.Node) {
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
		p.printBlockInline(block)
	}
}

// printFinallyClauseInline prints a finally clause without leading indent and with inline block
func (p *JavaPrettyPrinter) printFinallyClauseInline(node *parser.Node) {
	p.write("finally ")

	block := node.FirstChildOfKind(parser.KindBlock)
	if block != nil {
		p.printBlockInline(block)
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
		case parser.KindParenExpr, parser.KindFieldAccess:
			// Qualifier can be a parenthesized expression: (expr).super(...)
			p.printExpr(child)
			p.write(".")
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
