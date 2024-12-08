def main():
    values = eeaao_codegen.loadValues()
    print(values)
    specs = eeaao_codegen.loadSpecsGlob('json', '**')
    spec = specs["project.json"]
    eeaao_codegen.renderFile(
        values["javaPackage"].replace(".", "/") + "/ProjectInfoProvider.java",
        "ProjectInfoProvider.java.tmpl",
        values | spec,
    )
