package parser

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenKind
	}{
		{"", []TokenKind{TokenEOF}},
		{"class", []TokenKind{TokenClass, TokenEOF}},
		{"public class Main {}", []TokenKind{TokenPublic, TokenClass, TokenIdent, TokenLBrace, TokenRBrace, TokenEOF}},
		{"123", []TokenKind{TokenIntLiteral, TokenEOF}},
		{"3.14", []TokenKind{TokenFloatLiteral, TokenEOF}},
		{"\"hello\"", []TokenKind{TokenStringLiteral, TokenEOF}},
		{"'a'", []TokenKind{TokenCharLiteral, TokenEOF}},
		{"// comment\nclass", []TokenKind{TokenClass, TokenEOF}},
		{"/* block */ class", []TokenKind{TokenClass, TokenEOF}},
		{"+ - * / %", []TokenKind{TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenPercent, TokenEOF}},
		{"== != < <= > >=", []TokenKind{TokenEQ, TokenNE, TokenLT, TokenLE, TokenGT, TokenGE, TokenEOF}},
		{"&& || !", []TokenKind{TokenAnd, TokenOr, TokenNot, TokenEOF}},
		{"<< >> >>>", []TokenKind{TokenShl, TokenShr, TokenUShr, TokenEOF}},
		{"++ --", []TokenKind{TokenIncrement, TokenDecrement, TokenEOF}},
		{"->", []TokenKind{TokenArrow, TokenEOF}},
		{"::", []TokenKind{TokenColonColon, TokenEOF}},
		{"...", []TokenKind{TokenEllipsis, TokenEOF}},
		{"@", []TokenKind{TokenAt, TokenEOF}},
		{`"Hello \{name}"`, []TokenKind{TokenStringTemplate, TokenEOF}},
		{`"Hello world"`, []TokenKind{TokenStringLiteral, TokenEOF}},
		{"\"\"\"Hello \\{name}\"\"\"", []TokenKind{TokenTextBlockTemplate, TokenEOF}},
		{"\"\"\"Hello world\"\"\"", []TokenKind{TokenTextBlock, TokenEOF}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input), "test.java")
			var got []TokenKind
			for {
				tok := lexer.NextToken()
				if tok.Kind != TokenWhitespace && tok.Kind != TokenComment && tok.Kind != TokenLineComment {
					got = append(got, tok.Kind)
				}
				if tok.Kind == TokenEOF {
					break
				}
			}
			if len(got) != len(tt.expected) {
				t.Errorf("got %d tokens, want %d", len(got), len(tt.expected))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("token %d: got %v, want %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input string
		kind  NodeKind
	}{
		{"42", KindLiteral},
		{"x", KindIdentifier},
		{"x + y", KindBinaryExpr},
		{"x * y + z", KindBinaryExpr},
		{"-x", KindUnaryExpr},
		{"!x", KindUnaryExpr},
		{"x++", KindPostfixExpr},
		{"a ? b : c", KindTernaryExpr},
		{"x = 5", KindAssignExpr},
		{"(x)", KindParenExpr},
		{"obj.field", KindFieldAccess},
		{"obj.method()", KindCallExpr},
		{"arr[0]", KindArrayAccess},
		{"new Foo()", KindNewExpr},
		{"new int[10]", KindNewArrayExpr},
		{"x -> x + 1", KindLambdaExpr},
		{"(a, b) -> a + b", KindLambdaExpr},
		{"obj::method", KindMethodRef},
		{"x instanceof Foo", KindInstanceofExpr},
		{"(int) x", KindCastExpr},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := ParseExpression()
			p.Push([]byte(tt.input))
			node := p.Finish()
			if node.Kind != tt.kind {
				t.Errorf("got %v, want %v", node.Kind, tt.kind)
			}
		})
	}
}

func TestParseCompilationUnit(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"empty class",
			"class Foo {}",
		},
		{
			"class with package",
			"package com.example;\nclass Foo {}",
		},
		{
			"class with import",
			"import java.util.List;\nclass Foo {}",
		},
		{
			"class with field",
			"class Foo { int x; }",
		},
		{
			"class with method",
			"class Foo { void bar() {} }",
		},
		{
			"class with constructor",
			"class Foo { Foo() {} }",
		},
		{
			"public class",
			"public class Foo {}",
		},
		{
			"class extends",
			"class Foo extends Bar {}",
		},
		{
			"class implements",
			"class Foo implements Bar, Baz {}",
		},
		{
			"generic class",
			"class Foo<T> {}",
		},
		{
			"interface",
			"interface Foo {}",
		},
		{
			"enum",
			"enum Color { RED, GREEN, BLUE }",
		},
		{
			"record",
			"record Point(int x, int y) {}",
		},
		{
			"annotation",
			"@interface Override {}",
		},
		{
			"method with params",
			"class Foo { void bar(int x, String y) {} }",
		},
		{
			"method with throws",
			"class Foo { void bar() throws Exception {} }",
		},
		{
			"method with return type",
			"class Foo { int bar() { return 0; } }",
		},
		{
			"field with initializer",
			"class Foo { int x = 5; }",
		},
		{
			"static field",
			"class Foo { static int x; }",
		},
		{
			"annotated class",
			"@Deprecated public class Foo {}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ParseCompilationUnit(WithFile("test.java"))
			p.Push([]byte(tt.input))
			node := p.Finish()
			if node.Kind != KindCompilationUnit {
				t.Errorf("got %v, want CompilationUnit", node.Kind)
			}
			if hasError(node) {
				t.Errorf("parse error in: %s", tt.input)
			}
		})
	}
}

func TestParseStringTemplates(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"simple string template",
			`STR."Hello \{name}"`,
		},
		{
			"string template with expression",
			`STR."The sum is \{a + b}"`,
		},
		{
			"text block template",
			`STR."""
			Hello \{name}
			"""`,
		},
		{
			"nested template expression",
			`STR."Value: \{obj.getValue()}"`,
		},
		{
			"template with empty expression",
			`STR."Is \{} null?"`,
		},
		{
			"FMT template processor",
			`FMT."%-10s\{name}"`,
		},
		{
			"template in statement",
			`class Foo { void m() { String s = STR."Hello \{name}"; } }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node *Node
			if strings.HasPrefix(tt.input, "class") {
				p := ParseCompilationUnit()
				p.Push([]byte(tt.input))
				node = p.Finish()
			} else {
				p := ParseExpression()
				p.Push([]byte(tt.input))
				node = p.Finish()
			}
			if hasError(node) {
				t.Errorf("parse error in: %s", tt.input)
				printErrors(t, node, 0)
			}
		})
	}
}

func TestParseStringTemplateNodeKind(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantKind NodeKind
	}{
		{
			"simple string template produces TemplateExpr",
			`STR."Hello \{name}"`,
			KindTemplateExpr,
		},
		{
			"text block template produces TemplateExpr",
			`STR."""
			Hello \{name}
			"""`,
			KindTemplateExpr,
		},
		{
			"plain string is still FieldAccess",
			`STR."Hello"`,
			KindFieldAccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ParseExpression()
			p.Push([]byte(tt.input))
			node := p.Finish()
			if node == nil {
				t.Fatal("got nil node")
			}
			if node.Kind != tt.wantKind {
				t.Errorf("got %v, want %v", node.Kind, tt.wantKind)
			}
		})
	}
}

func TestParseStatements(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"if stmt", "class Foo { void m() { if (true) {} } }"},
		{"if-else stmt", "class Foo { void m() { if (true) {} else {} } }"},
		{"for stmt", "class Foo { void m() { for (int i = 0; i < 10; i++) {} } }"},
		{"enhanced for", "class Foo { void m() { for (var x : list) {} } }"},
		{"while stmt", "class Foo { void m() { while (true) {} } }"},
		{"do-while stmt", "class Foo { void m() { do {} while (true); } }"},
		{"switch stmt", "class Foo { void m() { switch (x) { case 1: break; default: break; } } }"},
		{"switch pattern", "class Foo { void m() { switch (x) { case Integer i: break; case String s when s.isEmpty(): break; default: break; } } }"},
		{"switch match-all", "class Foo { void m() { switch (x) { case Integer i: break; case _: break; } } }"},
		{"switch record pattern", "class Foo { void m() { switch (x) { case Point(int x, int y): break; case Box(Point p1, Point p2): break; } } }"},
		{"switch nested record", "class Foo { void m() { switch (x) { case Box(Point(int x, int y), _): break; } } }"},
		{"try-catch", "class Foo { void m() { try {} catch (Exception e) {} } }"},
		{"try-finally", "class Foo { void m() { try {} finally {} } }"},
		{"try-with-resources", "class Foo { void m() { try (var r = new R()) {} } }"},
		{"return stmt", "class Foo { int m() { return 1; } }"},
		{"throw stmt", "class Foo { void m() { throw new Exception(); } }"},
		{"assert stmt", "class Foo { void m() { assert x > 0; } }"},
		{"synchronized stmt", "class Foo { void m() { synchronized (this) {} } }"},
		{"labeled stmt", "class Foo { void m() { loop: for (;;) {} } }"},
		{"local var", "class Foo { void m() { int x = 5; } }"},
		{"local var infer", "class Foo { void m() { var x = 5; } }"},
		{"local class", "class Foo { void m() { class Inner {} } }"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ParseCompilationUnit()
			p.Push([]byte(tt.input))
			node := p.Finish()
			if hasError(node) {
				t.Errorf("parse error in: %s", tt.input)
			}
		})
	}
}

func TestPositionTracking(t *testing.T) {
	input := "class Foo {\n    int x;\n}"
	p := ParseCompilationUnit(WithFile("test.java"))
	p.Push([]byte(input))
	node := p.Finish()

	if node.Span.Start.Line != 1 {
		t.Errorf("start line: got %d, want 1", node.Span.Start.Line)
	}
	if node.Span.Start.Column != 1 {
		t.Errorf("start column: got %d, want 1", node.Span.Start.Column)
	}
}

func TestCompactCompilationUnit(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"simple main method",
			"void main() { println(\"Hello\"); }",
		},
		{
			"with import",
			"import java.util.List;\nvoid main() {}",
		},
		{
			"with field before method",
			"int x = 5;\nvoid main() {}",
		},
		{
			"with field after method",
			"void main() {}\nint x = 5;",
		},
		{
			"with nested class after method",
			"void main() {}\nclass Helper {}",
		},
		{
			"multiple methods",
			"void main() { helper(); }\nvoid helper() {}",
		},
		{
			"instance main with string array",
			"void main(String[] args) {}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ParseCompilationUnit(WithFile("test.java"))
			p.Push([]byte(tt.input))
			node := p.Finish()
			if node.Kind != KindCompilationUnit {
				t.Errorf("got %v, want CompilationUnit", node.Kind)
			}
			if hasError(node) {
				t.Errorf("parse error in: %s", tt.input)
				printErrors(t, node, 0)
			}
		})
	}
}

func TestComplexJavaFile(t *testing.T) {
	input := `
package com.example;

import java.util.List;
import java.util.ArrayList;

/**
 * A sample class.
 */
@SuppressWarnings("unchecked")
public class Example<T extends Comparable<T>> implements Runnable {
    private static final int MAX = 100;
    private List<T> items = new ArrayList<>();

    public Example() {
        this.items = new ArrayList<>();
    }

    public void add(T item) {
        if (items.size() < MAX) {
            items.add(item);
        }
    }

    @Override
    public void run() {
        for (T item : items) {
            System.out.println(item);
        }
    }

    public static void main(String[] args) {
        var example = new Example<String>();
        example.add("Hello");
        example.run();
    }
}
`
	p := ParseCompilationUnit(WithFile("Example.java"))
	p.Push([]byte(input))
	node := p.Finish()

	if node.Kind != KindCompilationUnit {
		t.Errorf("got %v, want CompilationUnit", node.Kind)
	}
	if hasError(node) {
		t.Error("parse error in complex file")
		printErrors(t, node, 0)
	}
}

func TestModularCompilationUnit(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"simple module",
			"module com.example {}",
		},
		{
			"open module",
			"open module com.example {}",
		},
		{
			"module with import",
			"import java.util.List;\nmodule com.example {}",
		},
		{
			"module with annotation",
			"@Deprecated\nmodule com.example {}",
		},
		{
			"module with requires",
			"module com.example {\n  requires java.base;\n}",
		},
		{
			"module with requires transitive",
			"module com.example {\n  requires transitive java.logging;\n}",
		},
		{
			"module with requires static",
			"module com.example {\n  requires static java.compiler;\n}",
		},
		{
			"module with exports",
			"module com.example {\n  exports com.example.api;\n}",
		},
		{
			"module with exports to",
			"module com.example {\n  exports com.example.internal to com.example.test;\n}",
		},
		{
			"module with opens",
			"module com.example {\n  opens com.example.internal;\n}",
		},
		{
			"module with opens to",
			"module com.example {\n  opens com.example.internal to com.example.test, com.example.other;\n}",
		},
		{
			"module with uses",
			"module com.example {\n  uses com.example.spi.Service;\n}",
		},
		{
			"module with provides",
			"module com.example {\n  provides com.example.spi.Service with com.example.impl.ServiceImpl;\n}",
		},
		{
			"module with provides multiple impls",
			"module com.example {\n  provides com.example.spi.Service with com.example.impl.Impl1, com.example.impl.Impl2;\n}",
		},
		{
			"complete module",
			`module com.example.app {
  requires java.base;
  requires transitive java.logging;
  requires static java.compiler;
  
  exports com.example.api;
  exports com.example.internal to com.example.test;
  
  opens com.example.model;
  opens com.example.internal to com.example.reflection;
  
  uses com.example.spi.Service;
  
  provides com.example.spi.Service with com.example.impl.ServiceImpl;
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ParseCompilationUnit(WithFile("module-info.java"))
			p.Push([]byte(tt.input))
			node := p.Finish()
			if node.Kind != KindCompilationUnit {
				t.Errorf("got %v, want CompilationUnit", node.Kind)
			}
			if hasError(node) {
				t.Errorf("parse error in: %s", tt.input)
				printErrors(t, node, 0)
			}
			// Verify we have a ModuleDecl child
			hasModuleDecl := false
			for _, child := range node.Children {
				if child.Kind == KindModuleDecl {
					hasModuleDecl = true
					break
				}
			}
			if !hasModuleDecl {
				t.Error("expected ModuleDecl child in CompilationUnit")
			}
		})
	}
}

func hasError(node *Node) bool {
	if node == nil {
		return false
	}
	if node.Kind == KindError {
		return true
	}
	for _, child := range node.Children {
		if hasError(child) {
			return true
		}
	}
	return false
}

func printErrors(t *testing.T, node *Node, depth int) {
	if node == nil {
		return
	}
	if node.Kind == KindError && node.Error != nil {
		t.Logf("%s error: %s at line %d", strings.Repeat("  ", depth), node.Error.Message, node.Span.Start.Line)
	}
	for _, child := range node.Children {
		printErrors(t, child, depth+1)
	}
}
