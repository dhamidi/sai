package codebase

import (
	"testing"

	"github.com/dhamidi/sai/java"
	"github.com/dhamidi/sai/java/parser"
)

func TestCompletionForRecordComponents(t *testing.T) {
	content := `package deps;

public class Main {
  public record Dependency(String name, String version) {}

  public static void main(String[] args) {
    Dependency dep = new Dependency("foo", "1.0");
    dep.
  }
}`

	c := New("/tmp/lsp_test")
	path := "/tmp/lsp_test/src/deps/Main.java"
	c.UpdateFile(path, []byte(content))

	f := c.GetFile(path)
	if f == nil {
		t.Fatal("GetFile returned nil")
	}
	if f.AST == nil {
		t.Fatal("AST is nil")
	}

	// line 8, after "dep."
	line := 8
	col := 8
	triggerCol := findTriggerPosition(f.Content, line, col)

	pos := parser.Position{Line: line, Column: triggerCol}
	typeName := java.TypeAtPoint(f.AST, pos, c.AllClasses())
	if typeName != "deps.Main.Dependency" {
		t.Errorf("TypeAtPoint = %q, want %q", typeName, "deps.Main.Dependency")
	}

	cls := c.FindClass("deps.Main.Dependency")
	if cls == nil {
		t.Fatal("FindClass returned nil")
	}
	if len(cls.RecordComponents) != 2 {
		t.Errorf("RecordComponents = %d, want 2", len(cls.RecordComponents))
	}

	completions := c.CompletionsAtPoint(path, line, triggerCol)
	if len(completions) != 2 {
		t.Errorf("Completions = %d, want 2", len(completions))
	}

	wantLabels := map[string]bool{"name": true, "version": true}
	for _, comp := range completions {
		if !wantLabels[comp.Label] {
			t.Errorf("unexpected completion label: %s", comp.Label)
		}
	}
}
