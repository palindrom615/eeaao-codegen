def main():
    spec = eeaao_codegen.loadSpecsGlob('json-schema', 'person.schema.json')["person.schema.json"]
    eeaao_codegen.renderFile(
        "__generated__/" + spec["title"].lower() + '.go',
        "struct.go.tmpl",
        spec
    )
