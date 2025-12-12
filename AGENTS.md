# Sai Development

## Building and Testing

```bash
go build ./...          # Build all packages
go test ./...           # Run all tests
```

## Tools

### Sai

The main Java toolchain. See README.md for usage.

### Ahi

**Ahi** (Estonian for "oven") provides development tools for sai. It's a companion tool for testing and working with language infrastructure.

#### Commands

- `ahi ebnf check <file>` — Parse and verify an EBNF grammar file
  - `--start <production>` — Specify start production for verification
- `ahi ebnf lex <grammar>` — Tokenize input based on an EBNF grammar, emitting tokens with source positions

#### Purpose

Ahi is used for:
- Validating EBNF grammar files used in language tooling
- Testing lexical analysis based on grammar definitions
- Debugging parser and lexer development

## Package Structure

- `cmd/sai/` — Main sai CLI
- `cmd/ahi/` — Ahi development tools CLI
- `ebnflex/` — EBNF-based lexer library
- `java/` — Java parsing and class model
- `java/parser/` — Java lexer and parser
- `classfile/` — Java classfile parsing
- `format/` — Java code formatting
- `pom/` — Maven POM parsing
- `ui/` — Web UI for code exploration
