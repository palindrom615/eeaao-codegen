def main():
    # load openapi specs
    openApiPlugin = eeaao_codegen.getPlugin("openapi")
    spec = openApiPlugin.loadSpecFile("../spec/petstore.json")

    # render a file
    eeaao_codegen.renderFile("petstoreApi.js", "Api.js.tmpl", spec)
