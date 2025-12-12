package parse

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/dhamidi/sai/ebnf/grammar"
	"github.com/dhamidi/sai/ebnf/lex"
)

var debugEarley = os.Getenv("DEBUG_EARLEY") != ""

// Tracer receives events during Earley parsing for debugging.
type Tracer interface {
	OnPredict(pos int, item *Item, production string)
	OnScan(pos int, item *Item, token lex.Token, matched bool)
	OnComplete(pos int, completed *Item)
	OnItemAdd(pos int, item *Item, reason string)
}

// EarleyParser implements Earley parsing for EBNF grammars.
// It can handle ambiguous grammars by producing a parse forest (SPPF).
type EarleyParser struct {
	grammar   grammar.Grammar
	tokens    []lex.Token
	skipKinds map[string]bool
	tracer    Tracer

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
	key := fmt.Sprintf("%s:%d:%d:%s", item.Name, item.Dot, item.Origin, exprKey(item.Expr))
	if s.itemSet[key] {
		return false
	}
	s.itemSet[key] = true
	s.items = append(s.items, item)
	return true
}

func exprKey(expr grammar.Expression) string {
	switch e := expr.(type) {
	case grammar.Sequence:
		if len(e) == 0 {
			return "seq[]"
		}
		parts := make([]string, len(e))
		for i, sub := range e {
			parts[i] = exprKey(sub)
		}
		return "seq[" + strings.Join(parts, ",") + "]"
	case grammar.Alternative:
		return "alt"
	case *grammar.Name:
		return "name:" + e.String
	case *grammar.Token:
		return "tok:" + e.String
	case *grammar.Repetition:
		return "rep:" + exprKey(e.Body)
	case *grammar.Option:
		return "opt:" + exprKey(e.Body)
	case *grammar.Group:
		return "grp:" + exprKey(e.Body)
	case *grammar.Range:
		return "range"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// Items returns all items in the set.
func (s *ItemSet) Items() []*Item {
	return s.items
}

// Position returns the chart position of this item set.
func (s *ItemSet) Position() int {
	return s.position
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

// SetTracer sets a tracer to receive parsing events.
func (p *EarleyParser) SetTracer(t Tracer) {
	p.tracer = t
}

// Chart returns the Earley chart after parsing.
func (p *EarleyParser) Chart() []ItemSet {
	return p.chart
}

// FilteredTokens returns the tokens after filtering trivia.
func (p *EarleyParser) FilteredTokens() []lex.Token {
	return p.filtered
}

// Parse parses starting from the given production and returns a parse forest.
func (p *EarleyParser) Parse(startProduction string) (*SPPFNode, error) {
	prod := p.grammar.Get(startProduction)
	if prod == nil || prod.Expr == nil {
		return nil, fmt.Errorf("production %q not found in grammar", startProduction)
	}

	// Filter trivia tokens (and exclude EOF - parsing succeeds when all non-EOF tokens are consumed)
	p.filtered = make([]lex.Token, 0, len(p.tokens))
	for _, tok := range p.tokens {
		if tok.Kind == "EOF" {
			continue // Don't include EOF in tokens to parse
		}
		if !p.skipKinds[tok.Kind] {
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
			if p.tracer != nil {
				p.tracer.OnPredict(pos, item, e.String)
			}
			p.addPredictedItems(pos, e.String, prod.Expr)
		}

	case *grammar.Group:
		// Inline the group - recursively predict if it contains a non-terminal
		p.predictOrAdd(pos, item, e.Body)

	case *grammar.Option:
		// Option: try matching the body, or skip (epsilon)
		p.predictOrAdd(pos, item, e.Body)
		// Also add item with dot advanced (epsilon case)
		p.chart[pos].Add(&Item{
			Name:   item.Name,
			Expr:   item.Expr,
			Dot:    item.Dot + 1,
			Origin: item.Origin,
		})

	case *grammar.Repetition:
		// Repetition: try matching the body, or skip (epsilon for zero matches)
		// Create a synthetic item to match the repetition body
		bodySeq, ok := e.Body.(grammar.Sequence)
		if !ok {
			bodySeq = grammar.Sequence{e.Body}
		}
		// Create helper item for matching the repetition body
		// When this completes, we'll advance the original item but allow more repetitions
		repName := fmt.Sprintf("$rep[%s@%d]", item.Name, item.Dot)
		p.chart[pos].Add(&Item{
			Name:   repName,
			Expr:   bodySeq,
			Dot:    0,
			Origin: pos,
		})
		// Predict non-terminals in the body
		p.predictOrAdd(pos, item, e.Body)
		// Also add item with dot advanced (epsilon case)
		p.chart[pos].Add(&Item{
			Name:   item.Name,
			Expr:   item.Expr,
			Dot:    item.Dot + 1,
			Origin: item.Origin,
		})

	case grammar.Alternative:
		// Add items for each alternative
		// We need to create actual items for each alternative so the scanner can match them
		for _, alt := range e {
			p.addPredictedItems(pos, item.Name, alt)
		}
	}
}

// predictOrAdd handles the body of Option/Repetition/Group/Alternative.
// If the body is a non-terminal, it expands it.
// For terminals and sequences, the caller is responsible for handling them.
func (p *EarleyParser) predictOrAdd(pos int, item *Item, body grammar.Expression) {
	switch b := body.(type) {
	case *grammar.Name:
		// Non-terminal: expand the production
		prod := p.grammar.Get(b.String)
		if prod != nil && prod.Expr != nil {
			p.addPredictedItems(pos, b.String, prod.Expr)
		}
	case *grammar.Group:
		p.predictOrAdd(pos, item, b.Body)
	case *grammar.Option:
		p.predictOrAdd(pos, item, b.Body)
		// Epsilon case handled by caller
	case *grammar.Repetition:
		p.predictOrAdd(pos, item, b.Body)
		// Epsilon case handled by caller
	case grammar.Alternative:
		for _, alt := range b {
			p.predictOrAdd(pos, item, alt)
		}
	case grammar.Sequence:
		// For sequences that are the body of a repetition/option/group,
		// we need to predict any non-terminals in the sequence.
		// The sequence itself will be matched by the scanner.
		if len(b) > 0 {
			first := b[0]
			p.predictOrAdd(pos, item, first)
		}
	default:
		// Terminal - nothing to predict, scanner will handle it
	}
}

func (p *EarleyParser) addPredictedItems(pos int, context string, expr grammar.Expression) {
	if debugEarley && context == "packageDeclaration" && pos == 2 {
		fmt.Fprintf(os.Stderr, "addPredictedItems: context=%s, pos=%d, expr=%T\n", context, pos, expr)
		// Print caller info
		for i := 1; i < 10; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			fn := runtime.FuncForPC(pc)
			fmt.Fprintf(os.Stderr, "  caller %d: %s:%d %s\n", i, file, line, fn.Name())
		}
	}
	switch e := expr.(type) {
	case grammar.Sequence:
		p.chart[pos].Add(&Item{
			Name:   context,
			Expr:   e,
			Dot:    0,
			Origin: pos,
		})
	case grammar.Alternative:
		// For each alternative, create a separate item
		// This allows proper completion when the alternative's body completes
		for _, alt := range e {
			p.addPredictedItems(pos, context, alt)
		}
	case *grammar.Name:
		// For a single Name reference, create an item that will complete when the named production completes
		// This is important for alternatives like: compilationUnit = A | B | C
		// We need an item [compilationUnit -> • A, 0] that completes when A completes
		p.chart[pos].Add(&Item{
			Name:   context,
			Expr:   grammar.Sequence{e},
			Dot:    0,
			Origin: pos,
		})
		// Also predict the named production's items
		prod := p.grammar.Get(e.String)
		if prod != nil && prod.Expr != nil {
			p.addPredictedItems(pos, e.String, prod.Expr)
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

	if debugEarley && item.Name == "packageDeclaration" && pos == 0 {
		fmt.Fprintf(os.Stderr, "SCAN: [%s] pos=%d, next=%T, tok=%s %q\n", item.Name, pos, next, tok.Kind, tok.Literal)
	}

	matched := false
	switch e := next.(type) {
	case *grammar.Token:
		// Match literal token
		literal := strings.Trim(e.String, "\"")
		if debugEarley && item.Name == "packageDeclaration" && pos == 0 {
			fmt.Fprintf(os.Stderr, "  Token match: literal=%q, tok.Literal=%q, match=%v\n", literal, tok.Literal, tok.Literal == literal)
		}
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

	if p.tracer != nil {
		p.tracer.OnScan(pos, item, tok, matched)
	}

	if matched {
		// Add item to next chart position with dot advanced
		if debugEarley && item.Name == "packageDeclaration" {
			fmt.Fprintf(os.Stderr, "SCAN matched [%s] at pos=%d, advancing dot from %d to %d, origin=%d\n", item.Name, pos, item.Dot, item.Dot+1, item.Origin)
		}
		newItem := &Item{
			Name:   item.Name,
			Expr:   item.Expr,
			Dot:    item.Dot + 1,
			Origin: item.Origin,
		}
		if p.chart[pos+1].Add(newItem) && p.tracer != nil {
			p.tracer.OnItemAdd(pos+1, newItem, "scan")
		}
	}
}

// complete handles completion of an item.
func (p *EarleyParser) complete(pos int, completed *Item) {
	if p.tracer != nil {
		p.tracer.OnComplete(pos, completed)
	}

	// Handle synthetic repetition items specially
	if strings.HasPrefix(completed.Name, "$rep[") {
		// Extract the parent item info from the name: $rep[itemName@dot]
		// Find the original item that created this repetition and allow more matches
		p.completeRepetition(pos, completed)
		return
	}

	// Find items at the origin that were waiting for this production
	origin := completed.Origin
	if debugEarley && (completed.Name == "packageDeclaration" || completed.Name == "ordinaryCompilationUnit" || completed.Name == "compilationUnit") {
		fmt.Fprintf(os.Stderr, "COMPLETE: [%s] at pos=%d, origin=%d\n", completed.Name, pos, origin)
		fmt.Fprintf(os.Stderr, "  Checking %d items at origin\n", len(p.chart[origin].items))
	}
	for _, item := range p.chart[origin].items {
		if p.isComplete(item) {
			continue
		}
		next := p.nextSymbol(item)
		if debugEarley && (completed.Name == "packageDeclaration" || completed.Name == "ordinaryCompilationUnit" || completed.Name == "compilationUnit") {
			fmt.Fprintf(os.Stderr, "  Item [%s] dot=%d, next=%T\n", item.Name, item.Dot, next)
		}
		if p.matchesCompleted(next, completed) {
			if debugEarley && (completed.Name == "packageDeclaration" || completed.Name == "ordinaryCompilationUnit" || completed.Name == "compilationUnit") {
				fmt.Fprintf(os.Stderr, "    MATCHED! Advancing [%s]\n", item.Name)
			}
			// Advance the waiting item
			newItem := &Item{
				Name:   item.Name,
				Expr:   item.Expr,
				Dot:    item.Dot + 1,
				Origin: item.Origin,
			}
			if p.chart[pos].Add(newItem) && p.tracer != nil {
				p.tracer.OnItemAdd(pos, newItem, "complete")
			}
		}
	}
}

// completeRepetition handles completion of a synthetic repetition item.
func (p *EarleyParser) completeRepetition(pos int, completed *Item) {
	origin := completed.Origin
	
	// Find items at the origin that have a Repetition at their current dot position
	for _, item := range p.chart[origin].items {
		if p.isComplete(item) {
			continue
		}
		next := p.nextSymbol(item)
		if rep, ok := next.(*grammar.Repetition); ok {
			// The repetition body matched once
			// Add item with dot advanced past the repetition
			p.chart[pos].Add(&Item{
				Name:   item.Name,
				Expr:   item.Expr,
				Dot:    item.Dot + 1,
				Origin: item.Origin,
			})
			// Also add another repetition item to allow more matches
			bodySeq, ok := rep.Body.(grammar.Sequence)
			if !ok {
				bodySeq = grammar.Sequence{rep.Body}
			}
			repName := fmt.Sprintf("$rep[%s@%d]", item.Name, item.Dot)
			p.chart[pos].Add(&Item{
				Name:   repName,
				Expr:   bodySeq,
				Dot:    0,
				Origin: pos,
			})
			// IMPORTANT: Also add a copy of the original item at this position
			// so future repetition completions can find it
			p.chart[pos].Add(&Item{
				Name:   item.Name,
				Expr:   item.Expr,
				Dot:    item.Dot, // Keep the same dot (pointing at Repetition)
				Origin: item.Origin,
			})
			// Predict non-terminals in the body for the new position
			p.predictOrAdd(pos, item, rep.Body)
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
	case grammar.Alternative:
		// Check if the completed item matches any alternative
		for _, alt := range e {
			if p.matchesCompleted(alt, completed) {
				return true
			}
		}
		return false
	case grammar.Sequence:
		// A sequence matches if its first element matches
		if len(e) > 0 {
			return p.matchesCompleted(e[0], completed)
		}
		return false
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
