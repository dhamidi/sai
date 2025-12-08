// Package parser provides a streaming, error-tolerant parser for Java source code.
//
// # Overview
//
// The parser consumes bytes incrementally, producing a concrete syntax tree (CST)
// that preserves all source information including whitespace and comments. It is
// designed for IDE-like tooling where incomplete or malformed input is common.
//
// # Architecture
//
//	┌─────────────┐     ┌─────────────┐     ┌─────────────┐
//	│   Input     │────▶│   Lexer     │────▶│   Parser    │
//	│  (bytes)    │     │  (tokens)   │     │   (CST)     │
//	└─────────────┘     └─────────────┘     └─────────────┘
//	                           │                   │
//	                           ▼                   ▼
//	                    ┌─────────────┐     ┌─────────────┐
//	                    │  Position   │     │  ErrorNode  │
//	                    │  Tracking   │     │  Recovery   │
//	                    └─────────────┘     └─────────────┘
//
// # Streaming Interface
//
// The parser implements a push-based streaming model:
//
//	type Parser struct {
//	    // unexported fields
//	}
//
//	// Push feeds bytes into the parser. May be called multiple times
//	// with chunks of input. Returns the number of bytes consumed.
//	func (p *Parser) Push(data []byte) int
//
//	// IsComplete reports whether it is safe to call Finish.
//	// Returns true when the input can be parsed to produce a complete
//	// node without blocking. For example, "1 + " returns false because
//	// the expression is incomplete, while "1 + 2" returns true.
//	func (p *Parser) IsComplete() bool
//
//	// Finish signals end of input and finalizes the parse tree.
//	// Must be called after all Push calls are complete.
//	// Use IsComplete to check if calling Finish is safe.
//	func (p *Parser) Finish() *Node
//
//	// Reset clears parser state for reuse with new input.
//	func (p *Parser) Reset()
//
// # Source Context
//
// Every node in the tree carries precise source location information:
//
//	type Position struct {
//	    File   string // source file path
//	    Offset int    // byte offset from start of file
//	    Line   int    // 1-based line number
//	    Column int    // 1-based column (in bytes, not runes)
//	}
//
//	type Span struct {
//	    Start Position
//	    End   Position
//	}
//
// The parser tracks position incrementally as bytes are pushed, handling
// newlines (LF, CR, CRLF) to maintain accurate line/column information.
//
// # Error Recovery
//
// The parser never panics on malformed input. Instead, it creates ErrorNode
// entries in the tree that capture:
//
//   - The span of unparsable text
//   - A diagnostic message describing the expected grammar
//   - The actual tokens encountered
//
// Recovery strategies:
//
//  1. Statement-level: Skip to next semicolon or closing brace
//  2. Declaration-level: Skip to next class/method/field keyword
//  3. Expression-level: Insert synthetic nodes for missing operands
//
// Example tree with error:
//
//	CompilationUnit
//	├── ClassDeclaration
//	│   └── MethodDeclaration
//	│       └── Block
//	│           ├── ExpressionStatement
//	│           └── ErrorNode("expected ';', got '}'")
//
// # Entry Points
//
// The parser supports multiple entry points for different use cases:
//
//	// ParseCompilationUnit parses a complete .java source file.
//	// This is the standard entry point for file-based parsing.
//	//
//	// Grammar: CompilationUnit → OrdinaryCompilationUnit
//	//                          | CompactCompilationUnit
//	//                          | ModularCompilationUnit
//	func ParseCompilationUnit(opts ...Option) *Parser
//
//	// ParseExpression parses a standalone expression.
//	// Useful for evaluating snippets, REPL input, or template expressions.
//	//
//	// Grammar: Expression → LambdaExpression | AssignmentExpression
//	func ParseExpression(opts ...Option) *Parser
//
// # Completion Support
//
// The parser can be queried for completion context at any position:
//
//	type CompletionContext struct {
//	    Position    Position       // cursor position
//	    Scope       *Scope         // enclosing scope (for local names)
//	    Expected    []TokenKind    // valid next tokens
//	    Incomplete  *Node          // partially parsed node at cursor
//	    InString    bool           // inside string literal
//	    InComment   bool           // inside comment
//	}
//
//	// CompletionAt returns context for code completion at the given offset.
//	// The parser must have processed input up to at least this offset.
//	func (p *Parser) CompletionAt(offset int) *CompletionContext
//
// # Node Types
//
// The CST uses a uniform node structure:
//
//	type Node struct {
//	    Kind     NodeKind   // e.g., KindClassDecl, KindMethodDecl, KindBinaryExpr
//	    Span     Span       // source location
//	    Children []*Node    // child nodes (for non-terminals)
//	    Token    *Token     // lexical token (for terminals)
//	    Error    *Error     // non-nil for error nodes
//	}
//
// Key node kinds (following JLS Chapter 19 grammar):
//
//	// Compilation unit level
//	KindCompilationUnit
//	KindPackageDecl
//	KindImportDecl
//
//	// Type declarations
//	KindClassDecl
//	KindInterfaceDecl
//	KindEnumDecl
//	KindRecordDecl
//	KindAnnotationDecl
//
//	// Members
//	KindFieldDecl
//	KindMethodDecl
//	KindConstructorDecl
//
//	// Statements
//	KindBlock
//	KindIfStmt
//	KindForStmt
//	KindWhileStmt
//	KindReturnStmt
//	KindSwitchStmt
//
//	// Expressions
//	KindBinaryExpr
//	KindUnaryExpr
//	KindCallExpr
//	KindFieldAccess
//	KindArrayAccess
//	KindLambdaExpr
//	KindLiteral
//	KindIdentifier
//
//	// Special
//	KindError
//
// # Configuration
//
//	type Option func(*Parser)
//
//	// WithFile sets the file path for position tracking.
//	func WithFile(path string) Option
//
//	// WithStartLine sets the initial line number (default 1).
//	func WithStartLine(line int) Option
//
// # Thread Safety
//
// A Parser instance is not safe for concurrent use. Create separate
// instances for concurrent parsing of different files.
//
// # Example Usage
//
//	// Parse a complete file
//	p := parser.ParseCompilationUnit(parser.WithFile("Main.java"))
//	p.Push([]byte("package com.example;\n"))
//	p.Push([]byte("public class Main {}\n"))
//	tree := p.Finish()
//
//	// Parse an expression
//	p := parser.ParseExpression()
//	p.Push([]byte("x + y * 2"))
//	tree := p.Finish()
//
//	// Get completion context
//	p := parser.ParseCompilationUnit(parser.WithFile("Main.java"))
//	p.Push([]byte("obj."))
//	ctx := p.CompletionAt(4) // after the dot
//	// ctx.Expected contains method/field tokens
package parser
