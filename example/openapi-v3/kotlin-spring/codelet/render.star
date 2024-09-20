def getTags(spec):
    tags = spec.get("tags", [])
    tagNames = [t["name"] for t in tags]
    for path, pathSpec in spec["paths"].items():
        for method, methodSpec in pathSpec.items():
            for tag in methodSpec["tags"]:
                if tag not in tagNames:
                    tagNames.append(tag)
                    tags.append({"name": tag})
    return tags


def main():
    spec = eeaao_codegen.loadSpecsGlob('json', 'petstore.json')["petstore.json"]
    tags = getTags(spec)
    for t in tags:
        eeaao_codegen.renderFile(
            t["name"].title() + "Api.kt",
            "Api.kt.tmpl",
            {"spec": spec, "tag": t}
        )