// Package v2 provides an alternative Java parser using EBNF grammar and Earley parsing.
package v2

import (
	_ "embed"
	"fmt"
	"io"
	"strings"

	"github.com/dhamidi/sai/ebnf/grammar"
	"github.com/dhamidi/sai/ebnf/lex"
	"github.com/dhamidi/sai/ebnf/parse"
)

//go:embed java25.ebnf
var javaGrammarSource []byte

var javaGrammar grammar.Grammar

func init() {
	g, err := grammar.Parse("java25.ebnf", strings.NewReader(string(javaGrammarSource)))
	if err != nil {
		panic(fmt.Sprintf("failed to parse java grammar: %v", err))
	}
	javaGrammar = g
}

// Parser wraps the Earley parser for Java source files.
type Parser struct {
	input    []byte
	filename string
	tokens   []lex.Token
	comments []lex.Token
	cst      *parse.Node
	err      error
}

// Parse parses Java source from a reader and returns a CST.
func Parse(r io.Reader, filename string) (*parse.Node, []lex.Token, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("read input: %w", err)
	}

	return ParseBytes(input, filename)
}

// ParseBytes parses Java source from bytes and returns a CST.
func ParseBytes(input []byte, filename string) (*parse.Node, []lex.Token, error) {
	lexer := lex.NewLexer(javaGrammar, input, filename)
	tokens, err := lexer.Tokenize()
	if err != nil && err != io.EOF {
		return nil, nil, fmt.Errorf("tokenize: %w", err)
	}

	var comments []lex.Token
	var filtered []lex.Token
	for _, tok := range tokens {
		if tok.Kind == "Comment" {
			comments = append(comments, tok)
		}
		filtered = append(filtered, tok)
	}

	parser := parse.NewEarleyParser(javaGrammar, filtered)
	parser.SetSkipKinds("WhiteSpace", "Comment", "LineTerminator")

	cst, err := parser.ParseToCST("compilationUnit")
	if err != nil {
		return nil, comments, fmt.Errorf("parse: %w", err)
	}

	return cst, comments, nil
}

// Grammar returns the Java grammar used by this parser.
func Grammar() grammar.Grammar {
	return javaGrammar
}
