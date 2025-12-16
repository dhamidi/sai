// Package javadoc provides a parser for Javadoc comments.
package javadoc

// Node is the interface implemented by all Javadoc AST nodes.
type Node interface {
	node()
}

// DocComment represents a complete Javadoc comment.
type DocComment struct {
	Body      []Node // Main description content
	BlockTags []Node // Block tags like @param, @return, etc.
}

func (DocComment) node() {}

// Text represents plain text content.
type Text struct {
	Content string
}

func (Text) node() {}

// Code represents an {@code ...} inline tag.
type Code struct {
	Content string
}

func (Code) node() {}

// Literal represents an {@literal ...} inline tag.
type Literal struct {
	Content string
}

func (Literal) node() {}

// Link represents an {@link ...} or {@linkplain ...} inline tag.
type Link struct {
	Reference string // The reference (e.g., "java.util.List#add")
	Label     []Node // Optional label content
	Plain     bool   // true for @linkplain, false for @link
}

func (Link) node() {}

// Value represents an {@value ...} inline tag.
type Value struct {
	Reference string
}

func (Value) node() {}

// DocRoot represents an {@docRoot} inline tag.
type DocRoot struct{}

func (DocRoot) node() {}

// InheritDoc represents an {@inheritDoc} inline tag.
type InheritDoc struct {
	Reference string // Optional superclass reference
}

func (InheritDoc) node() {}

// Index represents an {@index ...} inline tag.
type Index struct {
	Term        string
	Description []Node
}

func (Index) node() {}

// Summary represents an {@summary ...} inline tag.
type Summary struct {
	Content []Node
}

func (Summary) node() {}

// Return represents an {@return ...} inline tag or @return block tag.
type Return struct {
	Description []Node
	Inline      bool // true if {@return ...}, false if @return
}

func (Return) node() {}

// SystemProperty represents an {@systemProperty ...} inline tag.
type SystemProperty struct {
	Name string
}

func (SystemProperty) node() {}

// Snippet represents an {@snippet ...} inline tag.
type Snippet struct {
	Attributes map[string]string
	Body       string
}

func (Snippet) node() {}

// UnknownInlineTag represents an unknown inline tag.
type UnknownInlineTag struct {
	Name    string
	Content string
}

func (UnknownInlineTag) node() {}

// Param represents a @param block tag.
type Param struct {
	Name        string
	IsTypeParam bool // true if <T>, false if regular parameter
	Description []Node
}

func (Param) node() {}

// Throws represents a @throws or @exception block tag.
type Throws struct {
	Exception   string
	Description []Node
}

func (Throws) node() {}

// See represents a @see block tag.
type See struct {
	Reference []Node // Can be a reference, string literal, or HTML
}

func (See) node() {}

// Since represents a @since block tag.
type Since struct {
	Version []Node
}

func (Since) node() {}

// Deprecated represents a @deprecated block tag.
type Deprecated struct {
	Description []Node
}

func (Deprecated) node() {}

// Author represents an @author block tag.
type Author struct {
	Name []Node
}

func (Author) node() {}

// Version represents a @version block tag.
type Version struct {
	Version []Node
}

func (Version) node() {}

// Serial represents a @serial block tag.
type Serial struct {
	Description []Node
}

func (Serial) node() {}

// SerialData represents a @serialData block tag.
type SerialData struct {
	Description []Node
}

func (SerialData) node() {}

// SerialField represents a @serialField block tag.
type SerialField struct {
	Name        string
	Type        string
	Description []Node
}

func (SerialField) node() {}

// Hidden represents a @hidden block tag.
type Hidden struct {
	Description []Node
}

func (Hidden) node() {}

// Provides represents a @provides block tag (for module-info).
type Provides struct {
	ServiceType string
	Description []Node
}

func (Provides) node() {}

// Uses represents a @uses block tag (for module-info).
type Uses struct {
	ServiceType string
	Description []Node
}

func (Uses) node() {}

// Spec represents a @spec block tag.
type Spec struct {
	URL   string
	Title []Node
}

func (Spec) node() {}

// UnknownBlockTag represents an unknown block tag.
type UnknownBlockTag struct {
	Name    string
	Content []Node
}

func (UnknownBlockTag) node() {}

// StartElement represents the start of an HTML element.
type StartElement struct {
	Name       string
	Attributes []Attribute
	SelfClose  bool
}

func (StartElement) node() {}

// EndElement represents the end of an HTML element.
type EndElement struct {
	Name string
}

func (EndElement) node() {}

// Attribute represents an HTML attribute.
type Attribute struct {
	Name  string
	Value string
}

// Entity represents an HTML entity like &nbsp; or &#160;.
type Entity struct {
	Name string // The entity name without & and ;
}

func (Entity) node() {}

// Erroneous represents malformed content.
type Erroneous struct {
	Content string
	Message string
}

func (Erroneous) node() {}
