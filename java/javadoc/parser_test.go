package javadoc

import (
	"testing"
)

func TestParseSimpleText(t *testing.T) {
	doc := Parse("/** Simple text. */")

	if len(doc.Body) != 1 {
		t.Fatalf("expected 1 body node, got %d", len(doc.Body))
	}

	text, ok := doc.Body[0].(Text)
	if !ok {
		t.Fatalf("expected Text node, got %T", doc.Body[0])
	}

	if text.Content != "Simple text. " {
		t.Errorf("expected 'Simple text. ', got %q", text.Content)
	}
}

func TestParseCodeTag(t *testing.T) {
	doc := Parse("/** Use {@code Map<String, List<Integer>>} for this. */")

	// Should have: Text, Code, Text
	if len(doc.Body) != 3 {
		t.Fatalf("expected 3 body nodes, got %d: %+v", len(doc.Body), doc.Body)
	}

	code, ok := doc.Body[1].(Code)
	if !ok {
		t.Fatalf("expected Code node, got %T", doc.Body[1])
	}

	expected := "Map<String, List<Integer>>"
	if code.Content != expected {
		t.Errorf("expected %q, got %q", expected, code.Content)
	}
}

func TestParseCodeTagWithBraces(t *testing.T) {
	doc := Parse("/** Use {@code class Foo { int x; }} for this. */")

	if len(doc.Body) != 3 {
		t.Fatalf("expected 3 body nodes, got %d: %+v", len(doc.Body), doc.Body)
	}

	code, ok := doc.Body[1].(Code)
	if !ok {
		t.Fatalf("expected Code node, got %T", doc.Body[1])
	}

	expected := "class Foo { int x; }"
	if code.Content != expected {
		t.Errorf("expected %q, got %q", expected, code.Content)
	}
}

func TestParseLinkTag(t *testing.T) {
	doc := Parse("/** See {@link java.util.List} for more. */")

	if len(doc.Body) != 3 {
		t.Fatalf("expected 3 body nodes, got %d", len(doc.Body))
	}

	link, ok := doc.Body[1].(Link)
	if !ok {
		t.Fatalf("expected Link node, got %T", doc.Body[1])
	}

	if link.Reference != "java.util.List" {
		t.Errorf("expected 'java.util.List', got %q", link.Reference)
	}
	if link.Plain {
		t.Error("expected Plain to be false")
	}
}

func TestParseLinkTagWithLabel(t *testing.T) {
	doc := Parse("/** See {@link java.util.List the List interface}. */")

	link, ok := doc.Body[1].(Link)
	if !ok {
		t.Fatalf("expected Link node, got %T", doc.Body[1])
	}

	if link.Reference != "java.util.List" {
		t.Errorf("expected 'java.util.List', got %q", link.Reference)
	}

	if len(link.Label) != 1 {
		t.Fatalf("expected 1 label node, got %d", len(link.Label))
	}

	text, ok := link.Label[0].(Text)
	if !ok {
		t.Fatalf("expected Text label, got %T", link.Label[0])
	}

	if text.Content != "the List interface" {
		t.Errorf("expected 'the List interface', got %q", text.Content)
	}
}

func TestParseParamTag(t *testing.T) {
	doc := Parse(`/**
	 * Description.
	 * @param name the name of the thing
	 */`)

	if len(doc.BlockTags) != 1 {
		t.Fatalf("expected 1 block tag, got %d", len(doc.BlockTags))
	}

	param, ok := doc.BlockTags[0].(Param)
	if !ok {
		t.Fatalf("expected Param, got %T", doc.BlockTags[0])
	}

	if param.Name != "name" {
		t.Errorf("expected param name 'name', got %q", param.Name)
	}
	if param.IsTypeParam {
		t.Error("expected IsTypeParam to be false")
	}
}

func TestParseTypeParamTag(t *testing.T) {
	doc := Parse(`/**
	 * @param <T> the element type
	 */`)

	if len(doc.BlockTags) != 1 {
		t.Fatalf("expected 1 block tag, got %d", len(doc.BlockTags))
	}

	param, ok := doc.BlockTags[0].(Param)
	if !ok {
		t.Fatalf("expected Param, got %T", doc.BlockTags[0])
	}

	if param.Name != "T" {
		t.Errorf("expected param name 'T', got %q", param.Name)
	}
	if !param.IsTypeParam {
		t.Error("expected IsTypeParam to be true")
	}
}

func TestParseReturnTag(t *testing.T) {
	doc := Parse(`/**
	 * @return the computed value
	 */`)

	if len(doc.BlockTags) != 1 {
		t.Fatalf("expected 1 block tag, got %d", len(doc.BlockTags))
	}

	ret, ok := doc.BlockTags[0].(Return)
	if !ok {
		t.Fatalf("expected Return, got %T", doc.BlockTags[0])
	}

	if ret.Inline {
		t.Error("expected Inline to be false")
	}
}

func TestParseThrowsTag(t *testing.T) {
	doc := Parse(`/**
	 * @throws IllegalArgumentException if the argument is null
	 */`)

	if len(doc.BlockTags) != 1 {
		t.Fatalf("expected 1 block tag, got %d", len(doc.BlockTags))
	}

	throws, ok := doc.BlockTags[0].(Throws)
	if !ok {
		t.Fatalf("expected Throws, got %T", doc.BlockTags[0])
	}

	if throws.Exception != "IllegalArgumentException" {
		t.Errorf("expected 'IllegalArgumentException', got %q", throws.Exception)
	}
}

func TestParseHTMLEntity(t *testing.T) {
	doc := Parse("/** A &lt; B &amp;&amp; C &gt; D */")

	formatted := Format(doc)
	expected := "A < B && C > D"
	if formatted != expected {
		t.Errorf("expected %q, got %q", expected, formatted)
	}
}

func TestParseHTMLTags(t *testing.T) {
	doc := Parse("/** <p>First paragraph.</p><p>Second.</p> */")

	// Count the nodes
	var startCount, endCount, textCount int
	for _, node := range doc.Body {
		switch node.(type) {
		case StartElement:
			startCount++
		case EndElement:
			endCount++
		case Text:
			textCount++
		}
	}

	if startCount != 2 {
		t.Errorf("expected 2 StartElements, got %d", startCount)
	}
	if endCount != 2 {
		t.Errorf("expected 2 EndElements, got %d", endCount)
	}
}

func TestParseMultipleBlockTags(t *testing.T) {
	doc := Parse(`/**
	 * Description here.
	 *
	 * @param x the x coordinate
	 * @param y the y coordinate
	 * @return the distance
	 * @throws IllegalArgumentException if negative
	 */`)

	if len(doc.BlockTags) != 4 {
		t.Fatalf("expected 4 block tags, got %d", len(doc.BlockTags))
	}

	// Check types
	if _, ok := doc.BlockTags[0].(Param); !ok {
		t.Errorf("expected Param at 0, got %T", doc.BlockTags[0])
	}
	if _, ok := doc.BlockTags[1].(Param); !ok {
		t.Errorf("expected Param at 1, got %T", doc.BlockTags[1])
	}
	if _, ok := doc.BlockTags[2].(Return); !ok {
		t.Errorf("expected Return at 2, got %T", doc.BlockTags[2])
	}
	if _, ok := doc.BlockTags[3].(Throws); !ok {
		t.Errorf("expected Throws at 3, got %T", doc.BlockTags[3])
	}
}

func TestFormat(t *testing.T) {
	input := `/**
	 * This is a description with {@code some code} in it.
	 * And a {@link java.util.List} reference.
	 *
	 * @param name the name to use
	 * @return the result
	 */`

	doc := Parse(input)
	formatted := Format(doc)

	// Check that code is formatted with backticks
	if !contains(formatted, "`some code`") {
		t.Errorf("expected backtick-wrapped code in output: %s", formatted)
	}

	// Check that block tags are present
	if !contains(formatted, "@param name") {
		t.Errorf("expected @param in output: %s", formatted)
	}
	if !contains(formatted, "@return") {
		t.Errorf("expected @return in output: %s", formatted)
	}
}

func TestParseNestedBraces(t *testing.T) {
	// This is the problematic case that was causing truncation
	input := `/**
	 * Example:
	 * {@code
	 * class OneShotPublisher implements Publisher {
	 *   public void subscribe(Subscriber subscriber) {
	 *     if (subscribed)
	 *       subscriber.onError(new IllegalStateException());
	 *   }
	 * }
	 * }
	 */`

	doc := Parse(input)

	// Find the Code node
	var codeNode *Code
	for _, node := range doc.Body {
		if c, ok := node.(Code); ok {
			codeNode = &c
			break
		}
	}

	if codeNode == nil {
		t.Fatal("expected to find Code node")
	}

	// The code should contain all the braces
	if !contains(codeNode.Content, "class OneShotPublisher") {
		t.Errorf("code content missing class declaration: %s", codeNode.Content)
	}
	if !contains(codeNode.Content, "subscriber.onError") {
		t.Errorf("code content missing method body: %s", codeNode.Content)
	}
}

func TestParseSeeTag(t *testing.T) {
	doc := Parse(`/**
	 * @see java.util.List#add(Object)
	 */`)

	if len(doc.BlockTags) != 1 {
		t.Fatalf("expected 1 block tag, got %d", len(doc.BlockTags))
	}

	see, ok := doc.BlockTags[0].(See)
	if !ok {
		t.Fatalf("expected See, got %T", doc.BlockTags[0])
	}

	if len(see.Reference) == 0 {
		t.Error("expected non-empty reference")
	}
}

func TestParseSinceTag(t *testing.T) {
	doc := Parse(`/**
	 * @since 1.8
	 */`)

	if len(doc.BlockTags) != 1 {
		t.Fatalf("expected 1 block tag, got %d", len(doc.BlockTags))
	}

	since, ok := doc.BlockTags[0].(Since)
	if !ok {
		t.Fatalf("expected Since, got %T", doc.BlockTags[0])
	}

	formatted := formatNodes(since.Version)
	if !contains(formatted, "1.8") {
		t.Errorf("expected '1.8' in since version, got %q", formatted)
	}
}

func TestParsePreBlock(t *testing.T) {
	input := `/**
	 * Example:
	 * <pre>{@code
	 * public class Foo {
	 *     private int x;
	 * }
	 * }</pre>
	 */`

	doc := Parse(input)
	formatted := Format(doc)

	// The formatted output should contain the code
	if !contains(formatted, "public class Foo") {
		t.Errorf("expected class declaration in output: %s", formatted)
	}
	if !contains(formatted, "private int x") {
		t.Errorf("expected field declaration in output: %s", formatted)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
