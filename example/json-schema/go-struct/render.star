def validStructName(name):
    return name.title().replace(" ", "").replace("-", "")

def renderSpec(specs, specFile):
    spec = specs[specFile]
    location = "__generated__/" + specFile.replace('.schema.json', '')

    if "$defs" in spec:
        for name, schema in spec["$defs"].items():
            schema["title"] = name
            eeaao_codegen.renderFile(
                location + "/" + validStructName(name) + ".go",
                "root.go.tmpl",
                schema
            )
    eeaao_codegen.renderFile(
        location + "/" + validStructName(spec['title'])+ ".go",
        "root.go.tmpl",
        spec
    )

def main():
    specs = eeaao_codegen.loadSpecsGlob('json', '*.schema.json')
    renderSpec(specs, "person.schema.json")
    renderSpec(specs, "arrays.schema.json")
    renderSpec(specs, "regex-pattern.schema.json")
