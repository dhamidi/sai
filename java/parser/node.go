package parser

type NodeKind int

const (
	KindError NodeKind = iota

	// Compilation unit level
	KindCompilationUnit
	KindPackageDecl
	KindImportDecl

	// Type declarations
	KindClassDecl
	KindInterfaceDecl
	KindEnumDecl
	KindRecordDecl
	KindAnnotationDecl
	KindModuleDecl
	KindRequiresDirective
	KindExportsDirective
	KindOpensDirective
	KindUsesDirective
	KindProvidesDirective

	// Members
	KindFieldDecl
	KindMethodDecl
	KindConstructorDecl

	// Type and modifiers
	KindModifiers
	KindTypeParameters
	KindTypeParameter
	KindTypeArguments
	KindTypeArgument
	KindType
	KindArrayType
	KindParameterizedType
	KindWildcard
	KindAnnotation
	KindAnnotationElement

	// Method components
	KindParameters
	KindParameter
	KindThrowsList

	// Statements
	KindBlock
	KindEmptyStmt
	KindExprStmt
	KindIfStmt
	KindForStmt
	KindForInit
	KindForUpdate
	KindEnhancedForStmt
	KindWhileStmt
	KindDoStmt
	KindSwitchStmt
	KindSwitchCase
	KindSwitchLabel
	KindTypePattern
	KindGuard
	KindReturnStmt
	KindBreakStmt
	KindContinueStmt
	KindThrowStmt
	KindTryStmt
	KindCatchClause
	KindFinallyClause
	KindSynchronizedStmt
	KindAssertStmt
	KindLabeledStmt
	KindLocalVarDecl
	KindLocalClassDecl
	KindYieldStmt

	// Expressions
	KindAssignExpr
	KindTernaryExpr
	KindBinaryExpr
	KindUnaryExpr
	KindPostfixExpr
	KindCastExpr
	KindInstanceofExpr
	KindCallExpr
	KindMethodRef
	KindFieldAccess
	KindArrayAccess
	KindNewExpr
	KindNewArrayExpr
	KindArrayInit
	KindLambdaExpr
	KindParenExpr
	KindLiteral
	KindIdentifier
	KindQualifiedName
	KindThis
	KindSuper
	KindClassLiteral
	KindSwitchExpr
)

var nodeKindNames = map[NodeKind]string{
	KindError:             "Error",
	KindCompilationUnit:   "CompilationUnit",
	KindPackageDecl:       "PackageDecl",
	KindImportDecl:        "ImportDecl",
	KindClassDecl:         "ClassDecl",
	KindInterfaceDecl:     "InterfaceDecl",
	KindEnumDecl:          "EnumDecl",
	KindRecordDecl:        "RecordDecl",
	KindAnnotationDecl:    "AnnotationDecl",
	KindModuleDecl:        "ModuleDecl",
	KindRequiresDirective: "RequiresDirective",
	KindExportsDirective:  "ExportsDirective",
	KindOpensDirective:    "OpensDirective",
	KindUsesDirective:     "UsesDirective",
	KindProvidesDirective: "ProvidesDirective",
	KindFieldDecl:         "FieldDecl",
	KindMethodDecl:        "MethodDecl",
	KindConstructorDecl:   "ConstructorDecl",
	KindModifiers:         "Modifiers",
	KindTypeParameters:    "TypeParameters",
	KindTypeParameter:     "TypeParameter",
	KindTypeArguments:     "TypeArguments",
	KindTypeArgument:      "TypeArgument",
	KindType:              "Type",
	KindArrayType:         "ArrayType",
	KindParameterizedType: "ParameterizedType",
	KindWildcard:          "Wildcard",
	KindAnnotation:        "Annotation",
	KindAnnotationElement: "AnnotationElement",
	KindParameters:        "Parameters",
	KindParameter:         "Parameter",
	KindThrowsList:        "ThrowsList",
	KindBlock:             "Block",
	KindEmptyStmt:         "EmptyStmt",
	KindExprStmt:          "ExprStmt",
	KindIfStmt:            "IfStmt",
	KindForStmt:           "ForStmt",
	KindForInit:           "ForInit",
	KindForUpdate:         "ForUpdate",
	KindEnhancedForStmt:   "EnhancedForStmt",
	KindWhileStmt:         "WhileStmt",
	KindDoStmt:            "DoStmt",
	KindSwitchStmt:        "SwitchStmt",
	KindSwitchCase:        "SwitchCase",
	KindSwitchLabel:       "SwitchLabel",
	KindTypePattern:       "TypePattern",
	KindGuard:             "Guard",
	KindReturnStmt:        "ReturnStmt",
	KindBreakStmt:         "BreakStmt",
	KindContinueStmt:      "ContinueStmt",
	KindThrowStmt:         "ThrowStmt",
	KindTryStmt:           "TryStmt",
	KindCatchClause:       "CatchClause",
	KindFinallyClause:     "FinallyClause",
	KindSynchronizedStmt:  "SynchronizedStmt",
	KindAssertStmt:        "AssertStmt",
	KindLabeledStmt:       "LabeledStmt",
	KindLocalVarDecl:      "LocalVarDecl",
	KindLocalClassDecl:    "LocalClassDecl",
	KindYieldStmt:         "YieldStmt",
	KindAssignExpr:        "AssignExpr",
	KindTernaryExpr:       "TernaryExpr",
	KindBinaryExpr:        "BinaryExpr",
	KindUnaryExpr:         "UnaryExpr",
	KindPostfixExpr:       "PostfixExpr",
	KindCastExpr:          "CastExpr",
	KindInstanceofExpr:    "InstanceofExpr",
	KindCallExpr:          "CallExpr",
	KindMethodRef:         "MethodRef",
	KindFieldAccess:       "FieldAccess",
	KindArrayAccess:       "ArrayAccess",
	KindNewExpr:           "NewExpr",
	KindNewArrayExpr:      "NewArrayExpr",
	KindArrayInit:         "ArrayInit",
	KindLambdaExpr:        "LambdaExpr",
	KindParenExpr:         "ParenExpr",
	KindLiteral:           "Literal",
	KindIdentifier:        "Identifier",
	KindQualifiedName:     "QualifiedName",
	KindThis:              "This",
	KindSuper:             "Super",
	KindClassLiteral:      "ClassLiteral",
	KindSwitchExpr:        "SwitchExpr",
}

func (k NodeKind) String() string {
	if name, ok := nodeKindNames[k]; ok {
		return name
	}
	return "Unknown"
}

type Error struct {
	Message  string
	Expected []TokenKind
	Got      *Token
}

type Node struct {
	Kind     NodeKind
	Span     Span
	Children []*Node
	Token    *Token
	Error    *Error
}

func (n *Node) AddChild(child *Node) {
	if child != nil {
		n.Children = append(n.Children, child)
	}
}

func (n *Node) IsError() bool {
	return n.Kind == KindError
}

func (n *Node) FirstChildOfKind(kind NodeKind) *Node {
	for _, child := range n.Children {
		if child.Kind == kind {
			return child
		}
	}
	return nil
}

func (n *Node) ChildrenOfKind(kind NodeKind) []*Node {
	var result []*Node
	for _, child := range n.Children {
		if child.Kind == kind {
			result = append(result, child)
		}
	}
	return result
}

func (n *Node) TokenLiteral() string {
	if n.Token != nil {
		return n.Token.Literal
	}
	return ""
}
