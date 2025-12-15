# How to Fix a Roundtrip Test Failure

This document describes the process for fixing failures in the roundtrip formatter test (`format/roundtrip_test.go`).

## Architecture Overview

The codebase has two paths for converting Java source to output:

1. **AST → Pretty Printer** (`format/java_pretty.go`)
   - Works directly with parser AST nodes
   - Preserves source formatting as much as possible
   - Used by the roundtrip test
   - Entry point: `PrettyPrintJava(source []byte)`

2. **AST → Class Model → Serializer** (`java/from_source.go` → `format/java.go`)
   - Converts AST to a structured `ClassModel` (`java/model.go`)
   - Serializes the model to Java source
   - Used for generating simplified class representations
   - Entry point: `JavaModelEncoder.Encode(model *ClassModel)`

**Important**: "Formatting" in this context means serializing source code. The pretty printer and model serializer are both formatters, just at different abstraction levels.

## The Roundtrip Test

The test in `format/roundtrip_test.go`:

1. Parses original Java source → AST
2. Formats AST back to source using `PrettyPrintJava`
3. Parses formatted output → AST
4. Compares node counts between original and formatted ASTs

A failure means either:
- The formatted output doesn't parse (syntax error introduced)
- The formatted output parses but has different AST structure (nodes dropped/added)

## Step-by-Step Fix Process

### 1. Identify a Failing Test

Run the roundtrip tests with `-failfast` to stop at the first failure:

```bash
go test ./format -v -run 'TestRoundTrip_Testcases' -failfast
```

To run a specific file:

```bash
go test ./format -v -run 'TestRoundTrip_Testcases/filename_without_java'
```

To filter by substring:

```bash
go test ./format -v -run 'TestRoundTrip_Testcases' -filter=SomePattern
```

### 2. Analyze the Failure

The test output shows:
- **Parse errors**: The formatted output is syntactically invalid
- **Node count mismatches**: Lists which node kinds differ and by how much
- **Formatted output**: The actual problematic output (truncated)

For parse errors, look at the formatted output around the error line. The issue is usually obvious once you see the malformed output.

### 3. Find the Root Cause

Look at the original source to understand what construct is being mishandled:

```bash
# View specific lines of the original file
head -N /path/to/testcases/file.java | tail -M
```

Common causes:
- **Missing case in formatter**: A node kind isn't handled in `printNode()` or related functions
- **Incorrect AST structure assumption**: The parser creates a different tree than expected
- **Comment handling**: Comments stored separately need explicit emission

Check the parser to understand the AST structure:

```bash
# Search for how a construct is parsed
grep -n "KindSomething" java/parser/parser.go
```

### 4. Decide: Model Change or Formatter-Only Fix?

**Add to ClassModel if**:
- The construct represents a semantic concept (fields, methods, initializers, etc.)
- It should be preserved when round-tripping through the model
- Other tools might need to query/manipulate it

**Formatter-only fix if**:
- It's purely syntactic (operator precedence, parentheses)
- It's a special case of existing handling

### 5. Update the Class Model (if needed)

Edit `java/model.go` to add the new type:

```go
type ClassModel struct {
    // ... existing fields ...
    NewFeature []NewFeatureModel  // Add new field
}

type NewFeatureModel struct {
    // Fields that capture the semantic content
}
```

### 6. Extract from AST (if model was updated)

Edit `java/from_source.go` to extract the new construct:

1. Find `extractClassBodyMembers()` or the relevant extraction function
2. Add a case for the new node kind
3. Create a helper function if the extraction is complex

```go
case parser.KindSomething:
    if feature := featureFromNode(child); feature != nil {
        model.NewFeature = append(model.NewFeature, *feature)
    }
```

### 7. Update the Model Serializer (if model was updated)

Edit `format/java.go` to serialize the new construct:

1. Add a `writeNewFeature()` method to `JavaModelEncoder`
2. Call it from `MarshalText()` in the appropriate order

### 8. Fix the Pretty Printer

Edit `format/java_pretty.go`:

1. **For new node kinds**: Add a case in `printNode()` dispatching to a new print function
2. **For special cases of existing kinds**: Add detection logic and special handling

Example pattern for detecting special AST structures:

```go
func (p *JavaPrettyPrinter) isSpecialCase(node *parser.Node) bool {
    // Check node structure to identify the special case
    for _, child := range node.Children {
        if child.Kind == parser.KindSomething {
            return true
        }
    }
    return false
}

func (p *JavaPrettyPrinter) printSpecialCase(node *parser.Node) {
    // Handle the special case
}
```

### 9. Add a Test Case

Create a minimal test file in `testcases/` that exercises the construct:

```java
// testcases/NewFeatureTest.java
package test;

public class NewFeatureTest {
    // Minimal example of the construct
}
```

### 10. Verify the Fix

Run the test for your new test case:

```bash
go test ./format -v -run 'TestRoundTrip_Testcases/NewFeatureTest'
```

Run the original failing test:

```bash
go test ./format -v -run 'TestRoundTrip_Testcases/OriginalFailingTest'
```

Run all roundtrip tests to check for regressions:

```bash
go test ./format -v -run 'TestRoundTrip_Testcases'
```

### 11. Commit

```bash
git add -A
git commit -m "Add support for <feature> in class model and formatter

- Add <Type>Model to ClassModel
- Extract <feature> in from_source.go
- Serialize <feature> in JavaModelEncoder
- Fix java_pretty.go to handle <feature> nodes
- Add <Feature>Test.java test case"
```

## Common Patterns

### Comment Handling

Comments are stored separately from the AST. Use `emitCommentsBeforeLine(line)` to emit comments that appear before a given source line:

```go
func (p *JavaPrettyPrinter) printSomething(node *parser.Node) {
    for _, child := range node.Children {
        p.emitCommentsBeforeLine(child.Span.Start.Line)
        // ... print child ...
    }
    // Emit comments before closing brace
    p.emitCommentsBeforeLine(node.Span.End.Line)
}
```

### Detecting AST Patterns

The parser sometimes creates wrapper nodes. Check the actual structure:

```go
// Static initializer is: KindBlock containing [KindIdentifier("static"), KindBlock]
func (p *JavaPrettyPrinter) isStaticInitializer(node *parser.Node) bool {
    hasStatic := false
    hasBlock := false
    for _, child := range node.Children {
        if child.Kind == parser.KindIdentifier &&
           child.Token != nil &&
           child.Token.Literal == "static" {
            hasStatic = true
        } else if child.Kind == parser.KindBlock {
            hasBlock = true
        }
    }
    return hasStatic && hasBlock
}
```

## Files Reference

| File | Purpose |
|------|---------|
| `java/model.go` | Class model definitions |
| `java/from_source.go` | AST → ClassModel extraction |
| `format/java.go` | ClassModel → Java source serialization |
| `format/java_pretty.go` | AST → Java source pretty printing |
| `format/roundtrip_test.go` | Roundtrip test implementation |
| `java/parser/parser.go` | Parser (to understand AST structure) |
| `java/parser/node.go` | Node kind definitions |
| `testcases/` | Test Java files |
