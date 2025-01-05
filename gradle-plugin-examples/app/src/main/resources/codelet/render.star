def main():
    values = eeaao_codegen.loadValues()
    print(values)
    jsonPlugin = eeaao_codegen.getPlugin("json")
    spec = jsonPlugin.loadSpecFile('src/main/resources/spec/project.json')
    eeaao_codegen.renderFile(
        values["javaPackage"].replace(".", "/") + "/ProjectInfoProvider.java",
        "ProjectInfoProvider.java.tmpl",
        values | spec,
    )
