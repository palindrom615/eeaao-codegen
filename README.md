# eeaao-codegen

everything-everywhere-all-at-once code generator

eeaao-codegen is a highly customizable code generation tool that empowers developers to integrate schema-defined
interfaces (IDLs) into their projects seamlessly. Unlike traditional code generators, it does not impose predefined
generation logic. Instead, it allows you to define your own rules and templates using its `codelet`s.

## Motivation

When you use a language-agnostic Interface Definition language (IDL) in your project --such as gRPC, Thrift, GraphQL,
OpenAPI specification, AsyncAPI for pub-sub systems, SQL or prisma and so on-- more often than not, manually writing
schema-to-language
mappings or implementing a custom parser is a tedious and error-prone task. Code generators are a common solution to
this
problem. They automatically generate code from schema files, which can save you a lot of time and effort.

Typical code generators take schema files (e.g. openAPI specification, GraphQL schemas, Protobuf definitions...) as
input and
produce code in one of their predefined supported language. It works like a charm, until it doesn't. Because the
codegen lacks awareness of your codebase, they often have to make trade-offs between the following:

- to keep generated code be generic, codegen tries to cover wide variation of versions of language, dependencies and
  environments. sometimes it abuses reflection of language, or full of boilerplate codes that has overhead on runtime.
- to be more idiomatic, codegen should make bold assumptions over source code, or provide enormous number of options and
  flags to tune the code generation. You have to read the documentation and try-and-error to find the best fit for
  your project. On out of luck, you might find them unusable without substantial customization.

To avoid both situation, `eeaao-codegen` does **not** provide any predefined generation logic. Instead, it is a tool
that
helps you write your own code generation logic, organized into units called `codelet`.

## Usage

Assume you have openAPI specs in the `openapi` directory, a codelet defined in the `./codelet` directory, and
you want to generate code in the `__generated__` directory.

```
eeaao-codegen-cli --specdir=./openapi --codeletdir=./codelet --outputdir=./__generated__
```

## Codelet

`codelet` is a unit of code generation, which is a directory that contains below:

- `render.star` : an entrypoint of codelet
- `templates` : a directory that contains template files. `render.star` can use these files to render the code.
- `values.json`: a global key-value json file that can contains anything (optional)

`render.start` is a starlark script that defines the code generation process. It must have a `main()` function.
`eeaao-codegen`
calls `main()` function on code generation.

```starlark
def main():
    pass
```

of course, this is not useful. `eeaao_codegen` starlark module is provided to generate code. Belows are the functions
that you can use in `render.star`:

| Function                                                              | Description                                                                                          |
|-----------------------------------------------------------------------|------------------------------------------------------------------------------------------------------|
| `getPlugin(pluginName: str): EeaaoPlugin`                             | load plugin of given name and return loaded plugin                                                   |
| `renderFile(filePath: str, templatePath: str, data: Any): str`        | render `templates/{templatePath}` file with given `data` on `filePath` and return rendered file path |
| `loadValues(filePaths: List[str], templatePath: str, data: Any): Any` | load `values.json` and return loaded value                                                           |

`EeaaoPlugin` is an interface that defines below functions:

| Function                       | Description                                              |
|--------------------------------|----------------------------------------------------------|
| `loadSpec(filepath: str): Any` | load spec file of `filepath` and return loaded spec data |
| `loadSpecUrl(url: str): Any`   | load spec file of `url` and return loaded spec data      |

using these functions, you can write your own codelet. below is the example codelet that generates a simple `hello world`

```starlark
def main():
  # load openapi specs
  openApiPlugin = eeaao_codegen.getPlugin("openapi")
  spec = openApiPlugin.loadSpecFile("../spec/petstore.json")

  # render a file
  eeaao_codegen.renderFile("petstoreApi.js", "Api.js.tmpl", spec)
```

The codelet above loads `{specdir}/petstore.yaml` file, reads template file
`{codeletdir}/templates/Api.js.tmpl`, and
renders the template with loaded spec data. The rendered file will be saved as `{outputdir}/petstoreApi.js`.

The template file used in `renderFile` must be valid [go template](https://pkg.go.dev/text/template) files. below is the
example template file that generates a simple javascript code from openapi spec.

```gotemplate
// Api.js.tmpl
export const {{ .info.title | camelcase }}Api = {
  {{- range $path, $pathSpec := .paths }}
    {{- range $method, $methodSpec := $pathSpec }}
      {{ $methodSpec.operationId }}: (data) => {
        return fetch("{{ $path }}", {
          method: "{{ $method }}",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify(data)
        })
      },
    {{- end }}
  {{- end }}
}

```

you can find above example and more in [here](./example/openapi-v3/js)

Inside the template, `eeaao-codegen` exposes some variables and functions. You can
use [sprig](https://masterminds.github.io/sprig/)
and functions defined in [`HelperFuncs`](https://pkg.go.dev/github.com/palindrom615/eeaao-codegen)

## Plugins

`eeaao-codegen` itself does not have any knowledge about schema language. It delegates the schema parsing to
its internal plugins. Writing your own schema plugin is working in progress.

Currently, `eeaao-codegen` supports below plugins:

- json
- yaml
- openapi
- proto
