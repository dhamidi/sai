package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseJDKSourceFiles(t *testing.T) {
	if os.Getenv("IN_GIT_PRECOMMIT") != "" {
		t.Skip("skipping JDK tests during pre-commit")
	}
	jdkDir := "../../testcases/jdk"

	err := filepath.Walk(jdkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".java") {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			p := ParseCompilationUnit(strings.NewReader(string(content)), WithFile(path))
			node := p.Finish()

			if hasError(node) {
				t.Errorf("parse error in %s", path)
				t.Logf("Parse tree:\n%s", node.String())
			}
		})

		return nil
	})

	if err != nil {
		t.Fatalf("failed to walk jdk directory: %v", err)
	}
}
