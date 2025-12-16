# Sai - A Toasty Java Toolchain

<p align="center">
  <img src="./docs/sai-logo.png" width="256" height="256">
</p>

<p align="center">
  <em>Parse, format, and navigate Java code with ease</em>
</p>

---

Sai is a lightweight Java toolchain written in Go. It parses `.java` and `.class` files, builds navigable code models, and includes an LSP server for editor integration.

## Features

- **Parse** — Analyze `.java` source files and `.class` bytecode
- **Dump** — Extract class models in JSON, Java, or line-based formats  
- **Scan** — Batch process directories, JARs, and ZIP archives

- **Web UI** — Interactive browser-based code explorer

## Installation

```bash
go install github.com/dhamidi/sai/cmd/sai@latest
```

## Usage

### Parse a file

```bash
# Parse Java source
sai parse MyClass.java -f json

# Parse compiled bytecode
sai parse MyClass.class -f java
```

### Dump class models

```bash
# Output as JSON
sai dump MyClass.java -f json

# Output as Java-like declaration
sai dump MyClass.java -f java

# Output as line format (default)
sai dump MyClass.java
```

### Format Java source

```bash
# Format to stdout
sai fmt MyClass.java

# Format from stdin
cat MyClass.java | sai fmt

# Overwrite file in place
sai fmt -w MyClass.java
```

### Scan a codebase

```bash
# Scan a directory
sai scan ./src

# Scan a JAR file
sai scan library.jar

# With custom timeout
sai scan ./src -t 30s
```

### Launch the web UI

```bash
sai ui -a :8080
```

Then open http://localhost:8080 in your browser.

## Project Structure

A sai-managed Java project follows this structure:

```
myproject/
├── src/
│   └── myproject/
│       ├── core/                # Logic module
│       │   ├── module-info.java
│       │   └── *.java
│       └── main/                # Entrypoints module
│           ├── module-info.java
│           └── *.java
├── lib/                         # Dependencies (JAR files)
├── out/                         # Compiled classes
├── mlib/                        # Modular JARs for jlink/jpackage
└── dist/                        # Distribution output
```

### lib/ vs mlib/

- **`lib/`** contains downloaded dependencies and is used during compilation (`javac -p lib`)
- **`mlib/`** is a build artifact containing your project's modules packaged as JARs plus copies of dependencies from `lib/`. It serves as the module path for `jlink` and `jpackage`.

## License

MIT
