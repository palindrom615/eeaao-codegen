def struct_filename(schema):
    return "__generated__/" + schema["title"].lower() + '.go'

def main():
    spec = eeaao_codegen.loadSpecsGlob('json-schema', '*.schema.json')["arrays.schema.json"]
    if spec["$defs"] != None:
        for name, schema in spec["$defs"].items():
            schema["title"] = name
            eeaao_codegen.renderFile(
                struct_filename(schema),
                "struct.go.tmpl",
                schema
            )
    eeaao_codegen.renderFile(
        struct_filename(spec),
        "root.go.tmpl",
        spec
    )
