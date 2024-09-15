package eeaao_codegen

import (
	"encoding/json"
	"github.com/palindrom615/eeaao-codegen/plugin"
	"log"
	"os"
	"path/filepath"
)

func (a *App) RenderFile(filePath string, templatePath string, data any) string {
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

func (a *App) LoadSpecsGlob(pluginName string, glob string) (map[string]string, error) {
	p := plugin.GetPlugin(pluginName)
	matches, err := filepath.Glob(filepath.Join(a.specDir, glob))
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for _, match := range matches {
		doc, err := p.LoadSpecFile(match)
		if err != nil {
			log.Printf("Error loading spec file '%s': %v\n", match, err)
			continue
		}
		docStr, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshaling spec file '%s': %v\n", match, err)
			continue
		}
		p, _ := filepath.Rel(a.specDir, match)
		res[p] = string(docStr)
	}
	return res, nil
}

func (a *App) WithConfig() map[string]any {
	return a.Conf
}
