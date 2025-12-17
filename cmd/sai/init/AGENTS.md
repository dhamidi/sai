# Sai Java Project

This project uses **sai** instead of Maven or Gradle for Java development.

## Commands

### Adding Dependencies

```bash
sai add <groupId:artifactId:version>
# Example: sai add com.google.guava:guava:32.1.3-jre
```

Dependencies are downloaded to `lib/`.

### Listing and Searching Libraries

```bash
sai libs list                # List all dependencies in lib/
sai libs search <query>      # Search Maven Central for libraries
sai libs search -c <class>   # Search by class name
```

### Getting the Classpath

```bash
sai classpath
```

Returns a colon-separated list of JAR paths for use with `javac` or `java`.

### Compiling

```bash
sai compile
```

Compiles the application.

### Running

```bash
sai run
```

Runs the Cli.java entrypoint.

### Formatting

```bash
sai fmt src/Main.java        # Print formatted output
sai fmt -w src/Main.java     # Overwrite file in place
```

### Documentation

```bash
sai doc java.util            # List classes and subpackages in a package
sai doc java.util.List       # Show documentation for a class
sai doc java.util.List.add   # Show documentation for a method
sai doc String               # Short for java.lang.String
```

Use `sai doc` to explore the JDK API and look up method signatures and javadoc.

## Project Structure

```
src/
  {{PROJECT_ID}}/
    core/                    # Logic module
      module-info.java
    main/                    # Entrypoints module
      module-info.java
lib/                         # Dependencies (JARs) - JUnit 5 pre-installed
out/                         # Compiled classes
mlib/                        # Modular JARs for jlink/jpackage
```

## Build Commands

### Compile

```sh
mkdir -p out/{{PROJECT_ID}}.core out/{{PROJECT_ID}}.main mlib

# Compile core
javac -p lib -d out/{{PROJECT_ID}}.core \
  src/{{PROJECT_ID}}/core/module-info.java \
  src/{{PROJECT_ID}}/core/*.java

# Compile main
javac -p lib:out -d out/{{PROJECT_ID}}.main \
  src/{{PROJECT_ID}}/main/module-info.java \
  src/{{PROJECT_ID}}/main/*.java
```

### Create Modular JARs

```sh
# Core module
jar --create --file=mlib/{{PROJECT_ID}}.core.jar -C out/{{PROJECT_ID}}.core .

# Main module with entrypoint
jar --create --file=mlib/{{PROJECT_ID}}.main.jar --main-class={{PROJECT_ID}}.main.Main -C out/{{PROJECT_ID}}.main .

# Include dependencies
cp lib/*.jar mlib/
```

### Run During Development

```sh
# Using module path (after creating modular JARs)
java -p mlib -m {{PROJECT_ID}}.main

# Using exploded classes
java -p lib:out -m {{PROJECT_ID}}.main/{{PROJECT_ID}}.main.Main
```

## Notes for AI Agents

- Do NOT use Maven (`mvn`) or Gradle (`gradle`, `./gradlew`)
- **MUST use `sai compile` to compile** - do not invoke `javac` directly
- **MUST use `sai run` to run the application** - do not invoke `java` directly
- Use `sai add` to add dependencies
- Use `sai classpath` to get the classpath for compilation/execution
- Use `sai fmt -w` to format Java files
- The `lib/` directory contains all JAR dependencies
- Use `-v` flag with `sai compile -v` or `sai run -v` to see exact commands being executed
