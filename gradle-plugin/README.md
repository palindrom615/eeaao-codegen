# eeaao-codegen Gradle Plugin

This gradle plugin is a wrapper around the [eeaao-codegen](../README.md).

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

By default the plugin register a task named `eeaaoCodegen`. You can configure the task by using `eeaaoCodegen` extension.

```kotlin
eeaaoCodegen {
    specDir = "src/main/resources/spec"
    codeletDir = "src/main/resources/codelet"
    outDir = "src/__generated__/java"
}
```

internally, each of the option is mapped to a command line argument. The following is the list of options and their corresponding command line arguments:

| Option | Command Line Argument | Description |
|--------|--------------------|-------------|
| `specDir` | `--specdir` | The directory containing the specification files. |
| `codeletDir` | `--codeletdir` | The directory containing the codelet files. |
| `outDir` | `--outdir` | The directory where the generated code will be written. |

for further description about each options, check [main eeaao-codegen project](../README.md)

#### further configurations

By default, the plugin itself does not add `outDir` in your project's srcSets or make dependency on `build` task.
This is intentional not to mess up your project's layout or task dependency graph. 

If you want to add `outDir` to your project's srcSets, you can do it manually. 

```kotlin
// add generated code to sourceSets
// see https://docs.gradle.org/current/dsl/org.gradle.api.tasks.SourceSet.html
sourceSets {
    main {
        java {
            srcDir("build/__generated__/java")
        }
    }
}

// make compileJava task depend on the codegen task
// see https://docs.gradle.org/current/userguide/controlling_task_execution.html#sec:adding_dependencies_to_tasks
tasks["compileJava"].dependsOn(tasks.withType(GenerateEeaaoTask::class.java))
```

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
