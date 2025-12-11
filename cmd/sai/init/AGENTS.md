# Sai Java Project

This project uses **sai** instead of Maven or Gradle for Java development.

## Commands

### Adding Dependencies

```bash
sai add <groupId:artifactId:version>
# Example: sai add com.google.guava:guava:32.1.3-jre
```

Dependencies are downloaded to `lib/`.

### Getting the Classpath

```bash
sai classpath
```

Returns a colon-separated list of JAR paths for use with `javac` or `java`.

### Formatting

```bash
sai fmt src/Main.java        # Print formatted output
sai fmt -w src/Main.java     # Overwrite file in place
```

## Project Structure

```
src/
  {{PROJECT_ID}}/
    core/                    # Logic module
      module-info.java
    main/                    # Entrypoints module
      module-info.java
lib/                         # Dependencies (JARs)
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
- Use `sai add` to add dependencies
- Use `sai classpath` to get the classpath for compilation/execution
- Use `sai fmt -w` to format Java files
- The `lib/` directory contains all JAR dependencies
