package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func Render(specFile string, codeletDir string, outDir string) {
	// Read JSON specFile
	specData, err := os.ReadFile(specFile)
	if err != nil {
		fmt.Printf("Error reading spec file: %v\n", err)
		return
	}

	// Parse JSON into Spec struct
	var spec map[string]interface{}
	err = json.Unmarshal(specData, &spec)
	if err != nil {
		fmt.Printf("Error parsing spec file: %v\n", err)
		return
	}

	// Iterate over templates in codeletDir
	err = filepath.Walk(codeletDir, func(path string, info os.FileInfo, err error) error {
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
			tmpl, err := template.New(filepath.Base(path)).Parse(string(tmplData))
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
			err = tmpl.Execute(outFile, spec)
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
