package render

import (
	spec2 "github.com/palindrom615/eeaao-codegen/spec"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func Render(c *Config) {
	// Read JSON specFile

	s := spec2.SpecDir{Dir: specDir}

	// render `render.tmpl` with spec data
	tmplData, err := os.ReadFile(filepath.Join(codeletDir, "render.tmpl"))
	tmpl := template.New("render.tmpl")
	tmpl.Funcs(
		template.FuncMap{
			"loadSpecsGlob": s.LoadSpecsGlob,
		},
	)

	tmpl, err = tmpl.Parse(string(tmplData))
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create output file
	outFilePath := filepath.Join(outDir, info.Name())
	outFile, err := os.Create(outFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer outFile.Close()

	// Render template with spec data
	err = tmpl.Execute(outFile, struct{}{})
	if err != nil {
		log.Fatal(err)
		return
	}
}
