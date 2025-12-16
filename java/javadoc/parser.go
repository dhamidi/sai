package javadoc

import (
	"strings"
	"unicode"
)

// Parser is a recursive-descent parser for Javadoc comments.
type Parser struct {
	input []rune
	pos   int
	len   int
}

// Parse parses a Javadoc comment string and returns a DocComment AST.
func Parse(javadoc string) *DocComment {
	p := &Parser{
		input: []rune(javadoc),
	}
	p.len = len(p.input)
	return p.parseDocComment()
}

func (p *Parser) parseDocComment() *DocComment {
	p.skipCommentStart()

	doc := &DocComment{}
	doc.Body = p.parseBody()
	doc.BlockTags = p.parseBlockTags()

	return doc
}

// skipCommentStart skips the leading /** and any whitespace/asterisks.
func (p *Parser) skipCommentStart() {
	p.skipWhitespace()
	if p.match("/**") {
		p.advance(3)
	}
	p.skipLinePrefix()
}

// skipLinePrefix skips leading whitespace and a single asterisk at the start of a line.
func (p *Parser) skipLinePrefix() {
	p.skipHorizontalWhitespace()
	if p.peek() == '*' && p.peekAt(1) != '/' {
		p.advance(1)
		if p.peek() == ' ' {
			p.advance(1)
		}
	}
}

// parseBody parses the main description until a block tag or end of comment.
func (p *Parser) parseBody() []Node {
	return p.parseContent(false)
}

// parseContent parses rich text content (text, HTML, inline tags).
// If inInlineTag is true, parsing stops at an unmatched '}'.
func (p *Parser) parseContent(inInlineTag bool) []Node {
	var nodes []Node
	var textBuf strings.Builder
	depth := 0

	flushText := func() {
		if textBuf.Len() > 0 {
			nodes = append(nodes, Text{Content: textBuf.String()})
			textBuf.Reset()
		}
	}

	for p.pos < p.len {
		ch := p.peek()

		// Check for end of comment
		if ch == '*' && p.peekAt(1) == '/' {
			break
		}

		// Check for block tag at start of line
		if !inInlineTag && p.isAtBlockTag() {
			break
		}

		switch ch {
		case '\n', '\r':
			textBuf.WriteRune(ch)
			p.advance(1)
			if ch == '\r' && p.peek() == '\n' {
				textBuf.WriteRune('\n')
				p.advance(1)
			}
			p.skipLinePrefix()

		case '{':
			if p.peekAt(1) == '@' {
				flushText()
				node := p.parseInlineTag()
				if node != nil {
					nodes = append(nodes, node)
				}
			} else {
				if inInlineTag {
					depth++
				}
				textBuf.WriteRune(ch)
				p.advance(1)
			}

		case '}':
			if inInlineTag {
				if depth == 0 {
					flushText()
					return nodes
				}
				depth--
			}
			textBuf.WriteRune(ch)
			p.advance(1)

		case '<':
			flushText()
			node := p.parseHTML()
			if node != nil {
				nodes = append(nodes, node)
			}

		case '&':
			flushText()
			node := p.parseEntity()
			if node != nil {
				nodes = append(nodes, node)
			}

		default:
			textBuf.WriteRune(ch)
			p.advance(1)
		}
	}

	flushText()
	return nodes
}

// isAtBlockTag checks if we're at the start of a block tag (@ at start of line).
func (p *Parser) isAtBlockTag() bool {
	// Look back to see if we're at the start of a line
	if p.pos == 0 {
		return p.peek() == '@'
	}

	// Check if previous non-whitespace was a newline (accounting for * prefix)
	i := p.pos - 1
	for i >= 0 {
		ch := p.input[i]
		if ch == '\n' || ch == '\r' {
			return p.peek() == '@'
		}
		if ch == '*' {
			// Check if this is the line prefix asterisk
			j := i - 1
			for j >= 0 && (p.input[j] == ' ' || p.input[j] == '\t') {
				j--
			}
			if j < 0 || p.input[j] == '\n' || p.input[j] == '\r' {
				return p.peek() == '@'
			}
		}
		if ch != ' ' && ch != '\t' {
			return false
		}
		i--
	}
	return p.peek() == '@'
}

// parseInlineTag parses an inline tag like {@code ...} or {@link ...}.
func (p *Parser) parseInlineTag() Node {
	if !p.match("{@") {
		return nil
	}
	p.advance(2)

	tagName := p.readTagName()
	if tagName == "" {
		return Erroneous{Content: "{@", Message: "missing tag name"}
	}

	p.skipHorizontalWhitespace()

	var node Node
	switch tagName {
	case "code":
		node = p.parseCodeTag()
	case "literal":
		node = p.parseLiteralTag()
	case "link":
		node = p.parseLinkTag(false)
	case "linkplain":
		node = p.parseLinkTag(true)
	case "value":
		node = p.parseValueTag()
	case "docRoot":
		node = p.parseDocRootTag()
	case "inheritDoc":
		node = p.parseInheritDocTag()
	case "index":
		node = p.parseIndexTag()
	case "summary":
		node = p.parseSummaryTag()
	case "return":
		node = p.parseInlineReturnTag()
	case "systemProperty":
		node = p.parseSystemPropertyTag()
	case "snippet":
		node = p.parseSnippetTag()
	default:
		node = p.parseUnknownInlineTag(tagName)
	}

	// Consume closing brace if present
	if p.peek() == '}' {
		p.advance(1)
	}

	return node
}

// parseCodeTag parses the content of a {@code ...} tag.
func (p *Parser) parseCodeTag() Node {
	content := p.readBalancedContent()
	return Code{Content: content}
}

// parseLiteralTag parses the content of a {@literal ...} tag.
func (p *Parser) parseLiteralTag() Node {
	content := p.readBalancedContent()
	return Literal{Content: content}
}

// parseLinkTag parses the content of a {@link ...} or {@linkplain ...} tag.
func (p *Parser) parseLinkTag(plain bool) Node {
	ref := p.readReference()
	p.skipHorizontalWhitespace()

	var label []Node
	if p.peek() != '}' {
		label = p.parseContent(true)
	}

	return Link{Reference: ref, Label: label, Plain: plain}
}

// parseValueTag parses the content of a {@value ...} tag.
func (p *Parser) parseValueTag() Node {
	ref := p.readReference()
	return Value{Reference: ref}
}

// parseDocRootTag parses {@docRoot}.
func (p *Parser) parseDocRootTag() Node {
	return DocRoot{}
}

// parseInheritDocTag parses {@inheritDoc} or {@inheritDoc reference}.
func (p *Parser) parseInheritDocTag() Node {
	ref := p.readReference()
	return InheritDoc{Reference: ref}
}

// parseIndexTag parses {@index term description}.
func (p *Parser) parseIndexTag() Node {
	var term string
	if p.peek() == '"' {
		term = p.readQuotedString()
	} else {
		term = p.readWord()
	}

	p.skipHorizontalWhitespace()

	var desc []Node
	if p.peek() != '}' {
		desc = p.parseContent(true)
	}

	return Index{Term: term, Description: desc}
}

// parseSummaryTag parses {@summary ...}.
func (p *Parser) parseSummaryTag() Node {
	content := p.parseContent(true)
	return Summary{Content: content}
}

// parseInlineReturnTag parses {@return ...}.
func (p *Parser) parseInlineReturnTag() Node {
	content := p.parseContent(true)
	return Return{Description: content, Inline: true}
}

// parseSystemPropertyTag parses {@systemProperty name}.
func (p *Parser) parseSystemPropertyTag() Node {
	name := p.readWord()
	return SystemProperty{Name: name}
}

// parseSnippetTag parses {@snippet attributes : body}.
func (p *Parser) parseSnippetTag() Node {
	attrs := make(map[string]string)

	// Parse attributes until ':' or '}'
	for p.pos < p.len && p.peek() != ':' && p.peek() != '}' {
		p.skipHorizontalWhitespace()
		if !isJavaIdentifierStart(p.peek()) {
			break
		}

		name := p.readIdentifier()
		p.skipHorizontalWhitespace()

		var value string
		if p.peek() == '=' {
			p.advance(1)
			p.skipHorizontalWhitespace()
			if p.peek() == '"' || p.peek() == '\'' {
				value = p.readQuotedString()
			} else {
				value = p.readUnquotedAttrValue()
			}
		}
		attrs[name] = value
		p.skipHorizontalWhitespace()
	}

	var body string
	if p.peek() == ':' {
		p.advance(1)
		// Skip to end of line
		for p.pos < p.len && p.peek() != '\n' && p.peek() != '\r' {
			p.advance(1)
		}
		if p.peek() == '\r' {
			p.advance(1)
		}
		if p.peek() == '\n' {
			p.advance(1)
		}
		body = p.readBalancedContent()
	}

	return Snippet{Attributes: attrs, Body: body}
}

// parseUnknownInlineTag parses an unknown inline tag.
func (p *Parser) parseUnknownInlineTag(name string) Node {
	content := p.readBalancedContent()
	return UnknownInlineTag{Name: name, Content: content}
}

// parseHTML parses an HTML element or comment.
func (p *Parser) parseHTML() Node {
	if !p.match("<") {
		return nil
	}

	// Check for HTML comment
	if p.match("<!--") {
		return p.parseHTMLComment()
	}

	p.advance(1)

	// Check for end tag
	if p.peek() == '/' {
		p.advance(1)
		name := p.readHTMLTagName()
		p.skipHorizontalWhitespace()
		if p.peek() == '>' {
			p.advance(1)
		}
		return EndElement{Name: name}
	}

	// Start tag
	name := p.readHTMLTagName()
	if name == "" {
		return Text{Content: "<"}
	}

	attrs := p.parseHTMLAttributes()

	selfClose := false
	p.skipHorizontalWhitespace()
	if p.peek() == '/' {
		selfClose = true
		p.advance(1)
	}
	if p.peek() == '>' {
		p.advance(1)
	}

	return StartElement{Name: name, Attributes: attrs, SelfClose: selfClose}
}

// parseHTMLComment parses an HTML comment <!-- ... -->.
func (p *Parser) parseHTMLComment() Node {
	p.advance(4) // skip <!--
	start := p.pos

	for p.pos < p.len {
		if p.match("-->") {
			content := string(p.input[start:p.pos])
			p.advance(3)
			return Text{Content: "<!--" + content + "-->"}
		}
		p.advance(1)
	}

	return Text{Content: "<!--" + string(p.input[start:])}
}

// parseHTMLAttributes parses HTML attributes.
func (p *Parser) parseHTMLAttributes() []Attribute {
	var attrs []Attribute

	for {
		p.skipWhitespaceInTag()
		if p.peek() == '>' || p.peek() == '/' || p.pos >= p.len {
			break
		}

		name := p.readHTMLAttrName()
		if name == "" {
			break
		}

		p.skipWhitespaceInTag()

		var value string
		if p.peek() == '=' {
			p.advance(1)
			p.skipWhitespaceInTag()

			if p.peek() == '"' || p.peek() == '\'' {
				value = p.readQuotedString()
			} else {
				value = p.readUnquotedAttrValue()
			}
		}

		attrs = append(attrs, Attribute{Name: name, Value: value})
	}

	return attrs
}

// skipWhitespaceInTag skips whitespace including newlines and javadoc line prefixes within HTML tags.
func (p *Parser) skipWhitespaceInTag() {
	for p.pos < p.len {
		ch := p.peek()
		if ch == ' ' || ch == '\t' {
			p.advance(1)
		} else if ch == '\n' || ch == '\r' {
			p.advance(1)
			if ch == '\r' && p.peek() == '\n' {
				p.advance(1)
			}
			// Skip the javadoc line prefix
			p.skipHorizontalWhitespace()
			if p.peek() == '*' && p.peekAt(1) != '/' {
				p.advance(1)
				if p.peek() == ' ' {
					p.advance(1)
				}
			}
		} else {
			break
		}
	}
}

// parseEntity parses an HTML entity like &nbsp; or &#160;.
func (p *Parser) parseEntity() Node {
	if p.peek() != '&' {
		return nil
	}
	p.advance(1)

	start := p.pos
	if p.peek() == '#' {
		// Numeric entity
		p.advance(1)
		if p.peek() == 'x' || p.peek() == 'X' {
			p.advance(1)
			for isHexDigit(p.peek()) {
				p.advance(1)
			}
		} else {
			for isDigit(p.peek()) {
				p.advance(1)
			}
		}
	} else {
		// Named entity
		for isLetter(p.peek()) {
			p.advance(1)
		}
	}

	name := string(p.input[start:p.pos])

	if p.peek() == ';' {
		p.advance(1)
		return Entity{Name: name}
	}

	// Not a valid entity, return as text
	return Text{Content: "&" + name}
}

// parseBlockTags parses block tags until end of comment.
func (p *Parser) parseBlockTags() []Node {
	var tags []Node

	for p.pos < p.len {
		// Skip whitespace and line prefixes
		p.skipWhitespace()
		p.skipLinePrefix()

		// Check for end of comment
		if p.match("*/") {
			break
		}

		if p.peek() != '@' {
			// Not a block tag, skip this character
			p.advance(1)
			continue
		}

		p.advance(1)
		tagName := p.readTagName()
		if tagName == "" {
			continue
		}

		p.skipHorizontalWhitespace()

		var tag Node
		switch tagName {
		case "param":
			tag = p.parseParamTag()
		case "return":
			tag = p.parseReturnTag()
		case "throws", "exception":
			tag = p.parseThrowsTag()
		case "see":
			tag = p.parseSeeTag()
		case "since":
			tag = p.parseSinceTag()
		case "deprecated":
			tag = p.parseDeprecatedTag()
		case "author":
			tag = p.parseAuthorTag()
		case "version":
			tag = p.parseVersionTag()
		case "serial":
			tag = p.parseSerialTag()
		case "serialData":
			tag = p.parseSerialDataTag()
		case "serialField":
			tag = p.parseSerialFieldTag()
		case "hidden":
			tag = p.parseHiddenTag()
		case "provides":
			tag = p.parseProvidesTag()
		case "uses":
			tag = p.parseUsesTag()
		case "spec":
			tag = p.parseSpecTag()
		default:
			tag = p.parseUnknownBlockTag(tagName)
		}

		if tag != nil {
			tags = append(tags, tag)
		}
	}

	return tags
}

// parseParamTag parses a @param tag.
func (p *Parser) parseParamTag() Node {
	isTypeParam := false
	if p.peek() == '<' {
		isTypeParam = true
		p.advance(1)
	}

	name := p.readIdentifier()

	if isTypeParam && p.peek() == '>' {
		p.advance(1)
	}

	p.skipHorizontalWhitespace()
	desc := p.parseBlockContent()

	return Param{Name: name, IsTypeParam: isTypeParam, Description: desc}
}

// parseReturnTag parses a @return tag.
func (p *Parser) parseReturnTag() Node {
	desc := p.parseBlockContent()
	return Return{Description: desc, Inline: false}
}

// parseThrowsTag parses a @throws or @exception tag.
func (p *Parser) parseThrowsTag() Node {
	exc := p.readReference()
	p.skipHorizontalWhitespace()
	desc := p.parseBlockContent()
	return Throws{Exception: exc, Description: desc}
}

// parseSeeTag parses a @see tag.
func (p *Parser) parseSeeTag() Node {
	var ref []Node

	switch p.peek() {
	case '"':
		// String literal
		s := p.readQuotedString()
		ref = []Node{Text{Content: "\"" + s + "\""}}
	case '<':
		// HTML link
		ref = p.parseBlockContent()
	default:
		// Reference with optional label
		r := p.readReference()
		p.skipHorizontalWhitespace()
		if p.peek() != '@' && !p.isAtBlockTag() && !p.match("*/") {
			rest := p.parseBlockContent()
			ref = append([]Node{Text{Content: r}}, rest...)
		} else {
			ref = []Node{Text{Content: r}}
		}
	}

	return See{Reference: ref}
}

// parseSinceTag parses a @since tag.
func (p *Parser) parseSinceTag() Node {
	desc := p.parseBlockContent()
	return Since{Version: desc}
}

// parseDeprecatedTag parses a @deprecated tag.
func (p *Parser) parseDeprecatedTag() Node {
	desc := p.parseBlockContent()
	return Deprecated{Description: desc}
}

// parseAuthorTag parses an @author tag.
func (p *Parser) parseAuthorTag() Node {
	name := p.parseBlockContent()
	return Author{Name: name}
}

// parseVersionTag parses a @version tag.
func (p *Parser) parseVersionTag() Node {
	ver := p.parseBlockContent()
	return Version{Version: ver}
}

// parseSerialTag parses a @serial tag.
func (p *Parser) parseSerialTag() Node {
	desc := p.parseBlockContent()
	return Serial{Description: desc}
}

// parseSerialDataTag parses a @serialData tag.
func (p *Parser) parseSerialDataTag() Node {
	desc := p.parseBlockContent()
	return SerialData{Description: desc}
}

// parseSerialFieldTag parses a @serialField tag.
func (p *Parser) parseSerialFieldTag() Node {
	name := p.readIdentifier()
	p.skipHorizontalWhitespace()
	typ := p.readReference()
	p.skipHorizontalWhitespace()
	desc := p.parseBlockContent()
	return SerialField{Name: name, Type: typ, Description: desc}
}

// parseHiddenTag parses a @hidden tag.
func (p *Parser) parseHiddenTag() Node {
	desc := p.parseBlockContent()
	return Hidden{Description: desc}
}

// parseProvidesTag parses a @provides tag.
func (p *Parser) parseProvidesTag() Node {
	service := p.readReference()
	p.skipHorizontalWhitespace()
	desc := p.parseBlockContent()
	return Provides{ServiceType: service, Description: desc}
}

// parseUsesTag parses a @uses tag.
func (p *Parser) parseUsesTag() Node {
	service := p.readReference()
	p.skipHorizontalWhitespace()
	desc := p.parseBlockContent()
	return Uses{ServiceType: service, Description: desc}
}

// parseSpecTag parses a @spec tag.
func (p *Parser) parseSpecTag() Node {
	url := p.readWord()
	p.skipHorizontalWhitespace()
	title := p.parseBlockContent()
	return Spec{URL: url, Title: title}
}

// parseUnknownBlockTag parses an unknown block tag.
func (p *Parser) parseUnknownBlockTag(name string) Node {
	content := p.parseBlockContent()
	return UnknownBlockTag{Name: name, Content: content}
}

// parseBlockContent parses content until the next block tag or end of comment.
func (p *Parser) parseBlockContent() []Node {
	return p.parseContent(false)
}

// Helper methods for reading tokens

func (p *Parser) peek() rune {
	if p.pos >= p.len {
		return 0
	}
	return p.input[p.pos]
}

func (p *Parser) peekAt(offset int) rune {
	pos := p.pos + offset
	if pos >= p.len || pos < 0 {
		return 0
	}
	return p.input[pos]
}

func (p *Parser) advance(n int) {
	p.pos += n
	if p.pos > p.len {
		p.pos = p.len
	}
}

func (p *Parser) match(s string) bool {
	if p.pos+len(s) > p.len {
		return false
	}
	for i, ch := range s {
		if p.input[p.pos+i] != ch {
			return false
		}
	}
	return true
}

func (p *Parser) skipWhitespace() {
	for p.pos < p.len && isWhitespace(p.peek()) {
		p.advance(1)
	}
}

func (p *Parser) skipHorizontalWhitespace() {
	for p.pos < p.len && (p.peek() == ' ' || p.peek() == '\t') {
		p.advance(1)
	}
}

func (p *Parser) readTagName() string {
	start := p.pos
	for p.pos < p.len && isJavaIdentifierPart(p.peek()) {
		p.advance(1)
	}
	return string(p.input[start:p.pos])
}

func (p *Parser) readIdentifier() string {
	start := p.pos
	if p.pos < p.len && isJavaIdentifierStart(p.peek()) {
		p.advance(1)
		for p.pos < p.len && isJavaIdentifierPart(p.peek()) {
			p.advance(1)
		}
	}
	return string(p.input[start:p.pos])
}

func (p *Parser) readWord() string {
	start := p.pos
	for p.pos < p.len && !isWhitespace(p.peek()) && p.peek() != '}' {
		p.advance(1)
	}
	return string(p.input[start:p.pos])
}

func (p *Parser) readReference() string {
	start := p.pos
	// A reference can include package.class#member(params)
	// Stop at whitespace, '}', or end of content
	for p.pos < p.len {
		ch := p.peek()
		if isWhitespace(ch) || ch == '}' {
			break
		}
		p.advance(1)
	}
	return strings.TrimSpace(string(p.input[start:p.pos]))
}

func (p *Parser) readQuotedString() string {
	if p.peek() != '"' && p.peek() != '\'' {
		return ""
	}
	quote := p.peek()
	p.advance(1)

	start := p.pos
	for p.pos < p.len && p.peek() != quote {
		if p.peek() == '\\' && p.peekAt(1) == quote {
			p.advance(2)
		} else {
			p.advance(1)
		}
	}

	result := string(p.input[start:p.pos])
	if p.peek() == quote {
		p.advance(1)
	}
	return result
}

func (p *Parser) readUnquotedAttrValue() string {
	start := p.pos
	for p.pos < p.len {
		ch := p.peek()
		if isWhitespace(ch) || ch == '>' || ch == '}' || ch == ':' {
			break
		}
		p.advance(1)
	}
	return string(p.input[start:p.pos])
}

func (p *Parser) readHTMLTagName() string {
	start := p.pos
	for p.pos < p.len {
		ch := p.peek()
		if isLetter(ch) || isDigit(ch) || ch == '-' || ch == '_' || ch == ':' {
			p.advance(1)
		} else {
			break
		}
	}
	return string(p.input[start:p.pos])
}

func (p *Parser) readHTMLAttrName() string {
	start := p.pos
	for p.pos < p.len {
		ch := p.peek()
		if isLetter(ch) || isDigit(ch) || ch == '-' || ch == '_' || ch == ':' {
			p.advance(1)
		} else {
			break
		}
	}
	return string(p.input[start:p.pos])
}

// readBalancedContent reads content until a closing '}', handling nested braces.
func (p *Parser) readBalancedContent() string {
	start := p.pos
	depth := 0

	for p.pos < p.len {
		ch := p.peek()

		if ch == '{' {
			depth++
			p.advance(1)
		} else if ch == '}' {
			if depth == 0 {
				break
			}
			depth--
			p.advance(1)
		} else if ch == '*' && p.peekAt(1) == '/' {
			// End of comment
			break
		} else {
			p.advance(1)
		}
	}

	result := string(p.input[start:p.pos])

	// Trim leading space if first char is space
	if len(result) > 0 && result[0] == ' ' {
		result = result[1:]
	}

	return result
}

// Character classification helpers

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isJavaIdentifierStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_' || ch == '$'
}

func isJavaIdentifierPart(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '$'
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
