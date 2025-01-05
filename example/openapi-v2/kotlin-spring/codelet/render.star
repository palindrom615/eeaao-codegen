def main():
    jsonPlugin = eeaao_codegen.getPlugin("json")
    spec = jsonPlugin.loadSpecFile('petstorev2.json')
    tags = spec["tags"]
    for t in tags:
        print(t)
        eeaao_codegen.renderFile(
            t["name"].title() + "Api.kt",
            "Api.kt.tmpl",
            {"spec": spec, "tag": t}
        )