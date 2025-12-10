package format

import (
	"encoding/json"
	"io"

	"github.com/dhamidi/sai/java/parser"
)

type ASTJSONEncoder struct {
	w io.Writer
}

func NewASTJSONEncoder(w io.Writer) *ASTJSONEncoder {
	return &ASTJSONEncoder{w: w}
}

func (e *ASTJSONEncoder) Encode(node *parser.Node) error {
	text, err := e.MarshalText(node)
	if err != nil {
		return err
	}
	_, err = e.w.Write(text)
	return err
}

func (e *ASTJSONEncoder) MarshalText(node *parser.Node) ([]byte, error) {
	return json.MarshalIndent(nodeToJSON(node), "", "  ")
}

type astJSONNode struct {
	Kind     string         `json:"kind"`
	Span     *astJSONSpan   `json:"span,omitempty"`
	Token    string         `json:"token,omitempty"`
	Error    *astJSONError  `json:"error,omitempty"`
	Children []*astJSONNode `json:"children,omitempty"`
}

type astJSONSpan struct {
	Start astJSONPosition `json:"start"`
	End   astJSONPosition `json:"end"`
}

type astJSONPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type astJSONError struct {
	Message  string   `json:"message"`
	Expected []string `json:"expected,omitempty"`
	Got      string   `json:"got,omitempty"`
}

func nodeToJSON(n *parser.Node) *astJSONNode {
	jn := &astJSONNode{
		Kind: n.Kind.String(),
	}

	if n.Span.Start.Line != 0 || n.Span.End.Line != 0 {
		jn.Span = &astJSONSpan{
			Start: astJSONPosition{Line: n.Span.Start.Line, Column: n.Span.Start.Column},
			End:   astJSONPosition{Line: n.Span.End.Line, Column: n.Span.End.Column},
		}
	}

	if n.Token != nil {
		jn.Token = n.Token.Literal
	}

	if n.Error != nil {
		jn.Error = &astJSONError{
			Message: n.Error.Message,
		}
		for _, exp := range n.Error.Expected {
			jn.Error.Expected = append(jn.Error.Expected, exp.String())
		}
		if n.Error.Got != nil {
			jn.Error.Got = n.Error.Got.Literal
		}
	}

	if len(n.Children) > 0 {
		jn.Children = make([]*astJSONNode, len(n.Children))
		for i, child := range n.Children {
			jn.Children[i] = nodeToJSON(child)
		}
	}

	return jn
}
