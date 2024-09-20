equivalent openapi-generator command

```bash
openapi-generator generate -i ../spec/petstorev2.json -g kotlin-spring -o build-expected/openapi --additional-properties=annotationLibrary=swagger1,documentationProvider=springfox,requestMappingMode=none,interfaceOnly=true,useTags=true
```

```bash
go run github.com/palindrom615/eeaao-codegen/cmd/eeaao-codegen-cli --codeletdir ./codelet --outdir build --specdir ../spec
```
```bash
diff build build-expected
```
