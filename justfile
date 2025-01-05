target := "github.com/palindrom615/eeaao-codegen/cmd/eeaao-codegen-cli"

os_arch := "linux/amd64 linux/386 linux/arm linux/arm64 darwin/amd64 darwin/arm64 " + \
    "windows/amd64 windows/386 windows/arm windows/arm64 freebsd/amd64 freebsd/386 freebsd/arm openbsd/amd64 " + \
    "openbsd/386 openbsd/arm netbsd/amd64 netbsd/386 netbsd/arm solaris/amd64"

default:
    just --list

build:
    go build -o eeaao-codegen-cli {{target}}

build-all BUILDDIR:
    #!/usr/bin/env bash
    set -euxo pipefail
    for os_arch in {{os_arch}}; do \
        os=$(echo $os_arch | cut -d/ -f1); \
        arch=$(echo $os_arch | cut -d/ -f2); \
        just build-os-arch $os $arch {{BUILDDIR}}; \
    done

build-os-arch os arch BUILDDIR="build":
    #!/usr/bin/env bash
    set -euxo pipefail
    output_path="{{BUILDDIR}}/eeaao-codegen-cli-{{os}}-{{arch}}"
    GOOS={{os}} GOARCH={{arch}} go build -o $output_path {{target}}

gradle_opts := "--console=plain --no-daemon --quiet"

[group("gradle-plugin")]
build-gradle-plugin:
    #!/usr/bin/env bash
    set -euxo pipefail
    gradledir="./gradle-plugin/src/generated/resources/dev/palindrom615/eeaao"
    just build-all $gradledir
    cd ./gradle-plugin
    ./gradlew build {{gradle_opts}}

[group("gradle-plugin")]
publish-gradle-plugin: build-gradle-plugin
    #!/usr/bin/env bash
    set -euxo pipefail
    cd ./gradle-plugin
    ./gradlew publishPlugins {{gradle_opts}}
