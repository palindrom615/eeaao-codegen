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
    values = eeaao_codegen.loadValues()
    openapiPlugin = eeaao_codegen.getPlugin("openapi")
    predefinedComponentSchema = values.get("predefinedComponentSchema")
    baseJavaPackage = values.get("javaPackage")
    baseDirectory = baseJavaPackage.replace(".", "/")
    def genApi(api):
        directory = baseDirectory + "/" + api
        spec = openapiPlugin.loadSpecFile("../spec/" + api + '.json')
        schemas = spec.get("components", {}).get("schemas", {})
        for name, schema in schemas.items():
            if name in predefinedComponentSchema:
                continue
            className = name
            javaPackage = baseJavaPackage + "." + api + ".schema"
            renderFile = directory + "/schema/" + className + ".kt"
            eeaao_codegen.renderFile(
                directory + "/schema/" + className + ".kt",
                "Schema.kt.tmpl",
                {"className": className, "javaPackage": javaPackage, "schema": schema}
            )
            print(name)
        tags = getTags(spec)
        for t in tags:
            className = t["name"].title().replace("/", "_") + "Api"
            javaPackage = baseJavaPackage + "." + api
            eeaao_codegen.renderFile(
                directory + "/" + className + ".kt",
                "Api.kt.tmpl",
                {"className": className, "javaPackage": javaPackage, "spec": spec, "tag": t, }
            )
    genApi("petstore")
