# EEAAO Codegen Gradle Plugin

This gradle plugin is a wrapper around the [EEAAO Codegen](../README.md). It provides a way to run the codegen tool as a gradle task.

## Usage

### Applying the Plugin

```kotlin
// build.gradle.kts
plugins {
    id("dev.palindrom615.eeaao-codegen-plugin") version "0.1.0"
}
```

```groovy
// build.gradle
plugins {
    id 'dev.palindrom615.eeaao-codegen-plugin' version '0.1.0'
}
```

### Configuration

```kotlin
eeaaoCodegen {
    specDir = "src/main/resources/spec"
    codeletDir = "src/main/resources/codelet"
    outDir = "src/generated/kotlin"
}
```

internally, each of the option is mapped to a command line argument. The following is the list of options and their corresponding command line arguments:

| Option | Command Line Argument | Description |
|--------|--------------------|-------------|
| `specDir` | `--specdir` | The directory containing the specification files. |
| `codeletDir` | `--codeletdir` | The directory containing the codelet files. |
| `outDir` | `--outdir` | The directory where the generated code will be written. |

for further description about each options, check [main eeaao-codegen project](../README.md)

### Running the Plugin

```bash
./gradlew eeaaoCodegen
```

### Register additional tasks

Instead of built-in task registered by the plugin itself, you can register additional tasks by yourself.

```kotlin
tasks.register<GenerateEeaaoTask>("AnotherCodegenTask") {
    eeaaoCodegen {
        specDir = "src/main/resources/spec"
        codeletDir = "src/main/resources/anotherCodelet"
        outDir = "src/generated/kotlin"
    }
}
```

for more examples, check [../gradle-plugin-examples](../gradle-plugin-examples) directory.
