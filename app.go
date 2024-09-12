package eeaao_codegen

import (
	"github.com/palindrom615/eeaao-codegen/plugin"
	"os"
	"path/filepath"
	"text/template"
)

type App struct {
	SpecDir    string
	OutDir     string
	CodeletDir string
}

func (a *App) renderFile(filePath string, templatePath string, data any) string {
	dst := filepath.Join(a.OutDir, filePath)
	dstFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	tmplStr, err := os.ReadFile(filepath.Join(a.CodeletDir, "templates", templatePath))
	tmpl := template.New(filepath.Base(filepath.Join(a.CodeletDir, "templates", templatePath)))
	tmpl.Parse(string(tmplStr))

	tmpl.Execute(dstFile, data)

	return dst
}

func (a *App) LoadSpecsGlob(pluginName string, glob string) (res []plugin.SpecData) {
	p := plugin.GetPlugin(pluginName)
	matches, err := filepath.Glob(filepath.Join(a.SpecDir, glob))
	if err != nil {
		return nil
	}
	for _, match := range matches {
		if doc := p.LoadSpecFile(match); doc != nil {
			res = append(res, doc)
		}
	}
	return res
}
