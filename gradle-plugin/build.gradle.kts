/*
 * This file was generated by the Gradle 'init' task.
 *
 * This generated file contains a sample Java library project to get you started.
 * For more details on building Java & JVM projects, please refer to https://docs.gradle.org/8.9/userguide/building_java_projects.html in the Gradle documentation.
 */

plugins {
    `java-gradle-plugin`
    alias(libs.plugins.plugin.publish)
}

gradlePlugin {
    website = "https://github.com/palindrom615/eeaao-codegen/tree/main/gradle-plugin"
    vcsUrl = "https://github.com/palindrom615/eeaao-codegen.git"
    plugins {
        create("eeaaoCodegenPlugin") {
            id = "dev.palindrom615.eeaao-codegen-plugin"
            implementationClass = "dev.palindrom615.eeaao.EeaaoCodegenPlugin"
            displayName = "EEAAO Codegen Plugin"
            description = "A plugin to generate code from eeaao-codegen template"
            tags = listOf("codegen")
        }
    }
}

repositories {
    // Use Maven Central for resolving dependencies.
    mavenCentral()
    mavenLocal()
}

dependencies {
    // Use JUnit Jupiter for testing.
    testImplementation(libs.junit.jupiter)

    testRuntimeOnly(libs.junit.platform.launcher)

    // This dependency is exported to consumers, that is to say found on their compile classpath.
    api(libs.commons.math3)

    // This dependency is used internally, and not exposed to consumers on their own compile classpath.
    implementation(libs.guava)
    implementation(libs.osdetector)
}

// Apply a specific Java toolchain to ease working on different environments.
java {
    targetCompatibility = JavaVersion.VERSION_1_8
    sourceCompatibility = JavaVersion.VERSION_1_8
}

sourceSets {
    main {
        resources {
            srcDir("src/generated/resources")
        }
    }
}

tasks.named<Test>("test") {
    // Use JUnit Platform for unit tests.
    useJUnitPlatform()
}
