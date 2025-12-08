package parser

import (
	"testing"
)

func TestLexerNewLexer(t *testing.T) {
	lexer := NewLexer([]byte("class Foo {}"), "Test.java")
	pos := lexer.Position()

	if pos.File != "Test.java" {
		t.Errorf("File = %q, want %q", pos.File, "Test.java")
	}
	if pos.Line != 1 {
		t.Errorf("Line = %d, want %d", pos.Line, 1)
	}
	if pos.Column != 1 {
		t.Errorf("Column = %d, want %d", pos.Column, 1)
	}
	if pos.Offset != 0 {
		t.Errorf("Offset = %d, want %d", pos.Offset, 0)
	}
}

func TestLexerKeywords(t *testing.T) {
	tests := []struct {
		input string
		kind  TokenKind
	}{
		{"class", TokenClass},
		{"public", TokenPublic},
		{"private", TokenPrivate},
		{"protected", TokenProtected},
		{"static", TokenStatic},
		{"final", TokenFinal},
		{"abstract", TokenAbstract},
		{"interface", TokenInterface},
		{"extends", TokenExtends},
		{"implements", TokenImplements},
		{"void", TokenVoid},
		{"int", TokenInt},
		{"boolean", TokenBoolean},
		{"if", TokenIf},
		{"else", TokenElse},
		{"for", TokenFor},
		{"while", TokenWhile},
		{"return", TokenReturn},
		{"new", TokenNew},
		{"this", TokenThis},
		{"super", TokenSuper},
		{"true", TokenTrue},
		{"false", TokenFalse},
		{"null", TokenNull},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input), "test.java")
			tok := lexer.NextToken()
			if tok.Kind != tt.kind {
				t.Errorf("Kind = %v, want %v", tok.Kind, tt.kind)
			}
			if tok.Literal != tt.input {
				t.Errorf("Literal = %q, want %q", tok.Literal, tt.input)
			}
		})
	}
}

func TestLexerIdentifiers(t *testing.T) {
	tests := []string{
		"foo",
		"Bar",
		"_private",
		"$special",
		"camelCase",
		"SCREAMING_CASE",
		"with123Numbers",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			lexer := NewLexer([]byte(input), "test.java")
			tok := lexer.NextToken()
			if tok.Kind != TokenIdent {
				t.Errorf("Kind = %v, want %v", tok.Kind, TokenIdent)
			}
			if tok.Literal != input {
				t.Errorf("Literal = %q, want %q", tok.Literal, input)
			}
		})
	}
}

func TestLexerOperators(t *testing.T) {
	tests := []struct {
		input string
		kind  TokenKind
	}{
		{"(", TokenLParen},
		{")", TokenRParen},
		{"{", TokenLBrace},
		{"}", TokenRBrace},
		{"[", TokenLBracket},
		{"]", TokenRBracket},
		{";", TokenSemicolon},
		{",", TokenComma},
		{".", TokenDot},
		{"...", TokenEllipsis},
		{"@", TokenAt},
		{"::", TokenColonColon},
		{":", TokenColon},
		{"=", TokenAssign},
		{"==", TokenEQ},
		{"!=", TokenNE},
		{"<", TokenLT},
		{"<=", TokenLE},
		{">", TokenGT},
		{">=", TokenGE},
		{"&&", TokenAnd},
		{"||", TokenOr},
		{"!", TokenNot},
		{"&", TokenBitAnd},
		{"|", TokenBitOr},
		{"^", TokenBitXor},
		{"~", TokenBitNot},
		{"<<", TokenShl},
		{">>", TokenShr},
		{">>>", TokenUShr},
		{"+", TokenPlus},
		{"-", TokenMinus},
		{"*", TokenStar},
		{"/", TokenSlash},
		{"%", TokenPercent},
		{"++", TokenIncrement},
		{"--", TokenDecrement},
		{"?", TokenQuestion},
		{"->", TokenArrow},
		{"+=", TokenPlusAssign},
		{"-=", TokenMinusAssign},
		{"*=", TokenStarAssign},
		{"/=", TokenSlashAssign},
		{"%=", TokenPercentAssign},
		{"&=", TokenAndAssign},
		{"|=", TokenOrAssign},
		{"^=", TokenXorAssign},
		{"<<=", TokenShlAssign},
		{">>=", TokenShrAssign},
		{">>>=", TokenUShrAssign},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input), "test.java")
			tok := lexer.NextToken()
			if tok.Kind != tt.kind {
				t.Errorf("Kind = %v, want %v", tok.Kind, tt.kind)
			}
			if tok.Literal != tt.input {
				t.Errorf("Literal = %q, want %q", tok.Literal, tt.input)
			}
		})
	}
}

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		input string
		kind  TokenKind
	}{
		{"0", TokenIntLiteral},
		{"123", TokenIntLiteral},
		{"1_000_000", TokenIntLiteral},
		{"123L", TokenIntLiteral},
		{"0x1F", TokenIntLiteral},
		{"0xDEAD_BEEF", TokenIntLiteral},
		{"0b1010", TokenIntLiteral},
		{"0b1010_1010", TokenIntLiteral},
		{"3.14", TokenFloatLiteral},
		{"3.14f", TokenFloatLiteral},
		{"3.14d", TokenFloatLiteral},
		{"1e10", TokenFloatLiteral},
		{"1.5e-10", TokenFloatLiteral},
		{"1.5E+10", TokenFloatLiteral},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input), "test.java")
			tok := lexer.NextToken()
			if tok.Kind != tt.kind {
				t.Errorf("Kind = %v, want %v", tok.Kind, tt.kind)
			}
			if tok.Literal != tt.input {
				t.Errorf("Literal = %q, want %q", tok.Literal, tt.input)
			}
		})
	}
}

func TestLexerStrings(t *testing.T) {
	tests := []struct {
		input string
		kind  TokenKind
	}{
		{`"hello"`, TokenStringLiteral},
		{`"hello world"`, TokenStringLiteral},
		{`"with \"escapes\""`, TokenStringLiteral},
		{`"with\nnewline"`, TokenStringLiteral},
		{`""`, TokenStringLiteral},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input), "test.java")
			tok := lexer.NextToken()
			if tok.Kind != tt.kind {
				t.Errorf("Kind = %v, want %v", tok.Kind, tt.kind)
			}
			if tok.Literal != tt.input {
				t.Errorf("Literal = %q, want %q", tok.Literal, tt.input)
			}
		})
	}
}

func TestLexerCharLiterals(t *testing.T) {
	tests := []string{
		`'a'`,
		`'\n'`,
		`'\''`,
		`'\\'`,
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			lexer := NewLexer([]byte(input), "test.java")
			tok := lexer.NextToken()
			if tok.Kind != TokenCharLiteral {
				t.Errorf("Kind = %v, want %v", tok.Kind, TokenCharLiteral)
			}
			if tok.Literal != input {
				t.Errorf("Literal = %q, want %q", tok.Literal, input)
			}
		})
	}
}

func TestLexerComments(t *testing.T) {
	t.Run("line comment", func(t *testing.T) {
		lexer := NewLexer([]byte("// this is a comment"), "test.java")
		tok := lexer.NextToken()
		if tok.Kind != TokenLineComment {
			t.Errorf("Kind = %v, want %v", tok.Kind, TokenLineComment)
		}
		if tok.Literal != "// this is a comment" {
			t.Errorf("Literal = %q", tok.Literal)
		}
	})

	t.Run("block comment", func(t *testing.T) {
		lexer := NewLexer([]byte("/* block comment */"), "test.java")
		tok := lexer.NextToken()
		if tok.Kind != TokenComment {
			t.Errorf("Kind = %v, want %v", tok.Kind, TokenComment)
		}
		if tok.Literal != "/* block comment */" {
			t.Errorf("Literal = %q", tok.Literal)
		}
	})

	t.Run("multiline block comment", func(t *testing.T) {
		input := "/* line1\n   line2 */"
		lexer := NewLexer([]byte(input), "test.java")
		tok := lexer.NextToken()
		if tok.Kind != TokenComment {
			t.Errorf("Kind = %v, want %v", tok.Kind, TokenComment)
		}
	})
}

func TestLexerTextBlock(t *testing.T) {
	input := `"""
    hello
    world
    """`
	lexer := NewLexer([]byte(input), "test.java")
	tok := lexer.NextToken()
	if tok.Kind != TokenTextBlock {
		t.Errorf("Kind = %v, want %v", tok.Kind, TokenTextBlock)
	}
}

func TestLexerWhitespace(t *testing.T) {
	lexer := NewLexer([]byte("   \t\n  "), "test.java")
	tok := lexer.NextToken()
	if tok.Kind != TokenWhitespace {
		t.Errorf("Kind = %v, want %v", tok.Kind, TokenWhitespace)
	}
}

func TestLexerEOF(t *testing.T) {
	lexer := NewLexer([]byte(""), "test.java")
	tok := lexer.NextToken()
	if tok.Kind != TokenEOF {
		t.Errorf("Kind = %v, want %v", tok.Kind, TokenEOF)
	}
}

func TestLexerPositionTracking(t *testing.T) {
	lexer := NewLexer([]byte("foo\nbar"), "test.java")

	tok1 := lexer.NextToken()
	if tok1.Span.Start.Line != 1 || tok1.Span.Start.Column != 1 {
		t.Errorf("First token at (%d, %d), want (1, 1)", tok1.Span.Start.Line, tok1.Span.Start.Column)
	}

	lexer.NextToken() // newline whitespace

	tok2 := lexer.NextToken()
	if tok2.Span.Start.Line != 2 || tok2.Span.Start.Column != 1 {
		t.Errorf("Second token at (%d, %d), want (2, 1)", tok2.Span.Start.Line, tok2.Span.Start.Column)
	}
}

func TestLexerSequence(t *testing.T) {
	input := "public class Foo { }"
	lexer := NewLexer([]byte(input), "test.java")

	expected := []TokenKind{
		TokenPublic,
		TokenWhitespace,
		TokenClass,
		TokenWhitespace,
		TokenIdent,
		TokenWhitespace,
		TokenLBrace,
		TokenWhitespace,
		TokenRBrace,
		TokenEOF,
	}

	for i, want := range expected {
		tok := lexer.NextToken()
		if tok.Kind != want {
			t.Errorf("Token %d: Kind = %v, want %v", i, tok.Kind, want)
		}
	}
}

func TestLexerUnknownCharacter(t *testing.T) {
	lexer := NewLexer([]byte("#"), "test.java")
	tok := lexer.NextToken()
	if tok.Kind != TokenError {
		t.Errorf("Kind = %v, want %v", tok.Kind, TokenError)
	}
}
