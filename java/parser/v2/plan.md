# Java Parser v2 Implementation Plan

**Thread:** https://ampcode.com/threads/T-019b12da-f52e-72fc-b85d-bfb243294a37

## Overview

Implement `java/parser/v2` using the EBNF grammar (`java/java25.ebnf`) and Earley parser (`ebnf/parse`) as an alternative to the hand-written recursive descent parser in `java/parser`.

## Architecture

```
v1 (java/parser)                    v2 (java/parser/v2)
─────────────────                   ───────────────────
Lexer (lexer.go)                    java25.ebnf
    ↓                                   ↓
Parser (parser.go)                  ebnf/lex.Lexer
    ↓                                   ↓
parser.Node                         ebnf/parse.Earley
    ↓                                   ↓
from_source.go                      ebnf/parse.Node
    ↓                                   ↓
ClassModel ←────────────────────→  from_source.go (v2)
                                        ↓
                                    ClassModel (same)
```

Both implementations produce identical `ClassModel` output. The v2 parser uses `ebnf/parse.Node` directly — no adapter layer needed.

## Components

| Component | File | Effort | Description |
|-----------|------|--------|-------------|
| Parser wrapper | `parser.go` | S (1-2h) | Load grammar, tokenize, run Earley |
| ClassModel extraction | `from_source.go` | M-L (4-8h) | Walk CST to produce ClassModel |
| Comparison tests | `compare_test.go` | M (2-4h) | Verify v1 vs v2 produce same models |

## Implementation Details

### parser.go

```go
package v2

import (
    "io"
    "github.com/dhamidi/sai/ebnf/grammar"
    "github.com/dhamidi/sai/ebnf/lex"
    "github.com/dhamidi/sai/ebnf/parse"
)

//go:embed ../java25.ebnf
var javaGrammarSource []byte

var javaGrammar grammar.Grammar

func init() {
    // Parse embedded grammar
}

type Parser struct {
    input    []byte
    filename string
    tokens   []lex.Token
    comments []lex.Token // extracted before Earley filtering
}

func Parse(r io.Reader, filename string) (*parse.Node, error) {
    // 1. Read input
    // 2. Tokenize with lex.NewLexer(javaGrammar, input, filename)
    // 3. Extract comments from token stream
    // 4. Run Earley parser with start production "compilationUnit"
    // 5. Return CST
}
```

### from_source.go

Mirrors `java/from_source.go` but walks `ebnf/parse.Node` using grammar production names:

| v1 NodeKind | v2 CST Production |
|-------------|-------------------|
| `KindCompilationUnit` | `"compilationUnit"` |
| `KindPackageDecl` | `"packageDeclaration"` |
| `KindImportDecl` | `"importDeclaration"` |
| `KindClassDecl` | `"normalClassDeclaration"` |
| `KindInterfaceDecl` | `"normalInterfaceDeclaration"` |
| `KindEnumDecl` | `"enumDeclaration"` |
| `KindRecordDecl` | `"recordDeclaration"` |
| `KindAnnotationDecl` | `"annotationInterfaceDeclaration"` |
| `KindMethodDecl` | `"methodDeclaration"` |
| `KindFieldDecl` | `"fieldDeclaration"` |
| `KindConstructorDecl` | `"constructorDeclaration"` |

### compare_test.go

```go
func TestV1V2Equivalence(t *testing.T) {
    // For each .java file in testdata:
    //   v1Models := java.ClassModelsFromSource(source)
    //   v2Models := v2.ClassModelsFromSource(source)
    //   Compare and assert equality
}
```

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| **Performance** | High | Earley is O(n³) worst case. Benchmark on real JDK files. Start as opt-in. |
| **Grammar ambiguity** | Medium | Earley picks first derivation. May differ from v1 disambiguation. Golden tests catch this. |
| **CST shape mismatch** | Medium | Production names/nesting may differ from v1 expectations. Careful mapping needed. |
| **Comments/Javadoc** | Low | Extract comments from token stream before Earley filters them. |

## Success Criteria

1. `v2.ClassModelsFromSource(source)` returns identical `[]ClassModel` as `java.ClassModelsFromSource(source)` for all test files
2. Parses JDK source files without errors (same coverage as v1)
3. Performance within 10x of v1 on typical files (acceptable for experimental use)

## Parallel Evaluation

```go
// In java/ package or tests:
func CompareParserOutputs(source []byte) error {
    v1Models, err := ClassModelsFromSource(source)
    if err != nil {
        return err
    }
    v2Models, err := v2.ClassModelsFromSource(source)
    if err != nil {
        return err
    }
    return compareModels(v1Models, v2Models)
}
```

## Future Work

- Performance optimization if Earley proves too slow
- Error recovery for partial/invalid input
- Incremental parsing for editor integration
