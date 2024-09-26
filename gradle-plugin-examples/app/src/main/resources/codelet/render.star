def main():
    values = eeaao_codegen.loadValues()
    print(values)
    spec = eeaao_codegen.loadSpecFile('json', 'project.json')
    eeaao_codegen.renderFile(
        values["javaPackage"].replace(".", "/") + "/ProjectInfoProvider.kt",
        "ProjectInfoProvider.kt.tmpl",
        values | spec,
    )
