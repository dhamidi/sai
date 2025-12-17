# Column Position Tracking for Java Formatter

This document describes the implementation plan for adding column position tracking to the Java pretty printer, enabling line-length-aware formatting decisions.

## Problem Statement

The formatter needs to wrap long constructs to stay within a reasonable line length (80 characters). Currently, `JavaPrettyPrinter` has no concept of the current column position, so it cannot make wrapping decisions based on line length.

Three failing tests require this feature:
1. **TestLongRecordParameterWrapping** - Records with many parameters should wrap each on its own line
2. **TestLongTernaryWrapping** - Long ternary expressions should wrap `? value` and `: value` on separate lines  
3. **TestLongMethodChainWrapping** - Long method chains should wrap at each `.method()` call (partially working, but needs length awareness for arguments)

## 1. Struct Changes

Add these fields to `JavaPrettyPrinter` in [java_pretty.go](file:///Users/dhamidi/projects/sai/format/java_pretty.go#L11-L20):

```go
type JavaPrettyPrinter struct {
    // ... existing fields ...
    
    column    int  // Current column position (0-indexed)
    maxColumn int  // Maximum line length (default 80)
}
```

Initialize in `NewJavaPrettyPrinter`:

```go
func NewJavaPrettyPrinter(w io.Writer) *JavaPrettyPrinter {
    return &JavaPrettyPrinter{
        // ... existing ...
        column:    0,
        maxColumn: 80,
    }
}
```

## 2. Core Methods to Modify

### 2.1 Update `write()` to track column

Modify the existing [`write()`](file:///Users/dhamidi/projects/sai/format/java_pretty.go#L140-L142) method:

```go
func (p *JavaPrettyPrinter) write(s string) {
    p.w.Write([]byte(s))
    // Update column tracking
    if idx := strings.LastIndex(s, "\n"); idx >= 0 {
        p.column = len(s) - idx - 1
    } else {
        p.column += len(s)
    }
}
```

### 2.2 Update `writeIndent()` to track column

Modify [`writeIndent()`](file:///Users/dhamidi/projects/sai/format/java_pretty.go#L130-L138):

```go
func (p *JavaPrettyPrinter) writeIndent() {
    if !p.atLineStart {
        return
    }
    for i := 0; i < p.indent; i++ {
        p.write(p.indentStr)
    }
    p.atLineStart = false
}
```

The `write()` call already handles column tracking, so no additional changes needed here.

### 2.3 Reset column on newline

Ensure `p.column = 0` when `p.atLineStart = true` is set. Add a helper:

```go
func (p *JavaPrettyPrinter) newline() {
    p.write("\n")
    p.atLineStart = true
    p.column = 0
}
```

Consider refactoring existing `p.write("\n"); p.atLineStart = true` patterns to use this helper.

### 2.4 Add measurement helper

Add a method to measure the printed length of a node without actually printing:

```go
// measureExpr returns the approximate printed length of an expression
func (p *JavaPrettyPrinter) measureExpr(node *parser.Node) int {
    var buf bytes.Buffer
    mp := &JavaPrettyPrinter{
        w:           &buf,
        source:      p.source,
        indentStr:   p.indentStr,
        atLineStart: false,
        column:      0,
        maxColumn:   p.maxColumn,
    }
    mp.printExpr(node)
    return buf.Len()
}

// measureParameters returns the approximate printed length of parameters
func (p *JavaPrettyPrinter) measureParameters(node *parser.Node) int {
    var buf bytes.Buffer
    mp := &JavaPrettyPrinter{
        w:           &buf,
        source:      p.source,
        indentStr:   p.indentStr,
        atLineStart: false,
        column:      0,
        maxColumn:   p.maxColumn,
    }
    mp.printParameters(node)
    return buf.Len()
}
```

### 2.5 Add "would exceed" check

```go
func (p *JavaPrettyPrinter) wouldExceed(additionalChars int) bool {
    return p.column + additionalChars > p.maxColumn
}
```

## 3. Decision Points

### 3.1 Record Parameters (TestLongRecordParameterWrapping)

Location: [`printParameters()`](file:///Users/dhamidi/projects/sai/format/java_pretty_decl.go#L1134-L1159)

**Strategy**: Measure total parameter length. If > maxColumn, wrap each parameter.

```go
func (p *JavaPrettyPrinter) printParameters(node *parser.Node) {
    if node == nil {
        p.write("()")
        return
    }

    // Count parameters
    var params []*parser.Node
    for _, child := range node.Children {
        if child.Kind == parser.KindParameter {
            params = append(params, child)
        } else if child.Kind == parser.KindIdentifier && child.Token != nil {
            params = append(params, child)
        }
    }

    // Measure total length
    totalLen := p.measureParameters(node)
    shouldWrap := p.column + totalLen > p.maxColumn && len(params) > 1

    if shouldWrap {
        p.write("(\n")
        p.atLineStart = true
        p.indent++
        for i, param := range params {
            p.writeIndent()
            if param.Kind == parser.KindParameter {
                p.printParameter(param)
            } else {
                p.write(param.Token.Literal)
            }
            if i < len(params)-1 {
                p.write(",")
            }
            p.write("\n")
            p.atLineStart = true
        }
        p.indent--
        p.writeIndent()
        p.write(")")
    } else {
        // Existing single-line logic
        p.write("(")
        first := true
        for _, child := range node.Children {
            if child.Kind == parser.KindParameter {
                if !first {
                    p.write(", ")
                }
                p.printParameter(child)
                first = false
            } else if child.Kind == parser.KindIdentifier && child.Token != nil {
                if !first {
                    p.write(", ")
                }
                p.write(child.Token.Literal)
                first = false
            }
        }
        p.write(")")
    }
}
```

### 3.2 Ternary Expressions (TestLongTernaryWrapping)

Location: [`printTernaryExpr()`](file:///Users/dhamidi/projects/sai/format/java_pretty_expr.go#L112-L121)

**Strategy**: Measure total ternary length. If > maxColumn, wrap `?` and `:` branches.

```go
func (p *JavaPrettyPrinter) printTernaryExpr(node *parser.Node) {
    children := node.Children
    if len(children) < 3 {
        return
    }
    
    // Measure total length
    condLen := p.measureExpr(children[0])
    trueLen := p.measureExpr(children[1])
    falseLen := p.measureExpr(children[2])
    totalLen := condLen + 3 + trueLen + 3 + falseLen // " ? " and " : "
    
    shouldWrap := p.column + totalLen > p.maxColumn
    
    if shouldWrap {
        p.printExpr(children[0])
        p.write("\n")
        p.atLineStart = true
        p.indent++
        p.writeIndent()
        p.write("? ")
        p.printExpr(children[1])
        p.write("\n")
        p.atLineStart = true
        p.writeIndent()
        p.write(": ")
        p.printExpr(children[2])
        p.indent--
    } else {
        p.printExpr(children[0])
        p.write(" ? ")
        p.printExpr(children[1])
        p.write(" : ")
        p.printExpr(children[2])
    }
}
```

### 3.3 Method Chain Arguments (TestLongMethodChainWrapping)

Location: [`printArguments()`](file:///Users/dhamidi/projects/sai/format/java_pretty_expr.go#L142-L151)

The existing method chain wrapping in `printMethodChain()` handles breaking at `.method()` calls. What's missing is wrapping long argument lists (like the `Collectors.toMap(...)` call).

**Strategy**: If arguments would exceed line length, wrap each argument.

```go
func (p *JavaPrettyPrinter) printArguments(node *parser.Node) {
    if len(node.Children) == 0 {
        return
    }
    
    // Measure total arguments length
    var totalLen int
    for i, child := range node.Children {
        if i > 0 {
            totalLen += 2 // ", "
        }
        totalLen += p.measureExpr(child)
    }
    
    shouldWrap := p.column + totalLen > p.maxColumn && len(node.Children) > 1
    
    if shouldWrap {
        p.write("\n")
        p.atLineStart = true
        p.indent++
        for i, child := range node.Children {
            p.writeIndent()
            p.printExpr(child)
            if i < len(node.Children)-1 {
                p.write(",")
            }
            p.write("\n")
            p.atLineStart = true
        }
        p.indent--
        p.writeIndent()
    } else {
        first := true
        for _, child := range node.Children {
            if !first {
                p.write(", ")
            }
            p.printExpr(child)
            first = false
        }
    }
}
```

## 4. Implementation Order

Implement in this order to minimize risk and enable incremental testing:

### Phase 1: Column Tracking Infrastructure
1. Add `column` and `maxColumn` fields to struct
2. Modify `write()` to track column position
3. Add `newline()` helper (optional, for cleaner code)
4. Add `wouldExceed()` helper
5. **Test**: Verify column tracking works by adding debug prints

### Phase 2: Measurement Helpers
1. Add `measureExpr()` 
2. Add `measureParameters()`
3. **Test**: Unit test the measurement functions

### Phase 3: Record Parameter Wrapping
1. Modify `printParameters()` with wrapping logic
2. **Test**: Run `TestLongRecordParameterWrapping`

### Phase 4: Ternary Expression Wrapping
1. Modify `printTernaryExpr()` with wrapping logic
2. **Test**: Run `TestLongTernaryWrapping`

### Phase 5: Method Chain Argument Wrapping
1. Modify `printArguments()` with wrapping logic
2. **Test**: Run `TestLongMethodChainWrapping`

### Phase 6: Regression Testing
1. Run full test suite: `go test ./format/...`
2. Test against real-world Java files
3. Fix any regressions

## 5. Test Strategy

### Unit Tests for New Helpers

```go
func TestColumnTracking(t *testing.T) {
    var buf bytes.Buffer
    pp := NewJavaPrettyPrinter(&buf)
    
    pp.write("hello")
    if pp.column != 5 {
        t.Errorf("column = %d, want 5", pp.column)
    }
    
    pp.write("\n")
    if pp.column != 0 {
        t.Errorf("column after newline = %d, want 0", pp.column)
    }
    
    pp.write("ab\ncd")
    if pp.column != 2 {
        t.Errorf("column after embedded newline = %d, want 2", pp.column)
    }
}

func TestMeasureExpr(t *testing.T) {
    // Parse a simple expression and verify measurement
}
```

### Existing Tests

The three failing tests already define expected behavior:
- `TestLongRecordParameterWrapping`
- `TestLongTernaryWrapping`  
- `TestLongMethodChainWrapping`

### Manual Testing

Test with real project files:
```bash
sai fmt --check path/to/java/files
```

## 6. Edge Cases to Consider

1. **Already short lines**: Don't wrap if already under limit
2. **Deeply nested expressions**: Measure recursively  
3. **String literals with newlines**: These affect column tracking
4. **Comments**: May affect line length perception
5. **Indentation level**: Higher indent = less space for content
6. **Empty parameters/arguments**: `()` should never wrap
7. **Single parameter/argument**: Don't wrap a single item

## 7. Future Enhancements (Out of Scope)

- Configurable line length via options
- Smart wrapping at specific positions (e.g., after `=`)
- Binary expression wrapping
- Import statement wrapping
- Annotation argument wrapping
