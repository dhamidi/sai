package parse

import (
	"strings"
	"testing"

	"github.com/dhamidi/sai/ebnf/grammar"
	"github.com/dhamidi/sai/ebnf/lex"
)

func TestEarleyParser_AlternativeItems(t *testing.T) {
	g, err := grammar.Parse("test", strings.NewReader(`
		classModifier = annotation | "public" | "private" .
		annotation = "@" Identifier .
		Identifier = letter { letter } .
		letter = "a" … "z" .
	`))
	if err != nil {
		t.Fatalf("parse grammar: %v", err)
	}

	tokens := []lex.Token{
		{Kind: "Identifier", Literal: "public"},
	}

	parser := NewEarleyParser(g, tokens)
	parser.SetSkipKinds("WhiteSpace")

	_, err = parser.Parse("classModifier")

	chart := parser.Chart()
	if len(chart) == 0 {
		t.Fatal("chart is empty after parsing")
	}

	foundPublicItem := false
	for _, item := range chart[0].Items() {
		if item.Name == "classModifier" {
			if seq, ok := item.Expr.(grammar.Sequence); ok && len(seq) == 1 {
				if tok, ok := seq[0].(*grammar.Token); ok && strings.Contains(tok.String, "public") {
					foundPublicItem = true
				}
			}
		}
	}

	if !foundPublicItem {
		t.Errorf("expected to find classModifier item for 'public' alternative in chart[0]")
		t.Logf("chart[0] items:")
		for _, item := range chart[0].Items() {
			t.Logf("  %s (expr type: %T)", item, item.Expr)
		}
	}

	if err != nil {
		t.Errorf("parse failed: %v", err)
	}
}

func TestEarleyParser_MultipleRepetitions(t *testing.T) {
	g, err := grammar.Parse("test", strings.NewReader(`
		methodDeclaration = { methodModifier } result methodDeclarator .
		methodModifier = "public" | "static" | "final" .
		result = "void" | "int" .
		methodDeclarator = Identifier "(" ")" .
		Identifier = "a" … "z" { "a" … "z" } .
	`))
	if err != nil {
		t.Fatalf("parse grammar: %v", err)
	}

	tokens := []lex.Token{
		{Kind: "Identifier", Literal: "public"},
		{Kind: "Identifier", Literal: "static"},
		{Kind: "Identifier", Literal: "void"},
		{Kind: "Identifier", Literal: "main"},
		{Kind: "LPAREN", Literal: "("},
		{Kind: "RPAREN", Literal: ")"},
	}

	parser := NewEarleyParser(g, tokens)
	parser.SetSkipKinds("WhiteSpace")

	_, err = parser.Parse("methodDeclaration")
	if err != nil {
		t.Errorf("parse failed: %v", err)
	}
}

func TestEarleyParser_NestedRepetitions(t *testing.T) {
	g, err := grammar.Parse("test", strings.NewReader(`
		packageDeclaration = "package" Identifier { "." Identifier } ";" .
		Identifier = "a" … "z" { "a" … "z" } .
	`))
	if err != nil {
		t.Fatalf("parse grammar: %v", err)
	}

	tokens := []lex.Token{
		{Kind: "Identifier", Literal: "package"},
		{Kind: "Identifier", Literal: "com"},
		{Kind: "DOT", Literal: "."},
		{Kind: "Identifier", Literal: "example"},
		{Kind: "DOT", Literal: "."},
		{Kind: "Identifier", Literal: "foo"},
		{Kind: "SEMI", Literal: ";"},
	}

	parser := NewEarleyParser(g, tokens)
	parser.SetSkipKinds("WhiteSpace")

	_, err = parser.Parse("packageDeclaration")
	if err != nil {
		t.Errorf("parse failed: %v", err)
	}
}

func TestItemSetDeduplication(t *testing.T) {
	set := newItemSet(0)

	item1 := &Item{
		Name:   "test",
		Expr:   grammar.Sequence{&grammar.Token{String: "a"}},
		Dot:    0,
		Origin: 0,
	}
	item2 := &Item{
		Name:   "test",
		Expr:   grammar.Sequence{&grammar.Token{String: "b"}},
		Dot:    0,
		Origin: 0,
	}

	added1 := set.Add(item1)
	if !added1 {
		t.Error("first item should be added")
	}

	added2 := set.Add(item2)
	if !added2 {
		t.Error("second item with different Expr should be added")
	}

	if len(set.Items()) != 2 {
		t.Errorf("expected 2 items, got %d", len(set.Items()))
	}

	added3 := set.Add(item1)
	if added3 {
		t.Error("duplicate item should not be added")
	}
}
