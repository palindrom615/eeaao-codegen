def main():
    spec = eeaao_codegen.loadSpecsGlob('json', 'petstorev2.json')["petstorev2.json"]
    tags = spec["tags"]
    for t in tags:
        print(t)
        eeaao_codegen.renderFile(
            t["name"].title() + "Api.kt",
            "Api.kt.tmpl",
            {"spec": spec, "tag": t}
        )