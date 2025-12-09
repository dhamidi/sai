package parser

import "io"

type Option func(*Parser)

func WithFile(path string) Option {
	return func(p *Parser) {
		p.file = path
	}
}

func WithStartLine(line int) Option {
	return func(p *Parser) {
		p.startLine = line
	}
}

func WithComments() Option {
	return func(p *Parser) {
		p.includeComments = true
	}
}

func WithPositions() Option {
	return func(p *Parser) {
		p.includePositions = true
	}
}

type parseFunc func(*Parser) *Node

type Parser struct {
	file             string
	startLine        int
	includeComments  bool
	includePositions bool
	reader           io.Reader
	input            []byte
	lexer            *Lexer
	tokens           []Token
	comments         []Token
	pos              int
	entry            parseFunc
	incomplete       bool
}

func (p *Parser) IncludesPositions() bool {
	return p.includePositions
}

func (p *Parser) Comments() []Token {
	return p.comments
}

func ParseCompilationUnit(r io.Reader, opts ...Option) *Parser {
	p := &Parser{
		startLine: 1,
		reader:    r,
		entry:     (*Parser).parseCompilationUnit,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func ParseExpression(r io.Reader, opts ...Option) *Parser {
	p := &Parser{
		startLine: 1,
		reader:    r,
		entry:     (*Parser).parseExpression,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Parser) readAll() error {
	if p.input != nil {
		return nil
	}
	data, err := io.ReadAll(p.reader)
	if err != nil {
		return err
	}
	p.input = data
	return nil
}

// IsComplete reports whether it is safe to call Finish.
// Returns true when the input can be parsed to produce a complete node
// without blocking. For example, "1 + " returns false because the
// expression is incomplete.
func (p *Parser) IsComplete() bool {
	if err := p.readAll(); err != nil {
		return false
	}
	if len(p.input) == 0 {
		return false
	}
	// Save parser state
	savedLexer := p.lexer
	savedTokens := p.tokens
	savedPos := p.pos
	savedIncomplete := p.incomplete

	// Trial parse
	p.lexer = NewLexer(p.input, p.file)
	p.tokens = nil
	p.pos = 0
	p.incomplete = false
	p.tokenize()
	p.entry(p)

	complete := !p.incomplete

	// Restore parser state
	p.lexer = savedLexer
	p.tokens = savedTokens
	p.pos = savedPos
	p.incomplete = savedIncomplete

	return complete
}

func (p *Parser) Finish() *Node {
	if err := p.readAll(); err != nil {
		return nil
	}
	if len(p.input) == 0 {
		return nil
	}
	p.lexer = NewLexer(p.input, p.file)
	p.tokens = nil
	p.pos = 0
	p.incomplete = false
	p.tokenize()
	result := p.entry(p)
	if p.incomplete {
		return nil
	}
	return result
}

func (p *Parser) Reset(r io.Reader) {
	p.reader = r
	p.input = nil
	p.lexer = nil
	p.tokens = nil
	p.pos = 0
	p.incomplete = false
}

func (p *Parser) tokenize() {
	for {
		tok := p.lexer.NextToken()
		if tok.Kind == TokenWhitespace {
			continue
		}
		if tok.Kind == TokenComment || tok.Kind == TokenLineComment {
			if p.includeComments {
				p.comments = append(p.comments, tok)
			}
			continue
		}
		p.tokens = append(p.tokens, tok)
		if tok.Kind == TokenEOF {
			break
		}
	}
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Kind: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peekN(n int) Token {
	if p.pos+n >= len(p.tokens) {
		return Token{Kind: TokenEOF}
	}
	return p.tokens[p.pos+n]
}

func (p *Parser) advance() Token {
	tok := p.peek()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

func (p *Parser) expect(kind TokenKind) *Token {
	tok := p.peek()
	if tok.Kind == kind {
		p.advance()
		return &tok
	}
	return nil
}

func (p *Parser) expectIdentifier() *Token {
	if p.isIdentifierLike() {
		tok := p.advance()
		return &tok
	}
	return nil
}

func (p *Parser) check(kind TokenKind) bool {
	return p.peek().Kind == kind
}

// mustProgress returns a function that checks if the parser has advanced.
// Call it at the start of a loop iteration, then call the returned function
// at the end to break if no progress was made.
func (p *Parser) mustProgress() func() bool {
	saved := p.pos
	return func() bool {
		if p.pos == saved {
			if !p.check(TokenEOF) {
				p.advance()
			}
			return false
		}
		return true
	}
}

func (p *Parser) match(kinds ...TokenKind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			return true
		}
	}
	return false
}

func (p *Parser) isIdentifierLike() bool {
	switch p.peek().Kind {
	case TokenIdent,
		TokenModule, TokenOpen, TokenRequires, TokenTransitive,
		TokenExports, TokenOpens, TokenTo, TokenUses, TokenProvides, TokenWith,
		TokenVar, TokenYield, TokenRecord, TokenSealed, TokenNonSealed, TokenPermits:
		return true
	}
	return false
}

func (p *Parser) startNode(kind NodeKind) *Node {
	return &Node{
		Kind: kind,
		Span: Span{Start: p.peek().Span.Start},
	}
}

func (p *Parser) finishNode(n *Node) *Node {
	if p.pos > 0 && p.pos <= len(p.tokens) {
		n.Span.End = p.tokens[p.pos-1].Span.End
	} else if len(p.tokens) > 0 {
		n.Span.End = p.tokens[len(p.tokens)-1].Span.End
	}
	return n
}

func (p *Parser) errorNode(msg string, recoverTo []TokenKind, expected ...TokenKind) *Node {
	tok := p.peek()
	if tok.Kind == TokenEOF {
		p.incomplete = true
	}
	node := &Node{
		Kind: KindError,
		Span: Span{Start: tok.Span.Start, End: tok.Span.End},
		Error: &Error{
			Message:  msg,
			Expected: expected,
			Got:      &tok,
		},
	}
	p.recoverTo(recoverTo)
	return node
}

func (p *Parser) recoverTo(kinds []TokenKind) {
	if !p.check(TokenEOF) {
		p.advance()
	}
	if len(kinds) == 0 {
		return
	}
	for !p.check(TokenEOF) {
		for _, kind := range kinds {
			if p.check(kind) {
				return
			}
		}
		p.advance()
	}
}

func (p *Parser) parseCompilationUnit() *Node {
	node := p.startNode(KindCompilationUnit)

	if p.check(TokenPackage) || p.isAnnotatedPackage() {
		node.AddChild(p.parsePackageDecl())
	}

	for p.check(TokenImport) {
		node.AddChild(p.parseImportDecl())
	}

	if p.isModularCompilationUnit() {
		node.AddChild(p.parseModuleDecl())
	} else if p.isCompactCompilationUnit() {
		for !p.check(TokenEOF) {
			node.AddChild(p.parseClassMember())
		}
	} else {
		for !p.check(TokenEOF) {
			// Skip stray semicolons at top level (empty declarations)
			if p.check(TokenSemicolon) {
				p.advance()
				continue
			}
			node.AddChild(p.parseTypeDecl())
		}
	}

	return p.finishNode(node)
}

func (p *Parser) isCompactCompilationUnit() bool {
	if p.check(TokenEOF) {
		return false
	}

	save := p.pos

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	for p.match(TokenPublic, TokenProtected, TokenPrivate,
		TokenAbstract, TokenStatic, TokenFinal,
		TokenStrictfp, TokenNative, TokenSynchronized,
		TokenTransient, TokenVolatile, TokenDefault,
		TokenSealed, TokenNonSealed) {
		p.advance()
	}

	isTypeDecl := false
	switch p.peek().Kind {
	case TokenClass, TokenInterface, TokenEnum, TokenRecord:
		isTypeDecl = true
	case TokenAt:
		if p.peekN(1).Kind == TokenInterface {
			isTypeDecl = true
		}
	}

	p.pos = save
	return !isTypeDecl
}

func (p *Parser) isModularCompilationUnit() bool {
	if p.check(TokenEOF) {
		return false
	}

	save := p.pos

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	if p.check(TokenOpen) {
		p.advance()
	}

	isModule := p.check(TokenModule)
	p.pos = save
	return isModule
}

func (p *Parser) parseModuleDecl() *Node {
	node := p.startNode(KindModuleDecl)

	for p.check(TokenAt) {
		node.AddChild(p.parseAnnotation())
	}

	if p.check(TokenOpen) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
	}

	p.expect(TokenModule)
	node.AddChild(p.parseQualifiedName())

	p.expect(TokenLBrace)
	for !p.check(TokenRBrace) && !p.check(TokenEOF) {
		node.AddChild(p.parseModuleDirective())
	}
	p.expect(TokenRBrace)

	return p.finishNode(node)
}

func (p *Parser) parseModuleDirective() *Node {
	switch {
	case p.check(TokenRequires):
		return p.parseRequiresDirective()
	case p.check(TokenExports):
		return p.parseExportsDirective()
	case p.check(TokenOpens):
		return p.parseOpensDirective()
	case p.check(TokenUses):
		return p.parseUsesDirective()
	case p.check(TokenProvides):
		return p.parseProvidesDirective()
	default:
		return p.errorNode("expected module directive", []TokenKind{
			TokenRequires, TokenExports, TokenOpens, TokenUses, TokenProvides, TokenRBrace,
		})
	}
}

func (p *Parser) parseRequiresDirective() *Node {
	node := p.startNode(KindRequiresDirective)
	p.expect(TokenRequires)

	for p.check(TokenTransitive) || p.check(TokenStatic) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
	}

	node.AddChild(p.parseQualifiedName())
	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseExportsDirective() *Node {
	node := p.startNode(KindExportsDirective)
	p.expect(TokenExports)

	node.AddChild(p.parseQualifiedName())

	if p.check(TokenTo) {
		p.advance()
		node.AddChild(p.parseQualifiedName())
		for p.check(TokenComma) {
			p.advance()
			node.AddChild(p.parseQualifiedName())
		}
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseOpensDirective() *Node {
	node := p.startNode(KindOpensDirective)
	p.expect(TokenOpens)

	node.AddChild(p.parseQualifiedName())

	if p.check(TokenTo) {
		p.advance()
		node.AddChild(p.parseQualifiedName())
		for p.check(TokenComma) {
			p.advance()
			node.AddChild(p.parseQualifiedName())
		}
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseUsesDirective() *Node {
	node := p.startNode(KindUsesDirective)
	p.expect(TokenUses)
	node.AddChild(p.parseQualifiedName())
	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseProvidesDirective() *Node {
	node := p.startNode(KindProvidesDirective)
	p.expect(TokenProvides)
	node.AddChild(p.parseQualifiedName())

	p.expect(TokenWith)
	node.AddChild(p.parseQualifiedName())
	for p.check(TokenComma) {
		p.advance()
		node.AddChild(p.parseQualifiedName())
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) isAnnotatedPackage() bool {
	if !p.check(TokenAt) {
		return false
	}
	save := p.pos
	for p.check(TokenAt) {
		p.parseAnnotation()
	}
	result := p.check(TokenPackage)
	p.pos = save
	return result
}

func (p *Parser) parsePackageDecl() *Node {
	node := p.startNode(KindPackageDecl)

	for p.check(TokenAt) {
		node.AddChild(p.parseAnnotation())
	}

	p.expect(TokenPackage)
	node.AddChild(p.parseQualifiedName())
	p.expect(TokenSemicolon)

	return p.finishNode(node)
}

func (p *Parser) parseImportDecl() *Node {
	node := p.startNode(KindImportDecl)
	p.expect(TokenImport)

	if p.check(TokenModule) || (p.check(TokenIdent) && p.peek().Literal == "module") {
		node.Kind = KindModuleImportDecl
		p.advance()
		node.AddChild(p.parseQualifiedName())
		p.expect(TokenSemicolon)
		return p.finishNode(node)
	}

	if p.check(TokenStatic) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
	}

	node.AddChild(p.parseQualifiedName())

	if p.check(TokenDot) {
		p.advance()
		if tok := p.expect(TokenStar); tok != nil {
			node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
		}
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseQualifiedName() *Node {
	node := p.startNode(KindQualifiedName)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	} else {
		return p.errorNode("expected identifier", nil)
	}

	for p.check(TokenDot) && p.peekN(1).Kind == TokenIdent {
		p.advance()
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
	}

	return p.finishNode(node)
}

func (p *Parser) parseTypeDecl() *Node {
	modifiers := p.parseModifiers()

	switch p.peek().Kind {
	case TokenClass:
		return p.parseClassDecl(modifiers)
	case TokenInterface:
		return p.parseInterfaceDecl(modifiers)
	case TokenEnum:
		return p.parseEnumDecl(modifiers)
	case TokenRecord:
		return p.parseRecordDecl(modifiers)
	case TokenAt:
		if p.peekN(1).Kind == TokenInterface {
			return p.parseAnnotationDecl(modifiers)
		}
	}

	recoverTokens := []TokenKind{
		TokenAt, TokenPublic, TokenPrivate, TokenProtected,
		TokenAbstract, TokenStatic, TokenFinal, TokenStrictfp,
		TokenClass, TokenInterface, TokenEnum, TokenRecord,
	}
	if modifiers != nil && len(modifiers.Children) > 0 {
		return p.errorNode("expected class, interface, enum, record, or @interface", recoverTokens)
	}

	return p.errorNode("expected type declaration", recoverTokens)
}

func (p *Parser) parseModifiers() *Node {
	node := p.startNode(KindModifiers)

	for {
		switch p.peek().Kind {
		case TokenAt:
			if p.peekN(1).Kind == TokenInterface {
				return p.finishNode(node)
			}
			node.AddChild(p.parseAnnotation())
		case TokenPublic, TokenProtected, TokenPrivate,
			TokenAbstract, TokenStatic, TokenFinal,
			TokenStrictfp, TokenNative, TokenSynchronized,
			TokenTransient, TokenVolatile, TokenDefault,
			TokenSealed, TokenNonSealed:
			tok := p.advance()
			node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		default:
			return p.finishNode(node)
		}
	}
}

func (p *Parser) parseAnnotation() *Node {
	node := p.startNode(KindAnnotation)
	p.expect(TokenAt)
	node.AddChild(p.parseQualifiedName())

	if p.check(TokenLParen) {
		p.advance()
		if !p.check(TokenRParen) {
			if p.peekN(1).Kind == TokenAssign {
				for {
					progress := p.mustProgress()
					node.AddChild(p.parseAnnotationElement())
					if !p.check(TokenComma) {
						break
					}
					p.advance()
					if !progress() {
						break
					}
				}
			} else {
				node.AddChild(p.parseAnnotationValue())
			}
		}
		p.expect(TokenRParen)
	}

	return p.finishNode(node)
}

func (p *Parser) parseAnnotationElement() *Node {
	node := p.startNode(KindAnnotationElement)
	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}
	p.expect(TokenAssign)
	node.AddChild(p.parseAnnotationValue())
	return p.finishNode(node)
}

func (p *Parser) parseAnnotationValue() *Node {
	if p.check(TokenAt) {
		return p.parseAnnotation()
	}
	if p.check(TokenLBrace) {
		node := p.startNode(KindArrayInit)
		p.advance()
		for !p.check(TokenRBrace) && !p.check(TokenEOF) {
			node.AddChild(p.parseAnnotationValue())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
		}
		p.expect(TokenRBrace)
		return p.finishNode(node)
	}
	return p.parseExpression()
}

func (p *Parser) parseClassDecl(modifiers *Node) *Node {
	node := p.startNode(KindClassDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}

	p.expect(TokenClass)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	if p.check(TokenLT) {
		node.AddChild(p.parseTypeParameters())
	}

	if p.check(TokenExtends) {
		p.advance()
		node.AddChild(p.parseType())
	}

	if p.check(TokenImplements) {
		p.advance()
		for {
			progress := p.mustProgress()
			node.AddChild(p.parseType())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			if !progress() {
				break
			}
		}
	}

	if p.check(TokenPermits) {
		p.advance()
		for {
			progress := p.mustProgress()
			node.AddChild(p.parseType())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			if !progress() {
				break
			}
		}
	}

	node.AddChild(p.parseClassBody())
	return p.finishNode(node)
}

func (p *Parser) parseInterfaceDecl(modifiers *Node) *Node {
	node := p.startNode(KindInterfaceDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}

	p.expect(TokenInterface)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	if p.check(TokenLT) {
		node.AddChild(p.parseTypeParameters())
	}

	if p.check(TokenExtends) {
		p.advance()
		for {
			progress := p.mustProgress()
			node.AddChild(p.parseType())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			if !progress() {
				break
			}
		}
	}

	if p.check(TokenPermits) {
		p.advance()
		for {
			progress := p.mustProgress()
			node.AddChild(p.parseType())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			if !progress() {
				break
			}
		}
	}

	node.AddChild(p.parseClassBody())
	return p.finishNode(node)
}

func (p *Parser) parseEnumDecl(modifiers *Node) *Node {
	node := p.startNode(KindEnumDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}

	p.expect(TokenEnum)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	if p.check(TokenImplements) {
		p.advance()
		for {
			progress := p.mustProgress()
			node.AddChild(p.parseType())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			if !progress() {
				break
			}
		}
	}

	p.expect(TokenLBrace)

	for p.check(TokenIdent) || p.check(TokenAt) {
		node.AddChild(p.parseEnumConstant())
		if p.check(TokenComma) {
			p.advance()
		} else {
			break
		}
	}

	if p.check(TokenSemicolon) {
		p.advance()
		for !p.check(TokenRBrace) && !p.check(TokenEOF) {
			node.AddChild(p.parseClassMember())
		}
	}

	p.expect(TokenRBrace)
	return p.finishNode(node)
}

func (p *Parser) parseEnumConstant() *Node {
	node := p.startNode(KindFieldDecl)

	for p.check(TokenAt) {
		node.AddChild(p.parseAnnotation())
	}

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	if p.check(TokenLParen) {
		node.AddChild(p.parseArguments())
	}

	if p.check(TokenLBrace) {
		node.AddChild(p.parseClassBody())
	}

	return p.finishNode(node)
}

func (p *Parser) parseRecordDecl(modifiers *Node) *Node {
	node := p.startNode(KindRecordDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}

	p.expect(TokenRecord)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	if p.check(TokenLT) {
		node.AddChild(p.parseTypeParameters())
	}

	node.AddChild(p.parseParameters())

	if p.check(TokenImplements) {
		p.advance()
		for {
			progress := p.mustProgress()
			node.AddChild(p.parseType())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			if !progress() {
				break
			}
		}
	}

	node.AddChild(p.parseClassBody())
	return p.finishNode(node)
}

func (p *Parser) parseAnnotationDecl(modifiers *Node) *Node {
	node := p.startNode(KindAnnotationDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}

	p.expect(TokenAt)
	p.expect(TokenInterface)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	node.AddChild(p.parseClassBody())
	return p.finishNode(node)
}

func (p *Parser) parseTypeParameters() *Node {
	node := p.startNode(KindTypeParameters)
	p.expect(TokenLT)

	for {
		progress := p.mustProgress()
		node.AddChild(p.parseTypeParameter())
		if !p.check(TokenComma) {
			break
		}
		p.advance()
		if !progress() {
			break
		}
	}

	p.expectGT()
	return p.finishNode(node)
}

func (p *Parser) parseTypeParameter() *Node {
	node := p.startNode(KindTypeParameter)

	for p.check(TokenAt) {
		node.AddChild(p.parseAnnotation())
	}

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	if p.check(TokenExtends) {
		p.advance()
		for {
			node.AddChild(p.parseType())
			if !p.check(TokenBitAnd) {
				break
			}
			p.advance()
		}
	}

	return p.finishNode(node)
}

func (p *Parser) parseType() *Node {
	node := p.startNode(KindType)

	for p.check(TokenAt) {
		node.AddChild(p.parseAnnotation())
	}

	switch p.peek().Kind {
	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble, TokenVoid, TokenVar:
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
	case TokenIdent:
		node.AddChild(p.parseQualifiedName())
		if p.check(TokenLT) {
			node.AddChild(p.parseTypeArguments())
		}
		// Handle parameterized inner class types: Outer<T>.Inner or Outer<T>.Inner<U>
		for p.check(TokenDot) && p.peekN(1).Kind == TokenIdent {
			p.advance() // consume dot
			node.AddChild(p.parseQualifiedName())
			if p.check(TokenLT) {
				node.AddChild(p.parseTypeArguments())
			}
		}
	default:
		return p.errorNode("expected type", []TokenKind{TokenIdent, TokenSemicolon, TokenRParen, TokenComma, TokenRBrace})
	}

	for p.check(TokenAt) || p.check(TokenLBracket) {
		progress := p.mustProgress()
		wrapper := p.startNode(KindArrayType)
		for p.check(TokenAt) {
			wrapper.AddChild(p.parseAnnotation())
		}
		if !p.check(TokenLBracket) {
			break
		}
		p.advance()
		p.expect(TokenRBracket)
		wrapper.AddChild(node)
		node = p.finishNode(wrapper)
		if !progress() {
			break
		}
	}

	return p.finishNode(node)
}

func (p *Parser) parseTypeArguments() *Node {
	node := p.startNode(KindTypeArguments)
	p.expect(TokenLT)

	for {
		progress := p.mustProgress()
		node.AddChild(p.parseTypeArgument())
		if !p.check(TokenComma) {
			break
		}
		p.advance()
		if !progress() {
			break
		}
	}

	p.expectGT()
	return p.finishNode(node)
}

func (p *Parser) expectGT() bool {
	switch p.peek().Kind {
	case TokenGT:
		p.advance()
		return true
	case TokenShr:
		p.splitShiftToken(TokenGT)
		return true
	case TokenUShr:
		p.splitShiftToken(TokenShr)
		return true
	case TokenGE:
		p.splitCompareToken(TokenAssign)
		return true
	case TokenShrAssign:
		p.splitShiftToken(TokenGE)
		return true
	case TokenUShrAssign:
		p.splitShiftToken(TokenShrAssign)
		return true
	}
	return false
}

func (p *Parser) splitShiftToken(remainder TokenKind) {
	tok := p.tokens[p.pos]
	newTok := Token{
		Kind:    remainder,
		Literal: tok.Literal[1:],
		Span: Span{
			Start: Position{
				File:   tok.Span.Start.File,
				Offset: tok.Span.Start.Offset + 1,
				Line:   tok.Span.Start.Line,
				Column: tok.Span.Start.Column + 1,
			},
			End: tok.Span.End,
		},
	}
	p.tokens[p.pos] = newTok
}

func (p *Parser) splitCompareToken(remainder TokenKind) {
	tok := p.tokens[p.pos]
	newTok := Token{
		Kind:    remainder,
		Literal: tok.Literal[1:],
		Span: Span{
			Start: Position{
				File:   tok.Span.Start.File,
				Offset: tok.Span.Start.Offset + 1,
				Line:   tok.Span.Start.Line,
				Column: tok.Span.Start.Column + 1,
			},
			End: tok.Span.End,
		},
	}
	p.tokens[p.pos] = newTok
}

func (p *Parser) parseTypeArgument() *Node {
	if p.check(TokenQuestion) {
		return p.parseWildcard()
	}
	return p.parseType()
}

func (p *Parser) parseWildcard() *Node {
	node := p.startNode(KindWildcard)
	p.expect(TokenQuestion)

	if p.check(TokenExtends) || p.check(TokenSuper) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseType())
	}

	return p.finishNode(node)
}

func (p *Parser) parseClassBody() *Node {
	node := p.startNode(KindBlock)
	p.expect(TokenLBrace)

	for !p.check(TokenRBrace) && !p.check(TokenEOF) {
		node.AddChild(p.parseClassMember())
	}

	p.expect(TokenRBrace)
	return p.finishNode(node)
}

func (p *Parser) parseClassMember() *Node {
	if p.check(TokenLBrace) {
		return p.parseBlock()
	}

	if p.check(TokenStatic) && p.peekN(1).Kind == TokenLBrace {
		node := p.startNode(KindBlock)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		block := p.parseBlock()
		node.AddChild(block)
		return p.finishNode(node)
	}

	if p.check(TokenSemicolon) {
		node := p.startNode(KindEmptyStmt)
		p.advance()
		return p.finishNode(node)
	}

	modifiers := p.parseModifiers()

	switch p.peek().Kind {
	case TokenClass:
		return p.parseClassDecl(modifiers)
	case TokenInterface:
		return p.parseInterfaceDecl(modifiers)
	case TokenEnum:
		return p.parseEnumDecl(modifiers)
	case TokenRecord:
		return p.parseRecordDecl(modifiers)
	case TokenAt:
		if p.peekN(1).Kind == TokenInterface {
			return p.parseAnnotationDecl(modifiers)
		}
	}

	if p.check(TokenLT) {
		typeParams := p.parseTypeParameters()
		return p.parseMethodOrConstructor(modifiers, typeParams)
	}

	if p.isIdentifierLike() && p.peekN(1).Kind == TokenLParen {
		return p.parseConstructor(modifiers, nil)
	}

	// Compact constructor for records: public ClassName { ... }
	if p.isIdentifierLike() && p.peekN(1).Kind == TokenLBrace {
		return p.parseCompactConstructor(modifiers)
	}

	typ := p.parseType()

	if p.isIdentifierLike() {
		if p.peekN(1).Kind == TokenLParen {
			return p.parseMethod(modifiers, nil, typ)
		}
		return p.parseField(modifiers, typ)
	}

	return p.errorNode("expected member declaration", []TokenKind{
		TokenAt, TokenPublic, TokenPrivate, TokenProtected,
		TokenAbstract, TokenStatic, TokenFinal, TokenNative,
		TokenSynchronized, TokenTransient, TokenVolatile,
		TokenStrictfp, TokenDefault, TokenSealed, TokenNonSealed,
		TokenClass, TokenInterface, TokenEnum, TokenRecord,
		TokenIdent, TokenVoid, TokenBoolean, TokenByte,
		TokenChar, TokenShort, TokenInt, TokenLong,
		TokenFloat, TokenDouble, TokenLT, TokenRBrace,
	})
}

func (p *Parser) parseMethodOrConstructor(modifiers *Node, typeParams *Node) *Node {
	if p.isIdentifierLike() && p.peekN(1).Kind == TokenLParen {
		return p.parseConstructor(modifiers, typeParams)
	}

	typ := p.parseType()
	return p.parseMethod(modifiers, typeParams, typ)
}

func (p *Parser) parseConstructor(modifiers *Node, typeParams *Node) *Node {
	node := p.startNode(KindConstructorDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}
	if typeParams != nil {
		node.AddChild(typeParams)
	}

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	node.AddChild(p.parseParameters())

	if p.check(TokenThrows) {
		node.AddChild(p.parseThrowsList())
	}

	node.AddChild(p.parseConstructorBody())
	return p.finishNode(node)
}

// parseCompactConstructor parses a compact constructor for records.
// Compact constructors have no parameter list: public ClassName { ... }
func (p *Parser) parseCompactConstructor(modifiers *Node) *Node {
	node := p.startNode(KindConstructorDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	// Compact constructors have no parameters, but we add an empty parameters node
	paramsNode := p.startNode(KindParameters)
	node.AddChild(p.finishNode(paramsNode))

	// Parse the block body (not constructor body - no explicit constructor invocation check needed)
	node.AddChild(p.parseBlock())
	return p.finishNode(node)
}

func (p *Parser) parseConstructorBody() *Node {
	node := p.startNode(KindBlock)
	p.expect(TokenLBrace)

	if p.isExplicitConstructorInvocation() {
		node.AddChild(p.parseExplicitConstructorInvocation())
	}

	for !p.check(TokenRBrace) && !p.check(TokenEOF) {
		node.AddChild(p.parseStatement())
	}

	p.expect(TokenRBrace)
	return p.finishNode(node)
}

func (p *Parser) isExplicitConstructorInvocation() bool {
	save := p.pos

	if p.check(TokenLT) {
		p.skipTypeArguments()
	}

	if p.check(TokenThis) || p.check(TokenSuper) {
		p.advance()
		if p.check(TokenLParen) {
			p.pos = save
			return true
		}
	}

	p.pos = save

	// Check for qualified super: expr.super(...) or expr.<T>super(...)
	// This handles ExpressionName.super() and Primary.super()
	if p.isQualifiedSuperInvocation() {
		return true
	}

	return false
}

// isQualifiedSuperInvocation checks for patterns like:
// - outer.super(...)
// - outer.<T>super(...)
// - (expr).super(...)
func (p *Parser) isQualifiedSuperInvocation() bool {
	save := p.pos
	defer func() { p.pos = save }()

	// Try to parse qualifying expression (identifier chain or primary)
	if p.check(TokenIdent) {
		// Skip identifier chain: a.b.c
		for p.check(TokenIdent) {
			p.advance()
			if p.check(TokenDot) {
				p.advance()
			} else {
				return false
			}
		}
	} else if p.check(TokenLParen) {
		// Skip parenthesized expression
		p.advance()
		depth := 1
		for depth > 0 && !p.check(TokenEOF) {
			if p.check(TokenLParen) {
				depth++
			} else if p.check(TokenRParen) {
				depth--
			}
			p.advance()
		}
		if !p.check(TokenDot) {
			return false
		}
		p.advance()
	} else {
		return false
	}

	// Optional type arguments
	if p.check(TokenLT) {
		p.skipTypeArguments()
	}

	// Must be super followed by (
	if p.check(TokenSuper) {
		p.advance()
		if p.check(TokenLParen) {
			return true
		}
	}

	return false
}

func (p *Parser) parseExplicitConstructorInvocation() *Node {
	node := p.startNode(KindExplicitConstructorInvocation)

	// Check for qualified super: expr.super() or expr.<T>super()
	if !p.check(TokenLT) && !p.check(TokenThis) && !p.check(TokenSuper) {
		// Must be a qualified super invocation
		qualifier := p.parseQualifiedSuperQualifier()
		node.AddChild(qualifier)

		// Optional type arguments after the dot
		if p.check(TokenLT) {
			node.AddChild(p.parseTypeArguments())
		}

		// Must be super
		if p.check(TokenSuper) {
			tok := p.advance()
			node.AddChild(&Node{Kind: KindSuper, Token: &tok, Span: tok.Span})
		}
	} else {
		// Unqualified: [TypeArguments] this(...) or [TypeArguments] super(...)
		if p.check(TokenLT) {
			node.AddChild(p.parseTypeArguments())
		}

		if p.check(TokenThis) {
			tok := p.advance()
			node.AddChild(&Node{Kind: KindThis, Token: &tok, Span: tok.Span})
		} else if p.check(TokenSuper) {
			tok := p.advance()
			node.AddChild(&Node{Kind: KindSuper, Token: &tok, Span: tok.Span})
		}
	}

	node.AddChild(p.parseArguments())
	p.expect(TokenSemicolon)

	return p.finishNode(node)
}

// parseQualifiedSuperQualifier parses the qualifying expression before .super()
// Returns a KindIdentifier, KindQualifiedName, or expression node
func (p *Parser) parseQualifiedSuperQualifier() *Node {
	if p.check(TokenIdent) {
		// Parse identifier chain: a.b.c (stopping before .super)
		node := p.startNode(KindIdentifier)
		tok := p.advance()
		node.Token = &tok
		node.Span = tok.Span
		node = p.finishNode(node)

		for p.check(TokenDot) {
			// Peek ahead to see if next is super or <T>super
			save := p.pos
			p.advance() // consume dot

			if p.check(TokenLT) {
				// Could be type args before super, restore and return
				p.pos = save
				p.advance() // consume the dot before returning
				return node
			}

			if p.check(TokenSuper) {
				// Don't consume super, just the dot
				return node
			}

			// It's another identifier in the chain
			if p.check(TokenIdent) {
				qualNode := p.startNode(KindQualifiedName)
				qualNode.AddChild(node)
				identTok := p.advance()
				qualNode.AddChild(&Node{Kind: KindIdentifier, Token: &identTok, Span: identTok.Span})
				node = p.finishNode(qualNode)
			} else {
				// Unexpected, restore and return what we have
				p.pos = save
				return node
			}
		}

		// Consume trailing dot before super
		if p.check(TokenDot) {
			p.advance()
		}
		return node
	} else if p.check(TokenLParen) {
		// Parse parenthesized expression
		expr := p.parseParenExpr()
		p.expect(TokenDot)
		return expr
	}

	// Fallback: parse as expression
	expr := p.parsePrimaryExpr()
	p.expect(TokenDot)
	return expr
}

func (p *Parser) parseMethod(modifiers *Node, typeParams *Node, returnType *Node) *Node {
	node := p.startNode(KindMethodDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}
	if typeParams != nil {
		node.AddChild(typeParams)
	}
	if returnType != nil {
		node.AddChild(returnType)
	}

	if tok := p.expectIdentifier(); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	node.AddChild(p.parseParameters())

	for p.check(TokenLBracket) {
		p.advance()
		p.expect(TokenRBracket)
	}

	if p.check(TokenThrows) {
		node.AddChild(p.parseThrowsList())
	}

	if p.check(TokenLBrace) {
		node.AddChild(p.parseBlock())
	} else if p.check(TokenDefault) {
		p.advance()
		node.AddChild(p.parseAnnotationValue())
		p.expect(TokenSemicolon)
	} else {
		p.expect(TokenSemicolon)
	}

	return p.finishNode(node)
}

func (p *Parser) parseField(modifiers *Node, typ *Node) *Node {
	node := p.startNode(KindFieldDecl)
	if modifiers != nil {
		node.AddChild(modifiers)
	}
	if typ != nil {
		node.AddChild(typ)
	}

	for {
		progress := p.mustProgress()
		if tok := p.expect(TokenIdent); tok != nil {
			node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
		}

		for p.check(TokenLBracket) {
			p.advance()
			p.expect(TokenRBracket)
		}

		if p.check(TokenAssign) {
			p.advance()
			node.AddChild(p.parseVarInitializer())
		}

		if !p.check(TokenComma) {
			break
		}
		p.advance()
		if !progress() {
			break
		}
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseVarInitializer() *Node {
	if p.check(TokenLBrace) {
		return p.parseArrayInitializer()
	}
	return p.parseExpression()
}

func (p *Parser) parseArrayInitializer() *Node {
	node := p.startNode(KindArrayInit)
	p.expect(TokenLBrace)

	for !p.check(TokenRBrace) && !p.check(TokenEOF) {
		node.AddChild(p.parseVarInitializer())
		if !p.check(TokenComma) {
			break
		}
		p.advance()
		if p.check(TokenRBrace) {
			break
		}
	}

	p.expect(TokenRBrace)
	return p.finishNode(node)
}

func (p *Parser) parseParameters() *Node {
	node := p.startNode(KindParameters)
	p.expect(TokenLParen)

	if !p.check(TokenRParen) {
		if p.isReceiverParameter() {
			node.AddChild(p.parseReceiverParameter())
			if p.check(TokenComma) {
				p.advance()
			}
		}
		for !p.check(TokenRParen) && !p.check(TokenEOF) {
			node.AddChild(p.parseParameter())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
		}
	}

	p.expect(TokenRParen)
	return p.finishNode(node)
}

func (p *Parser) isReceiverParameter() bool {
	save := p.pos

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	switch p.peek().Kind {
	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble:
		p.advance()
	case TokenIdent:
		p.parseQualifiedName()
		if p.check(TokenLT) {
			p.skipTypeArguments()
		}
	default:
		p.pos = save
		return false
	}

	for p.check(TokenLBracket) {
		p.advance()
		if p.check(TokenRBracket) {
			p.advance()
		}
	}

	if p.check(TokenIdent) {
		p.advance()
		if p.check(TokenDot) {
			p.advance()
			if p.check(TokenThis) {
				p.pos = save
				return true
			}
		}
	} else if p.check(TokenThis) {
		p.pos = save
		return true
	}

	p.pos = save
	return false
}

func (p *Parser) parseReceiverParameter() *Node {
	node := p.startNode(KindReceiverParameter)

	for p.check(TokenAt) {
		node.AddChild(p.parseAnnotation())
	}

	node.AddChild(p.parseType())

	if p.check(TokenIdent) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		p.expect(TokenDot)
	}

	p.expect(TokenThis)
	return p.finishNode(node)
}

func (p *Parser) parseParameter() *Node {
	node := p.startNode(KindParameter)
	node.AddChild(p.parseModifiers())

	node.AddChild(p.parseType())

	if p.check(TokenEllipsis) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
	}

	if id := p.parseVariableDeclaratorId(); id != nil {
		node.AddChild(id)
	}

	for p.check(TokenLBracket) {
		p.advance()
		p.expect(TokenRBracket)
	}

	return p.finishNode(node)
}

func (p *Parser) parseThrowsList() *Node {
	node := p.startNode(KindThrowsList)
	p.expect(TokenThrows)

	for {
		progress := p.mustProgress()
		node.AddChild(p.parseType())
		if !p.check(TokenComma) {
			break
		}
		p.advance()
		if !progress() {
			break
		}
	}

	return p.finishNode(node)
}

func (p *Parser) parseBlock() *Node {
	node := p.startNode(KindBlock)
	p.expect(TokenLBrace)

	for !p.check(TokenRBrace) && !p.check(TokenEOF) {
		node.AddChild(p.parseStatement())
	}

	p.expect(TokenRBrace)
	return p.finishNode(node)
}

func (p *Parser) parseStatement() *Node {
	switch p.peek().Kind {
	case TokenLBrace:
		return p.parseBlock()
	case TokenSemicolon:
		node := p.startNode(KindEmptyStmt)
		p.advance()
		return p.finishNode(node)
	case TokenIf:
		return p.parseIfStmt()
	case TokenFor:
		return p.parseForStmt()
	case TokenWhile:
		return p.parseWhileStmt()
	case TokenDo:
		return p.parseDoStmt()
	case TokenSwitch:
		return p.parseSwitchStmt()
	case TokenReturn:
		return p.parseReturnStmt()
	case TokenBreak:
		return p.parseBreakStmt()
	case TokenContinue:
		return p.parseContinueStmt()
	case TokenThrow:
		return p.parseThrowStmt()
	case TokenTry:
		return p.parseTryStmt()
	case TokenSynchronized:
		return p.parseSynchronizedStmt()
	case TokenAssert:
		return p.parseAssertStmt()
	case TokenYield:
		return p.parseYieldStmt()
	case TokenClass, TokenInterface, TokenEnum, TokenRecord:
		return p.parseLocalClassDecl()
	case TokenFinal, TokenAt:
		return p.parseLocalVarOrExprStmt()
	case TokenIdent:
		if p.peekN(1).Kind == TokenColon {
			return p.parseLabeledStmt()
		}
		return p.parseLocalVarOrExprStmt()
	default:
		return p.parseLocalVarOrExprStmt()
	}
}

func (p *Parser) parseLocalVarOrExprStmt() *Node {
	if p.isLocalVarDecl() {
		return p.parseLocalVarDecl()
	}
	return p.parseExprStmt()
}

func (p *Parser) isLocalVarDecl() bool {
	save := p.pos

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	if p.check(TokenFinal) {
		p.advance()
	}

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	isType := false
	switch p.peek().Kind {
	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble, TokenVar:
		isType = true
	default:
		if p.isIdentifierLike() {
			p.parseQualifiedName()
			if p.check(TokenLT) {
				p.skipTypeArguments()
			}
			for p.check(TokenLBracket) {
				p.advance()
				if !p.check(TokenRBracket) {
					p.pos = save
					return false
				}
				p.advance()
			}
			isType = p.isIdentifierLike() || p.isUnnamedVariable()
		}
	}

	p.pos = save
	return isType
}

func (p *Parser) skipTypeArguments() {
	if !p.check(TokenLT) {
		return
	}
	p.advance()
	depth := 1
	for depth > 0 && !p.check(TokenEOF) {
		switch p.peek().Kind {
		case TokenLT:
			depth++
		case TokenGT:
			depth--
		case TokenShr:
			depth -= 2
		case TokenUShr:
			depth -= 3
		}
		p.advance()
	}
}

func (p *Parser) parseLocalVarDecl() *Node {
	node := p.startNode(KindLocalVarDecl)
	node.AddChild(p.parseModifiers())

	if p.check(TokenVar) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindType, Token: &tok, Span: tok.Span})
	} else {
		node.AddChild(p.parseType())
	}

	for {
		progress := p.mustProgress()
		if id := p.parseVariableDeclaratorId(); id != nil {
			node.AddChild(id)
		}

		for p.check(TokenLBracket) {
			p.advance()
			p.expect(TokenRBracket)
		}

		if p.check(TokenAssign) {
			p.advance()
			node.AddChild(p.parseVarInitializer())
		}

		if !p.check(TokenComma) {
			break
		}
		p.advance()
		if !progress() {
			break
		}
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseExprStmt() *Node {
	node := p.startNode(KindExprStmt)
	node.AddChild(p.parseExpression())
	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseLocalClassDecl() *Node {
	node := p.startNode(KindLocalClassDecl)
	modifiers := p.parseModifiers()
	switch p.peek().Kind {
	case TokenClass:
		node.AddChild(p.parseClassDecl(modifiers))
	case TokenInterface:
		node.AddChild(p.parseInterfaceDecl(modifiers))
	case TokenEnum:
		node.AddChild(p.parseEnumDecl(modifiers))
	case TokenRecord:
		node.AddChild(p.parseRecordDecl(modifiers))
	}
	return p.finishNode(node)
}

func (p *Parser) parseIfStmt() *Node {
	node := p.startNode(KindIfStmt)
	p.expect(TokenIf)
	p.expect(TokenLParen)
	node.AddChild(p.parseExpression())
	p.expect(TokenRParen)
	node.AddChild(p.parseStatement())

	if p.check(TokenElse) {
		p.advance()
		node.AddChild(p.parseStatement())
	}

	return p.finishNode(node)
}

func (p *Parser) parseForStmt() *Node {
	p.expect(TokenFor)
	p.expect(TokenLParen)

	if p.isEnhancedFor() {
		return p.parseEnhancedForStmt()
	}

	node := p.startNode(KindForStmt)

	initNode := p.startNode(KindForInit)
	if !p.check(TokenSemicolon) {
		if p.isLocalVarDecl() {
			initNode.AddChild(p.parseLocalVarDeclNoSemi())
		} else {
			for {
				initNode.AddChild(p.parseExpression())
				if !p.check(TokenComma) {
					break
				}
				p.advance()
			}
		}
	}
	node.AddChild(p.finishNode(initNode))
	p.expect(TokenSemicolon)

	if !p.check(TokenSemicolon) {
		node.AddChild(p.parseExpression())
	}
	p.expect(TokenSemicolon)

	updateNode := p.startNode(KindForUpdate)
	if !p.check(TokenRParen) {
		for {
			updateNode.AddChild(p.parseExpression())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
		}
	}
	node.AddChild(p.finishNode(updateNode))
	p.expect(TokenRParen)

	node.AddChild(p.parseStatement())
	return p.finishNode(node)
}

func (p *Parser) isEnhancedFor() bool {
	save := p.pos

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	if p.check(TokenFinal) {
		p.advance()
	}

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	switch p.peek().Kind {
	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble, TokenVar:
		p.advance()
	case TokenIdent:
		p.parseQualifiedName()
		if p.check(TokenLT) {
			p.skipTypeArguments()
		}
	default:
		p.pos = save
		return false
	}

	for p.check(TokenLBracket) {
		p.advance()
		if p.check(TokenRBracket) {
			p.advance()
		}
	}

	if !p.check(TokenIdent) {
		p.pos = save
		return false
	}
	p.advance()

	result := p.check(TokenColon)
	p.pos = save
	return result
}

func (p *Parser) isLocalVarDeclWithUnderscore() bool {
	return p.check(TokenIdent) && p.peek().Literal == "_" && p.peekN(1).Kind == TokenAssign
}

func (p *Parser) parseEnhancedForStmt() *Node {
	node := p.startNode(KindEnhancedForStmt)

	node.AddChild(p.parseModifiers())

	if p.check(TokenVar) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindType, Token: &tok, Span: tok.Span})
	} else {
		node.AddChild(p.parseType())
	}

	if id := p.parseVariableDeclaratorId(); id != nil {
		node.AddChild(id)
	}

	p.expect(TokenColon)
	node.AddChild(p.parseExpression())
	p.expect(TokenRParen)
	node.AddChild(p.parseStatement())

	return p.finishNode(node)
}

func (p *Parser) parseLocalVarDeclNoSemi() *Node {
	node := p.startNode(KindLocalVarDecl)
	node.AddChild(p.parseModifiers())

	if p.check(TokenVar) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindType, Token: &tok, Span: tok.Span})
	} else {
		node.AddChild(p.parseType())
	}

	for {
		progress := p.mustProgress()
		if id := p.parseVariableDeclaratorId(); id != nil {
			node.AddChild(id)
		}

		for p.check(TokenLBracket) {
			p.advance()
			p.expect(TokenRBracket)
		}

		if p.check(TokenAssign) {
			p.advance()
			node.AddChild(p.parseVarInitializer())
		}

		if !p.check(TokenComma) {
			break
		}
		p.advance()
		if !progress() {
			break
		}
	}

	return p.finishNode(node)
}

func (p *Parser) parseWhileStmt() *Node {
	node := p.startNode(KindWhileStmt)
	p.expect(TokenWhile)
	p.expect(TokenLParen)
	node.AddChild(p.parseExpression())
	p.expect(TokenRParen)
	node.AddChild(p.parseStatement())
	return p.finishNode(node)
}

func (p *Parser) parseDoStmt() *Node {
	node := p.startNode(KindDoStmt)
	p.expect(TokenDo)
	node.AddChild(p.parseStatement())
	p.expect(TokenWhile)
	p.expect(TokenLParen)
	node.AddChild(p.parseExpression())
	p.expect(TokenRParen)
	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseSwitchStmt() *Node {
	node := p.startNode(KindSwitchStmt)
	p.expect(TokenSwitch)
	p.expect(TokenLParen)
	node.AddChild(p.parseExpression())
	p.expect(TokenRParen)
	p.expect(TokenLBrace)

	for !p.check(TokenRBrace) && !p.check(TokenEOF) {
		node.AddChild(p.parseSwitchCase())
	}

	p.expect(TokenRBrace)
	return p.finishNode(node)
}

func (p *Parser) parseSwitchCase() *Node {
	node := p.startNode(KindSwitchCase)

	isArrowCase := false
	for p.check(TokenCase) || p.check(TokenDefault) {
		label := p.parseSwitchLabel()
		node.AddChild(label)
		if label.isArrowCase {
			isArrowCase = true
			break
		}
	}

	if isArrowCase {
		switch p.peek().Kind {
		case TokenLBrace:
			node.AddChild(p.parseBlock())
		case TokenThrow:
			node.AddChild(p.parseThrowStmt())
		default:
			exprNode := p.startNode(KindExprStmt)
			exprNode.AddChild(p.parseExpression())
			p.expect(TokenSemicolon)
			node.AddChild(p.finishNode(exprNode))
		}
	} else {
		for !p.check(TokenCase) && !p.check(TokenDefault) && !p.check(TokenRBrace) && !p.check(TokenEOF) {
			node.AddChild(p.parseStatement())
		}
	}

	return p.finishNode(node)
}

func (p *Parser) parseSwitchLabel() *Node {
	node := p.startNode(KindSwitchLabel)

	if p.check(TokenCase) {
		p.advance()
		for {
			progress := p.mustProgress()
			if p.looksLikePattern() {
				node.AddChild(p.parsePattern())
			} else {
				node.AddChild(p.parseCaseLabelExpression())
			}
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			// Java 21: case null, default -> ...
			if p.check(TokenDefault) {
				tok := p.advance()
				node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
				break
			}
			if !progress() {
				break
			}
		}
		if p.check(TokenWhen) {
			node.AddChild(p.parseGuard())
		}
	} else {
		p.expect(TokenDefault)
	}

	if p.check(TokenArrow) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.isArrowCase = true
	} else {
		p.expect(TokenColon)
	}

	return p.finishNode(node)
}

func (p *Parser) looksLikePattern() bool {
	if p.looksLikeMatchAllPattern() {
		return true
	}

	save := p.pos
	defer func() { p.pos = save }()

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	switch p.peek().Kind {
	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble:
		p.advance()
	case TokenIdent:
		p.parseQualifiedName()
		if p.check(TokenLT) {
			p.parseTypeArguments()
		}
	default:
		return false
	}

	for p.check(TokenLBracket) {
		p.advance()
		if !p.check(TokenRBracket) {
			return false
		}
		p.advance()
	}

	// TypePattern: Type identifier
	// RecordPattern: Type ( ... )
	return p.check(TokenIdent) || p.check(TokenLParen)
}

func (p *Parser) parsePattern() *Node {
	if p.looksLikeMatchAllPattern() {
		return p.parseMatchAllPattern()
	}

	// Parse the type first, then decide based on what follows
	typeNode := p.parseType()

	if p.check(TokenLParen) {
		// RecordPattern: Type ( ComponentPatternList )
		node := p.startNode(KindRecordPattern)
		node.AddChild(typeNode)
		p.advance() // consume (
		if !p.check(TokenRParen) {
			for {
				progress := p.mustProgress()
				node.AddChild(p.parsePattern())
				if !p.check(TokenComma) {
					break
				}
				p.advance()
				if !progress() {
					break
				}
			}
		}
		p.expect(TokenRParen)
		return p.finishNode(node)
	}

	// TypePattern: Type identifier
	node := p.startNode(KindTypePattern)
	node.AddChild(typeNode)
	if p.check(TokenIdent) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
	}
	return p.finishNode(node)
}

func (p *Parser) parseGuard() *Node {
	node := p.startNode(KindGuard)
	p.expect(TokenWhen)
	node.AddChild(p.parseExpression())
	return p.finishNode(node)
}

func (p *Parser) looksLikeMatchAllPattern() bool {
	if !p.check(TokenIdent) || p.peek().Literal != "_" {
		return false
	}
	next := p.peekN(1).Kind
	return next == TokenColon || next == TokenArrow || next == TokenComma || next == TokenRParen
}

func (p *Parser) parseMatchAllPattern() *Node {
	node := p.startNode(KindMatchAllPattern)
	p.advance() // consume _
	return p.finishNode(node)
}

func (p *Parser) isUnnamedVariable() bool {
	return p.check(TokenIdent) && p.peek().Literal == "_"
}

func (p *Parser) parseVariableDeclaratorId() *Node {
	if p.isUnnamedVariable() {
		node := p.startNode(KindUnnamedVariable)
		p.advance()
		return p.finishNode(node)
	}
	if p.isIdentifierLike() {
		tok := p.advance()
		return &Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span}
	}
	return nil
}

func (p *Parser) parseReturnStmt() *Node {
	node := p.startNode(KindReturnStmt)
	p.expect(TokenReturn)

	if !p.check(TokenSemicolon) {
		node.AddChild(p.parseExpression())
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseBreakStmt() *Node {
	node := p.startNode(KindBreakStmt)
	p.expect(TokenBreak)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseContinueStmt() *Node {
	node := p.startNode(KindContinueStmt)
	p.expect(TokenContinue)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseThrowStmt() *Node {
	node := p.startNode(KindThrowStmt)
	p.expect(TokenThrow)
	node.AddChild(p.parseExpression())
	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseTryStmt() *Node {
	node := p.startNode(KindTryStmt)
	p.expect(TokenTry)

	if p.check(TokenLParen) {
		p.advance()
		for !p.check(TokenRParen) && !p.check(TokenEOF) {
			node.AddChild(p.parseResource())
			if p.check(TokenSemicolon) {
				p.advance()
			}
			if p.check(TokenRParen) {
				break
			}
		}
		p.expect(TokenRParen)
	}

	node.AddChild(p.parseBlock())

	for p.check(TokenCatch) {
		node.AddChild(p.parseCatchClause())
	}

	if p.check(TokenFinally) {
		node.AddChild(p.parseFinallyClause())
	}

	return p.finishNode(node)
}

func (p *Parser) parseResource() *Node {
	if p.isLocalVarDecl() {
		node := p.startNode(KindLocalVarDecl)
		node.AddChild(p.parseModifiers())
		node.AddChild(p.parseType())
		if id := p.parseVariableDeclaratorId(); id != nil {
			node.AddChild(id)
		}
		if p.check(TokenAssign) {
			p.advance()
			node.AddChild(p.parseExpression())
		}
		return p.finishNode(node)
	}
	return p.parseExpression()
}

func (p *Parser) parseCatchClause() *Node {
	node := p.startNode(KindCatchClause)
	p.expect(TokenCatch)
	p.expect(TokenLParen)

	node.AddChild(p.parseModifiers())

	typeNode := p.startNode(KindType)
	typeNode.AddChild(p.parseType())
	for p.check(TokenBitOr) {
		p.advance()
		typeNode.AddChild(p.parseType())
	}
	node.AddChild(p.finishNode(typeNode))

	if id := p.parseVariableDeclaratorId(); id != nil {
		node.AddChild(id)
	}

	p.expect(TokenRParen)
	node.AddChild(p.parseBlock())

	return p.finishNode(node)
}

func (p *Parser) parseFinallyClause() *Node {
	node := p.startNode(KindFinallyClause)
	p.expect(TokenFinally)
	node.AddChild(p.parseBlock())
	return p.finishNode(node)
}

func (p *Parser) parseSynchronizedStmt() *Node {
	node := p.startNode(KindSynchronizedStmt)
	p.expect(TokenSynchronized)
	p.expect(TokenLParen)
	node.AddChild(p.parseExpression())
	p.expect(TokenRParen)
	node.AddChild(p.parseBlock())
	return p.finishNode(node)
}

func (p *Parser) parseAssertStmt() *Node {
	node := p.startNode(KindAssertStmt)
	p.expect(TokenAssert)
	node.AddChild(p.parseExpression())

	if p.check(TokenColon) {
		p.advance()
		node.AddChild(p.parseExpression())
	}

	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseYieldStmt() *Node {
	node := p.startNode(KindYieldStmt)
	p.expect(TokenYield)
	node.AddChild(p.parseExpression())
	p.expect(TokenSemicolon)
	return p.finishNode(node)
}

func (p *Parser) parseLabeledStmt() *Node {
	node := p.startNode(KindLabeledStmt)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}
	p.expect(TokenColon)
	node.AddChild(p.parseStatement())

	return p.finishNode(node)
}

func (p *Parser) parseExpression() *Node {
	return p.parseAssignmentExpr()
}

func (p *Parser) parseCaseLabelExpression() *Node {
	return p.parseTernaryExpr()
}

func (p *Parser) parseAssignmentExpr() *Node {
	if p.isLambda() {
		return p.parseLambdaExpr()
	}

	left := p.parseTernaryExpr()

	if p.isAssignOp() {
		node := p.startNode(KindAssignExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseAssignmentExpr())
		return p.finishNode(node)
	}

	return left
}

func (p *Parser) isAssignOp() bool {
	switch p.peek().Kind {
	case TokenAssign, TokenPlusAssign, TokenMinusAssign,
		TokenStarAssign, TokenSlashAssign, TokenPercentAssign,
		TokenAndAssign, TokenOrAssign, TokenXorAssign,
		TokenShlAssign, TokenShrAssign, TokenUShrAssign:
		return true
	}
	return false
}

func (p *Parser) isLambda() bool {
	if p.check(TokenIdent) && p.peekN(1).Kind == TokenArrow {
		return true
	}

	if !p.check(TokenLParen) {
		return false
	}

	save := p.pos
	p.advance()
	depth := 1

	for depth > 0 && !p.check(TokenEOF) {
		switch p.peek().Kind {
		case TokenLParen:
			depth++
		case TokenRParen:
			depth--
		}
		if depth > 0 {
			p.advance()
		}
	}

	if p.check(TokenRParen) {
		p.advance()
	}

	result := p.check(TokenArrow)
	p.pos = save
	return result
}

func (p *Parser) parseLambdaExpr() *Node {
	node := p.startNode(KindLambdaExpr)

	if p.check(TokenIdent) {
		tok := p.advance()
		paramNode := p.startNode(KindParameters)
		paramNode.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.finishNode(paramNode))
	} else {
		node.AddChild(p.parseLambdaParameters())
	}

	p.expect(TokenArrow)

	if p.check(TokenLBrace) {
		node.AddChild(p.parseBlock())
	} else {
		node.AddChild(p.parseExpression())
	}

	return p.finishNode(node)
}

func (p *Parser) parseLambdaParameters() *Node {
	node := p.startNode(KindParameters)
	p.expect(TokenLParen)

	if !p.check(TokenRParen) {
		for {
			progress := p.mustProgress()
			if p.isLambdaTypedParam() {
				node.AddChild(p.parseParameter())
			} else {
				if tok := p.expect(TokenIdent); tok != nil {
					node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
				}
			}
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			if !progress() {
				break
			}
		}
	}

	p.expect(TokenRParen)
	return p.finishNode(node)
}

func (p *Parser) isLambdaTypedParam() bool {
	switch p.peek().Kind {
	case TokenFinal, TokenAt:
		return true
	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble, TokenVar:
		return true
	case TokenIdent:
		return p.peekN(1).Kind == TokenIdent || p.peekN(1).Kind == TokenLT ||
			p.peekN(1).Kind == TokenDot || p.peekN(1).Kind == TokenLBracket
	}
	return false
}

func (p *Parser) parseTernaryExpr() *Node {
	cond := p.parseOrExpr()

	if p.check(TokenQuestion) {
		node := p.startNode(KindTernaryExpr)
		node.AddChild(cond)
		p.advance()
		node.AddChild(p.parseExpression())
		p.expect(TokenColon)
		if p.isLambda() {
			node.AddChild(p.parseLambdaExpr())
		} else {
			node.AddChild(p.parseTernaryExpr())
		}
		return p.finishNode(node)
	}

	return cond
}

func (p *Parser) parseOrExpr() *Node {
	left := p.parseAndExpr()

	for p.check(TokenOr) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseAndExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseAndExpr() *Node {
	left := p.parseBitOrExpr()

	for p.check(TokenAnd) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseBitOrExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseBitOrExpr() *Node {
	left := p.parseBitXorExpr()

	for p.check(TokenBitOr) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseBitXorExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseBitXorExpr() *Node {
	left := p.parseBitAndExpr()

	for p.check(TokenBitXor) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseBitAndExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseBitAndExpr() *Node {
	left := p.parseEqualityExpr()

	for p.check(TokenBitAnd) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseEqualityExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseEqualityExpr() *Node {
	left := p.parseRelationalExpr()

	for p.check(TokenEQ) || p.check(TokenNE) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseRelationalExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseRelationalExpr() *Node {
	left := p.parseShiftExpr()

	for {
		if p.check(TokenLT) || p.check(TokenLE) || p.check(TokenGT) || p.check(TokenGE) {
			node := p.startNode(KindBinaryExpr)
			node.AddChild(left)
			tok := p.advance()
			node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
			node.AddChild(p.parseShiftExpr())
			left = p.finishNode(node)
		} else if p.check(TokenInstanceof) {
			node := p.startNode(KindInstanceofExpr)
			node.AddChild(left)
			p.advance()
			node.AddChild(p.parseType())
			if p.check(TokenIdent) {
				tok := p.advance()
				node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
			}
			left = p.finishNode(node)
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parseShiftExpr() *Node {
	left := p.parseAdditiveExpr()

	for p.check(TokenShl) || p.check(TokenShr) || p.check(TokenUShr) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseAdditiveExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseAdditiveExpr() *Node {
	left := p.parseMultiplicativeExpr()

	for p.check(TokenPlus) || p.check(TokenMinus) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseMultiplicativeExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseMultiplicativeExpr() *Node {
	left := p.parseUnaryExpr()

	for p.check(TokenStar) || p.check(TokenSlash) || p.check(TokenPercent) {
		node := p.startNode(KindBinaryExpr)
		node.AddChild(left)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseUnaryExpr())
		left = p.finishNode(node)
	}

	return left
}

func (p *Parser) parseUnaryExpr() *Node {
	switch p.peek().Kind {
	case TokenIncrement, TokenDecrement:
		node := p.startNode(KindUnaryExpr)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseUnaryExpr())
		return p.finishNode(node)
	case TokenPlus, TokenMinus, TokenNot, TokenBitNot:
		node := p.startNode(KindUnaryExpr)
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		node.AddChild(p.parseUnaryExpr())
		return p.finishNode(node)
	case TokenLParen:
		if p.isCast() {
			return p.parseCastExpr()
		}
	}

	return p.parsePostfixExpr()
}

func (p *Parser) isCast() bool {
	if !p.check(TokenLParen) {
		return false
	}

	save := p.pos
	p.advance()

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	isType := false
	switch p.peek().Kind {
	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble:
		isType = true
	case TokenIdent:
		p.parseQualifiedName()
		if p.check(TokenLT) {
			p.skipTypeArguments()
		}
		for p.check(TokenLBracket) {
			p.advance()
			if p.check(TokenRBracket) {
				p.advance()
			}
		}
		// Handle intersection types: (Type & Type2)
		for p.check(TokenBitAnd) {
			p.advance()
			p.parseQualifiedName()
			if p.check(TokenLT) {
				p.skipTypeArguments()
			}
		}
		isType = p.check(TokenRParen)
		if isType {
			p.advance()
			switch p.peek().Kind {
			case TokenIdent, TokenThis, TokenSuper, TokenNew,
				TokenLParen, TokenNot, TokenBitNot,
				TokenIncrement, TokenDecrement,
				TokenIntLiteral, TokenFloatLiteral,
				TokenCharLiteral, TokenStringLiteral,
				TokenTextBlock, TokenTrue, TokenFalse, TokenNull:
			default:
				isType = false
			}
		}
	}

	p.pos = save
	return isType
}

func (p *Parser) parseCastExpr() *Node {
	node := p.startNode(KindCastExpr)
	p.expect(TokenLParen)

	typeNode := p.startNode(KindType)
	typeNode.AddChild(p.parseType())
	for p.check(TokenBitAnd) {
		p.advance()
		typeNode.AddChild(p.parseType())
	}
	node.AddChild(p.finishNode(typeNode))

	p.expect(TokenRParen)
	// Handle cast to lambda: (Supplier) () -> value
	if p.isLambda() {
		node.AddChild(p.parseLambdaExpr())
	} else {
		node.AddChild(p.parseUnaryExpr())
	}
	return p.finishNode(node)
}

func (p *Parser) parsePostfixExpr() *Node {
	expr := p.parsePrimaryExpr()
	return p.parsePostfixSuffix(expr)
}

func (p *Parser) parsePostfixSuffix(expr *Node) *Node {
	for {
		progress := p.mustProgress()
		switch p.peek().Kind {
		case TokenIncrement, TokenDecrement:
			node := p.startNode(KindPostfixExpr)
			node.AddChild(expr)
			tok := p.advance()
			node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
			expr = p.finishNode(node)
		case TokenDot:
			p.advance()
			if p.check(TokenNew) {
				expr = p.parseInnerNewExpr(expr)
			} else if p.match(TokenStringTemplate, TokenTextBlockTemplate) {
				node := p.startNode(KindTemplateExpr)
				node.AddChild(expr)
				tok := p.advance()
				node.AddChild(&Node{Kind: KindLiteral, Token: &tok, Span: tok.Span})
				expr = p.finishNode(node)
			} else if p.match(TokenStringLiteral, TokenTextBlock) {
				node := p.startNode(KindFieldAccess)
				node.AddChild(expr)
				tok := p.advance()
				node.AddChild(&Node{Kind: KindLiteral, Token: &tok, Span: tok.Span})
				expr = p.finishNode(node)
			} else if p.check(TokenLT) {
				typeArgs := p.parseTypeArguments()
				if p.isIdentifierLike() {
					tok := p.advance()
					node := p.startNode(KindFieldAccess)
					node.AddChild(expr)
					node.AddChild(typeArgs)
					node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
					expr = p.finishNode(node)
					if p.check(TokenLParen) {
						expr = p.parseMethodCall(expr)
					}
				}
			} else if p.check(TokenClass) {
				node := p.startNode(KindClassLiteral)
				node.AddChild(expr)
				p.advance()
				expr = p.finishNode(node)
			} else if p.check(TokenThis) {
				node := p.startNode(KindFieldAccess)
				node.AddChild(expr)
				tok := p.advance()
				node.AddChild(&Node{Kind: KindThis, Token: &tok, Span: tok.Span})
				expr = p.finishNode(node)
			} else if p.check(TokenSuper) {
				node := p.startNode(KindFieldAccess)
				node.AddChild(expr)
				tok := p.advance()
				node.AddChild(&Node{Kind: KindSuper, Token: &tok, Span: tok.Span})
				expr = p.finishNode(node)
			} else if p.isIdentifierLike() {
				tok := p.advance()
				node := p.startNode(KindFieldAccess)
				node.AddChild(expr)
				node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
				expr = p.finishNode(node)
				if p.check(TokenLParen) {
					expr = p.parseMethodCall(expr)
				}
			}
		case TokenLBracket:
			// Check if this is an array type class literal like String[].class
			// or an array type method reference like String[]::new
			if p.peekN(1).Kind == TokenRBracket {
				if result := p.tryParseArrayClassLiteralOrMethodRef(expr); result != nil {
					expr = result
					continue
				}
			}
			p.advance()
			node := p.startNode(KindArrayAccess)
			node.AddChild(expr)
			node.AddChild(p.parseExpression())
			p.expect(TokenRBracket)
			expr = p.finishNode(node)
		case TokenLParen:
			expr = p.parseMethodCall(expr)
		case TokenColonColon:
			expr = p.parseMethodRef(expr)
		case TokenLT:
			// Try to parse as parameterized type for Class<?>[]::new or Class<?>.class patterns
			if result := p.tryParseParameterizedTypeSpecialForm(expr); result != nil {
				expr = result
				continue
			}
			return expr
		default:
			return expr
		}
		if !progress() {
			return expr
		}
	}
}

func (p *Parser) parseMethodCall(target *Node) *Node {
	node := p.startNode(KindCallExpr)
	node.AddChild(target)
	node.AddChild(p.parseArguments())
	return p.finishNode(node)
}

func (p *Parser) parseArguments() *Node {
	node := p.startNode(KindParameters)
	p.expect(TokenLParen)

	if !p.check(TokenRParen) {
		for {
			progress := p.mustProgress()
			node.AddChild(p.parseExpression())
			if !p.check(TokenComma) {
				break
			}
			p.advance()
			if !progress() {
				break
			}
		}
	}

	p.expect(TokenRParen)
	return p.finishNode(node)
}

func (p *Parser) parseMethodRef(target *Node) *Node {
	node := p.startNode(KindMethodRef)
	node.AddChild(target)
	p.expect(TokenColonColon)

	if p.check(TokenLT) {
		node.AddChild(p.parseTypeArguments())
	}

	if p.check(TokenNew) {
		tok := p.advance()
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
	} else if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	return p.finishNode(node)
}

func (p *Parser) parsePrimaryExpr() *Node {
	switch p.peek().Kind {
	case TokenIntLiteral, TokenFloatLiteral, TokenCharLiteral,
		TokenStringLiteral, TokenTextBlock, TokenTrue, TokenFalse, TokenNull:
		tok := p.advance()
		return &Node{Kind: KindLiteral, Token: &tok, Span: tok.Span}

	case TokenThis:
		tok := p.advance()
		return &Node{Kind: KindThis, Token: &tok, Span: tok.Span}

	case TokenSuper:
		tok := p.advance()
		node := &Node{Kind: KindSuper, Token: &tok, Span: tok.Span}
		if p.check(TokenDot) || p.check(TokenLParen) {
			return p.parsePostfixSuffix(node)
		}
		return node

	case TokenNew:
		return p.parseNewExpr()

	case TokenLParen:
		return p.parseParenExpr()

	case TokenSwitch:
		return p.parseSwitchExpr()

	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble, TokenVoid:
		return p.parsePrimitiveClassLiteral()

	default:
		if p.isIdentifierLike() {
			tok := p.advance()
			return &Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span}
		}
		return p.errorNode("expected expression", []TokenKind{TokenSemicolon, TokenComma, TokenRParen, TokenRBrace, TokenRBracket})
	}
}

func (p *Parser) parseParenExpr() *Node {
	node := p.startNode(KindParenExpr)
	p.expect(TokenLParen)
	node.AddChild(p.parseExpression())
	p.expect(TokenRParen)
	return p.finishNode(node)
}

func (p *Parser) parseNewExpr() *Node {
	p.expect(TokenNew)

	if p.check(TokenLT) {
		p.parseTypeArguments()
	}

	for p.check(TokenAt) {
		p.parseAnnotation()
	}

	switch p.peek().Kind {
	case TokenBoolean, TokenByte, TokenChar, TokenShort,
		TokenInt, TokenLong, TokenFloat, TokenDouble:
		return p.parseNewArrayExpr()
	}

	qualName := p.parseQualifiedName()

	if p.check(TokenLT) {
		p.parseTypeArguments()
	}

	if p.check(TokenAt) || p.check(TokenLBracket) {
		node := p.startNode(KindNewArrayExpr)
		node.AddChild(qualName)
		for p.check(TokenAt) || p.check(TokenLBracket) {
			progress := p.mustProgress()
			for p.check(TokenAt) {
				node.AddChild(p.parseAnnotation())
			}
			if !p.check(TokenLBracket) {
				break
			}
			p.advance()
			if !p.check(TokenRBracket) {
				node.AddChild(p.parseExpression())
			}
			p.expect(TokenRBracket)
			if !progress() {
				break
			}
		}
		if p.check(TokenLBrace) {
			node.AddChild(p.parseArrayInitializer())
		}
		return p.finishNode(node)
	}

	node := p.startNode(KindNewExpr)
	node.AddChild(qualName)
	node.AddChild(p.parseArguments())

	if p.check(TokenLBrace) {
		node.AddChild(p.parseClassBody())
	}

	return p.finishNode(node)
}

func (p *Parser) parseNewArrayExpr() *Node {
	node := p.startNode(KindNewArrayExpr)
	tok := p.advance()
	node.AddChild(&Node{Kind: KindType, Token: &tok, Span: tok.Span})

	for p.check(TokenAt) || p.check(TokenLBracket) {
		progress := p.mustProgress()
		for p.check(TokenAt) {
			node.AddChild(p.parseAnnotation())
		}
		if !p.check(TokenLBracket) {
			break
		}
		p.advance()
		if !p.check(TokenRBracket) {
			node.AddChild(p.parseExpression())
		}
		p.expect(TokenRBracket)
		if !progress() {
			break
		}
	}

	if p.check(TokenLBrace) {
		node.AddChild(p.parseArrayInitializer())
	}

	return p.finishNode(node)
}

func (p *Parser) parseInnerNewExpr(outer *Node) *Node {
	p.expect(TokenNew)

	if p.check(TokenLT) {
		p.parseTypeArguments()
	}

	node := p.startNode(KindNewExpr)
	node.AddChild(outer)

	if tok := p.expect(TokenIdent); tok != nil {
		node.AddChild(&Node{Kind: KindIdentifier, Token: tok, Span: tok.Span})
	}

	if p.check(TokenLT) {
		node.AddChild(p.parseTypeArguments())
	}

	node.AddChild(p.parseArguments())

	if p.check(TokenLBrace) {
		node.AddChild(p.parseClassBody())
	}

	return p.finishNode(node)
}

func (p *Parser) parsePrimitiveClassLiteral() *Node {
	node := p.startNode(KindClassLiteral)
	tok := p.advance()
	typeNode := &Node{Kind: KindType, Token: &tok, Span: tok.Span}

	for p.check(TokenLBracket) {
		p.advance()
		p.expect(TokenRBracket)
		wrapper := p.startNode(KindArrayType)
		wrapper.AddChild(typeNode)
		typeNode = p.finishNode(wrapper)
	}

	node.AddChild(typeNode)
	p.expect(TokenDot)
	p.expect(TokenClass)
	return p.finishNode(node)
}

// tryParseArrayClassLiteralOrMethodRef attempts to parse an array type class literal like String[].class
// or an array type method reference like String[]::new.
// If successful, returns the ClassLiteral or MethodRef node. Otherwise returns nil (parser position unchanged).
func (p *Parser) tryParseArrayClassLiteralOrMethodRef(baseExpr *Node) *Node {
	save := p.pos

	// Count consecutive [] pairs
	dims := 0
	for p.check(TokenLBracket) && p.peekN(1).Kind == TokenRBracket {
		p.advance() // [
		p.advance() // ]
		dims++
	}

	if dims == 0 {
		p.pos = save
		return nil
	}

	// Build the array type node wrapping the base expression
	buildArrayType := func() *Node {
		typeNode := baseExpr
		for i := 0; i < dims; i++ {
			wrapper := p.startNode(KindArrayType)
			wrapper.AddChild(typeNode)
			typeNode = p.finishNode(wrapper)
		}
		return typeNode
	}

	// Check if .class follows
	if p.check(TokenDot) && p.peekN(1).Kind == TokenClass {
		p.advance() // .
		p.advance() // class

		node := p.startNode(KindClassLiteral)
		node.AddChild(buildArrayType())
		return p.finishNode(node)
	}

	// Check if ::new follows (array type method reference)
	if p.check(TokenColonColon) && p.peekN(1).Kind == TokenNew {
		p.advance()        // ::
		tok := p.advance() // new

		node := p.startNode(KindMethodRef)
		node.AddChild(buildArrayType())
		node.AddChild(&Node{Kind: KindIdentifier, Token: &tok, Span: tok.Span})
		return p.finishNode(node)
	}

	// Not an array class literal or method ref, restore position
	p.pos = save
	return nil
}

// tryParseParameterizedTypeSpecialForm attempts to parse parameterized type patterns like:
// - Class<?>[]::new (array type method reference with generic element type)
// - Class<?>.class (parameterized type class literal)
// If successful, returns the result node. Otherwise returns nil (parser position unchanged).
func (p *Parser) tryParseParameterizedTypeSpecialForm(baseExpr *Node) *Node {
	save := p.pos

	// Parse type arguments
	if !p.check(TokenLT) {
		return nil
	}
	typeArgs := p.parseTypeArguments()

	// Build parameterized type node
	paramType := p.startNode(KindType)
	paramType.AddChild(baseExpr)
	paramType.AddChild(typeArgs)
	paramType = p.finishNode(paramType)

	// Check for []::new or [].class pattern
	if p.check(TokenLBracket) && p.peekN(1).Kind == TokenRBracket {
		if result := p.tryParseArrayClassLiteralOrMethodRef(paramType); result != nil {
			return result
		}
	}

	// Check for .class pattern
	if p.check(TokenDot) && p.peekN(1).Kind == TokenClass {
		p.advance() // .
		p.advance() // class

		node := p.startNode(KindClassLiteral)
		node.AddChild(paramType)
		return p.finishNode(node)
	}

	// Not a special form, restore position
	p.pos = save
	return nil
}

func (p *Parser) parseSwitchExpr() *Node {
	node := p.startNode(KindSwitchExpr)
	p.expect(TokenSwitch)
	p.expect(TokenLParen)
	node.AddChild(p.parseExpression())
	p.expect(TokenRParen)
	p.expect(TokenLBrace)

	for !p.check(TokenRBrace) && !p.check(TokenEOF) {
		node.AddChild(p.parseSwitchCase())
	}

	p.expect(TokenRBrace)
	return p.finishNode(node)
}
