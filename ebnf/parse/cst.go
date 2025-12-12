// Package ebnfparse provides parsing based on EBNF grammars, producing concrete syntax trees.
package parse

import "github.com/dhamidi/sai/ebnf/lex"

// Span represents a range in source code.
type Span struct {
	Start lex.Position
	End   lex.Position
}

// Node represents a node in the concrete syntax tree.
// Leaf nodes have a non-nil Token; interior nodes have Children.
type Node struct {
	Kind     string     // Production name or token kind
	Children []*Node    // Child nodes (nil for terminals)
	Token    *lex.Token // The token (non-nil for terminals)
	Span     Span       // Source span covering this node
	Error    string     // Non-empty if this is an error node
}

// IsError returns true if this node represents a parse error.
func (n *Node) IsError() bool {
	return n.Error != ""
}

// IsTerminal returns true if this is a leaf node (token).
func (n *Node) IsTerminal() bool {
	return n.Token != nil
}

// Text returns the source text of this node.
// For terminals, returns the token literal.
// For non-terminals, returns empty string (caller should use span to extract text).
func (n *Node) Text() string {
	if n.Token != nil {
		return n.Token.Literal
	}
	return ""
}

// AddChild appends a child node and updates the span.
func (n *Node) AddChild(child *Node) {
	if child == nil {
		return
	}
	n.Children = append(n.Children, child)
	// Update span
	if len(n.Children) == 1 {
		n.Span.Start = child.Span.Start
	}
	n.Span.End = child.Span.End
}

// NewTerminal creates a terminal node from a token.
func NewTerminal(tok lex.Token) *Node {
	return &Node{
		Kind:  tok.Kind,
		Token: &tok,
		Span: Span{
			Start: tok.Position,
			End: lex.Position{
				Filename: tok.Position.Filename,
				Offset:   tok.Position.Offset + len(tok.Literal),
				Line:     tok.Position.Line,
				Column:   tok.Position.Column + len(tok.Literal),
			},
		},
	}
}

// NewNonTerminal creates a non-terminal node.
func NewNonTerminal(kind string) *Node {
	return &Node{
		Kind:     kind,
		Children: make([]*Node, 0),
	}
}

// NewError creates an error node.
func NewError(message string, pos lex.Position) *Node {
	return &Node{
		Kind:  "ERROR",
		Error: message,
		Span: Span{
			Start: pos,
			End:   pos,
		},
	}
}
