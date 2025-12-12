package parse

import (
	"github.com/dhamidi/sai/ebnf/grammar"
	"github.com/dhamidi/sai/ebnf/lex"
)

// Parser wraps EarleyParser for backward compatibility.
type Parser struct {
	*EarleyParser
}

// NewParser creates a new parser (Earley-based) for backward compatibility.
func NewParser(g grammar.Grammar, tokens []lex.Token) *Parser {
	return &Parser{EarleyParser: NewEarleyParser(g, tokens)}
}

// Parse parses starting from the given production.
// This is the backward-compatible interface that returns a CST node.
func (p *Parser) Parse(startProduction string) (*Node, error) {
	return p.EarleyParser.ParseToCST(startProduction)
}

// SetSkipKinds sets which token kinds to skip between terminals.
func (p *Parser) SetSkipKinds(kinds ...string) {
	p.EarleyParser.SetSkipKinds(kinds...)
}
