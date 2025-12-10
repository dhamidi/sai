package parser

import "encoding/json"

type jsonNode struct {
	Kind     string      `json:"kind"`
	Span     *jsonSpan   `json:"span,omitempty"`
	Token    string      `json:"token,omitempty"`
	Error    *jsonError  `json:"error,omitempty"`
	Children []*jsonNode `json:"children,omitempty"`
}

type jsonSpan struct {
	Start jsonPosition `json:"start"`
	End   jsonPosition `json:"end"`
}

type jsonPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type jsonError struct {
	Message  string   `json:"message"`
	Expected []string `json:"expected,omitempty"`
	Got      string   `json:"got,omitempty"`
}

func (n *Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.toJSON())
}

func (n *Node) toJSON() *jsonNode {
	jn := &jsonNode{
		Kind: n.Kind.String(),
	}

	if n.Span.Start.Line != 0 || n.Span.End.Line != 0 {
		jn.Span = &jsonSpan{
			Start: jsonPosition{Line: n.Span.Start.Line, Column: n.Span.Start.Column},
			End:   jsonPosition{Line: n.Span.End.Line, Column: n.Span.End.Column},
		}
	}

	if n.Token != nil {
		jn.Token = n.Token.Literal
	}

	if n.Error != nil {
		jn.Error = &jsonError{
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
		jn.Children = make([]*jsonNode, len(n.Children))
		for i, child := range n.Children {
			jn.Children[i] = child.toJSON()
		}
	}

	return jn
}
