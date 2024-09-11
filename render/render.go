package render

import (
	"fmt"
	spec2 "github.com/palindrom615/eeaao-codegen/spec"
	"os"
	"path/filepath"
	"text/template"
)

func Render(specDir string, codeletDir string, outDir string) {
	// Read JSON specFile

	s := spec2.SpecDir{Dir: specDir}

	// Iterate over templates in codeletDir
	err := filepath.Walk(codeletDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Read template file
			tmplData, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Parse template
			tmpl := template.New(filepath.Base(path))
			tmpl.Funcs(
				template.FuncMap{
					"loadSpecsGlob": s.LoadSpecsGlob,
				},
			)

			tmpl, err = tmpl.Parse(string(tmplData))
			if err != nil {
				return err
			}

			// Create output file
			outFilePath := filepath.Join(outDir, info.Name())
			outFile, err := os.Create(outFilePath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Render template with spec data
			err = tmpl.Execute(outFile, struct{}{})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error processing templates: %v\n", err)
	}
}
