package parser

import (
	"testing"
)

func TestTokenKindString(t *testing.T) {
	tests := []struct {
		kind TokenKind
		want string
	}{
		{TokenEOF, "EOF"},
		{TokenError, "Error"},
		{TokenIdent, "Identifier"},
		{TokenIntLiteral, "IntLiteral"},
		{TokenStringLiteral, "StringLiteral"},
		{TokenTrue, "true"},
		{TokenFalse, "false"},
		{TokenNull, "null"},
		{TokenClass, "class"},
		{TokenPublic, "public"},
		{TokenPrivate, "private"},
		{TokenStatic, "static"},
		{TokenFinal, "final"},
		{TokenVoid, "void"},
		{TokenInt, "int"},
		{TokenLParen, "("},
		{TokenRParen, ")"},
		{TokenLBrace, "{"},
		{TokenRBrace, "}"},
		{TokenSemicolon, ";"},
		{TokenDot, "."},
		{TokenEllipsis, "..."},
		{TokenAssign, "="},
		{TokenEQ, "=="},
		{TokenArrow, "->"},
		{TokenColonColon, "::"},
		{TokenKind(9999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.want {
				t.Errorf("TokenKind(%d).String() = %q, want %q", tt.kind, got, tt.want)
			}
		})
	}
}

func TestLookupKeyword(t *testing.T) {
	tests := []struct {
		ident string
		want  TokenKind
	}{
		{"class", TokenClass},
		{"public", TokenPublic},
		{"private", TokenPrivate},
		{"protected", TokenProtected},
		{"static", TokenStatic},
		{"final", TokenFinal},
		{"abstract", TokenAbstract},
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
		{"instanceof", TokenInstanceof},
		{"synchronized", TokenSynchronized},
		{"myVariable", TokenIdent},
		{"SomeClass", TokenIdent},
		{"notAKeyword", TokenIdent},
		{"", TokenIdent},
	}

	for _, tt := range tests {
		t.Run(tt.ident, func(t *testing.T) {
			if got := LookupKeyword(tt.ident); got != tt.want {
				t.Errorf("LookupKeyword(%q) = %v, want %v", tt.ident, got, tt.want)
			}
		})
	}
}

func TestPosition(t *testing.T) {
	pos := Position{
		File:   "Test.java",
		Offset: 100,
		Line:   5,
		Column: 10,
	}

	if pos.File != "Test.java" {
		t.Errorf("File = %q, want %q", pos.File, "Test.java")
	}
	if pos.Offset != 100 {
		t.Errorf("Offset = %d, want %d", pos.Offset, 100)
	}
	if pos.Line != 5 {
		t.Errorf("Line = %d, want %d", pos.Line, 5)
	}
	if pos.Column != 10 {
		t.Errorf("Column = %d, want %d", pos.Column, 10)
	}
}

func TestSpan(t *testing.T) {
	span := Span{
		Start: Position{File: "Test.java", Line: 1, Column: 1, Offset: 0},
		End:   Position{File: "Test.java", Line: 1, Column: 6, Offset: 5},
	}

	if span.Start.Line != 1 || span.Start.Column != 1 {
		t.Errorf("Start = (%d, %d), want (1, 1)", span.Start.Line, span.Start.Column)
	}
	if span.End.Line != 1 || span.End.Column != 6 {
		t.Errorf("End = (%d, %d), want (1, 6)", span.End.Line, span.End.Column)
	}
}

func TestToken(t *testing.T) {
	tok := Token{
		Kind:    TokenClass,
		Literal: "class",
		Span: Span{
			Start: Position{File: "Test.java", Line: 1, Column: 1, Offset: 0},
			End:   Position{File: "Test.java", Line: 1, Column: 6, Offset: 5},
		},
	}

	if tok.Kind != TokenClass {
		t.Errorf("Kind = %v, want %v", tok.Kind, TokenClass)
	}
	if tok.Literal != "class" {
		t.Errorf("Literal = %q, want %q", tok.Literal, "class")
	}
	if tok.Kind.String() != "class" {
		t.Errorf("Kind.String() = %q, want %q", tok.Kind.String(), "class")
	}
}
