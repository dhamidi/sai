package format

import (
	"bytes"
	"io"
	"sort"

	"github.com/dhamidi/sai/java/parser"
)

type JavaPrettyPrinter struct {
	w            io.Writer
	source       []byte
	comments     []parser.Token
	commentIndex int
	indent       int
	indentStr    string
	atLineStart  bool
	lastLine     int
}

func NewJavaPrettyPrinter(w io.Writer) *JavaPrettyPrinter {
	return &JavaPrettyPrinter{
		w:           w,
		indentStr:   "    ",
		atLineStart: true,
		lastLine:    1,
	}
}

func (p *JavaPrettyPrinter) Print(node *parser.Node, source []byte, comments []parser.Token) error {
	p.source = source
	p.comments = comments
	sort.Slice(p.comments, func(i, j int) bool {
		if p.comments[i].Span.Start.Line != p.comments[j].Span.Start.Line {
			return p.comments[i].Span.Start.Line < p.comments[j].Span.Start.Line
		}
		return p.comments[i].Span.Start.Column < p.comments[j].Span.Start.Column
	})
	p.commentIndex = 0

	p.printNode(node)
	p.emitRemainingComments()
	return nil
}

func (p *JavaPrettyPrinter) printNode(node *parser.Node) {
	switch node.Kind {
	case parser.KindCompilationUnit:
		p.printCompilationUnit(node)
	case parser.KindPackageDecl:
		p.printPackageDecl(node)
	case parser.KindImportDecl:
		p.printImportDecl(node)
	case parser.KindModuleDecl:
		p.printModuleDecl(node)
	case parser.KindRequiresDirective:
		p.printRequiresDirective(node)
	case parser.KindExportsDirective:
		p.printExportsDirective(node)
	case parser.KindOpensDirective:
		p.printOpensDirective(node)
	case parser.KindUsesDirective:
		p.printUsesDirective(node)
	case parser.KindProvidesDirective:
		p.printProvidesDirective(node)
	case parser.KindClassDecl:
		p.printClassDecl(node)
	case parser.KindInterfaceDecl:
		p.printInterfaceDecl(node)
	case parser.KindEnumDecl:
		p.printEnumDecl(node)
	case parser.KindRecordDecl:
		p.printRecordDecl(node)
	case parser.KindAnnotationDecl:
		p.printAnnotationDecl(node)
	case parser.KindFieldDecl:
		p.printFieldDecl(node)
	case parser.KindMethodDecl:
		p.printMethodDecl(node)
	case parser.KindConstructorDecl:
		p.printConstructorDecl(node)
	case parser.KindBlock:
		if p.isStaticInitializer(node) {
			p.printStaticInitializer(node)
		} else {
			p.printBlock(node)
		}
	default:
		p.printGenericNode(node)
	}
}

// isStaticInitializer checks if a KindBlock represents a static initializer
// (has structure: KindIdentifier("static") + KindBlock)
func (p *JavaPrettyPrinter) isStaticInitializer(node *parser.Node) bool {
	if node.Kind != parser.KindBlock {
		return false
	}
	hasStatic := false
	hasBlock := false
	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil && child.Token.Literal == "static" {
			hasStatic = true
		} else if child.Kind == parser.KindBlock {
			hasBlock = true
		}
	}
	return hasStatic && hasBlock
}

func (p *JavaPrettyPrinter) printStaticInitializer(node *parser.Node) {
	p.writeIndent()
	p.write("static ")
	for _, child := range node.Children {
		if child.Kind == parser.KindBlock {
			p.printBlock(child)
			return
		}
	}
}

func (p *JavaPrettyPrinter) printCompilationUnit(node *parser.Node) {
	for _, child := range node.Children {
		p.emitCommentsBeforeLine(child.Span.Start.Line)
		p.printNode(child)
	}
}

func (p *JavaPrettyPrinter) writeIndent() {
	if !p.atLineStart {
		return
	}
	for i := 0; i < p.indent; i++ {
		p.write(p.indentStr)
	}
	p.atLineStart = false
}

func (p *JavaPrettyPrinter) write(s string) {
	p.w.Write([]byte(s))
}

// isStatementKind returns true if the node kind represents a statement
func (p *JavaPrettyPrinter) isStatementKind(kind parser.NodeKind) bool {
	switch kind {
	case parser.KindBlock, parser.KindEmptyStmt, parser.KindExprStmt, parser.KindIfStmt,
		parser.KindForStmt, parser.KindEnhancedForStmt, parser.KindWhileStmt, parser.KindDoStmt,
		parser.KindSwitchStmt, parser.KindSwitchExpr, parser.KindReturnStmt, parser.KindBreakStmt,
		parser.KindContinueStmt, parser.KindThrowStmt, parser.KindTryStmt, parser.KindSynchronizedStmt,
		parser.KindAssertStmt, parser.KindYieldStmt, parser.KindLocalVarDecl, parser.KindLocalClassDecl,
		parser.KindLabeledStmt:
		return true
	}
	return false
}

func PrettyPrintJava(source []byte) ([]byte, error) {
	return PrettyPrintJavaFile(source, "")
}

func PrettyPrintJavaFile(source []byte, filename string) ([]byte, error) {
	opts := []parser.Option{parser.WithComments()}
	if filename != "" {
		opts = append(opts, parser.WithFile(filename))
	}
	pr := parser.ParseCompilationUnit(bytes.NewReader(source), opts...)
	node := pr.Finish()
	if node == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	pp := NewJavaPrettyPrinter(&buf)
	if err := pp.Print(node, source, pr.Comments()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
