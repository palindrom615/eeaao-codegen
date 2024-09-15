
print(eeaao_codegen.withConfig())
specs = eeaao_codegen.loadSpecsGlob("openapi", "*.json")
print(specs["petstorev2.json"]["swagger"])
for path in specs:
  print(specs[path]["info"]["title"])
  eeaao_codegen.renderFile("partial_headerFuck", "partial_header.tmpl", specs[path])