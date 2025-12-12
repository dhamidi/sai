// Package ebnflex provides lexical scanning based on EBNF grammars.
package ebnflex

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/exp/ebnf"
)

// Position represents a location in source code.
type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

func (p Position) String() string {
	if p.Filename != "" {
		return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line, p.Column)
	}
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Token represents a lexical token with its position.
type Token struct {
	Kind     string
	Literal  string
	Position Position
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %q", t.Position, t.Kind, t.Literal)
}

// memoKey is used for memoization of match results.
type memoKey struct {
	name   string
	offset int
}

// Lexer tokenizes input based on an EBNF grammar.
type Lexer struct {
	grammar  ebnf.Grammar
	input    []byte
	filename string
	pos      int
	line     int
	column   int
	memo     map[memoKey]int  // memoization cache: key -> match length (-1 = no match)
	visiting map[memoKey]bool // cycle detection
}

// NewLexer creates a lexer for the given grammar and input.
func NewLexer(grammar ebnf.Grammar, input []byte, filename string) *Lexer {
	return &Lexer{
		grammar:  grammar,
		input:    input,
		filename: filename,
		pos:      0,
		line:     1,
		column:   1,
		memo:     make(map[memoKey]int),
		visiting: make(map[memoKey]bool),
	}
}

// LoadGrammar loads an EBNF grammar from a file.
func LoadGrammar(filename string) (ebnf.Grammar, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open grammar: %w", err)
	}
	defer f.Close()

	grammar, err := ebnf.Parse(filename, f)
	if err != nil {
		return nil, fmt.Errorf("parse grammar: %w", err)
	}

	return grammar, nil
}

// Position returns the current position in the input.
func (l *Lexer) Position() Position {
	return Position{
		Filename: l.filename,
		Offset:   l.pos,
		Line:     l.line,
		Column:   l.column,
	}
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) advance() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return ch
}

// NextToken returns the next token from the input.
// It tries each production in the grammar marked as a token (uppercase first letter)
// and returns the longest match.
func (l *Lexer) NextToken() (Token, error) {
	if l.pos >= len(l.input) {
		return Token{Kind: "EOF", Position: l.Position()}, io.EOF
	}

	startPos := l.Position()
	startOffset := l.pos

	// Clear memoization cache for each new token (positions change)
	l.memo = make(map[memoKey]int)

	var bestMatch string
	var bestKind string
	var bestLen int

	// Try each production that looks like a token (starts with uppercase)
	for name, prod := range l.grammar {
		if prod.Expr == nil {
			continue
		}
		if len(name) == 0 || name[0] < 'A' || name[0] > 'Z' {
			continue // Skip non-terminal productions (lowercase)
		}

		l.visiting = make(map[memoKey]bool)
		matchLen := l.tryMatch(prod.Expr, startOffset, name)
		if matchLen > bestLen {
			bestLen = matchLen
			bestKind = name
			bestMatch = string(l.input[startOffset : startOffset+matchLen])
		}
	}

	if bestLen == 0 {
		// No match - emit single character as error token
		ch := l.advance()
		return Token{
			Kind:     "ERROR",
			Literal:  string(ch),
			Position: startPos,
		}, nil
	}

	// Advance past the match
	for i := 0; i < bestLen; i++ {
		l.advance()
	}

	return Token{
		Kind:     bestKind,
		Literal:  bestMatch,
		Position: startPos,
	}, nil
}

// tryMatch attempts to match an expression at the given offset.
// Returns the length of the match, or 0 if no match.
func (l *Lexer) tryMatch(expr ebnf.Expression, offset int, context string) int {
	switch e := expr.(type) {
	case *ebnf.Token:
		return l.tryMatchToken(e.String, offset)

	case *ebnf.Range:
		return l.tryMatchRange(e.Begin.String, e.End.String, offset)

	case ebnf.Sequence:
		total := 0
		pos := offset
		for _, item := range e {
			n := l.tryMatch(item, pos, context)
			if n == 0 {
				return 0
			}
			total += n
			pos += n
		}
		return total

	case ebnf.Alternative:
		best := 0
		for _, alt := range e {
			n := l.tryMatch(alt, offset, context)
			if n > best {
				best = n
			}
		}
		return best

	case *ebnf.Repetition:
		total := 0
		pos := offset
		for {
			n := l.tryMatch(e.Body, pos, context)
			if n == 0 {
				break
			}
			total += n
			pos += n
		}
		return total

	case *ebnf.Option:
		n := l.tryMatch(e.Body, offset, context)
		// Option always succeeds (returns 0 if body doesn't match)
		return n

	case *ebnf.Group:
		return l.tryMatch(e.Body, offset, context)

	case *ebnf.Name:
		return l.tryMatchName(e.String, offset)

	default:
		return 0
	}
}

// tryMatchName matches a named production with memoization and cycle detection.
func (l *Lexer) tryMatchName(name string, offset int) int {
	key := memoKey{name: name, offset: offset}

	// Check memo cache
	if result, ok := l.memo[key]; ok {
		if result == -1 {
			return 0
		}
		return result
	}

	// Cycle detection - if we're already visiting this production at this offset,
	// return 0 to break the cycle (left recursion)
	if l.visiting[key] {
		return 0
	}

	prod, ok := l.grammar[name]
	if !ok || prod.Expr == nil {
		l.memo[key] = -1
		return 0
	}

	// Mark as visiting
	l.visiting[key] = true

	result := l.tryMatch(prod.Expr, offset, name)

	// Unmark visiting
	delete(l.visiting, key)

	// Cache result
	if result == 0 {
		l.memo[key] = -1
	} else {
		l.memo[key] = result
	}

	return result
}

// tryMatchToken matches a literal string token.
func (l *Lexer) tryMatchToken(token string, offset int) int {
	// Token string includes quotes
	s := strings.Trim(token, "\"")
	if offset+len(s) > len(l.input) {
		return 0
	}
	if string(l.input[offset:offset+len(s)]) == s {
		return len(s)
	}
	return 0
}

// tryMatchRange matches a character range (e.g., "a"…"z").
func (l *Lexer) tryMatchRange(begin, end string, offset int) int {
	if offset >= len(l.input) {
		return 0
	}
	beginChar := strings.Trim(begin, "\"")
	endChar := strings.Trim(end, "\"")
	if len(beginChar) != 1 || len(endChar) != 1 {
		return 0
	}
	ch := l.input[offset]
	if ch >= beginChar[0] && ch <= endChar[0] {
		return 1
	}
	return 0
}

// Tokenize reads all tokens from input.
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for {
		tok, err := l.NextToken()
		if err == io.EOF {
			tokens = append(tokens, tok)
			break
		}
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, tok)
	}
	return tokens, nil
}
