package java

import (
	"github.com/dhamidi/javalyzer/java/parser"
)

// TypeAtPoint returns the fully qualified type name of the value at the given position.
// This is used for auto-completion: given a position after "list.", it returns the type of "list".
// The classes parameter is used to resolve star imports (e.g., import com.example.*).
func TypeAtPoint(root *parser.Node, pos parser.Position, classes []*ClassModel) string {
	if root == nil {
		return ""
	}

	pkg := packageFromCompilationUnit(root)
	resolver := newTypeResolver(pkg, importsFromCompilationUnit(root), classes)

	// Find the enclosing class and register its inner classes
	enclosingClass := findEnclosingClass(root, pos)
	if enclosingClass != nil {
		className := getClassName(enclosingClass, pkg)
		collectAndRegisterInnerClasses(enclosingClass, className, resolver)
	}

	node := findNodeAtPosition(root, pos)
	if node == nil {
		return ""
	}

	// We found an identifier - look up its declaration
	if node.Kind == parser.KindIdentifier && node.Token != nil {
		name := node.Token.Literal
		baseType := resolveVariableType(root, name, pos, resolver)

		// Check if we're inside an array access - if so, reduce array depth
		arrayAccessDepth := countArrayAccessDepth(root, pos, name)
		return reduceArrayDepth(baseType, arrayAccessDepth)
	}

	return ""
}

// countArrayAccessDepth counts how many ArrayAccess nodes wrap the identifier at this position.
func countArrayAccessDepth(root *parser.Node, pos parser.Position, name string) int {
	depth := 0
	var search func(node *parser.Node, inArrayAccess bool)
	search = func(node *parser.Node, inArrayAccess bool) {
		if node == nil {
			return
		}

		isArrayAccess := node.Kind == parser.KindArrayAccess

		// Check if this is the identifier we're looking for
		if node.Kind == parser.KindIdentifier && node.Token != nil && node.Token.Literal == name {
			if positionInSpan(pos, node.Span) && inArrayAccess {
				depth++
			}
		}

		for _, child := range node.Children {
			search(child, isArrayAccess || inArrayAccess)
		}
	}
	search(root, false)
	return depth
}

// reduceArrayDepth removes array brackets from the type string.
func reduceArrayDepth(typeName string, depth int) string {
	for i := 0; i < depth; i++ {
		if len(typeName) >= 2 && typeName[len(typeName)-2:] == "[]" {
			typeName = typeName[:len(typeName)-2]
		}
	}
	return typeName
}

// findEnclosingClass finds the class declaration that contains the given position.
func findEnclosingClass(node *parser.Node, pos parser.Position) *parser.Node {
	if node == nil {
		return nil
	}

	// Check children first for more specific matches
	for _, child := range node.Children {
		if result := findEnclosingClass(child, pos); result != nil {
			return result
		}
	}

	// Check if this is a class containing the position
	if (node.Kind == parser.KindClassDecl || node.Kind == parser.KindInterfaceDecl ||
		node.Kind == parser.KindEnumDecl || node.Kind == parser.KindRecordDecl) &&
		positionInSpan(pos, node.Span) {
		return node
	}

	return nil
}

// getClassName returns the fully qualified class name from a class declaration node.
func getClassName(classNode *parser.Node, pkg string) string {
	for _, child := range classNode.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			if pkg != "" {
				return pkg + "." + child.Token.Literal
			}
			return child.Token.Literal
		}
	}
	return ""
}

// resolveVariableType finds the declaration of a variable and returns its type.
func resolveVariableType(root *parser.Node, name string, pos parser.Position, resolver *typeResolver) string {
	// Find the enclosing method/block and search for local variable declarations
	decl := findVariableDeclaration(root, name, pos)
	if decl == nil {
		return ""
	}

	return typeFromDeclaration(decl, resolver)
}

// findVariableDeclaration finds the declaration of a variable visible at the given position.
// Prioritizes local variables > parameters > fields (shadowing rules).
func findVariableDeclaration(root *parser.Node, name string, pos parser.Position) *parser.Node {
	var localVar, param, field *parser.Node

	var search func(node *parser.Node)
	search = func(node *parser.Node) {
		if node == nil {
			return
		}

		// Check if this is a local variable declaration with our name
		if node.Kind == parser.KindLocalVarDecl {
			// Check if this declaration is before our position
			if positionAfter(pos, node.Span.End) {
				// Check if the declared variable name matches
				for _, child := range node.Children {
					if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == name {
						localVar = node
						return
					}
				}
			}
		}

		// Check parameters
		if node.Kind == parser.KindParameter {
			for _, child := range node.Children {
				if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == name {
					param = node
					return
				}
			}
		}

		// Check field declarations
		if node.Kind == parser.KindFieldDecl {
			for _, child := range node.Children {
				if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == name {
					field = node
				}
			}
		}

		// Check enhanced for loop variable
		if node.Kind == parser.KindEnhancedForStmt && positionInSpan(pos, node.Span) {
			for _, child := range node.Children {
				if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == name {
					localVar = node
					return
				}
			}
		}

		// Check catch clause exception variable
		if node.Kind == parser.KindCatchClause && positionInSpan(pos, node.Span) {
			for _, child := range node.Children {
				if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == name {
					localVar = node
					return
				}
			}
		}

		// Check instanceof pattern variable (Java 16+)
		if node.Kind == parser.KindInstanceofExpr {
			// Pattern variable is declared after the type
			children := node.Children
			for i, child := range children {
				if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == name {
					// Make sure it's after a Type (pattern variable, not the tested expression)
					if i > 0 && (children[i-1].Kind == parser.KindType || children[i-1].Kind == parser.KindTypePattern) {
						localVar = node
						return
					}
				}
			}
		}

		// Check TypePattern (switch pattern matching, Java 21+)
		if node.Kind == parser.KindTypePattern {
			for _, child := range node.Children {
				if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == name {
					localVar = node
					return
				}
			}
		}

		// Recurse into children
		for _, child := range node.Children {
			search(child)
		}
	}

	search(root)

	// Return in priority order: local > param > field
	if localVar != nil {
		return localVar
	}
	if param != nil {
		return param
	}
	return field
}

// positionAfter returns true if pos is after (or at) the given position.
func positionAfter(pos parser.Position, after parser.Position) bool {
	if pos.Line > after.Line {
		return true
	}
	if pos.Line == after.Line && pos.Column >= after.Column {
		return true
	}
	return false
}

// typeFromDeclaration extracts the type from a LocalVarDecl, Parameter, EnhancedForStmt, or InstanceofExpr node.
func typeFromDeclaration(decl *parser.Node, resolver *typeResolver) string {
	var typeNode *parser.Node
	var initExpr *parser.Node
	isVarargs := false

	// For instanceof pattern, the type comes before the pattern variable identifier
	if decl.Kind == parser.KindInstanceofExpr {
		for _, child := range decl.Children {
			if child.Kind == parser.KindType || child.Kind == parser.KindTypePattern {
				typeNode = child
			}
		}
		if typeNode != nil {
			return typeModelFromTypeNode(typeNode, resolver).Name
		}
		return ""
	}

	// For TypePattern (switch pattern), extract the type
	if decl.Kind == parser.KindTypePattern {
		for _, child := range decl.Children {
			if child.Kind == parser.KindType {
				return typeModelFromTypeNode(child, resolver).Name
			}
		}
		return ""
	}

	for _, child := range decl.Children {
		switch child.Kind {
		case parser.KindType, parser.KindArrayType, parser.KindParameterizedType:
			typeNode = child
		case parser.KindNewExpr, parser.KindCallExpr:
			initExpr = child
		case parser.KindIdentifier:
			// Check for varargs marker "..."
			if child.Token != nil && child.Token.Literal == "..." {
				isVarargs = true
			}
		}
	}

	if typeNode == nil {
		return ""
	}

	// Check for 'var' - need to infer type from initializer
	if typeNode.Token != nil && typeNode.Token.Literal == "var" {
		return typeFromInitializer(initExpr, resolver)
	}

	// Explicit type declaration
	model := typeModelFromTypeNode(typeNode, resolver)
	result := model.Name
	for i := 0; i < model.ArrayDepth; i++ {
		result += "[]"
	}
	if isVarargs {
		result += "[]"
	}
	return result
}

// typeFromInitializer infers the type from an initializer expression (e.g., new ArrayList<String>()).
func typeFromInitializer(expr *parser.Node, resolver *typeResolver) string {
	if expr == nil {
		return ""
	}

	switch expr.Kind {
	case parser.KindNewExpr:
		// new ClassName(...) - extract ClassName
		for _, child := range expr.Children {
			if child.Kind == parser.KindQualifiedName {
				return resolver.resolve(qualifiedNameToString(child))
			}
			if child.Kind == parser.KindIdentifier && child.Token != nil {
				return resolver.resolve(child.Token.Literal)
			}
		}
	case parser.KindCallExpr:
		// Method call like URI.create(...) - extract the method's return type
		return typeFromMethodCall(expr, resolver)
	}

	return ""
}

// typeFromMethodCall extracts the return type from a method call expression.
// For static calls like URI.create(...), it looks up the method in the class.
func typeFromMethodCall(expr *parser.Node, resolver *typeResolver) string {
	if len(expr.Children) == 0 {
		return ""
	}

	target := expr.Children[0]

	// Handle FieldAccess (e.g., URI.create, or object.method)
	if target.Kind == parser.KindFieldAccess && len(target.Children) >= 2 {
		// Get the class/object and method name
		classOrObject := target.Children[0]
		methodNode := target.Children[len(target.Children)-1]

		if methodNode.Kind != parser.KindIdentifier || methodNode.Token == nil {
			return ""
		}
		methodName := methodNode.Token.Literal

		// Get the class name (for static calls like URI.create)
		var className string
		if classOrObject.Kind == parser.KindIdentifier && classOrObject.Token != nil {
			className = resolver.resolve(classOrObject.Token.Literal)
		} else if classOrObject.Kind == parser.KindQualifiedName {
			className = resolver.resolve(qualifiedNameToString(classOrObject))
		}

		if className == "" {
			return ""
		}

		// Look up the method's return type in the class
		return lookupMethodReturnType(className, methodName, resolver.classes)
	}

	return ""
}

// lookupMethodReturnType finds a method in a class and returns its return type.
func lookupMethodReturnType(className, methodName string, classes []*ClassModel) string {
	for _, cls := range classes {
		if cls.Name == className {
			for _, method := range cls.Methods {
				if method.Name == methodName {
					result := method.ReturnType.Name
					for i := 0; i < method.ReturnType.ArrayDepth; i++ {
						result += "[]"
					}
					return result
				}
			}
		}
	}
	return ""
}

// findNodeAtPosition finds the most specific (deepest) node that contains the given position.
func findNodeAtPosition(node *parser.Node, pos parser.Position) *parser.Node {
	var bestMatch *parser.Node

	for _, child := range node.Children {
		if found := findNodeAtPosition(child, pos); found != nil {
			if bestMatch == nil || hasLargerSpan(found, bestMatch) {
				bestMatch = found
			}
		}
	}

	if bestMatch != nil {
		return bestMatch
	}

	if positionInSpan(pos, node.Span) {
		return node
	}

	return nil
}

func hasLargerSpan(a, b *parser.Node) bool {
	aSize := spanSize(a.Span)
	bSize := spanSize(b.Span)
	if aSize > 0 && bSize == 0 {
		return true
	}
	return false
}

func spanSize(span parser.Span) int {
	if span.Start.Line == span.End.Line {
		return span.End.Column - span.Start.Column
	}
	return (span.End.Line - span.Start.Line) * 1000
}

func positionInSpan(pos parser.Position, span parser.Span) bool {
	if pos.Line < span.Start.Line {
		return false
	}
	if pos.Line == span.Start.Line && pos.Column < span.Start.Column {
		return false
	}
	if pos.Line > span.End.Line {
		return false
	}
	if pos.Line == span.End.Line && pos.Column > span.End.Column {
		return false
	}
	return true
}
