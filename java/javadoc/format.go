package javadoc

import (
	"strings"
)

// Format formats a DocComment AST into a readable string.
func Format(doc *DocComment) string {
	if doc == nil {
		return ""
	}

	var sb strings.Builder

	// Format body
	body := formatNodes(doc.Body)
	body = normalizeWhitespace(body)
	sb.WriteString(body)

	// Format block tags
	if len(doc.BlockTags) > 0 && sb.Len() > 0 {
		sb.WriteString("\n")
	}

	for _, tag := range doc.BlockTags {
		s := formatBlockTag(tag)
		if s != "" {
			sb.WriteString("\n")
			sb.WriteString(s)
		}
	}

	return strings.TrimSpace(sb.String())
}

// FormatPlainText formats a DocComment AST into plain text without any formatting.
func FormatPlainText(doc *DocComment) string {
	if doc == nil {
		return ""
	}

	var sb strings.Builder

	body := formatNodesPlain(doc.Body)
	body = normalizeWhitespace(body)
	sb.WriteString(body)

	return strings.TrimSpace(sb.String())
}

func formatNodes(nodes []Node) string {
	var sb strings.Builder
	for i, node := range nodes {
		// Skip <pre> tags when they immediately precede {@code} with multiline content
		// because the Code formatter will handle the code block formatting
		if start, ok := node.(StartElement); ok && strings.ToLower(start.Name) == "pre" {
			if hasMultilineCodeNext(nodes, i) {
				continue
			}
		}
		// Skip </pre> tags when they immediately follow {@code} with multiline content
		if end, ok := node.(EndElement); ok && strings.ToLower(end.Name) == "pre" {
			if hasMultilineCodeBefore(nodes, i) {
				continue
			}
		}
		sb.WriteString(formatNode(node))
	}
	return sb.String()
}

func hasMultilineCodeNext(nodes []Node, idx int) bool {
	// Look ahead for Code node, skipping whitespace-only Text nodes
	for i := idx + 1; i < len(nodes); i++ {
		switch n := nodes[i].(type) {
		case Text:
			if strings.TrimSpace(n.Content) == "" {
				continue
			}
			return false
		case Code:
			return strings.Contains(n.Content, "\n")
		default:
			return false
		}
	}
	return false
}

func hasMultilineCodeBefore(nodes []Node, idx int) bool {
	// Look back for Code node, skipping whitespace-only Text nodes
	for i := idx - 1; i >= 0; i-- {
		switch n := nodes[i].(type) {
		case Text:
			if strings.TrimSpace(n.Content) == "" {
				continue
			}
			return false
		case Code:
			return strings.Contains(n.Content, "\n")
		default:
			return false
		}
	}
	return false
}

func formatNodesPlain(nodes []Node) string {
	var sb strings.Builder
	for _, node := range nodes {
		sb.WriteString(formatNodePlain(node))
	}
	return sb.String()
}

func formatNode(node Node) string {
	switch n := node.(type) {
	case Text:
		return n.Content
	case Code:
		content := stripJavadocLinePrefix(n.Content)
		content = strings.TrimSpace(content)
		// Use fenced code block for multi-line code, inline backticks for single line
		if strings.Contains(content, "\n") {
			return "\n```\n" + content + "\n```\n"
		}
		return "`" + content + "`"
	case Literal:
		return n.Content
	case Link:
		if len(n.Label) > 0 {
			return formatNodes(n.Label)
		}
		return formatReference(n.Reference)
	case Value:
		return formatReference(n.Reference)
	case DocRoot:
		return ""
	case InheritDoc:
		return "[inherited documentation]"
	case Index:
		return n.Term
	case Summary:
		return formatNodes(n.Content)
	case Return:
		if n.Inline {
			return formatNodes(n.Description)
		}
		return ""
	case SystemProperty:
		return n.Name
	case Snippet:
		return n.Body
	case UnknownInlineTag:
		return n.Content
	case StartElement:
		return formatStartElement(n)
	case EndElement:
		return formatEndElement(n)
	case Entity:
		return decodeEntity(n.Name)
	case Erroneous:
		return n.Content
	default:
		return ""
	}
}

func formatNodePlain(node Node) string {
	switch n := node.(type) {
	case Text:
		return n.Content
	case Code:
		return n.Content
	case Literal:
		return n.Content
	case Link:
		if len(n.Label) > 0 {
			return formatNodesPlain(n.Label)
		}
		return formatReference(n.Reference)
	case Value:
		return formatReference(n.Reference)
	case DocRoot:
		return ""
	case InheritDoc:
		return ""
	case Index:
		return n.Term
	case Summary:
		return formatNodesPlain(n.Content)
	case Return:
		if n.Inline {
			return formatNodesPlain(n.Description)
		}
		return ""
	case SystemProperty:
		return n.Name
	case Snippet:
		return n.Body
	case UnknownInlineTag:
		return n.Content
	case StartElement:
		return ""
	case EndElement:
		return ""
	case Entity:
		return decodeEntity(n.Name)
	case Erroneous:
		return ""
	default:
		return ""
	}
}

func formatReference(ref string) string {
	// Extract simple name from reference like java.util.List#add(E)
	if idx := strings.LastIndex(ref, "#"); idx >= 0 {
		member := ref[idx+1:]
		// Remove parameters from method reference
		if paren := strings.Index(member, "("); paren >= 0 {
			member = member[:paren]
		}
		return member
	}
	// Extract class simple name from fully qualified name
	if idx := strings.LastIndex(ref, "."); idx >= 0 {
		return ref[idx+1:]
	}
	return ref
}

func formatStartElement(e StartElement) string {
	tag := strings.ToLower(e.Name)

	// Handle specific HTML elements
	switch tag {
	case "p":
		return "\n\n"
	case "br":
		return "\n"
	case "pre":
		return "\n```\n"
	case "code":
		return "`"
	case "ul", "ol":
		return "\n"
	case "li":
		return "\n- "
	case "b", "strong":
		return ""
	case "i", "em":
		return ""
	case "a":
		// For links, try to extract href for display
		for _, attr := range e.Attributes {
			if strings.ToLower(attr.Name) == "href" && attr.Value != "" {
				// External link - will show as [text](url) if we had proper link handling
				// For now just skip the tag itself
				return ""
			}
		}
		return ""
	case "blockquote":
		return "\n> "
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return "\n\n"
	case "table", "thead", "tbody", "tr":
		return "\n"
	case "td", "th":
		return " "
	case "dl":
		return "\n"
	case "dt":
		return "\n"
	case "dd":
		return "\n  "
	default:
		return ""
	}
}

func formatEndElement(e EndElement) string {
	tag := strings.ToLower(e.Name)

	switch tag {
	case "pre":
		return "\n```\n"
	case "code":
		return "`"
	case "ul", "ol":
		return "\n"
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return "\n"
	default:
		return ""
	}
}

func formatBlockTag(node Node) string {
	switch n := node.(type) {
	case Param:
		desc := formatNodes(n.Description)
		if n.IsTypeParam {
			return "@param <" + n.Name + "> " + strings.TrimSpace(desc)
		}
		return "@param " + n.Name + " " + strings.TrimSpace(desc)
	case Return:
		desc := formatNodes(n.Description)
		return "@return " + strings.TrimSpace(desc)
	case Throws:
		desc := formatNodes(n.Description)
		return "@throws " + n.Exception + " " + strings.TrimSpace(desc)
	case See:
		return "@see " + strings.TrimSpace(formatNodes(n.Reference))
	case Since:
		return "@since " + strings.TrimSpace(formatNodes(n.Version))
	case Deprecated:
		desc := formatNodes(n.Description)
		return "@deprecated " + strings.TrimSpace(desc)
	case Author:
		return "@author " + strings.TrimSpace(formatNodes(n.Name))
	case Version:
		return "@version " + strings.TrimSpace(formatNodes(n.Version))
	case Serial:
		desc := formatNodes(n.Description)
		return "@serial " + strings.TrimSpace(desc)
	case SerialData:
		desc := formatNodes(n.Description)
		return "@serialData " + strings.TrimSpace(desc)
	case SerialField:
		desc := formatNodes(n.Description)
		return "@serialField " + n.Name + " " + n.Type + " " + strings.TrimSpace(desc)
	case Hidden:
		return "@hidden"
	case Provides:
		desc := formatNodes(n.Description)
		return "@provides " + n.ServiceType + " " + strings.TrimSpace(desc)
	case Uses:
		desc := formatNodes(n.Description)
		return "@uses " + n.ServiceType + " " + strings.TrimSpace(desc)
	case Spec:
		title := formatNodes(n.Title)
		return "@spec " + n.URL + " " + strings.TrimSpace(title)
	case UnknownBlockTag:
		content := formatNodes(n.Content)
		return "@" + n.Name + " " + strings.TrimSpace(content)
	default:
		return ""
	}
}

func decodeEntity(name string) string {
	switch name {
	case "lt", "#60":
		return "<"
	case "gt", "#62":
		return ">"
	case "amp", "#38":
		return "&"
	case "quot", "#34":
		return "\""
	case "apos", "#39":
		return "'"
	case "nbsp", "#160":
		return " "
	case "mdash", "#8212":
		return "—"
	case "ndash", "#8211":
		return "–"
	case "copy", "#169":
		return "©"
	case "reg", "#174":
		return "®"
	case "trade", "#8482":
		return "™"
	default:
		// Return the entity as-is if unknown
		return "&" + name + ";"
	}
}

func normalizeWhitespace(s string) string {
	// Replace multiple consecutive newlines with at most two
	lines := strings.Split(s, "\n")
	var result []string
	prevEmpty := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !prevEmpty {
				result = append(result, "")
				prevEmpty = true
			}
		} else {
			result = append(result, line)
			prevEmpty = false
		}
	}

	return strings.Join(result, "\n")
}

// stripJavadocLinePrefix removes the leading " * " or " *" from each line of Javadoc content.
func stripJavadocLinePrefix(s string) string {
	lines := strings.Split(s, "\n")
	var result []string

	for _, line := range lines {
		// Try to strip common Javadoc line prefix patterns
		trimmed := line
		if strings.HasPrefix(trimmed, " * ") {
			trimmed = trimmed[3:]
		} else if strings.HasPrefix(trimmed, " *") {
			trimmed = trimmed[2:]
		} else if strings.HasPrefix(trimmed, "* ") {
			trimmed = trimmed[2:]
		} else if strings.HasPrefix(trimmed, "*") && len(trimmed) > 1 && trimmed[1] != '/' {
			trimmed = trimmed[1:]
		}
		result = append(result, trimmed)
	}

	return strings.Join(result, "\n")
}
