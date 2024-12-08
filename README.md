# eeaao-codegen

everything-everywhere-all-at-once code generator

Suppose that you adopt a language-agnostic IDL(Interface Definition language) for your project. (`grpc`, `thrift`, or 
`graphql` for your network protocol, `OpenAPI` for RESTful api and `AsyncAPI` for pub-sub, SQL or `prisma` for database
schema, and so on) When using them on your code, probably writing schema-language mappings by hand or implementing your
own parser might be an option. You may consider using `protoc`, `openapi-generator` or similar code generator to use
them on your codebase.

Typical code generators takes schema files (e.g. openAPI specs, graphql schemas, protobuf definitions) as input and
produce code in one of the predefined supported language. It works like a charm, until it doesn't. Because the 
codegen has no knowledge about your codebase, they often have to do trade-off between below:
- to keep generated code be generic, codegen tries to cover wide variation of versions of language, dependencies and 
    environments. sometimes it abuse reflection of language, or full of boilerplate codes that has overhead on runtime.
- to be more idiomatic, codegen should make bold assumptions over source code, or provide enormous number of options and 
    flags to tune the code generation. You have to read the documentation and try-and-error to find the best fit for 
    your project. On out of luck, you just could not use them without substantial hacking.

To avoid both situation, `eeaao-codegen` does **not** provide any pre-defined generation logic. Instead, it is a tool that
helps you to write your own code generation logic, called `codelet`.

## Usage

Suppose you have a bunch of openAPI specs in `openapi` directory, you wrote your codelet in `./codelet` directory, and
you want to generate code in `__generated__` directory.

```
eeaao-codegen-cli --specdir=./openapi --codeletdir=./codelet --outputdir=./__generated__
```

## Codelet

`codelet` is a unit of code generation, which is a directory that contains below:
- `render.star` : an entrypoint of codelet
- `templates` : a directory that contains template files. `render.star` can use these files to render the code.
- `values.json`: a global key-value json file that can contains anything (optional)

`render.start` is a starlark script that defines how to render the code. It must have `main()` function. `eeaao-codegen`
calls `main()` function on code generation. and that is the only requirements! so below is the simplest codelet that does
nothing.

```starlark
def main():
    pass
```

of course, this is not useful. `eeaao-codegen` provides `eeaao_codegen` module to generate code. Belows are the functions
that you can use in `render.star`:

| Function                                                              | Description                                                                                          |
|-----------------------------------------------------------------------|------------------------------------------------------------------------------------------------------|
| `loadSpecFile(pluginName: str, filepath: str): Any`                   | load spec file of `filepath` with plugin of given name and return loaded spec data                   |
| `loadSpecsGlob(pluginName: str, glob: str): Any`                      | load spec files with globa patterns with plugin of given name and return loaded spec data            |
| `renderFile(filePath: str, templatePath: str, data: Any): str`        | render `templates/{templatePath}` file with given `data` on `filePath` and return rendered file path |
| `loadValues(filePaths: List[str], templatePath: str, data: Any): Any` | load `values.json` and return loaded value                                                           |

using these functions, you can write your own codelet. below is the example codelet that generates a simple `hello world`

```starlark

def main():
    # load openapi specs
    specs = eeaao_codegen.loadSpecFile("openapi", "petstore.yaml")

    # render a file
    eeaao_codegen.renderFile("petstoreApi.js", "petstoreApi.js.tmpl", specs)
```
above codelet loads `{specdir}/petstore.yaml` file, read template file `{codeletdir}/templates/petstoreApi.js.tmpl`, and
render the template with loaded spec data. the rendered file will be saved as `{outputdir}/petstoreApi.js`.

templates must be valid go template files. below is the example template file that generates a simple javascript code from openapi spec.

```gotemplate
// petstoreApi.js.tmpl
export const petstoreApi = {
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

## Plugins

`eeaao-codegen` itself does not have any knowledge about schema language. It delegates the schema parsing to
its internal plugins. Writing your own schema plugin is working in progress.

Currently, `eeaao-codegen` supports below plugins:
- json
- yaml
- openapi
- proto
