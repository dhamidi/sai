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

func TestPrintCatchBlockWithComment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "empty catch with comment",
			input: `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (Exception e) {
            // ignored
        }
    }
}`,
			expected: `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (Exception e) {
            // ignored
        }
    }
}
`,
		},
		{
			name: "empty catch with block comment",
			input: `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (Exception e) {
            /* intentionally empty */
        }
    }
}`,
			expected: `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (Exception e) {
            /* intentionally empty */
        }
    }
}
`,
		},
		{
			name: "catch with comment before statement",
			input: `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (Exception e) {
            // log the error
            log(e);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (Exception e) {
            // log the error
            log(e);
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("catch block formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintMethodChainFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "short chain stays on one line",
			input: `class X {
    void foo() {
        obj.method().call();
    }
}`,
			expected: `class X {

    void foo() {
        obj.method().call();
    }
}
`,
		},
		{
			name: "long chain splits to separate lines",
			input: `class X {
    void foo() {
        builder.first().second().third().fourth();
    }
}`,
			expected: `class X {

    void foo() {
        builder
            .first()
            .second()
            .third()
            .fourth();
    }
}
`,
		},
		{
			name: "nanojson style with begin/end indentation",
			input: `class X {
    void foo() {
        String json = JsonWriter.string().object().array("a").value(1).value(2).end().value("b", false).end().done();
    }
}`,
			expected: `class X {

    void foo() {
        String json = JsonWriter.string()
            .object()
                .array("a")
                    .value(1)
                    .value(2)
                .end()
                .value("b", false)
            .end()
            .done();
    }
}
`,
		},
		{
			name: "nested object and array builders",
			input: `class X {
    void foo() {
        String json = JsonWriter.string().object().value("name", "test").array("items").value(1).value(2).value(3).end().object("nested").value("key", "value").end().end().done();
    }
}`,
			expected: `class X {

    void foo() {
        String json = JsonWriter.string()
            .object()
                .value("name", "test")
                .array("items")
                    .value(1)
                    .value(2)
                    .value(3)
                .end()
                .object("nested")
                    .value("key", "value")
                .end()
            .end()
            .done();
    }
}
`,
		},
		{
			name: "stream builder pattern",
			input: `class X {
    void foo() {
        list.stream().filter(x -> x > 0).map(x -> x * 2).collect(Collectors.toList());
    }
}`,
			expected: `class X {

    void foo() {
        list.stream()
            .filter(x -> x > 0)
            .map(x -> x * 2)
            .collect(Collectors.toList());
    }
}
`,
		},
		{
			name: "assignment with long chain",
			input: `class X {
    void foo() {
        Result result = builder.configure().withOption("a").withOption("b").withOption("c").build();
    }
}`,
			expected: `class X {

    void foo() {
        Result result = builder
            .configure()
            .withOption("a")
            .withOption("b")
            .withOption("c")
            .build();
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("method chain formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
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

func TestPrintLongPermitsClause(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "short permits stays on one line",
			input: `sealed interface Shape permits Circle, Square, Triangle {}`,
			expected: `sealed interface Shape permits Circle, Square, Triangle {
}
`,
		},
		{
			name:  "four permits splits to separate line",
			input: `sealed interface Shape permits Circle, Square, Triangle, Rectangle {}`,
			expected: `sealed interface Shape
        permits Circle, Square, Triangle,
                Rectangle {
}
`,
		},
		{
			name:  "six permits in groups of three",
			input: `sealed interface Shape permits Circle, Square, Triangle, Rectangle, Pentagon, Hexagon {}`,
			expected: `sealed interface Shape
        permits Circle, Square, Triangle,
                Rectangle, Pentagon, Hexagon {
}
`,
		},
		{
			name:  "seven permits in groups of three",
			input: `sealed interface Shape permits Circle, Square, Triangle, Rectangle, Pentagon, Hexagon, Heptagon {}`,
			expected: `sealed interface Shape
        permits Circle, Square, Triangle,
                Rectangle, Pentagon, Hexagon,
                Heptagon {
}
`,
		},
		{
			name:  "sealed class with long permits",
			input: `sealed class Animal permits Dog, Cat, Bird, Fish, Snake {}`,
			expected: `sealed class Animal
        permits Dog, Cat, Bird,
                Fish, Snake {
}
`,
		},
		{
			name:  "sealed interface with extends and long permits",
			input: `sealed interface Container extends Iterable permits List, Set, Map, Queue, Deque {}`,
			expected: `sealed interface Container extends Iterable
        permits List, Set, Map,
                Queue, Deque {
}
`,
		},
		{
			name:  "sealed class with implements and long permits",
			input: `sealed class Widget implements Serializable, Cloneable permits Button, Label, TextField, TextArea, ComboBox {}`,
			expected: `sealed class Widget implements Serializable, Cloneable
        permits Button, Label, TextField,
                TextArea, ComboBox {
}
`,
		},
		{
			name:  "permits with qualified type names",
			input: `sealed interface Buffer permits ByteBuffer, CharBuffer, ShortBuffer, IntBuffer, LongBuffer, FloatBuffer, DoubleBuffer {}`,
			expected: `sealed interface Buffer
        permits ByteBuffer, CharBuffer, ShortBuffer,
                IntBuffer, LongBuffer, FloatBuffer,
                DoubleBuffer {
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("permits clause formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintEnumInlineComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "enum with inline comments on each element",
			input: `enum Color {
    RED,    // Primary color
    GREEN,  // Primary color
    BLUE    // Primary color
}`,
			expected: `enum Color {
    RED, // Primary color
    GREEN, // Primary color
    BLUE // Primary color
}
`,
		},
		{
			name: "enum with inline comment on some elements",
			input: `enum Status {
    PENDING,   // Waiting for processing
    ACTIVE,
    COMPLETED  // Successfully finished
}`,
			expected: `enum Status {
    PENDING, // Waiting for processing
    ACTIVE,
    COMPLETED // Successfully finished
}
`,
		},
		{
			name: "enum with inline comments and methods",
			input: `enum Priority {
    LOW,     // Lowest priority
    MEDIUM,  // Default priority
    HIGH;    // Highest priority

    public boolean isUrgent() {
        return this == HIGH;
    }
}`,
			expected: `enum Priority {
    LOW, // Lowest priority
    MEDIUM, // Default priority
    HIGH; // Highest priority

    public boolean isUrgent() {
        return this == HIGH;
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("enum inline comment formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintInstanceofPatternEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Type patterns
		{
			name:     "simple type pattern",
			input:    "x instanceof String s",
			expected: "x instanceof String s",
		},
		{
			name:     "type pattern with final",
			input:    "x instanceof final String s",
			expected: "x instanceof final String s",
		},
		{
			name:     "type pattern with generic type",
			input:    "x instanceof List<String> list",
			expected: "x instanceof List<String> list",
		},
		{
			name:     "type pattern with array type",
			input:    "x instanceof int[] arr",
			expected: "x instanceof int[] arr",
		},
		{
			name:     "type pattern with 2D array",
			input:    "x instanceof String[][] matrix",
			expected: "x instanceof String[][] matrix",
		},
		// Record patterns
		{
			name:     "simple record pattern",
			input:    "x instanceof Point(var a, var b)",
			expected: "x instanceof Point(var a, var b)",
		},
		{
			name:     "record pattern with typed components",
			input:    "x instanceof Point(int a, int b)",
			expected: "x instanceof Point(int a, int b)",
		},
		{
			name:     "record pattern with mixed var and typed",
			input:    "x instanceof Pair(String s, var x)",
			expected: "x instanceof Pair(String s, var x)",
		},
		{
			name:     "nested record pattern",
			input:    "x instanceof Outer(Inner(var a), var b)",
			expected: "x instanceof Outer(Inner(var a), var b)",
		},
		{
			name:     "deeply nested record pattern",
			input:    "x instanceof A(B(C(var x)))",
			expected: "x instanceof A(B(C(var x)))",
		},
		{
			name:     "record pattern with generic type",
			input:    "x instanceof Box<String>(var s)",
			expected: "x instanceof Box<String>(var s)",
		},
		{
			name:     "record pattern single component",
			input:    "x instanceof Wrapper(var value)",
			expected: "x instanceof Wrapper(var value)",
		},
		{
			name:     "record pattern three components",
			input:    "x instanceof Triple(var a, var b, var c)",
			expected: "x instanceof Triple(var a, var b, var c)",
		},
		{
			name:     "record pattern with qualified type",
			input:    "x instanceof java.awt.Point(var x, var y)",
			expected: "x instanceof java.awt.Point(var x, var y)",
		},
		// Edge cases with expressions
		{
			name:     "instanceof in method call",
			input:    "handle(obj instanceof String s)",
			expected: "handle(obj instanceof String s)",
		},
		{
			name:     "instanceof with method call receiver",
			input:    "getObject() instanceof String s",
			expected: "getObject() instanceof String s",
		},
		{
			name:     "instanceof in ternary",
			input:    "obj instanceof String s ? s.length() : 0",
			expected: "obj instanceof String s ? s.length() : 0",
		},
		{
			name:     "instanceof with logical and",
			input:    "obj instanceof String s && s.length() > 0",
			expected: "obj instanceof String s && s.length() > 0",
		},
		{
			name:     "chained instanceof with or",
			input:    "x instanceof String s || x instanceof Integer i",
			expected: "x instanceof String s || x instanceof Integer i",
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

func TestPrintRecordPatternInClassContext(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "empty record pattern",
			input: `class X {
    void foo(Object o) {
        if (o instanceof EmptyRecord()) {
            process();
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof EmptyRecord()) {
            process();
        }
    }
}
`,
		},
		{
			name: "record pattern with underscore wildcard",
			input: `class X {
    void foo(Object o) {
        if (o instanceof Point(var x, _)) {
            useX(x);
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof Point(var x, _)) {
            useX(x);
        }
    }
}
`,
		},
		{
			name: "multiple underscore wildcards",
			input: `class X {
    void foo(Object o) {
        if (o instanceof Triple(_, var y, _)) {
            useY(y);
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof Triple(_, var y, _)) {
            useY(y);
        }
    }
}
`,
		},
		{
			name: "final on record pattern in if",
			input: `class X {
    void foo(Object o) {
        if (o instanceof final Point(var x, var y)) {
            System.out.println(x + y);
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof final Point(var x, var y)) {
            System.out.println(x + y);
        }
    }
}
`,
		},
		{
			name: "nested record with final inner patterns",
			input: `class X {
    void foo(Object o) {
        if (o instanceof Outer(final Inner(var x), var y)) {
            process(x, y);
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof Outer(final Inner(var x), var y)) {
            process(x, y);
        }
    }
}
`,
		},
		{
			name: "chained else-if with different patterns",
			input: `class X {
    void foo(Object o) {
        if (o instanceof String s) {
            handleString(s);
        } else if (o instanceof Point(var x, var y)) {
            handlePoint(x, y);
        } else if (o instanceof Circle(Point(var cx, var cy), var r)) {
            handleCircle(cx, cy, r);
        } else if (o instanceof Integer i) {
            handleInt(i);
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof String s) {
            handleString(s);
        } else if (o instanceof Point(var x, var y)) {
            handlePoint(x, y);
        } else if (o instanceof Circle(Point(var cx, var cy), var r)) {
            handleCircle(cx, cy, r);
        } else if (o instanceof Integer i) {
            handleInt(i);
        }
    }
}
`,
		},
		{
			name: "record pattern in while condition",
			input: `class X {
    void process(Queue<Object> queue) {
        Object item;
        while ((item = queue.poll()) != null && item instanceof Point(var x, var y)) {
            plot(x, y);
        }
    }
}`,
			expected: `class X {

    void process(Queue<Object> queue) {
        Object item;
        while ((item = queue.poll()) != null && item instanceof Point(var x, var y)) {
            plot(x, y);
        }
    }
}
`,
		},
		{
			name: "record pattern with array component type",
			input: `class X {
    void foo(Object o) {
        if (o instanceof ArrayHolder(int[] data)) {
            processArray(data);
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof ArrayHolder(int[] data)) {
            processArray(data);
        }
    }
}
`,
		},
		{
			name: "record pattern with generic component",
			input: `class X {
    void foo(Object o) {
        if (o instanceof Container(List<String> items)) {
            processItems(items);
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof Container(List<String> items)) {
            processItems(items);
        }
    }
}
`,
		},
		{
			name: "deeply nested generic record pattern",
			input: `class X {
    void foo(Object o) {
        if (o instanceof Box<Pair<Integer, String>>(Pair(Integer i, String s))) {
            handle(i, s);
        }
    }
}`,
			expected: `class X {

    void foo(Object o) {
        if (o instanceof Box<Pair<Integer, String>>(Pair(Integer i, String s))) {
            handle(i, s);
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintSwitchWithPatternMatching(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "switch with type patterns",
			input: `class X {
    String describe(Object o) {
        return switch (o) {
            case String s -> "String: " + s;
            case Integer i -> "Integer: " + i;
            case null -> "null";
            default -> "unknown";
        };
    }
}`,
			expected: `class X {

    String describe(Object o) {
        return switch (o) {
            case String s -> "String: " + s;
            case Integer i -> "Integer: " + i;
            case null -> "null";
            default -> "unknown";
        };
    }
}
`,
		},
		{
			name: "switch with record patterns",
			input: `class X {
    void process(Shape s) {
        switch (s) {
            case Circle(Point(var x, var y), var r) -> drawCircle(x, y, r);
            case Rectangle(Point(var x, var y), var w, var h) -> drawRect(x, y, w, h);
            default -> throw new IllegalArgumentException();
        }
    }
}`,
			expected: `class X {

    void process(Shape s) {
        switch (s) {
            case Circle(Point(var x, var y), var r) -> drawCircle(x, y, r);
            case Rectangle(Point(var x, var y), var w, var h) -> drawRect(
                x,
                y,
                w,
                h
            );
            default -> throw new IllegalArgumentException();
        }
    }
}
`,
		},
		{
			name: "switch with guarded patterns",
			input: `class X {
    String classify(Object o) {
        return switch (o) {
            case String s when s.isEmpty() -> "empty string";
            case String s when s.length() > 10 -> "long string";
            case String s -> "string";
            case Integer i when i > 0 -> "positive";
            case Integer i when i < 0 -> "negative";
            case Integer i -> "zero";
            default -> "other";
        };
    }
}`,
			expected: `class X {

    String classify(Object o) {
        return switch (o) {
            case String s when s.isEmpty() -> "empty string";
            case String s when s.length() > 10 -> "long string";
            case String s -> "string";
            case Integer i when i > 0 -> "positive";
            case Integer i when i < 0 -> "negative";
            case Integer i -> "zero";
            default -> "other";
        };
    }
}
`,
		},
		{
			name: "switch with record pattern and guard",
			input: `class X {
    void process(Object o) {
        switch (o) {
            case Point(var x, var y) when x == y -> handleDiagonal(x);
            case Point(var x, var y) when x > y -> handleAboveDiagonal(x, y);
            case Point(var x, var y) -> handleBelowDiagonal(x, y);
            default -> {}
        }
    }
}`,
			expected: `class X {

    void process(Object o) {
        switch (o) {
            case Point(var x, var y) when x == y -> handleDiagonal(x);
            case Point(var x, var y) when x > y -> handleAboveDiagonal(x, y);
            case Point(var x, var y) -> handleBelowDiagonal(x, y);
            default -> {}
        }
    }
}
`,
		},
		{
			name: "switch with underscore pattern",
			input: `class X {
    int count(Object o) {
        return switch (o) {
            case Point(var x, _) -> 1;
            case Circle(_, var r) -> 2;
            case null, default -> 0;
        };
    }
}`,
			expected: `class X {

    int count(Object o) {
        return switch (o) {
            case Point(var x, _) -> 1;
            case Circle(_, var r) -> 2;
            case null, default -> 0;
        };
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintForLoopEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "for loop with multiple initializers",
			input: `class X {
    void foo() {
        for (int i = 0, j = 10; i < j; i++, j--) {
            process(i, j);
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (int i = 0, j = 10; i < j; i++, j--) {
            process(i, j);
        }
    }
}
`,
		},
		{
			name: "for loop with empty initializer",
			input: `class X {
    void foo() {
        int i = 0;
        for (; i < 10; i++) {
            process(i);
        }
    }
}`,
			expected: `class X {

    void foo() {
        int i = 0;
        for (; i < 10; i++) {
            process(i);
        }
    }
}
`,
		},
		{
			name: "for loop with empty condition",
			input: `class X {
    void foo() {
        for (int i = 0; ; i++) {
            if (i > 10) break;
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (int i = 0; ; i++) {
            if (i > 10) break;
        }
    }
}
`,
		},
		{
			name: "infinite for loop",
			input: `class X {
    void foo() {
        for (;;) {
            process();
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (; ; ) {
            process();
        }
    }
}
`,
		},
		{
			name: "for loop with single statement body",
			input: `class X {
    void foo() {
        for (int i = 0; i < 10; i++)
            process(i);
    }
}`,
			expected: "class X {\n\n    void foo() {\n        for (int i = 0; i < 10; i++) \n            process(i);\n    }\n}\n",
		},
		{
			name: "for loop with empty body",
			input: `class X {
    void foo() {
        for (int i = 0; i < 10; i++);
    }
}`,
			expected: "class X {\n\n    void foo() {\n        for (int i = 0; i < 10; i++) ;\n    }\n}\n",
		},
		{
			name: "nested for loops",
			input: `class X {
    void foo() {
        for (int i = 0; i < 10; i++) {
            for (int j = 0; j < 10; j++) {
                matrix[i][j] = i * j;
            }
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (int i = 0; i < 10; i++) {
            for (int j = 0; j < 10; j++) {
                matrix[i][j] = i * j;
            }
        }
    }
}
`,
		},
		{
			name: "for loop initializer from variable",
			input: `class X {
    void foo(int lineEnd) {
        for (int i = lineEnd; i > 0; i--) {
            process(i);
        }
    }
}`,
			expected: `class X {

    void foo(int lineEnd) {
        for (int i = lineEnd; i > 0; i--) {
            process(i);
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintEnhancedForLoopEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "enhanced for with array",
			input: `class X {
    void foo() {
        for (int x : array) {
            process(x);
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (int x : array) {
            process(x);
        }
    }
}
`,
		},
		{
			name: "enhanced for with final modifier",
			input: `class X {
    void foo() {
        for (final String item : items) {
            process(item);
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (final String item : items) {
            process(item);
        }
    }
}
`,
		},
		{
			name: "enhanced for with generic type",
			input: `class X {
    void foo() {
        for (Map.Entry<String, Integer> entry : map.entrySet()) {
            process(entry.getKey(), entry.getValue());
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (Map.Entry<String, Integer> entry : map.entrySet()) {
            process(entry.getKey(), entry.getValue());
        }
    }
}
`,
		},
		{
			name: "enhanced for with var keyword",
			input: `class X {
    void foo() {
        for (var item : items) {
            process(item);
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (var item : items) {
            process(item);
        }
    }
}
`,
		},
		{
			name: "enhanced for with single statement body",
			input: `class X {
    void foo() {
        for (String s : strings)
            System.out.println(s);
    }
}`,
			expected: "class X {\n\n    void foo() {\n        for (String s : strings) \n            System.out.println(s);\n    }\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintWhileDoWhileEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "while with single statement body",
			input: `class X {
    void foo() {
        while (condition)
            process();
    }
}`,
			expected: "class X {\n\n    void foo() {\n        while (condition) \n            process();\n    }\n}\n",
		},
		{
			name: "while with complex condition",
			input: `class X {
    void foo() {
        while ((line = reader.readLine()) != null && !line.isEmpty()) {
            process(line);
        }
    }
}`,
			expected: `class X {

    void foo() {
        while ((line = reader.readLine()) != null && !line.isEmpty()) {
            process(line);
        }
    }
}
`,
		},
		{
			name: "do-while basic",
			input: `class X {
    void foo() {
        do {
            process();
        } while (condition);
    }
}`,
			expected: `class X {

    void foo() {
        do {
            process();
        }
        while (condition);
    }
}
`,
		},
		{
			name: "do-while with complex condition",
			input: `class X {
    void foo() {
        do {
            attempt++;
            result = tryOperation();
        } while (!result.isSuccess() && attempt < maxAttempts);
    }
}`,
			expected: `class X {

    void foo() {
        do {
            attempt++;
            result = tryOperation();
        }
        while (!result.isSuccess() && attempt < maxAttempts);
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintTryWithResourcesEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "try with single resource",
			input: `class X {
    void foo() {
        try (InputStream is = new FileInputStream(file)) {
            process(is);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try (InputStream is = new FileInputStream(file)) {
            process(is);
        }
    }
}
`,
		},
		{
			name: "try with multiple resources",
			input: `class X {
    void foo() {
        try (InputStream is = new FileInputStream(file); OutputStream os = new FileOutputStream(out)) {
            copy(is, os);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try (InputStream is = new FileInputStream(file); OutputStream os = new FileOutputStream(out)) {
            copy(is, os);
        }
    }
}
`,
		},
		{
			name: "try with resources and catch",
			input: `class X {
    void foo() {
        try (BufferedReader reader = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = reader.readLine()) != null) {
                process(line);
            }
        } catch (IOException e) {
            handleError(e);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try (BufferedReader reader = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = reader.readLine()) != null) {
                process(line);
            }
        } catch (IOException e) {
            handleError(e);
        }
    }
}
`,
		},
		{
			name: "try with resources catch and finally",
			input: `class X {
    void foo() {
        try (Connection conn = getConnection()) {
            executeQuery(conn);
        } catch (SQLException e) {
            log(e);
        } finally {
            cleanup();
        }
    }
}`,
			expected: `class X {

    void foo() {
        try (Connection conn = getConnection()) {
            executeQuery(conn);
        } catch (SQLException e) {
            log(e);
        } finally {
            cleanup();
        }
    }
}
`,
		},
		{
			name: "try with var keyword in resource",
			input: `class X {
    void foo() {
        try (var reader = new BufferedReader(new FileReader(file))) {
            process(reader);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try (var reader = new BufferedReader(new FileReader(file))) {
            process(reader);
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintMultiCatchEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "multi-catch two exceptions",
			input: `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (IOException | SQLException e) {
            handleError(e);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (IOException | SQLException e) {
            handleError(e);
        }
    }
}
`,
		},
		{
			name: "multi-catch three exceptions",
			input: `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (IOException | SQLException | ClassNotFoundException e) {
            handleError(e);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (IOException | SQLException | ClassNotFoundException e) {
            handleError(e);
        }
    }
}
`,
		},
		{
			name: "multi-catch with final",
			input: `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (final IOException | SQLException e) {
            handleError(e);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (final IOException | SQLException e) {
            handleError(e);
        }
    }
}
`,
		},
		{
			name: "multiple catch blocks with multi-catch",
			input: `class X {
    void foo() {
        try {
            riskyOperation();
        } catch (IOException | SQLException e) {
            handleDatabaseError(e);
        } catch (RuntimeException e) {
            handleRuntimeError(e);
        }
    }
}`,
			expected: `class X {

    void foo() {
        try {
            riskyOperation();
        } catch (IOException | SQLException e) {
            handleDatabaseError(e);
        } catch (RuntimeException e) {
            handleRuntimeError(e);
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintSynchronizedStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "synchronized on this",
			input: `class X {
    void foo() {
        synchronized (this) {
            counter++;
        }
    }
}`,
			expected: `class X {

    void foo() {
        synchronized (this) {
            counter++;
        }
    }
}
`,
		},
		{
			name: "synchronized on lock object",
			input: `class X {
    void foo() {
        synchronized (lock) {
            sharedData.add(item);
        }
    }
}`,
			expected: `class X {

    void foo() {
        synchronized (lock) {
            sharedData.add(item);
        }
    }
}
`,
		},
		{
			name: "synchronized on class literal",
			input: `class X {
    void foo() {
        synchronized (X.class) {
            staticCounter++;
        }
    }
}`,
			expected: `class X {

    void foo() {
        synchronized (X.class) {
            staticCounter++;
        }
    }
}
`,
		},
		{
			name: "nested synchronized blocks",
			input: `class X {
    void foo() {
        synchronized (lock1) {
            synchronized (lock2) {
                transferData();
            }
        }
    }
}`,
			expected: `class X {

    void foo() {
        synchronized (lock1) {
            synchronized (lock2) {
                transferData();
            }
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintAssertStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple assert",
			input: `class X {
    void foo() {
        assert condition;
    }
}`,
			expected: `class X {

    void foo() {
        assert condition;
    }
}
`,
		},
		{
			name: "assert with message",
			input: `class X {
    void foo() {
        assert condition : "Condition failed";
    }
}`,
			expected: `class X {

    void foo() {
        assert condition : "Condition failed";
    }
}
`,
		},
		{
			name: "assert with complex condition",
			input: `class X {
    void foo() {
        assert x > 0 && x < 100 : "x must be between 0 and 100";
    }
}`,
			expected: `class X {

    void foo() {
        assert x > 0 && x < 100 : "x must be between 0 and 100";
    }
}
`,
		},
		{
			name: "assert with method call message",
			input: `class X {
    void foo() {
        assert isValid() : getErrorMessage();
    }
}`,
			expected: `class X {

    void foo() {
        assert isValid() : getErrorMessage();
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintLabeledStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "labeled for loop with break",
			input: `class X {
    void foo() {
        outer: for (int i = 0; i < 10; i++) {
            for (int j = 0; j < 10; j++) {
                if (found) break outer;
            }
        }
    }
}`,
			expected: `class X {

    void foo() {
        outer:
        for (int i = 0; i < 10; i++) {
            for (int j = 0; j < 10; j++) {
                if (found) break outer;
            }
        }
    }
}
`,
		},
		{
			name: "labeled for loop with continue",
			input: `class X {
    void foo() {
        outer: for (int i = 0; i < 10; i++) {
            for (int j = 0; j < 10; j++) {
                if (skip) continue outer;
            }
        }
    }
}`,
			expected: `class X {

    void foo() {
        outer:
        for (int i = 0; i < 10; i++) {
            for (int j = 0; j < 10; j++) {
                if (skip) continue outer;
            }
        }
    }
}
`,
		},
		{
			name: "labeled while loop",
			input: `class X {
    void foo() {
        loop: while (true) {
            if (done) break loop;
        }
    }
}`,
			expected: `class X {

    void foo() {
        loop:
        while (true) {
            if (done) break loop;
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintStaticInitializers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "static initializer block",
			input: `class X {
    static Map<String, Integer> map;
    static {
        map = new HashMap<>();
        map.put("one", 1);
        map.put("two", 2);
    }
}`,
			expected: "class X {\n    static Map<String, Integer> map;\n    static {\n        map = new HashMap<>();\n        map.put(\"one\", 1);\n        map.put(\"two\", 2);\n    }\n}\n",
		},
		{
			name: "multiple static initializers",
			input: `class X {
    static int a;
    static {
        a = computeA();
    }
    static int b;
    static {
        b = computeB();
    }
}`,
			expected: "class X {\n    static int a;\n    static {\n        a = computeA();\n    }\n    static int b;\n    static {\n        b = computeB();\n    }\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintTraditionalSwitch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "traditional switch with fallthrough",
			input: `class X {
    void foo(int x) {
        switch (x) {
            case 1:
            case 2:
                handleOneOrTwo();
                break;
            case 3:
                handleThree();
                break;
            default:
                handleDefault();
        }
    }
}`,
			expected: `class X {

    void foo(int x) {
        switch (x) {
            case 1:
            case 2:
                handleOneOrTwo();
                break;
            case 3:
                handleThree();
                break;
            default:
                handleDefault();
        }
    }
}
`,
		},
		{
			name: "switch with enum",
			input: `class X {
    void foo(Day day) {
        switch (day) {
            case MONDAY:
            case TUESDAY:
            case WEDNESDAY:
            case THURSDAY:
            case FRIDAY:
                System.out.println("Weekday");
                break;
            case SATURDAY:
            case SUNDAY:
                System.out.println("Weekend");
                break;
        }
    }
}`,
			expected: `class X {

    void foo(Day day) {
        switch (day) {
            case MONDAY:
            case TUESDAY:
            case WEDNESDAY:
            case THURSDAY:
            case FRIDAY:
                System.out.println("Weekday");
                break;
            case SATURDAY:
            case SUNDAY:
                System.out.println("Weekend");
                break;
        }
    }
}
`,
		},
		{
			name: "switch with string",
			input: `class X {
    void foo(String s) {
        switch (s) {
            case "hello":
                greet();
                break;
            case "goodbye":
                farewell();
                break;
            default:
                unknown();
        }
    }
}`,
			expected: `class X {

    void foo(String s) {
        switch (s) {
            case "hello":
                greet();
                break;
            case "goodbye":
                farewell();
                break;
            default:
                unknown();
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintSwitchExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "switch expression with yield",
			input: `class X {
    int foo(int x) {
        return switch (x) {
            case 1 -> 10;
            case 2 -> 20;
            default -> {
                int result = compute(x);
                yield result;
            }
        };
    }
}`,
			expected: `class X {

    int foo(int x) {
        return switch (x) {
            case 1 -> 10;
            case 2 -> 20;
            default -> {
                int result = compute(x);
                yield result;
            }
        };
    }
}
`,
		},
		{
			name: "switch expression with multiple case labels",
			input: `class X {
    String foo(int x) {
        return switch (x) {
            case 1, 2, 3 -> "small";
            case 4, 5, 6 -> "medium";
            default -> "large";
        };
    }
}`,
			expected: `class X {

    String foo(int x) {
        return switch (x) {
            case 1, 2, 3 -> "small";
            case 4, 5, 6 -> "medium";
            default -> "large";
        };
    }
}
`,
		},
		{
			name: "switch expression assigned to variable",
			input: `class X {
    void foo(Day day) {
        int numLetters = switch (day) {
            case MONDAY, FRIDAY, SUNDAY -> 6;
            case TUESDAY -> 7;
            case THURSDAY, SATURDAY -> 8;
            case WEDNESDAY -> 9;
        };
    }
}`,
			expected: `class X {

    void foo(Day day) {
        int numLetters = switch (day) {
            case MONDAY, FRIDAY, SUNDAY -> 6;
            case TUESDAY -> 7;
            case THURSDAY, SATURDAY -> 8;
            case WEDNESDAY -> 9;
        };
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintComplexLambdas(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "lambda with block body",
			input: `class X {
    void foo() {
        list.forEach(item -> {
            process(item);
            log(item);
        });
    }
}`,
			expected: `class X {

    void foo() {
        list.forEach(item -> {
            process(item);
            log(item);
        });
    }
}
`,
		},
		{
			name: "SwingUtilities invokeLater with try-catch",
			input: `class X {
    public static void main(String[] args) {
        SwingUtilities.invokeLater(() -> {
            try {
                var app = new AmpVisor();
                app.run();
            } catch (Exception e) {
                e.printStackTrace();
                System.exit(1);
            }
        });
    }
}`,
			expected: `class X {

    public static void main(String[] args) {
        SwingUtilities.invokeLater(() -> {
            try {
                var app = new AmpVisor();
                app.run();
            } catch (Exception e) {
                e.printStackTrace();
                System.exit(1);
            }
        });
    }
}
`,
		},
		{
			name: "lambda as second argument",
			input: `class X {
    void foo() {
        executor.submit("task", () -> {
            doWork();
        });
    }
}`,
			expected: `class X {

    void foo() {
        executor.submit("task", () -> {
            doWork();
        });
    }
}
`,
		},
		{
			name: "chained call after lambda block",
			input: `class X {
    void foo() {
        CompletableFuture.runAsync(() -> {
            doWork();
        }).thenRun(() -> {
            cleanup();
        });
    }
}`,
			expected: `class X {

    void foo() {
        CompletableFuture.runAsync(() -> {
            doWork();
        }).thenRun(() -> {
            cleanup();
        });
    }
}
`,
		},
		{
			name: "lambda with multiple parameters",
			input: `class X {
    void foo() {
        map.forEach((key, value) -> System.out.println(key + ": " + value));
    }
}`,
			expected: `class X {

    void foo() {
        map.forEach((key, value) -> System.out.println(key + ": " + value));
    }
}
`,
		},
		{
			name: "lambda with typed parameters",
			input: `class X {
    void foo() {
        BiFunction<Integer, Integer, Integer> add = (Integer a, Integer b) -> a + b;
    }
}`,
			expected: `class X {

    void foo() {
        BiFunction<Integer, Integer, Integer> add = (Integer a, Integer b) -> a + b;
    }
}
`,
		},
		{
			name: "nested lambdas",
			input: `class X {
    void foo() {
        list.stream().map(x -> y -> x + y).collect(Collectors.toList());
    }
}`,
			expected: `class X {

    void foo() {
        list.stream()
            .map(x -> y -> x + y)
            .collect(Collectors.toList());
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintVarKeyword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "var with simple initialization",
			input: `class X {
    void foo() {
        var x = 10;
        var s = "hello";
    }
}`,
			expected: `class X {

    void foo() {
        var x = 10;
        var s = "hello";
    }
}
`,
		},
		{
			name: "var with complex type inference",
			input: `class X {
    void foo() {
        var list = new ArrayList<String>();
        var map = new HashMap<String, List<Integer>>();
    }
}`,
			expected: `class X {

    void foo() {
        var list = new ArrayList<String>();
        var map = new HashMap<String, List<Integer>>();
    }
}
`,
		},
		{
			name: "var in for loop",
			input: `class X {
    void foo() {
        for (var i = 0; i < 10; i++) {
            process(i);
        }
    }
}`,
			expected: `class X {

    void foo() {
        for (var i = 0; i < 10; i++) {
            process(i);
        }
    }
}
`,
		},
		{
			name: "var with method call",
			input: `class X {
    void foo() {
        var result = someMethod();
        var stream = list.stream();
    }
}`,
			expected: `class X {

    void foo() {
        var result = someMethod();
        var stream = list.stream();
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintRecordDeclarations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple record",
			input:    `record Point(int x, int y) {}`,
			expected: "record Point(int x, int y) {}\n",
		},
		{
			name: "record with compact constructor",
			input: `record Range(int lo, int hi) {
    Range {
        if (lo > hi) throw new IllegalArgumentException();
    }
}`,
			expected: `record Range(int lo, int hi) {

    Range() {
        if (lo > hi) throw new IllegalArgumentException();
    }
}
`,
		},
		{
			name: "record with method",
			input: `record Point(int x, int y) {
    public double distance() {
        return Math.sqrt(x * x + y * y);
    }
}`,
			expected: `record Point(int x, int y) {

    public double distance() {
        return Math.sqrt(x * x + y * y);
    }
}
`,
		},
		{
			name:     "record with generic type",
			input:    `record Pair<T, U>(T first, U second) {}`,
			expected: "record Pair<T, U>(T first, U second) {}\n",
		},
		{
			name:     "record implementing interface",
			input:    `record Point(int x, int y) implements Serializable {}`,
			expected: "record Point(int x, int y) implements Serializable {}\n",
		},
		{
			name: "record with static field",
			input: `record Point(int x, int y) {
    static Point ORIGIN = new Point(0, 0);
}`,
			expected: "record Point(int x, int y) {\n    static Point ORIGIN = new Point(0, 0);\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintInterfaceDefaultStaticMethods(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "interface with default method",
			input: `interface Greeter {
    default void greet() {
        System.out.println("Hello!");
    }
}`,
			expected: `interface Greeter {

    default void greet() {
        System.out.println("Hello!");
    }
}
`,
		},
		{
			name: "interface with static method",
			input: `interface Util {
    static int add(int a, int b) {
        return a + b;
    }
}`,
			expected: `interface Util {

    static int add(int a, int b) {
        return a + b;
    }
}
`,
		},
		{
			name: "interface with default and abstract",
			input: `interface Shape {
    double area();
    default String describe() {
        return "A shape with area " + area();
    }
}`,
			expected: `interface Shape {

    double area();

    default String describe() {
        return "A shape with area " + area();
    }
}
`,
		},
		{
			name: "interface with private method",
			input: `interface Logger {
    default void log(String msg) {
        logInternal("[LOG] " + msg);
    }
    private void logInternal(String msg) {
        System.out.println(msg);
    }
}`,
			expected: `interface Logger {

    default void log(String msg) {
        logInternal("[LOG] " + msg);
    }

    private void logInternal(String msg) {
        System.out.println(msg);
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintEnumWithConstructorAndMethods(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "enum with constructor",
			input: `enum Planet {
    MERCURY(3.303e23, 2.4397e6),
    VENUS(4.869e24, 6.0518e6),
    EARTH(5.976e24, 6.37814e6);

    private final double mass;
    private final double radius;

    Planet(double mass, double radius) {
        this.mass = mass;
        this.radius = radius;
    }
}`,
			expected: "enum Planet {\n    MERCURY(3.303e23, 2.4397e6),\n    VENUS(4.869e24, 6.0518e6),\n    EARTH(5.976e24, 6.37814e6);\n    private final double mass;\n    private final double radius;\n\n    Planet(double mass, double radius) {\n        this.mass = mass;\n        this.radius = radius;\n    }\n}\n",
		},
		{
			name: "enum with abstract method",
			input: `enum Operation {
    PLUS {
        double apply(double x, double y) { return x + y; }
    },
    MINUS {
        double apply(double x, double y) { return x - y; }
    };

    abstract double apply(double x, double y);
}`,
			expected: `enum Operation {
    PLUS {

        double apply(double x, double y) {
            return x + y;
        }
    },
    MINUS {

        double apply(double x, double y) {
            return x - y;
        }
    };

    abstract double apply(double x, double y);
}
`,
		},
		{
			name: "enum implementing interface",
			input: `enum Direction implements Rotatable {
    NORTH, EAST, SOUTH, WEST;

    public Direction rotate() {
        return values()[(ordinal() + 1) % 4];
    }
}`,
			expected: `enum Direction implements Rotatable {
    NORTH,
    EAST,
    SOUTH,
    WEST;

    public Direction rotate() {
        return values()[(ordinal() + 1) % 4];
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintTextBlocks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple text block",
			input: `class X {
    String html = """
        <html>
            <body>
                <p>Hello</p>
            </body>
        </html>
        """;
}`,
			expected: "class X {\n    String html = \"\"\"\n        <html>\n            <body>\n                <p>Hello</p>\n            </body>\n        </html>\n        \"\"\";\n}\n",
		},
		{
			name: "text block with json",
			input: `class X {
    String json = """
        {
            "name": "test",
            "value": 42
        }
        """;
}`,
			expected: "class X {\n    String json = \"\"\"\n        {\n            \"name\": \"test\",\n            \"value\": 42\n        }\n        \"\"\";\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintAnnotationsOnVariousElements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "annotation on local variable",
			input: `class X {
    void foo() {
        @SuppressWarnings("unchecked")
        List<String> list = (List<String>) obj;
    }
}`,
			expected: `class X {

    void foo() {
        @SuppressWarnings("unchecked")
        List<String> list = (List<String>) obj;
    }
}
`,
		},
		{
			name: "annotation on method parameter",
			input: `class X {
    void foo(@NonNull String s, @Nullable Integer i) {
        process(s, i);
    }
}`,
			expected: "class X {\n\n    void foo(@NonNull\n    String s, @Nullable\n    Integer i) {\n        process(s, i);\n    }\n}\n",
		},
		{
			name: "multiple annotations on method",
			input: `class X {
    @Override
    @Deprecated
    @SuppressWarnings("deprecation")
    public void foo() {
        oldMethod();
    }
}`,
			expected: `class X {

    @Override
    @Deprecated
    @SuppressWarnings("deprecation")
    public void foo() {
        oldMethod();
    }
}
`,
		},
		{
			name: "annotation with array value",
			input: `class X {
    @SuppressWarnings({"unchecked", "deprecation", "rawtypes"})
    void foo() {
        process();
    }
}`,
			expected: `class X {

    @SuppressWarnings({"unchecked", "deprecation", "rawtypes"})
    void foo() {
        process();
    }
}
`,
		},
		{
			name: "annotation with named parameters",
			input: `class X {
    @RequestMapping(value = "/api", method = RequestMethod.GET)
    void foo() {
        process();
    }
}`,
			expected: `class X {

    @RequestMapping(value = "/api", method = RequestMethod.GET)
    void foo() {
        process();
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintInnerClasses(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "static nested class",
			input: `class Outer {
    static class Inner {
        void foo() {
            process();
        }
    }
}`,
			expected: "class Outer {\n    static class Inner {\n\n        void foo() {\n            process();\n        }\n    }\n}\n",
		},
		{
			name: "non-static inner class",
			input: `class Outer {
    class Inner {
        void foo() {
            Outer.this.bar();
        }
    }
}`,
			expected: "class Outer {\n    class Inner {\n\n        void foo() {\n            Outer.this.bar();\n        }\n    }\n}\n",
		},
		{
			name: "local class in method",
			input: `class X {
    void foo() {
        class LocalHelper {
            void help() {
                process();
            }
        }
        new LocalHelper().help();
    }
}`,
			expected: "class X {\n\n    void foo() {\n        class LocalHelper {\n\n            void help() {\n                process();\n            }\n        }\n;\n        new LocalHelper().help();\n    }\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintGenericMethodsAndConstructors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "generic method",
			input: `class X {
    <T> T identity(T value) {
        return value;
    }
}`,
			expected: `class X {

    <T> T identity(T value) {
        return value;
    }
}
`,
		},
		{
			name: "generic method with bounds",
			input: `class X {
    <T extends Comparable<T>> T max(T a, T b) {
        return a.compareTo(b) > 0 ? a : b;
    }
}`,
			expected: `class X {

    <T extends Comparable<T>> T max(T a, T b) {
        return a.compareTo(b) > 0 ? a : b;
    }
}
`,
		},
		{
			name: "generic method with multiple type params",
			input: `class X {
    <K, V> Map<K, V> createMap(K key, V value) {
        Map<K, V> map = new HashMap<>();
        map.put(key, value);
        return map;
    }
}`,
			expected: `class X {

    <K, V> Map<K, V> createMap(K key, V value) {
        Map<K, V> map = new HashMap<>();
        map.put(key, value);
        return map;
    }
}
`,
		},
		{
			name: "generic constructor",
			input: `class Box<T> {
    T value;
    <U extends T> Box(U value) {
        this.value = value;
    }
}`,
			expected: "class Box<T> {\n    T value;\n\n    <U extends T> Box(U value) {\n        this.value = value;\n    }\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintComplexGenerics(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "nested generic types",
			input: `class X {
    Map<String, List<Set<Integer>>> complex;
}`,
			expected: "class X {\n    Map<String, List<Set<Integer>>> complex;\n}\n",
		},
		{
			name: "bounded wildcard types",
			input: `class X {
    void foo(List<? extends Number> nums, List<? super Integer> ints) {
        process(nums, ints);
    }
}`,
			expected: `class X {

    void foo(List<? extends Number> nums, List<? super Integer> ints) {
        process(nums, ints);
    }
}
`,
		},
		{
			name: "recursive type bounds",
			input: `class Enum<E extends Enum<E>> {
    int ordinal;
}`,
			expected: "class Enum<E extends Enum<E>> {\n    int ordinal;\n}\n",
		},
		{
			name: "intersection types",
			input: `class X {
    <T extends Serializable & Comparable<T>> void foo(T value) {
        process(value);
    }
}`,
			expected: `class X {

    <T extends Serializable & Comparable<T>> void foo(T value) {
        process(value);
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintVarargsAndArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "varargs method",
			input: `class X {
    void foo(String... args) {
        for (String arg : args) {
            process(arg);
        }
    }
}`,
			expected: `class X {

    void foo(String... args) {
        for (String arg : args) {
            process(arg);
        }
    }
}
`,
		},
		{
			name: "varargs with preceding params",
			input: `class X {
    void foo(int count, String format, Object... args) {
        process(count, format, args);
    }
}`,
			expected: `class X {

    void foo(int count, String format, Object... args) {
        process(count, format, args);
    }
}
`,
		},
		{
			name: "array field declarations",
			input: `class X {
    int[] singleDim;
    int[][] doubleDim;
    String[][][] tripleDim;
}`,
			expected: "class X {\n    int[] singleDim;\n    int[][] doubleDim;\n    String[][][] tripleDim;\n}\n",
		},
		{
			name: "array initializers",
			input: `class X {
    int[] nums = {1, 2, 3, 4, 5};
    String[] strings = {"a", "b", "c"};
    int[][] matrix = {{1, 2}, {3, 4}};
}`,
			expected: "class X {\n    int[] nums = {1, 2, 3, 4, 5};\n    String[] strings = {\"a\", \"b\", \"c\"};\n    int[][] matrix = {{1, 2}, {3, 4}};\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("formatting mismatch:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPreserveIntentionalBlankLines(t *testing.T) {
	// Blank lines that create logical sections in the input should be preserved
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "blank lines between statement groups",
			input: `class X {
    void setup() {
        // Initialize resources
        resource1 = createResource();
        resource2 = createResource();

        // Configure settings
        config.setOption("a", true);
        config.setOption("b", false);

        // Start processing
        processor.start();
    }
}`,
			expected: `class X {

    void setup() {
        // Initialize resources
        resource1 = createResource();
        resource2 = createResource();

        // Configure settings
        config.setOption("a", true);
        config.setOption("b", false);

        // Start processing
        processor.start();
    }
}
`,
		},
		{
			name: "no blank line added before comment in if block",
			input: `class X {
    void process() {
        if (children.length == 0) {
            // Leaf node - only add if it has actual data
            if (node.keys().length > 0) {
                keys.add(prefix);
            }
        }
    }
}`,
			expected: `class X {

    void process() {
        if (children.length == 0) {
            // Leaf node - only add if it has actual data
            if (node.keys().length > 0) {
                keys.add(prefix);
            }
        }
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := PrettyPrintJava([]byte(tt.input))
			if err != nil {
				t.Fatalf("PrettyPrintJava error: %v", err)
			}
			if string(output) != tt.expected {
				t.Errorf("Intentional blank lines not preserved:\ngot:\n%s\nwant:\n%s", output, tt.expected)
			}
		})
	}
}

func TestPrintFlowSubscriberInlineDefinition(t *testing.T) {
	input := `class Example {
    void subscribe() {
        publisher.subscribe(new Flow.Subscriber<Object>() {
            @Override
            public void onSubscribe(Flow.Subscription subscription) {
                subscription.request(1);
            }

            @Override
            public void onNext(Object item) {
                if (item instanceof Point(var x, var y)) {
                    System.out.println(x + ", " + y);
                } else if (item instanceof String s) {
                    var upper = (String) s.toUpperCase();
                    System.out.println(upper);
                }
            }

            @Override
            public void onError(Throwable throwable) {
                throwable.printStackTrace();
            }

            @Override
            public void onComplete() {
                System.out.println("Done");
            }
        });
    }
}`
	expected := `class Example {

    void subscribe() {
        publisher.subscribe(new Flow.Subscriber<Object>() {

            @Override
            public void onSubscribe(Flow.Subscription subscription) {
                subscription.request(1);
            }

            @Override
            public void onNext(Object item) {
                if (item instanceof Point(var x, var y)) {
                    System.out.println(x + ", " + y);
                } else if (item instanceof String s) {
                    var upper = (String) s.toUpperCase();
                    System.out.println(upper);
                }
            }

            @Override
            public void onError(Throwable throwable) {
                throwable.printStackTrace();
            }

            @Override
            public void onComplete() {
                System.out.println("Done");
            }
        });
    }
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("Flow.Subscriber inline definition formatting mismatch:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

// Tests for formatting issues found in HTTPInMemoryCache.java, AppStore.java, and RootState.java
func TestSwitchCaseArrowIndentation(t *testing.T) {
	// From AppStore.java: switch case arrow blocks should be properly indented
	input := `class X {
    void accept(Object action) {
        switch (action) {
        case Action.One(var x) ->
            {
                process(x);
            }
        case Action.Two(var y) ->
            {
                handle(y);
            }
        }
    }
}`
	expected := `class X {

    void accept(Object action) {
        switch (action) {
            case Action.One(var x) -> {
                process(x);
            }
            case Action.Two(var y) -> {
                handle(y);
            }
        }
    }
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("switch case arrow indentation mismatch:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

func TestLongMethodChainWrapping(t *testing.T) {
	// From AppStore.java: long method chains should be wrapped for readability
	input := `class X {
    List<Item> getAllItems() {
        return snapshots.values().stream().flatMap(List::stream).collect(java.util.stream.Collectors.toMap(Item::key, item -> item, (a, b) -> a.timestamp().isAfter(b.timestamp()) ? a : b)).values().stream().sorted(Comparator.comparing(Item::timestamp).reversed()).limit(MAX_ITEMS).toList();
    }
}`
	expected := `class X {

    List<Item> getAllItems() {
        return snapshots.values().stream()
            .flatMap(List::stream)
            .collect(java.util.stream.Collectors.toMap(
                Item::key,
                item -> item,
                (a, b) -> a.timestamp().isAfter(b.timestamp()) ? a : b
            ))
            .values().stream()
            .sorted(Comparator.comparing(Item::timestamp).reversed())
            .limit(MAX_ITEMS)
            .toList();
    }
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("long method chain wrapping mismatch:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

func TestLongRecordParameterWrapping(t *testing.T) {
	// From RootState.java: records with many parameters should wrap each parameter
	input := `public record RootState(AccountFormData addAccount, View currentView, Optional<ThreadSummary> tickerItem, boolean loading, Flow.Publisher<AppEvent> events, Consumer<AppAction> dispatch, FeedState feedState, Set<String> followedUsernames) {
}`
	expected := `public record RootState(
    AccountFormData addAccount,
    View currentView,
    Optional<ThreadSummary> tickerItem,
    boolean loading,
    Flow.Publisher<AppEvent> events,
    Consumer<AppAction> dispatch,
    FeedState feedState,
    Set<String> followedUsernames
) {}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("long record parameter wrapping mismatch:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

func TestLongTernaryWrapping(t *testing.T) {
	// From HTTPInMemoryCache.java: long ternary expressions should be wrapped
	input := `class X {
    void foo() {
        Optional<String> ifNoneMatch = entry.etag() != null && !entry.etag().isEmpty() ? Optional.of(entry.etag()) : Optional.empty();
    }
}`
	expected := `class X {

    void foo() {
        Optional<String> ifNoneMatch = entry.etag() != null && !entry.etag().isEmpty()
            ? Optional.of(entry.etag())
            : Optional.empty();
    }
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("long ternary wrapping mismatch:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

func TestEmptyMethodBodyFormatting(t *testing.T) {
	// From AppStore.java: empty method bodies should be on same line
	input := `class X {
    @Override
    public void onComplete() {
    }
}`
	expected := `class X {

    @Override
    public void onComplete() {}
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("empty method body formatting mismatch:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}

func TestSingleLineIfReturn(t *testing.T) {
	// From RootState.java: simple single-statement if should stay on one line
	input := `class X {
    List<Item> bufferItems() {
        if (items.isEmpty())
            return List.of();
        return items.subList(0, 1);
    }
}`
	expected := `class X {

    List<Item> bufferItems() {
        if (items.isEmpty()) return List.of();
        return items.subList(0, 1);
    }
}
`
	output, err := PrettyPrintJava([]byte(input))
	if err != nil {
		t.Fatalf("PrettyPrintJava error: %v", err)
	}
	if string(output) != expected {
		t.Errorf("single line if return mismatch:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}
