package eeaao_codegen

import (
	"github.com/Masterminds/sprig"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"maps"
	"os"
	"path/filepath"
	"text/template"
)

func (a *App) makeFuncmap() template.FuncMap {
	funcmap := template.FuncMap{
		"loadSpecsGlob": a.loadSpecsGlob,
		"renderFile":    a.renderFile,
		"withConfig":    func() map[string]any { return a.Conf },
	}
	maps.Copy(funcmap, sprig.FuncMap())
	return funcmap
}

func (a *App) renderFile(filePath string, templatePath string, data any) string {
	if !filepath.IsLocal(filePath) {
		panic("filePath must be local")
	}
	if !filepath.IsLocal(templatePath) {
		panic("templatePath must be local")
	}
	dst := filepath.Join(a.OutDir, filePath)
	os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	dstFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	a.tmpl.ExecuteTemplate(dstFile, templatePath, data)

	return dst
}

func (a *App) loadSpecsGlob(pluginName string, glob string) (res []plugin.SpecData) {
	p := plugin.GetPlugin(pluginName)
	matches, err := filepath.Glob(filepath.Join(a.specDir, glob))
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
