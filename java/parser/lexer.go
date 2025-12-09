package parser

import (
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input        []byte
	file         string
	pos          int
	line         int
	column       int
	isModuleInfo bool
}

func NewLexer(input []byte, file string) *Lexer {
	return &Lexer{
		input:        input,
		file:         file,
		pos:          0,
		line:         1,
		column:       1,
		isModuleInfo: isModuleInfoFile(file),
	}
}

func isModuleInfoFile(file string) bool {
	if len(file) < 16 {
		return file == "module-info.java"
	}
	return file[len(file)-16:] == "module-info.java"
}

func (l *Lexer) Position() Position {
	return Position{
		File:   l.file,
		Offset: l.pos,
		Line:   l.line,
		Column: l.column,
	}
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekN(n int) byte {
	if l.pos+n >= len(l.input) {
		return 0
	}
	return l.input[l.pos+n]
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

func (l *Lexer) advanceN(n int) {
	for i := 0; i < n; i++ {
		l.advance()
	}
}

func (l *Lexer) skipWhitespace() bool {
	start := l.pos
	for {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.advance()
		} else {
			break
		}
	}
	return l.pos > start
}

func (l *Lexer) NextToken() Token {
	startPos := l.Position()

	if l.pos >= len(l.input) {
		return Token{Kind: TokenEOF, Span: Span{Start: startPos, End: startPos}}
	}

	ch := l.peek()

	if ch == '/' && l.peekN(1) == '/' {
		return l.scanLineComment(startPos)
	}
	if ch == '/' && l.peekN(1) == '*' {
		return l.scanBlockComment(startPos)
	}

	if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
		return l.scanWhitespace(startPos)
	}

	if isJavaLetter(ch) {
		return l.scanIdentOrKeyword(startPos)
	}

	if isDigit(ch) {
		return l.scanNumber(startPos)
	}

	if ch == '\'' {
		return l.scanCharLiteral(startPos)
	}

	if ch == '"' {
		if l.peekN(1) == '"' && l.peekN(2) == '"' {
			return l.scanTextBlock(startPos)
		}
		return l.scanStringLiteral(startPos)
	}

	return l.scanOperator(startPos)
}

func (l *Lexer) scanWhitespace(start Position) Token {
	for {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.advance()
		} else {
			break
		}
	}
	end := l.Position()
	return Token{
		Kind:    TokenWhitespace,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) scanLineComment(start Position) Token {
	l.advanceN(2)
	for l.peek() != 0 && l.peek() != '\n' {
		l.advance()
	}
	end := l.Position()
	return Token{
		Kind:    TokenLineComment,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) scanBlockComment(start Position) Token {
	l.advanceN(2)
	for {
		if l.peek() == 0 {
			break
		}
		if l.peek() == '*' && l.peekN(1) == '/' {
			l.advanceN(2)
			break
		}
		l.advance()
	}
	end := l.Position()
	return Token{
		Kind:    TokenComment,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) scanIdentOrKeyword(start Position) Token {
	for isJavaLetterOrDigit(l.peek()) {
		l.advance()
	}
	end := l.Position()
	literal := string(l.input[start.Offset:end.Offset])

	// Handle "non-sealed" contextual keyword (Java 17+)
	if literal == "non" && l.peek() == '-' {
		// Check if followed by "sealed"
		remaining := l.input[l.pos:]
		if len(remaining) >= 7 && string(remaining[:7]) == "-sealed" {
			// Check that "sealed" is not followed by more identifier chars
			if len(remaining) == 7 || !isJavaLetterOrDigit(remaining[7]) {
				l.advanceN(7)
				end = l.Position()
				return Token{
					Kind:    TokenNonSealed,
					Span:    Span{Start: start, End: end},
					Literal: "non-sealed",
				}
			}
		}
	}

	kind := LookupKeyword(literal, l.isModuleInfo)
	return Token{
		Kind:    kind,
		Span:    Span{Start: start, End: end},
		Literal: literal,
	}
}

func (l *Lexer) scanNumber(start Position) Token {
	if l.peek() == '0' && (l.peekN(1) == 'x' || l.peekN(1) == 'X') {
		return l.scanHexNumber(start)
	}
	if l.peek() == '0' && (l.peekN(1) == 'b' || l.peekN(1) == 'B') {
		return l.scanBinaryNumber(start)
	}

	isFloat := false
	for isDigit(l.peek()) || l.peek() == '_' {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekN(1)) {
		isFloat = true
		l.advance()
		for isDigit(l.peek()) || l.peek() == '_' {
			l.advance()
		}
	}

	if l.peek() == 'e' || l.peek() == 'E' {
		isFloat = true
		l.advance()
		if l.peek() == '+' || l.peek() == '-' {
			l.advance()
		}
		for isDigit(l.peek()) || l.peek() == '_' {
			l.advance()
		}
	}

	ch := l.peek()
	if ch == 'f' || ch == 'F' || ch == 'd' || ch == 'D' {
		isFloat = true
		l.advance()
	} else if ch == 'l' || ch == 'L' {
		l.advance()
	}

	end := l.Position()
	kind := TokenIntLiteral
	if isFloat {
		kind = TokenFloatLiteral
	}
	return Token{
		Kind:    kind,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) scanHexNumber(start Position) Token {
	l.advanceN(2)
	for isHexDigit(l.peek()) || l.peek() == '_' {
		l.advance()
	}
	isFloat := false
	if l.peek() == '.' {
		isFloat = true
		l.advance()
		for isHexDigit(l.peek()) || l.peek() == '_' {
			l.advance()
		}
	}
	if l.peek() == 'p' || l.peek() == 'P' {
		isFloat = true
		l.advance()
		if l.peek() == '+' || l.peek() == '-' {
			l.advance()
		}
		for isDigit(l.peek()) || l.peek() == '_' {
			l.advance()
		}
	}
	if isFloat {
		if l.peek() == 'f' || l.peek() == 'F' || l.peek() == 'd' || l.peek() == 'D' {
			l.advance()
		}
	} else {
		if l.peek() == 'l' || l.peek() == 'L' {
			l.advance()
		}
	}
	end := l.Position()
	kind := TokenIntLiteral
	if isFloat {
		kind = TokenFloatLiteral
	}
	return Token{
		Kind:    kind,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) scanBinaryNumber(start Position) Token {
	l.advanceN(2)
	for l.peek() == '0' || l.peek() == '1' || l.peek() == '_' {
		l.advance()
	}
	if l.peek() == 'l' || l.peek() == 'L' {
		l.advance()
	}
	end := l.Position()
	return Token{
		Kind:    TokenIntLiteral,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) scanCharLiteral(start Position) Token {
	l.advance()
	for l.peek() != 0 && l.peek() != '\'' {
		if l.peek() == '\\' {
			l.advance()
		}
		l.advance()
	}
	if l.peek() == '\'' {
		l.advance()
	}
	end := l.Position()
	return Token{
		Kind:    TokenCharLiteral,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) scanStringLiteral(start Position) Token {
	l.advance()
	hasEmbeddedExpr := false
	for l.peek() != 0 && l.peek() != '"' && l.peek() != '\n' {
		if l.peek() == '\\' {
			l.advance()
			if l.peek() == '{' {
				hasEmbeddedExpr = true
				l.advance()
				l.skipEmbeddedExpression()
				continue
			}
		}
		l.advance()
	}
	if l.peek() == '"' {
		l.advance()
	}
	end := l.Position()
	kind := TokenStringLiteral
	if hasEmbeddedExpr {
		kind = TokenStringTemplate
	}
	return Token{
		Kind:    kind,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) scanTextBlock(start Position) Token {
	l.advanceN(3)
	hasEmbeddedExpr := false
	for l.peek() != 0 {
		if l.peek() == '"' && l.peekN(1) == '"' && l.peekN(2) == '"' {
			l.advanceN(3)
			break
		}
		if l.peek() == '\\' {
			l.advance()
			if l.peek() == '{' {
				hasEmbeddedExpr = true
				l.advance()
				l.skipEmbeddedExpression()
				continue
			}
		}
		l.advance()
	}
	end := l.Position()
	kind := TokenTextBlock
	if hasEmbeddedExpr {
		kind = TokenTextBlockTemplate
	}
	return Token{
		Kind:    kind,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) skipEmbeddedExpression() {
	depth := 1
	for l.peek() != 0 && depth > 0 {
		ch := l.peek()
		switch ch {
		case '{':
			depth++
			l.advance()
		case '}':
			depth--
			if depth > 0 {
				l.advance()
			}
		case '"':
			if l.peekN(1) == '"' && l.peekN(2) == '"' {
				l.advanceN(3)
				for l.peek() != 0 {
					if l.peek() == '"' && l.peekN(1) == '"' && l.peekN(2) == '"' {
						l.advanceN(3)
						break
					}
					if l.peek() == '\\' {
						l.advance()
					}
					l.advance()
				}
			} else {
				l.advance()
				for l.peek() != 0 && l.peek() != '"' && l.peek() != '\n' {
					if l.peek() == '\\' {
						l.advance()
					}
					l.advance()
				}
				if l.peek() == '"' {
					l.advance()
				}
			}
		case '\'':
			l.advance()
			for l.peek() != 0 && l.peek() != '\'' {
				if l.peek() == '\\' {
					l.advance()
				}
				l.advance()
			}
			if l.peek() == '\'' {
				l.advance()
			}
		case '/':
			if l.peekN(1) == '/' {
				for l.peek() != 0 && l.peek() != '\n' {
					l.advance()
				}
			} else if l.peekN(1) == '*' {
				l.advanceN(2)
				for l.peek() != 0 {
					if l.peek() == '*' && l.peekN(1) == '/' {
						l.advanceN(2)
						break
					}
					l.advance()
				}
			} else {
				l.advance()
			}
		default:
			l.advance()
		}
	}
}

func (l *Lexer) scanOperator(start Position) Token {
	ch := l.peek()

	switch ch {
	case '(':
		l.advance()
		return l.token(TokenLParen, start)
	case ')':
		l.advance()
		return l.token(TokenRParen, start)
	case '{':
		l.advance()
		return l.token(TokenLBrace, start)
	case '}':
		l.advance()
		return l.token(TokenRBrace, start)
	case '[':
		l.advance()
		return l.token(TokenLBracket, start)
	case ']':
		l.advance()
		return l.token(TokenRBracket, start)
	case ';':
		l.advance()
		return l.token(TokenSemicolon, start)
	case ',':
		l.advance()
		return l.token(TokenComma, start)
	case '@':
		l.advance()
		return l.token(TokenAt, start)
	case '~':
		l.advance()
		return l.token(TokenBitNot, start)
	case '?':
		l.advance()
		return l.token(TokenQuestion, start)

	case '.':
		if l.peekN(1) == '.' && l.peekN(2) == '.' {
			l.advanceN(3)
			return l.token(TokenEllipsis, start)
		}
		l.advance()
		return l.token(TokenDot, start)

	case ':':
		if l.peekN(1) == ':' {
			l.advanceN(2)
			return l.token(TokenColonColon, start)
		}
		l.advance()
		return l.token(TokenColon, start)

	case '=':
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenEQ, start)
		}
		l.advance()
		return l.token(TokenAssign, start)

	case '!':
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenNE, start)
		}
		l.advance()
		return l.token(TokenNot, start)

	case '<':
		if l.peekN(1) == '<' {
			if l.peekN(2) == '=' {
				l.advanceN(3)
				return l.token(TokenShlAssign, start)
			}
			l.advanceN(2)
			return l.token(TokenShl, start)
		}
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenLE, start)
		}
		l.advance()
		return l.token(TokenLT, start)

	case '>':
		if l.peekN(1) == '>' {
			if l.peekN(2) == '>' {
				if l.peekN(3) == '=' {
					l.advanceN(4)
					return l.token(TokenUShrAssign, start)
				}
				l.advanceN(3)
				return l.token(TokenUShr, start)
			}
			if l.peekN(2) == '=' {
				l.advanceN(3)
				return l.token(TokenShrAssign, start)
			}
			l.advanceN(2)
			return l.token(TokenShr, start)
		}
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenGE, start)
		}
		l.advance()
		return l.token(TokenGT, start)

	case '&':
		if l.peekN(1) == '&' {
			l.advanceN(2)
			return l.token(TokenAnd, start)
		}
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenAndAssign, start)
		}
		l.advance()
		return l.token(TokenBitAnd, start)

	case '|':
		if l.peekN(1) == '|' {
			l.advanceN(2)
			return l.token(TokenOr, start)
		}
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenOrAssign, start)
		}
		l.advance()
		return l.token(TokenBitOr, start)

	case '^':
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenXorAssign, start)
		}
		l.advance()
		return l.token(TokenBitXor, start)

	case '+':
		if l.peekN(1) == '+' {
			l.advanceN(2)
			return l.token(TokenIncrement, start)
		}
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenPlusAssign, start)
		}
		l.advance()
		return l.token(TokenPlus, start)

	case '-':
		if l.peekN(1) == '-' {
			l.advanceN(2)
			return l.token(TokenDecrement, start)
		}
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenMinusAssign, start)
		}
		if l.peekN(1) == '>' {
			l.advanceN(2)
			return l.token(TokenArrow, start)
		}
		l.advance()
		return l.token(TokenMinus, start)

	case '*':
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenStarAssign, start)
		}
		l.advance()
		return l.token(TokenStar, start)

	case '/':
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenSlashAssign, start)
		}
		l.advance()
		return l.token(TokenSlash, start)

	case '%':
		if l.peekN(1) == '=' {
			l.advanceN(2)
			return l.token(TokenPercentAssign, start)
		}
		l.advance()
		return l.token(TokenPercent, start)
	}

	l.advance()
	end := l.Position()
	return Token{
		Kind:    TokenError,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func (l *Lexer) token(kind TokenKind, start Position) Token {
	end := l.Position()
	return Token{
		Kind:    kind,
		Span:    Span{Start: start, End: end},
		Literal: string(l.input[start.Offset:end.Offset]),
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isJavaLetter(ch byte) bool {
	if ch >= 128 {
		r, _ := utf8.DecodeRune([]byte{ch})
		return unicode.IsLetter(r) || r == '_' || r == '$'
	}
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$'
}

func isJavaLetterOrDigit(ch byte) bool {
	if ch >= 128 {
		r, _ := utf8.DecodeRune([]byte{ch})
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$'
	}
	return isJavaLetter(ch) || isDigit(ch)
}
