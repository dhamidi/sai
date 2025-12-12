package parse

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dhamidi/sai/ebnf/grammar"
	"github.com/dhamidi/sai/ebnf/lex"
)

// EarleyParser implements Earley parsing for EBNF grammars.
// It can handle ambiguous grammars by producing a parse forest (SPPF).
type EarleyParser struct {
	grammar   grammar.Grammar
	tokens    []lex.Token
	skipKinds map[string]bool

	// Internal state
	chart    []ItemSet
	filtered []lex.Token // tokens after filtering trivia
}

// Item represents an Earley item: a production with a dot position and origin.
type Item struct {
	Name   string             // Production name
	Expr   grammar.Expression // The expression being parsed
	Dot    int                // Position in the flattened expression
	Origin int                // Chart position where this item started

	// For building the parse forest
	Node *SPPFNode
}

func (item *Item) String() string {
	return fmt.Sprintf("[%s → •%d, %d]", item.Name, item.Dot, item.Origin)
}

// ItemSet is a set of Earley items at a particular chart position.
type ItemSet struct {
	items    []*Item
	itemSet  map[string]bool // for deduplication
	position int
}

func newItemSet(pos int) *ItemSet {
	return &ItemSet{
		items:    make([]*Item, 0),
		itemSet:  make(map[string]bool),
		position: pos,
	}
}

func (s *ItemSet) Add(item *Item) bool {
	key := fmt.Sprintf("%s:%d:%d", item.Name, item.Dot, item.Origin)
	if s.itemSet[key] {
		return false
	}
	s.itemSet[key] = true
	s.items = append(s.items, item)
	return true
}

// SPPFNode represents a node in the Shared Packed Parse Forest.
// This allows representing ambiguous parses compactly.
type SPPFNode struct {
	Label    string        // Production name or terminal
	Start    int           // Start position in token stream
	End      int           // End position in token stream
	Token    *lex.Token    // Non-nil for terminal nodes
	Children [][]*SPPFNode // Each slice is one possible derivation (packed)
}

func (n *SPPFNode) IsTerminal() bool {
	return n.Token != nil
}

func (n *SPPFNode) IsAmbiguous() bool {
	return len(n.Children) > 1
}

// NewEarleyParser creates a new Earley parser.
func NewEarleyParser(g grammar.Grammar, tokens []lex.Token) *EarleyParser {
	return &EarleyParser{
		grammar:   g,
		tokens:    tokens,
		skipKinds: map[string]bool{"WhiteSpace": true, "Comment": true},
	}
}

// SetSkipKinds sets which token kinds to skip.
func (p *EarleyParser) SetSkipKinds(kinds ...string) {
	p.skipKinds = make(map[string]bool)
	for _, k := range kinds {
		p.skipKinds[k] = true
	}
}

// Parse parses starting from the given production and returns a parse forest.
func (p *EarleyParser) Parse(startProduction string) (*SPPFNode, error) {
	prod := p.grammar.Get(startProduction)
	if prod == nil || prod.Expr == nil {
		return nil, fmt.Errorf("production %q not found in grammar", startProduction)
	}

	// Filter trivia tokens
	p.filtered = make([]lex.Token, 0, len(p.tokens))
	for _, tok := range p.tokens {
		if tok.Kind == "EOF" || !p.skipKinds[tok.Kind] {
			p.filtered = append(p.filtered, tok)
		}
	}

	// Initialize chart
	n := len(p.filtered)
	p.chart = make([]ItemSet, n+1)
	for i := range p.chart {
		p.chart[i] = *newItemSet(i)
	}

	// Seed with start production
	p.chart[0].Add(&Item{
		Name:   startProduction,
		Expr:   prod.Expr,
		Dot:    0,
		Origin: 0,
	})

	// Main Earley loop
	for i := 0; i <= n; i++ {
		// Process items at position i
		// Note: items may be added during iteration
		for j := 0; j < len(p.chart[i].items); j++ {
			item := p.chart[i].items[j]

			if p.isComplete(item) {
				p.complete(i, item)
			} else {
				next := p.nextSymbol(item)
				if p.isTerminal(next) {
					p.scan(i, item, next)
				} else {
					p.predict(i, item, next)
				}
			}
		}
	}

	// Find completed start items at position n
	var roots []*Item
	for _, item := range p.chart[n].items {
		if item.Name == startProduction && item.Origin == 0 && p.isComplete(item) {
			roots = append(roots, item)
		}
	}

	if len(roots) == 0 {
		// Find furthest position reached for error reporting
		furthest := 0
		for i := n; i >= 0; i-- {
			if len(p.chart[i].items) > 0 {
				furthest = i
				break
			}
		}
		if furthest < len(p.filtered) && p.filtered[furthest].Kind != "EOF" {
			return nil, fmt.Errorf("parse error at %s: unexpected %q",
				p.filtered[furthest].Position, p.filtered[furthest].Literal)
		}
		return nil, fmt.Errorf("parse error: incomplete parse")
	}

	// Build SPPF from completed items
	return p.buildForest(roots, startProduction, 0, n), nil
}

// isComplete checks if an item has the dot at the end.
func (p *EarleyParser) isComplete(item *Item) bool {
	return item.Dot >= p.exprLength(item.Expr)
}

// exprLength returns the "length" of an expression for dot positioning.
func (p *EarleyParser) exprLength(expr grammar.Expression) int {
	switch e := expr.(type) {
	case grammar.Sequence:
		return len(e)
	case grammar.Alternative:
		return 1 // alternatives are handled specially
	case *grammar.Repetition:
		return 1
	case *grammar.Option:
		return 1
	case *grammar.Group:
		return 1
	default:
		return 1
	}
}

// nextSymbol returns the symbol after the dot, or nil if complete.
func (p *EarleyParser) nextSymbol(item *Item) grammar.Expression {
	switch e := item.Expr.(type) {
	case grammar.Sequence:
		if item.Dot < len(e) {
			return e[item.Dot]
		}
		return nil
	case grammar.Alternative:
		// For alternatives, we expand each alternative as a separate item
		if item.Dot == 0 {
			return item.Expr
		}
		return nil
	default:
		if item.Dot == 0 {
			return item.Expr
		}
		return nil
	}
}

// isTerminal checks if an expression is a terminal (token literal or token name).
func (p *EarleyParser) isTerminal(expr grammar.Expression) bool {
	switch e := expr.(type) {
	case *grammar.Token:
		return true
	case *grammar.Name:
		// Uppercase names are terminals (token kinds)
		if len(e.String) > 0 && e.String[0] >= 'A' && e.String[0] <= 'Z' {
			// But check if it's actually a production in the grammar
			if p.grammar.Has(e.String) {
				// It's a production - could be either terminal or non-terminal
				// Convention: uppercase = token (terminal)
				return true
			}
		}
		return false
	case *grammar.Range:
		return true
	default:
		return false
	}
}

// predict adds new items for a non-terminal.
func (p *EarleyParser) predict(pos int, item *Item, next grammar.Expression) {
	switch e := next.(type) {
	case *grammar.Name:
		// Add items for the named production
		prod := p.grammar.Get(e.String)
		if prod != nil && prod.Expr != nil {
			p.addPredictedItems(pos, e.String, prod.Expr)
		}

	case *grammar.Group:
		// Inline the group
		p.addPredictedItems(pos, item.Name, e.Body)

	case *grammar.Option:
		// Option: try matching the body, or skip (epsilon)
		p.addPredictedItems(pos, item.Name, e.Body)
		// Also add item with dot advanced (epsilon case)
		p.chart[pos].Add(&Item{
			Name:   item.Name,
			Expr:   item.Expr,
			Dot:    item.Dot + 1,
			Origin: item.Origin,
		})

	case *grammar.Repetition:
		// Repetition: try matching the body, or skip (epsilon for zero matches)
		p.addPredictedItems(pos, item.Name, e.Body)
		// Also add item with dot advanced (epsilon case)
		p.chart[pos].Add(&Item{
			Name:   item.Name,
			Expr:   item.Expr,
			Dot:    item.Dot + 1,
			Origin: item.Origin,
		})

	case grammar.Alternative:
		// Add items for each alternative
		for _, alt := range e {
			p.addPredictedItems(pos, item.Name, alt)
		}
	}
}

func (p *EarleyParser) addPredictedItems(pos int, context string, expr grammar.Expression) {
	switch e := expr.(type) {
	case grammar.Sequence:
		p.chart[pos].Add(&Item{
			Name:   context,
			Expr:   e,
			Dot:    0,
			Origin: pos,
		})
	case grammar.Alternative:
		for _, alt := range e {
			p.addPredictedItems(pos, context, alt)
		}
	default:
		// Wrap single expression in implicit sequence
		p.chart[pos].Add(&Item{
			Name:   context,
			Expr:   grammar.Sequence{e},
			Dot:    0,
			Origin: pos,
		})
	}
}

// scan handles terminal matching.
func (p *EarleyParser) scan(pos int, item *Item, next grammar.Expression) {
	if pos >= len(p.filtered) {
		return
	}
	tok := p.filtered[pos]

	matched := false
	switch e := next.(type) {
	case *grammar.Token:
		// Match literal token
		literal := strings.Trim(e.String, "\"")
		if tok.Literal == literal {
			matched = true
		}
	case *grammar.Name:
		// Match token kind
		if tok.Kind == e.String {
			matched = true
		}
		// Also try matching by literal for keywords
		if tok.Literal == e.String {
			matched = true
		}
	case *grammar.Range:
		// Match character range
		if len(tok.Literal) == 1 {
			begin := strings.Trim(e.Begin.String, "\"")
			end := strings.Trim(e.End.String, "\"")
			if len(begin) == 1 && len(end) == 1 {
				ch := tok.Literal[0]
				if ch >= begin[0] && ch <= end[0] {
					matched = true
				}
			}
		}
	}

	if matched {
		// Add item to next chart position with dot advanced
		p.chart[pos+1].Add(&Item{
			Name:   item.Name,
			Expr:   item.Expr,
			Dot:    item.Dot + 1,
			Origin: item.Origin,
		})
	}
}

// complete handles completion of an item.
func (p *EarleyParser) complete(pos int, completed *Item) {
	// Find items at the origin that were waiting for this production
	origin := completed.Origin
	for _, item := range p.chart[origin].items {
		if p.isComplete(item) {
			continue
		}
		next := p.nextSymbol(item)
		if p.matchesCompleted(next, completed) {
			// Advance the waiting item
			p.chart[pos].Add(&Item{
				Name:   item.Name,
				Expr:   item.Expr,
				Dot:    item.Dot + 1,
				Origin: item.Origin,
			})
		}
	}
}

// matchesCompleted checks if a symbol matches a completed item.
func (p *EarleyParser) matchesCompleted(sym grammar.Expression, completed *Item) bool {
	switch e := sym.(type) {
	case *grammar.Name:
		return e.String == completed.Name
	case *grammar.Group:
		// Groups are inlined, so check if the completed item matches the group's body
		return p.matchesCompleted(e.Body, completed)
	case *grammar.Option:
		return p.matchesCompleted(e.Body, completed)
	case *grammar.Repetition:
		return p.matchesCompleted(e.Body, completed)
	}
	return false
}

// buildForest constructs an SPPF node from completed items.
func (p *EarleyParser) buildForest(items []*Item, label string, start, end int) *SPPFNode {
	node := &SPPFNode{
		Label: label,
		Start: start,
		End:   end,
	}

	// For now, build a simple tree (not fully packed forest)
	// A complete implementation would merge nodes with same label/extent
	if len(items) > 0 {
		// Use first completed item to build children
		item := items[0]
		children := p.buildChildren(item, start, end)
		if len(children) > 0 {
			node.Children = append(node.Children, children)
		}
	}

	return node
}

// buildChildren reconstructs children for a completed item.
func (p *EarleyParser) buildChildren(item *Item, start, end int) []*SPPFNode {
	var children []*SPPFNode

	seq, ok := item.Expr.(grammar.Sequence)
	if !ok {
		seq = grammar.Sequence{item.Expr}
	}

	pos := start
	for _, elem := range seq {
		child, newPos := p.buildChild(elem, pos, end)
		if child != nil {
			children = append(children, child)
			pos = newPos
		}
	}

	return children
}

// buildChild builds an SPPF node for a single expression element.
func (p *EarleyParser) buildChild(expr grammar.Expression, start, end int) (*SPPFNode, int) {
	switch e := expr.(type) {
	case *grammar.Token:
		if start < len(p.filtered) {
			tok := p.filtered[start]
			literal := strings.Trim(e.String, "\"")
			if tok.Literal == literal {
				return &SPPFNode{
					Label: tok.Kind,
					Start: start,
					End:   start + 1,
					Token: &tok,
				}, start + 1
			}
		}
		return nil, start

	case *grammar.Name:
		if start < len(p.filtered) {
			tok := p.filtered[start]
			// Check if it's a terminal match
			if tok.Kind == e.String || tok.Literal == e.String {
				return &SPPFNode{
					Label: tok.Kind,
					Start: start,
					End:   start + 1,
					Token: &tok,
				}, start + 1
			}
		}
		// Try as non-terminal
		for checkEnd := start + 1; checkEnd <= end; checkEnd++ {
			if p.hasCompletedItem(e.String, start, checkEnd) {
				node := &SPPFNode{
					Label: e.String,
					Start: start,
					End:   checkEnd,
				}
				// Recursively build children
				items := p.getCompletedItems(e.String, start, checkEnd)
				if len(items) > 0 {
					children := p.buildChildren(items[0], start, checkEnd)
					if len(children) > 0 {
						node.Children = append(node.Children, children)
					}
				}
				return node, checkEnd
			}
		}
		return nil, start

	case *grammar.Option:
		child, newPos := p.buildChild(e.Body, start, end)
		if child != nil {
			return child, newPos
		}
		// Epsilon case - return empty node
		return &SPPFNode{Label: "ε", Start: start, End: start}, start

	case *grammar.Repetition:
		var children []*SPPFNode
		pos := start
		for pos < end {
			child, newPos := p.buildChild(e.Body, pos, end)
			if child == nil || newPos == pos {
				break
			}
			children = append(children, child)
			pos = newPos
		}
		if len(children) == 0 {
			return &SPPFNode{Label: "ε", Start: start, End: start}, start
		}
		node := &SPPFNode{Label: "*", Start: start, End: pos}
		node.Children = append(node.Children, children)
		return node, pos

	case *grammar.Group:
		return p.buildChild(e.Body, start, end)

	case grammar.Sequence:
		pos := start
		var children []*SPPFNode
		for _, elem := range e {
			child, newPos := p.buildChild(elem, pos, end)
			if child == nil {
				return nil, start
			}
			children = append(children, child)
			pos = newPos
		}
		if len(children) == 1 {
			return children[0], pos
		}
		node := &SPPFNode{Label: "seq", Start: start, End: pos}
		node.Children = append(node.Children, children)
		return node, pos

	case grammar.Alternative:
		for _, alt := range e {
			child, newPos := p.buildChild(alt, start, end)
			if child != nil {
				return child, newPos
			}
		}
		return nil, start
	}

	return nil, start
}

// hasCompletedItem checks if there's a completed item for a production at given extent.
func (p *EarleyParser) hasCompletedItem(name string, start, end int) bool {
	if end >= len(p.chart) {
		return false
	}
	for _, item := range p.chart[end].items {
		if item.Name == name && item.Origin == start && p.isComplete(item) {
			return true
		}
	}
	return false
}

// getCompletedItems returns all completed items for a production at given extent.
func (p *EarleyParser) getCompletedItems(name string, start, end int) []*Item {
	var items []*Item
	if end >= len(p.chart) {
		return items
	}
	for _, item := range p.chart[end].items {
		if item.Name == name && item.Origin == start && p.isComplete(item) {
			items = append(items, item)
		}
	}
	return items
}

// ToCST converts an SPPF node to a concrete syntax tree Node.
// For ambiguous parses, it uses the first derivation.
func (p *EarleyParser) ToCST(sppf *SPPFNode) *Node {
	if sppf == nil {
		return nil
	}

	if sppf.IsTerminal() {
		return NewTerminal(*sppf.Token)
	}

	node := NewNonTerminal(sppf.Label)

	if len(sppf.Children) > 0 {
		// Use first derivation
		for _, child := range sppf.Children[0] {
			childNode := p.ToCST(child)
			if childNode != nil && childNode.Kind != "ε" {
				node.AddChild(childNode)
			}
		}
	}

	// Set span from tokens
	if sppf.Start < len(p.filtered) {
		node.Span.Start = p.filtered[sppf.Start].Position
	}
	if sppf.End > 0 && sppf.End <= len(p.filtered) {
		lastTok := p.filtered[sppf.End-1]
		node.Span.End = lex.Position{
			Filename: lastTok.Position.Filename,
			Offset:   lastTok.Position.Offset + len(lastTok.Literal),
			Line:     lastTok.Position.Line,
			Column:   lastTok.Position.Column + len(lastTok.Literal),
		}
	}

	return node
}

// ParseToCST parses and returns a CST directly.
func (p *EarleyParser) ParseToCST(startProduction string) (*Node, error) {
	forest, err := p.Parse(startProduction)
	if err != nil {
		return nil, err
	}
	return p.ToCST(forest), nil
}

// LoadGrammar loads an EBNF grammar from a file.
func LoadGrammar(filename string) (grammar.Grammar, error) {
	f, err := os.Open(filename)
	if err != nil {
		return grammar.Grammar{}, fmt.Errorf("open grammar: %w", err)
	}
	defer f.Close()

	g, err := grammar.Parse(filename, f)
	if err != nil {
		return grammar.Grammar{}, fmt.Errorf("parse grammar: %w", err)
	}

	return g, nil
}

// ParseTokens is a convenience function to parse tokens with a grammar using Earley parsing.
func ParseTokens(g grammar.Grammar, tokens []lex.Token, start string) (*Node, error) {
	parser := NewEarleyParser(g, tokens)
	return parser.ParseToCST(start)
}

// ParseFile parses a file using a lexer grammar and parser grammar.
func ParseFile(lexerGrammar, parserGrammar grammar.Grammar, input []byte, filename, start string) (*Node, error) {
	// Tokenize with lexer grammar
	lexer := lex.NewLexer(lexerGrammar, input, filename)
	tokens, err := lexer.Tokenize()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("tokenize: %w", err)
	}

	// Parse with Earley parser
	parser := NewEarleyParser(parserGrammar, tokens)
	return parser.ParseToCST(start)
}
