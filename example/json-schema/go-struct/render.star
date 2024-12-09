def struct_filename(schema):
    return "__generated__/" + schema["title"].lower() + '.go'

def renderSpec(spec):
    if "$defs" in spec:
        for name, schema in spec["$defs"].items():
            schema["title"] = name
            eeaao_codegen.renderFile(
                struct_filename(schema),
                "root.go.tmpl",
                schema
            )
    eeaao_codegen.renderFile(
        struct_filename(spec),
        "root.go.tmpl",
        spec
    )

def main():
    specs = eeaao_codegen.loadSpecsGlob('json', '*.schema.json')
    renderSpec(specs["person.schema.json"])
    renderSpec(specs["arrays.schema.json"])
