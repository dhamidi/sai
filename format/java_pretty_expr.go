package format

import (
	"github.com/dhamidi/sai/java/parser"
)

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
	case parser.KindUnnamedVariable:
		p.write("_")
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
	if len(children) < 3 {
		return
	}

	condLen := p.measureExpr(children[0])
	trueLen := p.measureExpr(children[1])
	falseLen := p.measureExpr(children[2])
	totalLen := condLen + 3 + trueLen + 3 + falseLen // " ? " and " : "

	shouldWrap := p.column+totalLen > p.maxColumn

	if shouldWrap {
		p.printExpr(children[0])
		p.write("\n")
		p.atLineStart = true
		p.indent++
		p.writeIndent()
		p.write("? ")
		p.printExpr(children[1])
		p.write("\n")
		p.atLineStart = true
		p.writeIndent()
		p.write(": ")
		p.printExpr(children[2])
		p.indent--
	} else {
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
	if len(node.Children) == 0 {
		return
	}

	var totalLen int
	for i, child := range node.Children {
		if i > 0 {
			totalLen += 2 // ", "
		}
		totalLen += p.measureExpr(child)
	}

	shouldWrap := p.column+totalLen > p.maxColumn && len(node.Children) > 1

	if shouldWrap {
		p.write("\n")
		p.atLineStart = true
		p.indent++
		for i, child := range node.Children {
			p.writeIndent()
			p.printExpr(child)
			if i < len(node.Children)-1 {
				p.write(",")
			}
			p.write("\n")
			p.atLineStart = true
		}
		p.indent--
		p.writeIndent()
	} else {
		first := true
		for _, child := range node.Children {
			if !first {
				p.write(", ")
			}
			p.printExpr(child)
			first = false
		}
	}
}

// chainElement represents one element in a method chain
type chainElement struct {
	methodName string       // The method name (empty for base identifier)
	args       *parser.Node // Arguments node (nil for field access or base)
	typeArgs   *parser.Node // Type arguments for generic method calls
	isBase     bool         // True if this is the base (e.g., "obj" in obj.foo().bar())
	baseNode   *parser.Node // The base node if isBase is true
}

// collectMethodChain collects all elements of a method chain from an expression.
// Returns elements in order from base to final call, and the count of actual method calls.
func (p *JavaPrettyPrinter) collectMethodChain(node *parser.Node) ([]chainElement, int) {
	var elements []chainElement
	callCount := 0
	p.collectChainElements(node, &elements, &callCount)
	// Reverse to get base-to-end order
	for i, j := 0, len(elements)-1; i < j; i, j = i+1, j-1 {
		elements[i], elements[j] = elements[j], elements[i]
	}
	return elements, callCount
}

func (p *JavaPrettyPrinter) collectChainElements(node *parser.Node, elements *[]chainElement, callCount *int) {
	switch node.Kind {
	case parser.KindCallExpr:
		// CallExpr has children: [target, Parameters]
		if len(node.Children) < 2 {
			return
		}
		target := node.Children[0]
		args := node.Children[1]

		// If target is a FieldAccess, we need to get the method name from it
		if target.Kind == parser.KindFieldAccess {
			// FieldAccess for method call: [..., methodName] or [..., TypeArgs, methodName]
			var methodName string
			var typeArgs *parser.Node
			children := target.Children
			if len(children) > 0 {
				lastChild := children[len(children)-1]
				if lastChild.Kind == parser.KindIdentifier && lastChild.Token != nil {
					methodName = lastChild.Token.Literal
				}
				// Check for type arguments (generic method call)
				if len(children) >= 2 {
					secondLast := children[len(children)-2]
					if secondLast.Kind == parser.KindTypeArguments {
						typeArgs = secondLast
					}
				}
			}

			*elements = append(*elements, chainElement{
				methodName: methodName,
				args:       args,
				typeArgs:   typeArgs,
			})
			*callCount++

			// Continue collecting from the rest of the field access chain
			// We need to find what precedes the method name
			if len(children) > 1 {
				// Find the receiver (everything before typeArgs and methodName)
				endIdx := len(children) - 1
				if typeArgs != nil {
					endIdx--
				}
				if endIdx > 0 {
					// Build or traverse the receiver
					if endIdx == 1 {
						// Single receiver
						p.collectChainElements(children[0], elements, callCount)
					} else {
						// Multiple parts - create a synthetic receiver traversal
						for i := endIdx - 1; i >= 0; i-- {
							child := children[i]
							if child.Kind == parser.KindCallExpr {
								p.collectChainElements(child, elements, callCount)
								break
							} else if child.Kind == parser.KindIdentifier {
								// This is a field access or base identifier
								if i == 0 {
									// This is the base
									*elements = append(*elements, chainElement{
										isBase:   true,
										baseNode: child,
									})
								}
								// Otherwise it's part of a qualified name, handled below
							} else {
								// Other expression types (This, Super, etc.)
								p.collectChainElements(child, elements, callCount)
								break
							}
						}
						// If we have a chain of identifiers (qualified name), add as base
						if endIdx > 0 && children[0].Kind == parser.KindIdentifier {
							// Check if all are identifiers (qualified name)
							allIdent := true
							for i := 0; i < endIdx; i++ {
								if children[i].Kind != parser.KindIdentifier {
									allIdent = false
									break
								}
							}
							if allIdent {
								// Reconstruct the qualified name
								var name string
								for i := 0; i < endIdx; i++ {
									if children[i].Kind == parser.KindIdentifier && children[i].Token != nil {
										if i > 0 {
											name += "."
										}
										name += children[i].Token.Literal
									}
								}
								*elements = append(*elements, chainElement{
									isBase:     true,
									methodName: name,
								})
							}
						}
					}
				}
			}
		} else if target.Kind == parser.KindIdentifier {
			// Direct call: foo() - this is a base element (no receiver)
			*elements = append(*elements, chainElement{
				isBase:   true,
				baseNode: node, // Store the whole CallExpr as the base
			})
			*callCount++
		} else {
			// Other target (e.g., new Foo(), parenthesized expression)
			*elements = append(*elements, chainElement{
				args: args,
			})
			*callCount++
			p.collectChainElements(target, elements, callCount)
		}

	case parser.KindFieldAccess:
		// Pure field access (not a method call)
		children := node.Children
		if len(children) > 0 {
			lastChild := children[len(children)-1]
			if lastChild.Kind == parser.KindIdentifier && lastChild.Token != nil {
				*elements = append(*elements, chainElement{
					methodName: lastChild.Token.Literal,
				})
			} else if lastChild.Kind == parser.KindThis {
				// Handle qualified this: Outer.this
				*elements = append(*elements, chainElement{
					methodName: "this",
				})
			}
			// Continue with the receiver
			if len(children) > 1 {
				p.collectChainElements(children[0], elements, callCount)
			} else if len(children) == 1 {
				*elements = append(*elements, chainElement{
					isBase:   true,
					baseNode: children[0],
				})
			}
		}

	case parser.KindIdentifier:
		*elements = append(*elements, chainElement{
			isBase:   true,
			baseNode: node,
		})

	case parser.KindThis:
		*elements = append(*elements, chainElement{
			isBase:     true,
			methodName: "this",
		})

	case parser.KindSuper:
		*elements = append(*elements, chainElement{
			isBase:     true,
			methodName: "super",
		})

	case parser.KindNewExpr:
		*elements = append(*elements, chainElement{
			isBase:   true,
			baseNode: node,
		})

	default:
		*elements = append(*elements, chainElement{
			isBase:   true,
			baseNode: node,
		})
	}
}

// isBuilderBeginMethod returns true if the method name indicates the start of a nested builder block
func isBuilderBeginMethod(name string) bool {
	switch name {
	case "object", "array", "begin", "group", "block", "nest", "start":
		return true
	}
	return false
}

// isBuilderEndMethod returns true if the method name indicates the end of a nested builder block
func isBuilderEndMethod(name string) bool {
	switch name {
	case "end", "done", "close", "finish", "complete":
		return true
	}
	return false
}

// isChainStarterMethod returns true if the method name is a "starter" that should stay on the same line as the base
func isChainStarterMethod(name string) bool {
	switch name {
	case "stream", "parallelStream", "string", "builder", "of", "from", "create", "newBuilder", "values":
		return true
	}
	return false
}

// printMethodChain prints a method chain with proper formatting
func (p *JavaPrettyPrinter) printMethodChain(node *parser.Node, baseIndent int) {
	elements, _ := p.collectMethodChain(node)
	if len(elements) == 0 {
		p.printExpr(node)
		return
	}

	// Mark elements that should stay on the same line as the previous element
	// This happens when:
	// 1. The current method is a "starter" method (stream, values, etc.)
	// 2. The previous element was also a starter method (chaining starters)
	stayOnSameLine := make([]bool, len(elements))
	for i := 1; i < len(elements); i++ {
		elem := elements[i]
		if elem.isBase {
			continue
		}
		// Stay on same line if this is a starter method
		if isChainStarterMethod(elem.methodName) {
			stayOnSameLine[i] = true
		}
	}

	// Calculate indentation levels for begin/end methods
	indentLevels := make([]int, len(elements))
	currentLevel := 0
	for i := 1; i < len(elements); i++ {
		elem := elements[i]
		if isBuilderEndMethod(elem.methodName) {
			currentLevel--
			if currentLevel < 0 {
				currentLevel = 0
			}
		}
		indentLevels[i] = currentLevel
		if isBuilderBeginMethod(elem.methodName) {
			currentLevel++
		}
	}

	// Print the chain
	prevHadWrappedArgs := false
	prevWasStarterOnNewLine := false // Track if prev starter was forced to new line
	for i, elem := range elements {
		if elem.isBase {
			if elem.baseNode != nil {
				p.printExpr(elem.baseNode)
			} else if elem.methodName != "" {
				p.write(elem.methodName)
			}
			prevWasStarterOnNewLine = false
		} else if stayOnSameLine[i] && (!prevHadWrappedArgs || prevWasStarterOnNewLine) {
			// Stay on the same line as previous element when:
			// 1. This is a starter and the previous call didn't have wrapped args, OR
			// 2. This is a starter and the previous element was also a starter on a new line
			p.write(".")
			if elem.typeArgs != nil {
				p.printTypeArguments(elem.typeArgs)
			}
			p.write(elem.methodName)
			if elem.args != nil {
				p.write("(")
				p.printArguments(elem.args)
				p.write(")")
			}
			prevHadWrappedArgs = false
			prevWasStarterOnNewLine = false
		} else {
			// Method calls go on new lines
			p.write("\n")
			p.atLineStart = true
			// Track if this is a starter going to a new line
			prevWasStarterOnNewLine = stayOnSameLine[i]
			// Reset since we already forced a new line
			prevHadWrappedArgs = false
			// Base indentation + 1 for chain + additional for nested builders
			chainIndent := baseIndent + 1 + indentLevels[i]
			for j := 0; j < chainIndent; j++ {
				p.write(p.indentStr)
			}
			p.atLineStart = false // We've written content now
			p.write(".")
			if elem.typeArgs != nil {
				p.printTypeArguments(elem.typeArgs)
			}
			p.write(elem.methodName)
			if elem.args != nil {
				// Save and adjust indent for arguments
				savedIndent := p.indent
				p.indent = chainIndent
				columnBefore := p.column
				p.write("(")
				p.printArguments(elem.args)
				p.write(")")
				p.indent = savedIndent
				// If column decreased, we wrapped to new lines
				prevHadWrappedArgs = p.column < columnBefore || p.atLineStart
				if prevHadWrappedArgs {
					prevWasStarterOnNewLine = false
				}
			}
		}
	}
}

// shouldFormatAsChain returns true if the expression should be formatted as a multi-line method chain
func (p *JavaPrettyPrinter) shouldFormatAsChain(node *parser.Node) bool {
	_, callCount := p.collectMethodChain(node)
	return callCount > 2
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
		case parser.KindIdentifier, parser.KindFieldAccess, parser.KindThis, parser.KindCallExpr, parser.KindParenExpr, parser.KindNewExpr:
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

		// The second child is now a pattern node (TypePattern or RecordPattern)
		pattern := children[1]
		switch pattern.Kind {
		case parser.KindTypePattern:
			p.printTypePattern(pattern)
		case parser.KindRecordPattern:
			p.printRecordPattern(pattern)
		default:
			// Fallback for old-style parsing (type + optional identifier)
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
			// Check for optional pattern variable
			if idx < len(children) && children[idx].Kind == parser.KindIdentifier {
				p.write(" ")
				p.write(children[idx].Token.Literal)
			}
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
	var paramsNode *parser.Node

	if params != nil {
		// Check if this is a single inferred-type parameter (no Parameter children, just Identifier)
		// In that case, we print without parentheses: x -> body
		hasParameterChildren := false
		var singleIdentifier *parser.Node
		identifierCount := 0
		for _, child := range params.Children {
			if child.Kind == parser.KindParameter {
				hasParameterChildren = true
				break
			} else if child.Kind == parser.KindIdentifier && child.Token != nil {
				singleIdentifier = child
				identifierCount++
			}
		}

		if !hasParameterChildren && identifierCount == 1 && singleIdentifier != nil {
			// Single inferred-type parameter: print without parentheses
			p.write(singleIdentifier.Token.Literal)
			paramsNode = params
		} else {
			// Regular parameters with types or multiple params: use parentheses
			p.printParameters(params)
		}
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
		if child == paramsNode {
			continue
		}
		// This is the body
		if child.Kind == parser.KindBlock {
			p.printBlockInline(child)
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
