equivalent openapi-generator command

```bash
openapi-generator generate -i ../spec/petstore.json -g kotlin-spring -o build-expected/openapi --additional-properties=annotationLibrary=swagger2,documentationProvider=springdoc,requestMappingMode=none,interfaceOnly=true,useTags=true
```


```bash
go run github.com/palindrom615/eeaao-codegen/cmd/eeaao-codegen-cli --codeletdir ./codelet --outdir build
```
```bash
diff build build-expected
```
