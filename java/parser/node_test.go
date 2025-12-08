package parser

import (
	"testing"
)

func TestNodeKindString(t *testing.T) {
	tests := []struct {
		kind NodeKind
		want string
	}{
		{KindError, "Error"},
		{KindCompilationUnit, "CompilationUnit"},
		{KindPackageDecl, "PackageDecl"},
		{KindImportDecl, "ImportDecl"},
		{KindClassDecl, "ClassDecl"},
		{KindInterfaceDecl, "InterfaceDecl"},
		{KindEnumDecl, "EnumDecl"},
		{KindMethodDecl, "MethodDecl"},
		{KindFieldDecl, "FieldDecl"},
		{KindBlock, "Block"},
		{KindIfStmt, "IfStmt"},
		{KindForStmt, "ForStmt"},
		{KindWhileStmt, "WhileStmt"},
		{KindReturnStmt, "ReturnStmt"},
		{KindBinaryExpr, "BinaryExpr"},
		{KindCallExpr, "CallExpr"},
		{KindLiteral, "Literal"},
		{KindIdentifier, "Identifier"},
		{NodeKind(9999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.want {
				t.Errorf("NodeKind(%d).String() = %q, want %q", tt.kind, got, tt.want)
			}
		})
	}
}

func TestNodeAddChild(t *testing.T) {
	parent := &Node{Kind: KindClassDecl}
	child1 := &Node{Kind: KindMethodDecl}
	child2 := &Node{Kind: KindFieldDecl}

	parent.AddChild(child1)
	parent.AddChild(child2)
	parent.AddChild(nil)

	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}
	if parent.Children[0] != child1 {
		t.Error("First child mismatch")
	}
	if parent.Children[1] != child2 {
		t.Error("Second child mismatch")
	}
}

func TestNodeIsError(t *testing.T) {
	errorNode := &Node{Kind: KindError}
	normalNode := &Node{Kind: KindClassDecl}

	if !errorNode.IsError() {
		t.Error("Expected IsError() to be true for error node")
	}
	if normalNode.IsError() {
		t.Error("Expected IsError() to be false for non-error node")
	}
}

func TestNodeFirstChildOfKind(t *testing.T) {
	method1 := &Node{Kind: KindMethodDecl, Token: &Token{Literal: "method1"}}
	method2 := &Node{Kind: KindMethodDecl, Token: &Token{Literal: "method2"}}
	field := &Node{Kind: KindFieldDecl}

	parent := &Node{
		Kind:     KindClassDecl,
		Children: []*Node{field, method1, method2},
	}

	t.Run("finds first match", func(t *testing.T) {
		got := parent.FirstChildOfKind(KindMethodDecl)
		if got != method1 {
			t.Error("Expected to find first method")
		}
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		got := parent.FirstChildOfKind(KindIfStmt)
		if got != nil {
			t.Error("Expected nil for non-existent kind")
		}
	})
}

func TestNodeChildrenOfKind(t *testing.T) {
	method1 := &Node{Kind: KindMethodDecl}
	method2 := &Node{Kind: KindMethodDecl}
	field := &Node{Kind: KindFieldDecl}

	parent := &Node{
		Kind:     KindClassDecl,
		Children: []*Node{field, method1, method2},
	}

	t.Run("finds all matches", func(t *testing.T) {
		methods := parent.ChildrenOfKind(KindMethodDecl)
		if len(methods) != 2 {
			t.Errorf("Expected 2 methods, got %d", len(methods))
		}
	})

	t.Run("returns empty slice when not found", func(t *testing.T) {
		got := parent.ChildrenOfKind(KindIfStmt)
		if len(got) != 0 {
			t.Errorf("Expected empty slice, got %d elements", len(got))
		}
	})

	t.Run("finds single match", func(t *testing.T) {
		fields := parent.ChildrenOfKind(KindFieldDecl)
		if len(fields) != 1 {
			t.Errorf("Expected 1 field, got %d", len(fields))
		}
	})
}

func TestNodeTokenLiteral(t *testing.T) {
	t.Run("with token", func(t *testing.T) {
		node := &Node{
			Kind:  KindIdentifier,
			Token: &Token{Literal: "myVar"},
		}
		if got := node.TokenLiteral(); got != "myVar" {
			t.Errorf("TokenLiteral() = %q, want %q", got, "myVar")
		}
	})

	t.Run("without token", func(t *testing.T) {
		node := &Node{Kind: KindBlock}
		if got := node.TokenLiteral(); got != "" {
			t.Errorf("TokenLiteral() = %q, want empty string", got)
		}
	})
}

func TestError(t *testing.T) {
	tok := &Token{Kind: TokenIdent, Literal: "foo"}
	err := &Error{
		Message:  "unexpected token",
		Expected: []TokenKind{TokenClass, TokenInterface},
		Got:      tok,
	}

	if err.Message != "unexpected token" {
		t.Errorf("Message = %q, want %q", err.Message, "unexpected token")
	}
	if len(err.Expected) != 2 {
		t.Errorf("Expected 2 expected tokens, got %d", len(err.Expected))
	}
	if err.Got != tok {
		t.Error("Got token mismatch")
	}
}
