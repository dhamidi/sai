package format

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dhamidi/sai/java/parser"
)

// Helper function to parse Java expression and pretty print it
func formatExpr(t *testing.T, input string) string {
	t.Helper()
	p := parser.ParseExpression(strings.NewReader(input))
	node := p.Finish()
	if node.Kind == parser.KindError {
		t.Fatalf("parse error for input %q", input)
	}

	var buf bytes.Buffer
	printer := NewJavaPrettyPrinter(&buf)
	printer.printExpr(node)
	return buf.String()
}

func TestPrintBinaryExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple addition",
			input:    "a + b",
			expected: "a + b",
		},
		{
			name:     "simple subtraction",
			input:    "x - y",
			expected: "x - y",
		},
		{
			name:     "multiplication",
			input:    "a * b",
			expected: "a * b",
		},
		{
			name:     "division",
			input:    "a / b",
			expected: "a / b",
		},
		{
			name:     "modulo",
			input:    "a % b",
			expected: "a % b",
		},
		{
			name:     "bitwise and",
			input:    "a & b",
			expected: "a & b",
		},
		{
			name:     "bitwise or",
			input:    "a | b",
			expected: "a | b",
		},
		{
			name:     "bitwise xor",
			input:    "a ^ b",
			expected: "a ^ b",
		},
		{
			name:     "left shift",
			input:    "a << b",
			expected: "a << b",
		},
		{
			name:     "right shift",
			input:    "a >> b",
			expected: "a >> b",
		},
		{
			name:     "unsigned right shift",
			input:    "a >>> b",
			expected: "a >>> b",
		},
		{
			name:     "logical and",
			input:    "a && b",
			expected: "a && b",
		},
		{
			name:     "logical or",
			input:    "a || b",
			expected: "a || b",
		},
		{
			name:     "equality",
			input:    "a == b",
			expected: "a == b",
		},
		{
			name:     "inequality",
			input:    "a != b",
			expected: "a != b",
		},
		{
			name:     "less than",
			input:    "a < b",
			expected: "a < b",
		},
		{
			name:     "less than or equal",
			input:    "a <= b",
			expected: "a <= b",
		},
		{
			name:     "greater than",
			input:    "a > b",
			expected: "a > b",
		},
		{
			name:     "greater than or equal",
			input:    "a >= b",
			expected: "a >= b",
		},
		{
			name:     "chained addition",
			input:    "a + b + c",
			expected: "a + b + c",
		},
		{
			name:     "mixed operators",
			input:    "a + b * c",
			expected: "a + b * c",
		},
		{
			name:     "integer literals",
			input:    "1 + 2",
			expected: "1 + 2",
		},
		{
			name:     "string concatenation",
			input:    `"hello" + "world"`,
			expected: `"hello" + "world"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintUnaryExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "negation",
			input:    "-x",
			expected: "-x",
		},
		{
			name:     "positive",
			input:    "+x",
			expected: "+x",
		},
		{
			name:     "logical not",
			input:    "!x",
			expected: "!x",
		},
		{
			name:     "bitwise complement",
			input:    "~x",
			expected: "~x",
		},
		{
			name:     "prefix increment",
			input:    "++x",
			expected: "++x",
		},
		{
			name:     "prefix decrement",
			input:    "--x",
			expected: "--x",
		},
		{
			name:     "negation of literal",
			input:    "-42",
			expected: "-42",
		},
		{
			name:     "double negation",
			input:    "!!x",
			expected: "!!x",
		},
		{
			name:     "negation of expression",
			input:    "-(a + b)",
			expected: "-(a + b)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintPostfixExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "postfix increment",
			input:    "x++",
			expected: "x++",
		},
		{
			name:     "postfix decrement",
			input:    "x--",
			expected: "x--",
		},
		{
			name:     "postfix on array access",
			input:    "arr[i]++",
			expected: "arr[i]++",
		},
		{
			name:     "postfix on field access",
			input:    "obj.count++",
			expected: "obj.count++",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintAssignExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple assignment",
			input:    "x = 5",
			expected: "x = 5",
		},
		{
			name:     "add assign",
			input:    "x += 5",
			expected: "x += 5",
		},
		{
			name:     "subtract assign",
			input:    "x -= 5",
			expected: "x -= 5",
		},
		{
			name:     "multiply assign",
			input:    "x *= 5",
			expected: "x *= 5",
		},
		{
			name:     "divide assign",
			input:    "x /= 5",
			expected: "x /= 5",
		},
		{
			name:     "modulo assign",
			input:    "x %= 5",
			expected: "x %= 5",
		},
		{
			name:     "and assign",
			input:    "x &= 5",
			expected: "x &= 5",
		},
		{
			name:     "or assign",
			input:    "x |= 5",
			expected: "x |= 5",
		},
		{
			name:     "xor assign",
			input:    "x ^= 5",
			expected: "x ^= 5",
		},
		{
			name:     "left shift assign",
			input:    "x <<= 2",
			expected: "x <<= 2",
		},
		{
			name:     "right shift assign",
			input:    "x >>= 2",
			expected: "x >>= 2",
		},
		{
			name:     "unsigned right shift assign",
			input:    "x >>>= 2",
			expected: "x >>>= 2",
		},
		{
			name:     "assign to array element",
			input:    "arr[0] = 5",
			expected: "arr[0] = 5",
		},
		{
			name:     "assign to field",
			input:    "obj.field = 5",
			expected: "obj.field = 5",
		},
		{
			name:     "assign expression",
			input:    "x = y + z",
			expected: "x = y + z",
		},
		{
			name:     "chained assignment",
			input:    "x = y = z",
			expected: "x = y = z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintCallExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple call no args",
			input:    "foo()",
			expected: "foo()",
		},
		{
			name:     "simple call one arg",
			input:    "foo(x)",
			expected: "foo(x)",
		},
		{
			name:     "simple call multiple args",
			input:    "foo(x, y, z)",
			expected: "foo(x, y, z)",
		},
		{
			name:     "method call on object",
			input:    "obj.method()",
			expected: "obj.method()",
		},
		{
			name:     "method call with args",
			input:    "obj.method(a, b)",
			expected: "obj.method(a, b)",
		},
		{
			name:     "chained method calls",
			input:    "obj.foo().bar()",
			expected: "obj.foo().bar()",
		},
		{
			name:     "System.out.println",
			input:    `System.out.println("hello")`,
			expected: `System.out.println("hello")`,
		},
		{
			name:     "nested method calls",
			input:    "outer(inner())",
			expected: "outer(inner())",
		},
		{
			name:     "call with expression arg",
			input:    "foo(a + b)",
			expected: "foo(a + b)",
		},
		{
			name:     "call with lambda arg",
			input:    "list.forEach(() -> 42)",
			expected: "list.forEach(() -> 42)",
		},
		{
			name:     "static method call",
			input:    "Math.max(a, b)",
			expected: "Math.max(a, b)",
		},
		{
			name:     "method call on new expression",
			input:    "new Foo().bar()",
			expected: "new Foo().bar()",
		},
		{
			name:     "method call on string literal",
			input:    `"hello".length()`,
			expected: `"hello".length()`,
		},
		{
			name:     "method call on this",
			input:    "this.method()",
			expected: "this.method()",
		},
		{
			name:     "method call on super",
			input:    "super.method()",
			expected: "super.method()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintTernaryExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple ternary",
			input:    "a ? b : c",
			expected: "a ? b : c",
		},
		{
			name:     "ternary with comparison",
			input:    "x > 0 ? x : -x",
			expected: "x > 0 ? x : -x",
		},
		{
			name:     "nested ternary",
			input:    "a ? b ? c : d : e",
			expected: "a ? b ? c : d : e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintCastExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "cast to int",
			input:    "(int) x",
			expected: "(int) x",
		},
		{
			name:     "cast to String",
			input:    "(String) obj",
			expected: "(String) obj",
		},
		{
			name:     "cast to array type",
			input:    "(int[]) arr",
			expected: "(int[]) arr",
		},
		{
			name:     "cast to generic type",
			input:    "(List<String>) obj",
			expected: "(List<String>) obj",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintInstanceofExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple instanceof",
			input:    "x instanceof String",
			expected: "x instanceof String",
		},
		{
			name:     "instanceof negated",
			input:    "!(x instanceof String)",
			expected: "!(x instanceof String)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintNewExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple new",
			input:    "new Foo()",
			expected: "new Foo()",
		},
		{
			name:     "new with args",
			input:    "new Foo(1, 2)",
			expected: "new Foo(1, 2)",
		},
		{
			name:     "new with generic type",
			input:    "new ArrayList<String>()",
			expected: "new ArrayList<String>()",
		},
		{
			name:     "new with diamond",
			input:    "new ArrayList<>()",
			expected: "new ArrayList<>()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintArrayAccess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple array access",
			input:    "arr[0]",
			expected: "arr[0]",
		},
		{
			name:     "array access with variable",
			input:    "arr[i]",
			expected: "arr[i]",
		},
		{
			name:     "multidimensional array",
			input:    "matrix[i][j]",
			expected: "matrix[i][j]",
		},
		{
			name:     "array access with expression",
			input:    "arr[i + 1]",
			expected: "arr[i + 1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintFieldAccess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple field access",
			input:    "obj.field",
			expected: "obj.field",
		},
		{
			name:     "chained field access",
			input:    "obj.inner.field",
			expected: "obj.inner.field",
		},
		{
			name:     "System.out",
			input:    "System.out",
			expected: "System.out",
		},
		{
			name:     "field on this",
			input:    "this.field",
			expected: "this.field",
		},
		{
			name:     "field on super",
			input:    "super.field",
			expected: "super.field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintLambdaExpr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no param lambda",
			input:    "() -> 42",
			expected: "() -> 42",
		},
		{
			name:     "typed param lambda",
			input:    "(int x) -> x * 2",
			expected: "(int x) -> x * 2",
		},
		{
			name:     "lambda with method call body",
			input:    "() -> foo()",
			expected: "() -> foo()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintMethodRef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "static method reference",
			input:    "String::valueOf",
			expected: "String::valueOf",
		},
		{
			name:     "instance method reference",
			input:    "obj::method",
			expected: "obj::method",
		},
		{
			name:     "constructor reference",
			input:    "Foo::new",
			expected: "Foo::new",
		},
		{
			name:     "this method reference",
			input:    "this::method",
			expected: "this::method",
		},
		{
			name:     "super method reference",
			input:    "super::method",
			expected: "super::method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "integer literal",
			input:    "42",
			expected: "42",
		},
		{
			name:     "long literal",
			input:    "42L",
			expected: "42L",
		},
		{
			name:     "float literal",
			input:    "3.14f",
			expected: "3.14f",
		},
		{
			name:     "double literal",
			input:    "3.14",
			expected: "3.14",
		},
		{
			name:     "string literal",
			input:    `"hello"`,
			expected: `"hello"`,
		},
		{
			name:     "char literal",
			input:    "'a'",
			expected: "'a'",
		},
		{
			name:     "true literal",
			input:    "true",
			expected: "true",
		},
		{
			name:     "false literal",
			input:    "false",
			expected: "false",
		},
		{
			name:     "null literal",
			input:    "null",
			expected: "null",
		},
		{
			name:     "hex literal",
			input:    "0xFF",
			expected: "0xFF",
		},
		{
			name:     "binary literal",
			input:    "0b1010",
			expected: "0b1010",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "complex arithmetic",
			input:    "(a + b) * (c - d)",
			expected: "(a + b) * (c - d)",
		},
		{
			name:     "method call chain with operations",
			input:    "list.size() + 1",
			expected: "list.size() + 1",
		},
		{
			name:     "assignment with method call",
			input:    "x = obj.getValue()",
			expected: "x = obj.getValue()",
		},
		{
			name:     "ternary with method calls",
			input:    "list.isEmpty() ? 0 : list.size()",
			expected: "list.isEmpty() ? 0 : list.size()",
		},
		{
			name:     "nested new and method call",
			input:    "new StringBuilder().append(x).toString()",
			expected: "new StringBuilder().append(x).toString()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// Tests for wildcards in generic types - use valid expression forms like .class literals
func TestPrintWildcardGenerics(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "wildcard unbounded",
			input:    "List<?>.class",
			expected: "List<?>.class",
		},
		{
			name:     "wildcard extends",
			input:    "List<? extends Number>.class",
			expected: "List<? extends Number>.class",
		},
		{
			name:     "wildcard super",
			input:    "List<? super Integer>.class",
			expected: "List<? super Integer>.class",
		},
		{
			name:     "map with wildcards",
			input:    "Map<String, ?>.class",
			expected: "Map<String, ?>.class",
		},
		{
			name:     "class literal with wildcard",
			input:    "Class<?>.class",
			expected: "Class<?>.class",
		},
		{
			name:     "wildcard cast",
			input:    "(List<?>) x",
			expected: "(List<?>) x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// Helper to format annotation by parsing it in a class context
func formatAnnotation(t *testing.T, annotation string) string {
	t.Helper()
	// Wrap annotation in a class to make it valid Java
	input := annotation + "\nclass X {}"
	p := parser.ParseCompilationUnit(strings.NewReader(input))
	node := p.Finish()
	if node == nil || node.Kind == parser.KindError {
		t.Fatalf("parse error for input %q", input)
	}

	// Find the annotation in the AST
	classDecl := node.FirstChildOfKind(parser.KindClassDecl)
	if classDecl == nil {
		t.Fatalf("no class decl found")
	}
	modifiers := classDecl.FirstChildOfKind(parser.KindModifiers)
	if modifiers == nil {
		t.Fatalf("no modifiers found")
	}
	annotNode := modifiers.FirstChildOfKind(parser.KindAnnotation)
	if annotNode == nil {
		t.Fatalf("no annotation found")
	}

	var buf bytes.Buffer
	printer := NewJavaPrettyPrinter(&buf)
	printer.printAnnotation(annotNode)
	return buf.String()
}

func TestPrintAnnotationArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "annotation with single value array",
			input:    `@SuppressWarnings({"unchecked"})`,
			expected: `@SuppressWarnings({"unchecked"})`,
		},
		{
			name:     "annotation with multiple values",
			input:    `@SuppressWarnings({"deprecation", "rawtypes", "unchecked"})`,
			expected: `@SuppressWarnings({"deprecation", "rawtypes", "unchecked"})`,
		},
		{
			name:     "annotation with value key",
			input:    `@SuppressWarnings(value = {"unchecked"})`,
			expected: `@SuppressWarnings(value = {"unchecked"})`, // preserve explicit "value = " as written
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAnnotation(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatAnnotation(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintDiamondOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "diamond operator in new",
			input:    "new ArrayList<>()",
			expected: "new ArrayList<>()",
		},
		{
			name:     "diamond with arguments",
			input:    "new HashMap<>(16)",
			expected: "new HashMap<>(16)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintGenericMethodCall(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "explicit type argument",
			input:    "Collections.<String>emptyList()",
			expected: "Collections.<String>emptyList()",
		},
		{
			name:     "multiple type arguments",
			input:    "Collections.<String, Object>emptyMap()",
			expected: "Collections.<String, Object>emptyMap()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintNewArrayWithInit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "new array with initializer",
			input:    `new String[]{"a", "b"}`,
			expected: `new String[]{"a", "b"}`,
		},
		{
			name:     "new int array with values",
			input:    "new int[]{1, 2, 3}",
			expected: "new int[]{1, 2, 3}",
		},
		{
			name:     "new byte array with initializer",
			input:    "new byte[]{byteCode}",
			expected: "new byte[]{byteCode}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatExpr(t, tt.input)
			if got != tt.expected {
				t.Errorf("formatExpr(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPrintElseIfOnOneLine(t *testing.T) {
	input := `class X {
    void foo() {
        if (a) {
            x();
        } else if (b) {
            y();
        } else {
            z();
        }
    }
}`
	expected := `class X {

    void foo() {
        if (a) {
            x();
        } else if (b) {
            y();
        } else {
            z();
        }
    }
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("else if not on one line:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

func TestPrintTryCatchOnOneLine(t *testing.T) {
	input := `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (Exception e) {
            handleError(e);
        }
    }
}`
	expected := `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (Exception e) {
            handleError(e);
        }
    }
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("try/catch not on one line:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

func TestPrintTryCatchFinallyOnOneLine(t *testing.T) {
	input := `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (IOException e) {
            handleIO(e);
        } catch (Exception e) {
            handleError(e);
        } finally {
            cleanup();
        }
    }
}`
	expected := `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (IOException e) {
            handleIO(e);
        } catch (Exception e) {
            handleError(e);
        } finally {
            cleanup();
        }
    }
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("try/catch/finally not on one line:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

func TestPrintModuleInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple module",
			input:    "module com.example {}",
			expected: "module com.example {\n}\n",
		},
		{
			name:     "open module",
			input:    "open module com.example {}",
			expected: "open module com.example {\n}\n",
		},
		{
			name:  "module with requires",
			input: "module com.example {\n  requires java.base;\n}",
			expected: `module com.example {
    requires java.base;
}
`,
		},
		{
			name:  "module with requires transitive",
			input: "module com.example {\n  requires transitive java.logging;\n}",
			expected: `module com.example {
    requires transitive java.logging;
}
`,
		},
		{
			name:  "module with requires static",
			input: "module com.example {\n  requires static java.compiler;\n}",
			expected: `module com.example {
    requires static java.compiler;
}
`,
		},
		{
			name:  "module with exports",
			input: "module com.example {\n  exports com.example.api;\n}",
			expected: `module com.example {
    exports com.example.api;
}
`,
		},
		{
			name:  "module with exports to",
			input: "module com.example {\n  exports com.example.internal to com.example.test;\n}",
			expected: `module com.example {
    exports com.example.internal to com.example.test;
}
`,
		},
		{
			name:  "module with opens",
			input: "module com.example {\n  opens com.example.internal;\n}",
			expected: `module com.example {
    opens com.example.internal;
}
`,
		},
		{
			name:  "module with opens to multiple",
			input: "module com.example {\n  opens com.example.internal to com.example.test, com.example.other;\n}",
			expected: `module com.example {
    opens com.example.internal to com.example.test, com.example.other;
}
`,
		},
		{
			name:  "module with uses",
			input: "module com.example {\n  uses com.example.spi.Service;\n}",
			expected: `module com.example {
    uses com.example.spi.Service;
}
`,
		},
		{
			name:  "module with provides",
			input: "module com.example {\n  provides com.example.spi.Service with com.example.impl.ServiceImpl;\n}",
			expected: `module com.example {
    provides com.example.spi.Service with com.example.impl.ServiceImpl;
}
`,
		},
		{
			name:  "module with provides multiple impls",
			input: "module com.example {\n  provides com.example.spi.Service with com.example.impl.Impl1, com.example.impl.Impl2;\n}",
			expected: `module com.example {
    provides com.example.spi.Service with com.example.impl.Impl1, com.example.impl.Impl2;
}
`,
		},
		{
			name: "complete module",
			input: `module com.example.app {
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
			expected: `module com.example.app {
    requires java.base;
    requires transitive java.logging;
    requires static java.compiler;
    exports com.example.api;
    exports com.example.internal to com.example.test;
    opens com.example.model;
    opens com.example.internal to com.example.reflection;
    uses com.example.spi.Service;
    provides com.example.spi.Service with com.example.impl.ServiceImpl;
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJavaFile([]byte(tt.input), "module-info.java")
			if err != nil {
				t.Fatalf("PrettyPrintJavaFile error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("module formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}
