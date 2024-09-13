package eeaao_codegen

import (
	"github.com/Masterminds/sprig"
	"log"
	"maps"
	"os"
	"path/filepath"
	"text/template"
)

func Render(c *App) string {
	// Read JSON specFile

	os.Mkdir(c.OutDir, os.ModePerm)

	// render `render.tmpl` with spec data
	tmplData, err := os.ReadFile(filepath.Join(c.CodeletDir, "render.tmpl"))
	tmpl := template.New("render.tmpl")
	funcmap := template.FuncMap{
		"loadSpecsGlob": c.LoadSpecsGlob,
		"renderFile":    c.renderFile,
		"withConfig":    func() map[string]any { return c.Conf },
	}
	maps.Copy(funcmap, sprig.FuncMap())
	tmpl.Funcs(funcmap)

	tmpl, err = tmpl.Parse(string(tmplData))
	if err != nil {
		log.Fatal(err)
		return ""
	}

	// Create output file
	outFilePath := filepath.Join(c.OutDir, "render")
	outFile, err := os.Create(outFilePath)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer outFile.Close()

	// Render template with spec data
	err = tmpl.Execute(outFile, struct{}{})
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return outFilePath
}
