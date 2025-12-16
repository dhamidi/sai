# Skill: Developing and Testing `sai doc`

## When to Use

Use this skill when working on the `sai doc` subcommand, which provides documentation lookup for Java classes, methods, and fields.

## Prerequisites

A **concrete reference Java project** must be specified for testing. The project should have:
- A `src/` directory with Java source files
- A `lib/` directory with JAR dependencies (optional, for testing class file parsing)

Example reference project: `~/projects/java-playground/visor`

## Key Files

- `cmd/sai/cmd_doc.go` - Main implementation of the doc command
- `java/from_source.go` - Parses Java source to ClassModel
- `java/from_classfile.go` - Parses compiled .class files to ClassModel
- `java/model.go` - ClassModel, MethodModel, FieldModel definitions

## Development Workflow

1. Make changes to the relevant files
2. Install the updated binary:
   ```bash
   go install ./...
   ```
3. Test in the reference project directory:
   ```bash
   cd <reference-project>
   sai doc <test-case>
   ```

## Test Cases

When testing `sai doc`, verify these scenarios:

### Source-based lookups (src/)
```bash
sai doc <ClassName>              # Local class by simple name
sai doc com.example.MyClass      # Local class by FQN
sai doc com.example.MyClass.method  # Method lookup
```

### Library lookups (lib/ JARs)
```bash
sai doc com.grack.nanojson.JsonParser  # Class from JAR
sai doc com.grack.nanojson             # Package listing from JAR
```

### JDK lookups
```bash
sai doc java.util.List           # JDK class
sai doc java.util.List.add       # JDK method
sai doc java.util                # JDK package listing
```

### Package exploration
```bash
sai doc                          # List all available packages
sai doc com.example              # List package contents
```

## Verification

After changes, run the pre-commit hook:
```bash
githooks/pre-commit
```
