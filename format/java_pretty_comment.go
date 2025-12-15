package format

import (
	"github.com/dhamidi/sai/java/parser"
)

func (p *JavaPrettyPrinter) emitCommentsBeforeLine(line int) {
	p.emitCommentsBeforeLineSkippingAnnotationLines(line, nil)
}

func (p *JavaPrettyPrinter) emitCommentsBeforeLineSkippingAnnotationLines(line int, skipLineCommentLines map[int]bool) {
	for p.commentIndex < len(p.comments) {
		comment := p.comments[p.commentIndex]
		if comment.Span.Start.Line >= line {
			break
		}
		// Skip line comments on annotation lines - they'll be emitted as trailing comments
		if comment.Kind == parser.TokenLineComment && skipLineCommentLines[comment.Span.Start.Line] {
			break
		}
		if comment.Span.Start.Line > p.lastLine+1 {
			p.write("\n")
		}
		p.writeIndent()
		p.write(comment.Literal)
		p.write("\n")
		p.atLineStart = true
		p.lastLine = comment.Span.End.Line
		p.commentIndex++
	}
}

func (p *JavaPrettyPrinter) emitRemainingComments() {
	for p.commentIndex < len(p.comments) {
		comment := p.comments[p.commentIndex]
		if comment.Span.Start.Line > p.lastLine+1 {
			p.write("\n")
		}
		p.writeIndent()
		p.write(comment.Literal)
		p.write("\n")
		p.atLineStart = true
		p.lastLine = comment.Span.End.Line
		p.commentIndex++
	}
}

// emitTrailingLineComment emits a line comment on the given line if one exists.
// This is used to emit trailing comments that appear at the end of a line after code.
func (p *JavaPrettyPrinter) emitTrailingLineComment(line int) {
	if p.commentIndex >= len(p.comments) {
		return
	}
	comment := p.comments[p.commentIndex]
	// Only emit if it's a line comment on the exact same line
	if comment.Kind == parser.TokenLineComment && comment.Span.Start.Line == line {
		p.write(" ")
		p.write(comment.Literal)
		p.lastLine = comment.Span.End.Line
		p.commentIndex++
	}
}
